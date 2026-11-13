package release

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func inputs(t *testing.T, platform string) Options {
	t.Helper()
	root := t.TempDir()
	input := filepath.Join(root, "input")
	for _, dir := range []string{input, filepath.Join(root, "licenses")} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	binary := "storm"
	if platform == "windows" {
		binary += ".exe"
	}
	files := map[string]string{filepath.Join(input, binary): "synthetic executable", filepath.Join(input, "build-info.txt"): "synthetic public build metadata", filepath.Join(root, "LICENSE"): "synthetic license", filepath.Join(root, ".env"): "SENSITIVE_CANARY_DO_NOT_PACKAGE"}
	for _, name := range licenseNames {
		files[filepath.Join(root, "licenses", name)] = "synthetic notice"
	}
	for path, data := range files {
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return Options{Root: root, Input: input, Output: filepath.Join(root, "output"), Version: "v0.1.0-alpha.1", Revision: strings.Repeat("a", 40), OS: platform, Arch: "amd64"}
}
func unpack(t *testing.T, path string) map[string][]byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{}
	if strings.HasSuffix(path, ".zip") {
		z, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range z.File {
			r, err := file.Open()
			if err != nil {
				t.Fatal(err)
			}
			files[file.Name], err = io.ReadAll(r)
			_ = r.Close()
			if err != nil {
				t.Fatal(err)
			}
			if file.Name == "storm.exe" && file.Mode().Perm() != 0755 {
				t.Fatal("binary mode lost")
			}
		}
	} else {
		g, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		defer g.Close()
		r := tar.NewReader(g)
		for {
			header, err := r.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			files[header.Name], err = io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			if header.Name == "storm" && header.Mode != 0755 {
				t.Fatal("binary mode lost")
			}
		}
	}
	return files
}
func TestPackagesAreCompleteAndDeterministic(t *testing.T) {
	for _, platform := range []string{"linux", "darwin", "windows"} {
		t.Run(platform, func(t *testing.T) {
			o := inputs(t, platform)
			name, err := Package(o)
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(o.Output, name)
			files := unpack(t, path)
			if len(files) != 9 {
				t.Fatalf("unexpected package entries: %v", files)
			}
			var manifest Manifest
			if err = json.Unmarshal(files["build-manifest.json"], &manifest); err != nil {
				t.Fatal(err)
			}
			if manifest.APIVersion != manifestVersion || manifest.Version != o.Version || manifest.SourceRevision != o.Revision || manifest.OS != platform || manifest.Arch != o.Arch || manifest.GoVersion == "" {
				t.Fatal("manifest identity incorrect")
			}
			for _, item := range manifest.Files {
				digest := sha256.Sum256(files[item.Path])
				if item.SHA256 != hex.EncodeToString(digest[:]) || item.Bytes != len(files[item.Path]) {
					t.Fatal("manifest integrity incorrect")
				}
			}
			for name, data := range files {
				if strings.Contains(name, "..") || bytes.Contains(data, []byte("SENSITIVE_CANARY")) || bytes.Contains(data, []byte(o.Root)) {
					t.Fatal("private input escaped allowlist")
				}
			}
			original, _ := os.ReadFile(path)
			o.Output = filepath.Join(o.Root, "another-output")
			again, err := Package(o)
			if err != nil {
				t.Fatal(err)
			}
			second, _ := os.ReadFile(filepath.Join(o.Output, again))
			if !bytes.Equal(original, second) {
				t.Fatal("package is not deterministic")
			}
			if _, err = Package(o); err == nil {
				t.Fatal("existing artifact was overwritten")
			}
			after, _ := os.ReadFile(filepath.Join(o.Output, again))
			if !bytes.Equal(second, after) {
				t.Fatal("failed rerun changed artifact")
			}
		})
	}
}
func TestInvalidAndMissingInputsFailSafely(t *testing.T) {
	base := inputs(t, "linux")
	cases := []Options{base, base, base, base, base, base, base}
	cases[0].Version = "v0.1.0;SENSITIVE_CANARY"
	cases[1].Version = "v1.0.0"
	cases[2].Revision = "SENSITIVE_CANARY"
	cases[3].OS = "unknown"
	cases[4].Arch = "unknown"
	cases[5].Root = ""
	cases[6].Version = "v0.1.0-alpha.0"
	leadingZero := base
	leadingZero.Version = "v0.01.0-alpha.1"
	cases = append(cases, leadingZero)
	for _, o := range cases {
		_, err := Package(o)
		var safe *Failure
		if !errors.As(err, &safe) || safe.Code != "INVALID_OPTIONS" || safe.Category != "validation" || safe.Key != "release.INVALID_OPTIONS" {
			t.Fatal("invalid option was not safely rejected")
		}
	}
	for _, target := range []string{filepath.Join(base.Input, "storm"), filepath.Join(base.Input, "build-info.txt"), filepath.Join(base.Root, "LICENSE"), filepath.Join(base.Root, "licenses", licenseNames[0])} {
		original, err := os.ReadFile(target)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.Remove(target); err != nil {
			t.Fatal(err)
		}
		if _, err = Package(base); err == nil {
			t.Fatal("missing required input accepted")
		}
		if err = os.WriteFile(target, original, 0600); err != nil {
			t.Fatal(err)
		}
	}
	badOutput := filepath.Join(base.Root, "blocked")
	if err := os.WriteFile(badOutput, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	base.Output = badOutput
	if _, err := Package(base); err == nil {
		t.Fatal("file used as output directory")
	}
}
func TestBoundedFiles(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "data")
	for _, data := range [][]byte{nil, []byte("oversized")} {
		if err := os.WriteFile(file, data, 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := readBounded(file, 2); err == nil {
			t.Fatal("empty/oversized file accepted")
		}
	}
	if _, err := readBounded(dir, 10); err == nil {
		t.Fatal("directory accepted")
	}
	if _, err := readBounded(filepath.Join(dir, "absent"), 10); err == nil {
		t.Fatal("absent file accepted")
	}
	if runtime.GOOS != "windows" && os.Geteuid() != 0 {
		if err := os.Chmod(file, 0000); err != nil {
			t.Fatal(err)
		}
		if _, err := readBounded(file, 100); err == nil {
			t.Fatal("unreadable input accepted")
		}
		if err := os.Chmod(file, 0600); err != nil {
			t.Fatal(err)
		}
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(file, link); err == nil {
		if _, err = readBounded(link, 100); err == nil {
			t.Fatal("symlink accepted")
		}
	}
}

type limitWriter struct{ left int }

func (w *limitWriter) Write(data []byte) (int, error) {
	if w.left < len(data) {
		n := w.left
		w.left = 0
		return n, io.ErrShortWrite
	}
	w.left -= len(data)
	return len(data), nil
}

type failingFile struct {
	*os.File
	failClose bool
}

func (f failingFile) Write(data []byte) (int, error) {
	if !f.failClose {
		return 0, io.ErrShortWrite
	}
	return f.File.Write(data)
}
func (f failingFile) Close() error {
	_ = f.File.Close()
	if f.failClose {
		return io.ErrClosedPipe
	}
	return nil
}

func TestPartialArchiveFailuresAreNotPublished(t *testing.T) {
	payload := make([]byte, 128<<10)
	_, _ = rand.New(rand.NewSource(42)).Read(payload)
	for _, windows := range []bool{false, true} {
		for _, quota := range []int{0, 1, 10, 50, 1000, 8000, 32000, 64000, 120000, 150000} {
			err := archive(&limitWriter{left: quota}, windows, []entry{{"storm", payload, 0755}, {"license", payload, 0644}})
			if err == nil {
				t.Fatal("partial archive write reported success")
			}
		}
	}
	for _, failClose := range []bool{false, true} {
		o := inputs(t, "linux")
		_, err := packageWith(o, func(name string) (outputFile, error) {
			f, e := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
			return failingFile{f, failClose}, e
		})
		if err == nil {
			t.Fatal("partial package reported success")
		}
		items, e := os.ReadDir(o.Output)
		if e != nil || len(items) != 0 {
			t.Fatal("partial package was retained")
		}
	}
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, errors.New("SENSITIVE_WRITE_DETAIL") }
func TestArchiveAndCommandFailures(t *testing.T) {
	for _, windows := range []bool{false, true} {
		if err := archive(brokenWriter{}, windows, []entry{{"storm", []byte("binary"), 0755}}); err == nil {
			t.Fatal("write failure ignored")
		}
	}
	for _, args := range [][]string{{"--unknown=SENSITIVE_CANARY"}, {"unexpected"}, {"--version=bad"}} {
		var stdout, stderr bytes.Buffer
		if code := Run(args, &stdout, &stderr); code == 0 || stdout.Len() != 0 || strings.Contains(stderr.String(), "SENSITIVE_CANARY") {
			t.Fatal("unsafe command failure")
		}
	}
	o := inputs(t, "linux")
	args := []string{"--root", o.Root, "--input", o.Input, "--output", o.Output, "--version", o.Version, "--revision", o.Revision, "--os", o.OS, "--arch", o.Arch}
	var out, stderr bytes.Buffer
	if code := Run(args, &out, &stderr); code != 0 || stderr.Len() != 0 || !strings.HasPrefix(out.String(), "storm_v0.1.0-alpha.1_linux_amd64") {
		t.Fatal("command failed")
	}
	o.Output = filepath.Join(o.Root, "write-error")
	args[5] = o.Output
	if code := Run(args, brokenWriter{}, &stderr); code != 74 {
		t.Fatal("command output failure ignored")
	}
}
