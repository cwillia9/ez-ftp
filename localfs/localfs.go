package localfs

// T holds the local file system struct
type T struct {
	rootDir string
}

type LocalFile struct {
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

func (f *LocalFile) Close() error {
	return nil
}

func (f *LocalFile) Read(p []byte) (int, error) {
	return 0, nil
}

func (f *LocalFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *LocalFile) Stat()

func (t *T) Open(path string)
