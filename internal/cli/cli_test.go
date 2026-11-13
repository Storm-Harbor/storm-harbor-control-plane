package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	recovery "github.com/Storm-Harbor/storm-lib-recovery-contracts"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "examples", name+".yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return data
}
func args() []string {
	return []string{"validate", "--file", "fixture.yaml", "--organization", "org-example", "--project", "project-example", "--output", "json"}
}
func invoke(t *testing.T, ctx context.Context, command []string, options Options) (int, Report, string, string) {
	t.Helper()
	var out, err bytes.Buffer
	code := Run(ctx, command, &out, &err, options)
	var report Report
	data := out.Bytes()
	if code != 0 {
		data = err.Bytes()
	}
	if len(data) > 0 && data[0] == '{' {
		if e := json.Unmarshal(data, &report); e != nil {
			t.Fatal(e)
		}
	}
	return code, report, out.String(), err.String()
}
func TestValidConfigurations(t *testing.T) {
	for _, name := range []string{"exercise", "restore"} {
		t.Run(name, func(t *testing.T) {
			data := fixture(t, name)
			calls := 0
			code, r, out, err := invoke(t, context.Background(), args(), Options{Load: func(_ context.Context, path string) ([]byte, error) {
				calls++
				if path != "fixture.yaml" {
					t.Fatal("wrong input")
				}
				return data, nil
			}})
			digest := sha256.Sum256(data)
			if code != 0 || calls != 1 || err != "" || r.Status != "valid" || r.APIVersion != ReportVersion || r.Boundary != "offline-local-expectation" || r.ContractVersion != recovery.APIVersion || r.InputSHA256 != hex.EncodeToString(digest[:]) {
				t.Fatalf("unexpected outcome: %d %#v %s", code, r, err)
			}
			for _, s := range []string{"org-example", "project-example", "recovery-database", "primary-environment"} {
				if strings.Contains(out, s) {
					t.Fatal("configuration value leaked")
				}
			}
		})
	}
	data := fixture(t, "exercise")
	_, _, out, _ := invoke(t, context.Background(), args()[:7], Options{Load: func(context.Context, string) ([]byte, error) { return data, nil }})
	if !strings.Contains(out, "raw bytes") || !strings.Contains(out, "No authentication") {
		t.Fatal("missing boundary explanation")
	}
}
func TestValidationFailuresAreSafe(t *testing.T) {
	base := string(fixture(t, "exercise"))
	canary := "SENSITIVE_CANARY_DO_NOT_ECHO"
	cases := []struct{ name, input, code string }{
		{"malformed", "[" + canary, "INVALID_DOCUMENT"},
		{"unknown", base + "unexpected: " + canary + "\n", "INVALID_DOCUMENT"},
		{"duplicate", base + "kind: " + canary + "\n", "INVALID_DOCUMENT"},
		{"multiple", base + "---\nsecret: " + canary, "INVALID_DOCUMENT"},
		{"scope", strings.ReplaceAll(base, "org-example", "org-other"), "SCOPE_MISMATCH"},
		{"lifecycle", strings.Replace(base, "cleanup: run-owned", "cleanup: retain", 1), "INVALID_CONFIGURATION"},
		{"too-large", strings.Repeat("x", recovery.MaxDocumentBytes+1), "INVALID_DOCUMENT"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, r, out, err := invoke(t, context.Background(), args(), Options{Load: func(context.Context, string) ([]byte, error) { return []byte(tc.input), nil }})
			if code != ExitInvalid || r.Status != "invalid" || r.Error == nil || r.Error.Code != tc.code || !strings.HasPrefix(r.Error.Key, "contracts.") || out != "" || strings.Contains(err, canary) || strings.Contains(err, "org-other") {
				t.Fatalf("unsafe/unexpected result: %d %#v", code, r)
			}
		})
	}
	command := args()
	command[4] = "not a valid scope"
	code, r, _, _ := invoke(t, context.Background(), command, Options{Load: func(context.Context, string) ([]byte, error) {
		t.Fatal("invalid scope must precede file access")
		return nil, nil
	}})
	if code != ExitInvalid || r.Error.Code != "INVALID_SCOPE" {
		t.Fatal("missing scope rejection")
	}
}
func TestCommandsAndArguments(t *testing.T) {
	for _, command := range [][]string{nil, {"help"}, {"--help"}, {"-h"}} {
		code, _, out, err := invoke(t, context.Background(), command, Options{})
		if code != 0 || err != "" || !strings.Contains(out, "storm validate") {
			t.Fatal("help failed")
		}
	}
	for _, command := range [][]string{{"version"}, {"version", "--output=text"}, {"version", "--output", "json"}} {
		code, r, out, err := invoke(t, context.Background(), command, Options{})
		if code != 0 || err != "" || out == "" {
			t.Fatal("version failed")
		}
		if r.Build != nil && (r.Build.Version != "dev" || r.Build.Revision != "unknown" || r.Build.SDKVersion != SDKVersion) {
			t.Fatal("bad build defaults")
		}
	}
	code, r, _, _ := invoke(t, context.Background(), []string{"version", "--output=json"}, Options{Version: "alpha-test", Revision: "0123456789"})
	if code != 0 || r.Build.Version != "alpha-test" || r.Build.Revision != "0123456789" {
		t.Fatal("bad build metadata")
	}
	canary := "SENSITIVE_CANARY_DO_NOT_ECHO"
	cases := [][]string{
		{"run", canary, "--output=json"}, {"help", canary}, {"version", "--file", canary},
		{"validate"}, {"validate", "--file"}, {"validate", "--file="}, {"validate", "--file", "--project"},
		{"version", "--output=xml"}, {"version", "--output"}, {"version", "--output=json", "--output=text"},
		append(args(), "--organization", canary), append(args(), canary), append(args(), "--unknown="+canary),
	}
	for _, command := range cases {
		t.Run(strings.Join(command, " "), func(t *testing.T) {
			code, r, out, err := invoke(t, context.Background(), command, Options{})
			if code != ExitUsage || out != "" || err == "" || strings.Contains(err, canary) || (r.Error != nil && r.Error.Code != "INVALID_ARGUMENTS") {
				t.Fatalf("unsafe usage result %d", code)
			}
		})
	}
	data := fixture(t, "exercise")
	command := []string{"validate", "--file=fixture", "--organization=org-example", "--project=project-example", "--output=text"}
	code, _, _, _ = invoke(t, context.Background(), command, Options{Load: func(context.Context, string) ([]byte, error) { return data, nil }})
	if code != 0 {
		t.Fatal("equals flags rejected")
	}
}
func TestCancellationTimeoutAndLoaderFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	code, r, _, _ := invoke(t, ctx, args(), Options{})
	if code != ExitCancelled || r.Status != "cancelled" {
		t.Fatal("cancel failed")
	}
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	code, r, _, _ = invoke(t, ctx, args(), Options{})
	if code != ExitTimeout || r.Error.Code != "TIMEOUT" {
		t.Fatal("deadline failed")
	}
	for _, e := range []error{context.Canceled, context.DeadlineExceeded, errors.New("PRIVATE_INTERNAL_ERROR"), fault("INPUT_UNREADABLE", "io"), fault("INPUT_TYPE", "io"), fault("INPUT_TOO_LARGE", "validation")} {
		code, r, out, err := invoke(t, context.Background(), args(), Options{Load: func(context.Context, string) ([]byte, error) { return nil, e }})
		if code == 0 || r.Error == nil || out != "" || strings.Contains(err, "PRIVATE_INTERNAL_ERROR") {
			t.Fatal("unsafe failure")
		}
	}
	started, finished := make(chan struct{}), make(chan struct{})
	ctx, cancel = context.WithCancel(context.Background())
	result := make(chan int, 1)
	var out, err bytes.Buffer
	go func() {
		result <- Run(ctx, args(), &out, &err, Options{Load: func(ctx context.Context, _ string) ([]byte, error) {
			close(started)
			<-ctx.Done()
			defer close(finished)
			return nil, ctx.Err()
		}})
	}()
	<-started
	cancel()
	if <-result != ExitCancelled {
		t.Fatal("in-flight cancel failed")
	}
	<-finished
	code, _, _, _ = invoke(t, context.Background(), args(), Options{Timeout: time.Millisecond, Load: func(ctx context.Context, _ string) ([]byte, error) { <-ctx.Done(); return nil, ctx.Err() }})
	if code != ExitTimeout {
		t.Fatal("bounded timeout failed")
	}
	ctx, cancel = context.WithCancel(context.Background())
	code, _, _, _ = invoke(t, ctx, args(), Options{Load: func(context.Context, string) ([]byte, error) { cancel(); return fixture(t, "exercise"), nil }})
	if code != ExitCancelled {
		t.Fatal("post-read cancellation ignored")
	}
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, errors.New("PRIVATE_WRITE_ERROR") }
func TestOutputFailure(t *testing.T) {
	for _, command := range [][]string{{"help"}, {"version"}, {"version", "--output=json"}, {"run"}} {
		var err bytes.Buffer
		code := Run(context.Background(), command, failedWriter{}, &err, Options{})
		if command[0] == "run" {
			code = Run(context.Background(), command, io.Discard, failedWriter{}, Options{})
		}
		if code != ExitOutput || strings.Contains(err.String(), "PRIVATE_WRITE_ERROR") {
			t.Fatal("write failure was not safe")
		}
		if outputFormat(command) == "json" {
			var report Report
			if e := json.Unmarshal(err.Bytes(), &report); e != nil || report.Error == nil || report.Error.Code != "OUTPUT_FAILED" {
				t.Fatal("JSON write failure lost its structured envelope")
			}
		}
	}
	if fault("INTERNAL_ERROR", "internal").Error() != "INTERNAL_ERROR" {
		t.Fatal("typed error not stable")
	}
}

func TestSafeTextDiagnostics(t *testing.T) {
	command := args()
	command[len(command)-1] = "text"
	for _, e := range []error{context.Canceled, context.DeadlineExceeded, errors.New("PRIVATE_DIAGNOSTIC"), fault("INPUT_UNREADABLE", "io"), fault("INPUT_TYPE", "io"), fault("INPUT_TOO_LARGE", "validation")} {
		code, _, out, err := invoke(t, context.Background(), command, Options{Load: func(context.Context, string) ([]byte, error) { return nil, e }})
		if code == 0 || out != "" || err == "" || strings.Contains(err, "PRIVATE_DIAGNOSTIC") {
			t.Fatal("unsafe text diagnostic")
		}
	}
	code, _, _, err := invoke(t, context.Background(), command, Options{Load: func(context.Context, string) ([]byte, error) { return []byte("not a configuration"), nil }})
	if code != ExitInvalid || !strings.Contains(err, "INVALID_DOCUMENT") {
		t.Fatal("missing safe validation text")
	}
}

func TestProcessEntryPoint(t *testing.T) {
	var out, stderr bytes.Buffer
	code := Execute([]string{"version", "--output=json"}, &out, &stderr, Options{Version: "entry-point-test"})
	var report Report
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if code != 0 || stderr.Len() != 0 || report.Build.Version != "entry-point-test" {
		t.Fatal("process wiring failed")
	}
}
func TestRealFileAdapter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "exercise.yaml")
	if err := os.WriteFile(path, fixture(t, "exercise"), 0600); err != nil {
		t.Fatal(err)
	}
	command := args()
	command[2] = path
	code, r, _, _ := invoke(t, context.Background(), command, Options{})
	if code != 0 || r.Status != "valid" {
		t.Fatal("real input failed")
	}
	for _, path := range []string{filepath.Join(dir, "missing"), dir} {
		if _, err := LoadFile(context.Background(), path); err == nil {
			t.Fatal("unsafe input accepted")
		}
	}
	large := filepath.Join(dir, "large")
	if err := os.WriteFile(large, make([]byte, recovery.MaxDocumentBytes+1), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFile(context.Background(), large); err == nil {
		t.Fatal("large file accepted")
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(path, link); err == nil {
		if _, err = LoadFile(context.Background(), link); err == nil {
			t.Fatal("symlink accepted")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := LoadFile(ctx, path); !errors.Is(err, context.Canceled) {
		t.Fatal("cancel ignored")
	}
}
