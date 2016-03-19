package localfs

import (
	"os"
	"path/filepath"
)

// T holds the local file system struct
type T struct {
	rootDir string
}

// Config config expectation
type Config interface {
	GetRootDir() string
}

// New returns a new initialized local file system instance
func New(cfg Config) (*T, error) {
	localfs := &T{}
	localfs.rootDir = cfg.GetRootDir()
	return localfs, nil
}

// Opens file on local file system. The path given will be relative
// to the rootDir
func (t *T) Open(name string, flags int) (*os.File, error) {
	// TODO(cwilliams): It's a security concern to not check if
	// name has any '..' in it
	path := filepath.Join(t.rootDir, name)
	file, err := os.OpenFile(path, flags, 0755)
	return file, err
}
