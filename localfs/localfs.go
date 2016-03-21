package localfs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cwillia9/ez-ftp/domain"
)

// T holds the local file system struct
type T struct {
	rootDir string
}

// Config config expectation
type Config interface {
	GetRootDir() string
}

type localFile struct {
	f           *os.File
	outwardName string
}

func (l *localFile) Close() error {
	return l.f.Close()
}

func (l *localFile) Name() string {
	return l.f.Name()
}

func (l *localFile) Read(b []byte) (int, error) {
	return l.f.Read(b)
}

func (l *localFile) Write(b []byte) (int, error) {
	return l.f.Write(b)
}

func (l *localFile) FileName() string {
	return l.outwardName
}

// New returns a new initialized local file system instance
func New(cfg Config) (*T, error) {
	localfs := &T{}
	localfs.rootDir = cfg.GetRootDir()
	return localfs, nil
}

// Opens file on local file system. The path given will be relative
// to the rootDir
func (t *T) Open(name string, flags int) (domain.File, error) {
	// TODO(cwilliams): It's a security concern to not check if
	// name has any '..' in it
	path := filepath.Join(t.rootDir, name)
	file, err := os.OpenFile(path, flags, 0755)
	if err != nil {
		return nil, err
	}
	f := &localFile{}
	f.f = file
	splitName := strings.Split(name, "/")
	f.outwardName = splitName[len(splitName)-1]

	return f, nil
}
