// Package cli implements a credential-free presentation and file boundary over
// the immutable recovery contract. Local expected scope is not authorization.
package cli

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	recovery "github.com/Storm-Harbor/storm-lib-recovery-contracts"
)

const ReportVersion = "cli.stormharbor.io/v1alpha1"
const SDKVersion = "v0.1.0-alpha.1"

const (
	ExitOK        = 0
	ExitInvalid   = 2
	ExitInput     = 3
	ExitInternal  = 70
	ExitUsage     = 64
	ExitOutput    = 74
	ExitTimeout   = 124
	ExitCancelled = 130
)

// Options is internal application wiring, never an authority supplied by YAML.
type Options struct {
	Version  string
	Revision string
	Timeout  time.Duration
	Load     func(context.Context, string) ([]byte, error)
}

type Failure struct {
	Code      string            `json:"code"`
	Module    string            `json:"module"`
	Category  string            `json:"category"`
	Key       string            `json:"errorKey"`
	Arguments map[string]string `json:"arguments,omitempty"`
}

func (f *Failure) Error() string { return f.Code }

type Build struct {
	Version    string `json:"version"`
	Revision   string `json:"revision"`
	GoVersion  string `json:"goVersion"`
	SDKVersion string `json:"sdkVersion"`
}

// Report never includes supplied configuration, flag values or filesystem paths.
type Report struct {
	APIVersion      string   `json:"apiVersion"`
	Command         string   `json:"command"`
	Status          string   `json:"status"`
	Boundary        string   `json:"boundary,omitempty"`
	ContractVersion string   `json:"configurationApiVersion,omitempty"`
	InputSHA256     string   `json:"inputSHA256,omitempty"`
	Build           *Build   `json:"build,omitempty"`
	Error           *Failure `json:"error,omitempty"`
}

type arguments struct{ file, organization, project, output string }
type outcome struct {
	report Report
	code   int
}

func fault(code, category string) *Failure {
	return &Failure{Code: code, Module: "storm-cli", Category: category, Key: "cli." + code}
}

func failure(command string, err error) outcome {
	r := Report{APIVersion: ReportVersion, Command: command, Status: "error"}
	code := ExitInternal
	var violation *recovery.Violation
	var safe *Failure
	switch {
	case errors.Is(err, context.Canceled):
		r.Status = "cancelled"
		r.Error = fault("CANCELLED", "cancellation")
		code = ExitCancelled
	case errors.Is(err, context.DeadlineExceeded):
		r.Error = fault("TIMEOUT", "timeout")
		code = ExitTimeout
	case errors.As(err, &violation):
		r.Status = "invalid"
		code = ExitInvalid
		r.Error = &Failure{Code: string(violation.Code), Module: violation.Module, Category: violation.Category, Key: violation.Key, Arguments: violation.Arguments}
	case errors.As(err, &safe):
		r.Error = safe
		switch safe.Code {
		case "INVALID_ARGUMENTS":
			code = ExitUsage
		case "INPUT_UNREADABLE", "INPUT_TYPE":
			code = ExitInput
		case "OUTPUT_FAILED":
			code = ExitOutput
		case "INPUT_TOO_LARGE":
			r.Status = "invalid"
			code = ExitInvalid
		}
	default:
		r.Error = fault("INTERNAL_ERROR", "internal")
	}
	return outcome{r, code}
}

// outputFormat chooses JSON for argument errors only when requested unambiguously.
func outputFormat(args []string) string {
	mode, count := "text", 0
	for i, arg := range args {
		if strings.HasPrefix(arg, "--output=") {
			mode = strings.TrimPrefix(arg, "--output=")
			count++
		}
		if arg == "--output" && i+1 < len(args) {
			mode = args[i+1]
			count++
		}
	}
	if count == 1 && mode == "json" {
		return mode
	}
	return "text"
}

func parse(command string, args []string) (arguments, error) {
	result := arguments{output: outputFormat(args)}
	seen := map[string]bool{}
	for i := 0; i < len(args); i++ {
		key, value, hasValue := strings.Cut(args[i], "=")
		if key != "--output" && (command != "validate" || (key != "--file" && key != "--organization" && key != "--project")) {
			return result, fault("INVALID_ARGUMENTS", "usage")
		}
		if seen[key] {
			return result, fault("INVALID_ARGUMENTS", "usage")
		}
		seen[key] = true
		if !hasValue {
			i++
			if i >= len(args) {
				return result, fault("INVALID_ARGUMENTS", "usage")
			}
			value = args[i]
		}
		if value == "" || strings.HasPrefix(value, "--") {
			return result, fault("INVALID_ARGUMENTS", "usage")
		}
		switch key {
		case "--file":
			result.file = value
		case "--organization":
			result.organization = value
		case "--project":
			result.project = value
		case "--output":
			if value != "json" && value != "text" {
				return result, fault("INVALID_ARGUMENTS", "usage")
			}
			result.output = value
		}
	}
	if command == "validate" && (result.file == "" || result.organization == "" || result.project == "") {
		return result, fault("INVALID_ARGUMENTS", "usage")
	}
	return result, nil
}

func validate(ctx context.Context, args arguments, load func(context.Context, string) ([]byte, error)) outcome {
	scope := recovery.Scope{OrganizationID: args.organization, ProjectID: args.project}
	if err := scope.Validate(); err != nil {
		return failure("validate", err)
	}
	data, err := load(ctx, args.file)
	if err != nil {
		return failure("validate", err)
	}
	if err = ctx.Err(); err != nil {
		return failure("validate", err)
	}
	if _, err = recovery.DecodeYAML(data, scope); err != nil {
		return failure("validate", err)
	}
	digest := sha256.Sum256(data)
	return outcome{Report{APIVersion: ReportVersion, Command: "validate", Status: "valid", Boundary: "offline-local-expectation", ContractVersion: recovery.APIVersion, InputSHA256: hex.EncodeToString(digest[:])}, ExitOK}
}

const help = `Storm Harbor offline CLI

storm version [--output text|json]
storm validate --file <regular-file> --organization <expected-id> --project <expected-id> [--output text|json]
storm help

Validation compares a document to an explicit local scope expectation.
It does not authenticate, authorize, resolve secrets, plan or execute recovery.
Only bounded regular files are accepted; stdin and symlinks are unsupported.
`

// Execute owns process-signal wiring and delegates to the context-aware boundary.
func Execute(args []string, stdout, stderr io.Writer, options Options) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	return Run(ctx, args, stdout, stderr, options)
}

// Run returns stable exit codes. It makes no network call or configuration write.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer, options Options) int {
	mode := outputFormat(args)
	command := "cli"
	if len(args) == 0 || (len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h")) {
		if _, err := io.WriteString(stdout, help); err != nil {
			return emit(stderr, stderr, "text", failure("help", fault("OUTPUT_FAILED", "io")), ExitOutput)
		}
		return ExitOK
	}
	if args[0] == "validate" || args[0] == "version" {
		command = args[0]
	} else {
		return emit(stdout, stderr, mode, failure(command, fault("INVALID_ARGUMENTS", "usage")), 0)
	}
	opts, err := parse(command, args[1:])
	if err != nil {
		return emit(stdout, stderr, mode, failure(command, err), 0)
	}
	if err = ctx.Err(); err != nil {
		return emit(stdout, stderr, mode, failure(command, err), 0)
	}
	if command == "version" {
		if options.Version == "" {
			options.Version = "dev"
		}
		if options.Revision == "" {
			options.Revision = "unknown"
		}
		return emit(stdout, stderr, mode, outcome{Report{APIVersion: ReportVersion, Command: command, Status: "ok", Build: &Build{Version: options.Version, Revision: options.Revision, GoVersion: runtime.Version(), SDKVersion: SDKVersion}}, 0}, 0)
	}
	if options.Load == nil {
		options.Load = LoadFile
	}
	if options.Timeout <= 0 {
		options.Timeout = 5 * time.Second
	}
	bounded, cancel := context.WithTimeout(ctx, options.Timeout)
	defer cancel()
	completed := make(chan outcome, 1)
	go func() { completed <- validate(bounded, opts, options.Load) }()
	var result outcome
	select {
	case <-bounded.Done():
		result = failure(command, bounded.Err())
	case result = <-completed:
		if err = bounded.Err(); err != nil {
			result = failure(command, err)
		}
	}
	return emit(stdout, stderr, mode, result, 0)
}

func emit(stdout, stderr io.Writer, mode string, result outcome, override int) int {
	target := stdout
	if result.code != 0 {
		target = stderr
	}
	var err error
	if mode == "json" {
		err = json.NewEncoder(target).Encode(result.report)
	} else {
		switch {
		case result.report.Error != nil:
			_, err = fmt.Fprintf(target, "%s: %s\n", result.report.Error.Code, message(result.report.Error.Code))
		case result.report.Build != nil:
			_, err = fmt.Fprintf(target, "storm %s\nRevision: %s\nGo: %s\nRecovery contracts: %s\n", result.report.Build.Version, result.report.Build.Revision, result.report.Build.GoVersion, result.report.Build.SDKVersion)
		default:
			_, err = fmt.Fprintf(target, "Valid against the local scope expectation.\nInput SHA-256 (raw bytes): %s\nNo authentication, authorization or recovery was performed.\n", result.report.InputSHA256)
		}
	}
	if err != nil {
		if mode == "json" {
			_ = json.NewEncoder(stderr).Encode(failure(result.report.Command, fault("OUTPUT_FAILED", "io")).report)
		} else {
			_, _ = io.WriteString(stderr, "OUTPUT_FAILED: Cannot write CLI output.\n")
		}
		return ExitOutput
	}
	if override != 0 {
		return override
	}
	return result.code
}

func message(code string) string {
	switch code {
	case "INVALID_ARGUMENTS":
		return "Use storm help for supported commands and required flags."
	case "INPUT_UNREADABLE":
		return "Cannot read the input file."
	case "INPUT_TYPE":
		return "Use a regular file; symlinks and special files are unsupported."
	case "INPUT_TOO_LARGE":
		return "The input exceeds the configuration size limit."
	case "CANCELLED":
		return "Operation cancelled."
	case "TIMEOUT":
		return "The local validation deadline expired."
	case "OUTPUT_FAILED":
		return "Cannot write CLI output."
	case "INTERNAL_ERROR":
		return "The operation failed safely."
	default:
		return "Configuration validation failed; inspect the structured error code and safe field."
	}
}
