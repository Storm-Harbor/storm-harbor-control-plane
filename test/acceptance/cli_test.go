package acceptance

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestExecutable(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	binary := filepath.Join(dir, "storm")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	if candidate := os.Getenv("STORM_TEST_BINARY"); candidate != "" {
		binary = candidate
		if !filepath.IsAbs(binary) {
			binary = filepath.Join(root, binary)
		}
	} else {
		build := exec.CommandContext(ctx, "go", "build", "-trimpath", "-o", binary, "./cmd/storm")
		build.Dir = root
		if out, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build: %v %s", err, out)
		}
	}
	canary := "SENSITIVE_CANARY_DO_NOT_ECHO"
	bad := filepath.Join(dir, canary+".yaml")
	if err := os.WriteFile(bad, []byte("secret: "+canary), 0600); err != nil {
		t.Fatal(err)
	}
	large := filepath.Join(dir, "large.yaml")
	if err := os.WriteFile(large, make([]byte, 65537), 0600); err != nil {
		t.Fatal(err)
	}
	valid := []string{"validate", "--file", filepath.Join(root, "examples", "exercise.yaml"), "--organization", "org-example", "--project", "project-example", "--output=json"}
	restore := append([]string{}, valid...)
	restore[2] = filepath.Join(root, "examples", "restore.yaml")
	wrong := append([]string{}, valid...)
	wrong[4] = "org-other"
	invalid := append([]string{}, valid...)
	invalid[2] = bad
	missing := append([]string{}, valid...)
	missing[2] = filepath.Join(dir, canary+"-missing")
	oversized := append([]string{}, valid...)
	oversized[2] = large
	cases := []struct {
		name   string
		args   []string
		code   int
		status string
	}{
		{"version", []string{"version", "--output=json"}, 0, "ok"},
		{"exercise", valid, 0, "valid"}, {"restore", restore, 0, "valid"},
		{"wrong-scope", wrong, 2, "invalid"}, {"malformed", invalid, 2, "invalid"},
		{"oversized", oversized, 2, "invalid"}, {"missing", missing, 3, "error"},
		{"unsupported", []string{"run", canary, "--output=json"}, 64, "error"},
		{"duplicate", append(valid, "--project", canary), 64, "error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.CommandContext(ctx, binary, tc.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			code := 0
			if err != nil {
				if exit, ok := err.(*exec.ExitError); ok {
					code = exit.ExitCode()
				} else {
					t.Fatal(err)
				}
			}
			if code != tc.code {
				t.Fatalf("exit %d expected %d: %s", code, tc.code, stderr.String())
			}
			data := stdout.Bytes()
			if code != 0 {
				data = stderr.Bytes()
				if stdout.Len() != 0 {
					t.Fatal("error wrote stdout")
				}
			} else if stderr.Len() != 0 {
				t.Fatal("success wrote stderr")
			}
			var report map[string]any
			if err = json.Unmarshal(data, &report); err != nil {
				t.Fatalf("invalid structured output: %v", err)
			}
			if report["apiVersion"] != "cli.stormharbor.io/v1alpha1" || report["status"] != tc.status || bytes.Contains(data, []byte(canary)) || strings.Contains(string(data), root) {
				t.Fatal("unsafe/unexpected report")
			}
			if tc.name == "version" {
				build, ok := report["build"].(map[string]any)
				if !ok {
					t.Fatal("version metadata missing")
				}
				for name, key := range map[string]string{"STORM_TEST_VERSION": "version", "STORM_TEST_REVISION": "revision"} {
					if expected := os.Getenv(name); expected != "" && build[key] != expected {
						t.Fatalf("release build metadata %s does not match expected value", key)
					}
				}
			}
		})
	}
	cmd := exec.CommandContext(ctx, binary, "help")
	if out, err := cmd.CombinedOutput(); err != nil || !bytes.Contains(out, []byte("does not authenticate")) {
		t.Fatal("help boundary missing")
	}
}
