// Package release creates bounded, deterministic maintainer release packages.
// It does not authenticate a build, sign it, create a Git tag or publish assets.
package release

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

const manifestVersion = "release.stormharbor.io/v1alpha1"
const maxBinaryBytes = 64 << 20

var alphaVersion = regexp.MustCompile(`^v0\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)-alpha\.[1-9][0-9]*$`)
var commitDigest = regexp.MustCompile(`^[a-f0-9]{40}$`)

type Failure struct{ Code, Category, Key string }

func (e *Failure) Error() string          { return e.Code }
func fail(code, category string) *Failure { return &Failure{code, category, "release." + code} }

type Options struct{ Root, Input, Output, Version, Revision, OS, Arch string }
type Digest struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Bytes  int    `json:"bytes"`
}
type Manifest struct {
	APIVersion     string   `json:"apiVersion"`
	Version        string   `json:"version"`
	SourceRevision string   `json:"sourceRevision"`
	OS             string   `json:"os"`
	Arch           string   `json:"architecture"`
	GoVersion      string   `json:"goVersion"`
	Files          []Digest `json:"files"`
}
type entry struct {
	name string
	data []byte
	mode int64
}

var licenseNames = []string{"go-LICENSE.txt", "recovery-contracts-LICENSE.txt", "recovery-contracts-NOTICE.txt", "yaml-v4-LICENSE.txt", "yaml-v4-NOTICE.txt"}

func valid(o Options) bool {
	return o.Root != "" && o.Input != "" && o.Output != "" && len(o.Version) <= 64 && alphaVersion.MatchString(o.Version) && commitDigest.MatchString(o.Revision) &&
		(o.OS == "linux" || o.OS == "darwin" || o.OS == "windows") && (o.Arch == "amd64" || o.Arch == "arm64")
}

func readBounded(path string, limit int64) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > limit {
		return nil, fail("INPUT_INVALID", "io")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fail("INPUT_INVALID", "io")
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil || len(data) == 0 || int64(len(data)) > limit {
		return nil, fail("INPUT_INVALID", "io")
	}
	return data, nil
}

func collect(o Options) ([]entry, error) {
	binary := "storm"
	if o.OS == "windows" {
		binary += ".exe"
	}
	sources := []struct {
		name, path  string
		mode, limit int64
	}{
		{binary, filepath.Join(o.Input, binary), 0755, maxBinaryBytes},
		{"build-info.txt", filepath.Join(o.Input, "build-info.txt"), 0644, 256 << 10},
		{"LICENSE", filepath.Join(o.Root, "LICENSE"), 0644, 1 << 20},
	}
	for _, name := range licenseNames {
		sources = append(sources, struct {
			name, path  string
			mode, limit int64
		}{"licenses/" + name, filepath.Join(o.Root, "licenses", name), 0644, 1 << 20})
	}
	manifest := Manifest{APIVersion: manifestVersion, Version: o.Version, SourceRevision: o.Revision, OS: o.OS, Arch: o.Arch, GoVersion: runtime.Version()}
	entries := make([]entry, 0, len(sources)+1)
	for _, source := range sources {
		data, err := readBounded(source.path, source.limit)
		if err != nil {
			return nil, err
		}
		digest := sha256.Sum256(data)
		manifest.Files = append(manifest.Files, Digest{source.name, hex.EncodeToString(digest[:]), len(data)})
		entries = append(entries, entry{source.name, data, source.mode})
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fail("MANIFEST_INVALID", "internal")
	}
	return append(entries, entry{"build-manifest.json", append(data, '\n'), 0644}), nil
}

func archive(w io.Writer, windows bool, entries []entry) error {
	epoch := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	if windows {
		z := zip.NewWriter(w)
		for _, item := range entries {
			header := &zip.FileHeader{Name: item.name, Method: zip.Deflate, Modified: epoch}
			header.SetMode(os.FileMode(item.mode))
			file, err := z.CreateHeader(header)
			if err != nil {
				return err
			}
			if _, err = file.Write(item.data); err != nil {
				return err
			}
		}
		return z.Close()
	}
	g := gzip.NewWriter(w)
	g.Header.ModTime = epoch
	t := tar.NewWriter(g)
	for _, item := range entries {
		if err := t.WriteHeader(&tar.Header{Name: item.name, Mode: item.mode, Size: int64(len(item.data)), ModTime: epoch, Typeflag: tar.TypeReg}); err != nil {
			return err
		}
		if _, err := t.Write(item.data); err != nil {
			return err
		}
	}
	if err := t.Close(); err != nil {
		return err
	}
	return g.Close()
}

// Package never overwrites an existing artifact. Partial output is removed only
// if this call created it. Names and metadata contain no filesystem source paths.
type outputFile interface {
	io.Writer
	Close() error
}

func Package(o Options) (string, error) {
	return packageWith(o, func(name string) (outputFile, error) {
		return os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	})
}

func packageWith(o Options, create func(string) (outputFile, error)) (asset string, err error) {
	if !valid(o) {
		return "", fail("INVALID_OPTIONS", "validation")
	}
	entries, err := collect(o)
	if err != nil {
		return "", err
	}
	if err = os.MkdirAll(o.Output, 0755); err != nil {
		return "", fail("OUTPUT_FAILED", "io")
	}
	ext := ".tar.gz"
	if o.OS == "windows" {
		ext = ".zip"
	}
	asset = "storm_" + o.Version + "_" + o.OS + "_" + o.Arch + ext
	target := filepath.Join(o.Output, asset)
	f, err := create(target)
	if err != nil {
		return "", fail("OUTPUT_FAILED", "io")
	}
	defer func() {
		_ = f.Close()
		if err != nil {
			_ = os.Remove(target)
		}
	}()
	if err = archive(f, o.OS == "windows", entries); err != nil {
		return "", fail("OUTPUT_FAILED", "io")
	}
	if err = f.Close(); err != nil {
		return "", fail("OUTPUT_FAILED", "io")
	}
	return asset, nil
}

// Run is the maintainer command boundary. It suppresses raw flag/OS errors.
func Run(args []string, stdout, stderr io.Writer) int {
	var o Options
	flags := flag.NewFlagSet("storm-package", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&o.Root, "root", ".", "")
	flags.StringVar(&o.Input, "input", ".artifacts/release", "")
	flags.StringVar(&o.Output, "output", ".artifacts/packages", "")
	flags.StringVar(&o.Version, "version", "", "")
	flags.StringVar(&o.Revision, "revision", "", "")
	flags.StringVar(&o.OS, "os", "", "")
	flags.StringVar(&o.Arch, "arch", "", "")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		_, _ = fmt.Fprintln(stderr, "release.INVALID_OPTIONS")
		return 64
	}
	asset, err := Package(o)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "release."+err.Error())
		return 1
	}
	if _, err = fmt.Fprintln(stdout, asset); err != nil {
		_, _ = fmt.Fprintln(stderr, "release.OUTPUT_FAILED")
		return 74
	}
	return 0
}
