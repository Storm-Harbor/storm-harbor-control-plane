package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	recovery "github.com/Storm-Harbor/storm-lib-recovery-contracts"
)

// These port tests exercise disappearance/replacement/read failures that are
// unsafe to race nondeterministically against the operating system.
type scenarioFiles struct {
	before  os.FileInfo
	file    inputFile
	openErr error
}

func (f scenarioFiles) Lstat(string) (os.FileInfo, error) { return f.before, nil }
func (f scenarioFiles) Open(string) (inputFile, error)    { return f.file, f.openErr }

type scenarioFile struct {
	reader  io.Reader
	info    os.FileInfo
	statErr error
	closed  bool
}

func (f *scenarioFile) Read(data []byte) (int, error) { return f.reader.Read(data) }
func (f *scenarioFile) Stat() (os.FileInfo, error)    { return f.info, f.statErr }
func (f *scenarioFile) Close() error                  { f.closed = true; return nil }

type badReader struct{ cancel context.CancelFunc }

func (r badReader) Read([]byte) (int, error) {
	if r.cancel != nil {
		r.cancel()
	}
	return 0, errors.New("PRIVATE_FS_DETAIL")
}
func TestInputPartialFailures(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a")
	b := filepath.Join(dir, "b")
	for _, p := range []string{a, b} {
		if err := os.WriteFile(p, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	first, _ := os.Stat(a)
	second, _ := os.Stat(b)
	directory, _ := os.Stat(dir)
	cases := []struct {
		name     string
		fs       scenarioFiles
		expected string
	}{
		{"disappeared", scenarioFiles{before: first, openErr: errors.New("PRIVATE_PATH")}, "INPUT_UNREADABLE"},
		{"stat-denied", scenarioFiles{before: first, file: &scenarioFile{statErr: errors.New("PRIVATE_STAT")}}, "INPUT_TYPE"},
		{"replaced", scenarioFiles{before: first, file: &scenarioFile{info: second}}, "INPUT_TYPE"},
		{"became-directory", scenarioFiles{before: first, file: &scenarioFile{info: directory}}, "INPUT_TYPE"},
		{"read-error", scenarioFiles{before: first, file: &scenarioFile{info: first, reader: badReader{}}}, "INPUT_UNREADABLE"},
		{"grew", scenarioFiles{before: first, file: &scenarioFile{info: first, reader: bytes.NewReader(make([]byte, recovery.MaxDocumentBytes+10))}}, "INPUT_TOO_LARGE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := loadFile(context.Background(), "unused", tc.fs)
			var safe *Failure
			if !errors.As(err, &safe) || safe.Code != tc.expected {
				t.Fatalf("unexpected safe error %v", err)
			}
			if file, ok := tc.fs.file.(*scenarioFile); ok && !file.closed {
				t.Fatal("handle leaked")
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	file := &scenarioFile{info: first, reader: badReader{cancel: cancel}}
	_, err := loadFile(ctx, "unused", scenarioFiles{before: first, file: file})
	if !errors.Is(err, context.Canceled) {
		t.Fatal("read cancellation lost")
	}
	reader := contextReader{ctx: ctx, source: bytes.NewReader(nil)}
	if _, err = reader.Read(make([]byte, 1)); !errors.Is(err, context.Canceled) {
		t.Fatal("cancelled reader accessed source")
	}
}
