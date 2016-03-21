package domain

import "io"

// A FileSystem implements access to a collection of named files.
// The elements in a file path are separated by slash ('/', U+002F)
// characters, regardless of host operating system convention.
type FileSystem interface {
	Open(name string, perms int) (File, error)
}

type File interface {
	Close() error
	Name() string     // Full path to file.
	FileName() string // Desired name of the file to be returned as
	io.Writer
	io.Reader
}
