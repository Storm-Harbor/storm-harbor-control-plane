package cli

import (
	"context"
	"io"
	"os"

	recovery "github.com/Storm-Harbor/storm-lib-recovery-contracts"
)

// LoadFile accepts a bounded regular file. The caller bounds the whole operation
// separately, including filesystem stalls. No underlying OS error is disclosed.
func LoadFile(ctx context.Context, name string) ([]byte, error) {
	return loadFile(ctx, name, localFiles{})
}

type inputFile interface {
	io.Reader
	Stat() (os.FileInfo, error)
	Close() error
}
type files interface {
	Lstat(string) (os.FileInfo, error)
	Open(string) (inputFile, error)
}
type localFiles struct{}

func (localFiles) Lstat(name string) (os.FileInfo, error) { return os.Lstat(name) }
func (localFiles) Open(name string) (inputFile, error)    { return os.Open(name) }

func loadFile(ctx context.Context, name string, fs files) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	before, err := fs.Lstat(name)
	if err != nil {
		return nil, fault("INPUT_UNREADABLE", "io")
	}
	if !before.Mode().IsRegular() {
		return nil, fault("INPUT_TYPE", "io")
	}
	if before.Size() > recovery.MaxDocumentBytes {
		return nil, fault("INPUT_TOO_LARGE", "validation")
	}
	file, err := fs.Open(name)
	if err != nil {
		return nil, fault("INPUT_UNREADABLE", "io")
	}
	defer file.Close()
	after, err := file.Stat()
	if err != nil || !after.Mode().IsRegular() || !os.SameFile(before, after) {
		return nil, fault("INPUT_TYPE", "io")
	}
	data, err := io.ReadAll(io.LimitReader(&contextReader{ctx: ctx, source: file}, recovery.MaxDocumentBytes+1))
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fault("INPUT_UNREADABLE", "io")
	}
	if len(data) > recovery.MaxDocumentBytes {
		return nil, fault("INPUT_TOO_LARGE", "validation")
	}
	return data, nil
}

type contextReader struct {
	ctx    context.Context
	source io.Reader
}

func (r *contextReader) Read(data []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.source.Read(data)
}
