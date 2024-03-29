package hdfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cwillia9/ez-ftp/domain"
)

type t struct {
	rootDir  string
	hdfsHost string
	hdfsPort int
	hdfsUser string
}

type Config struct {
	rootDir  string
	hdfsHost string
	hdfsPort int
	hdfsUser string
}

type hdfsFile struct {
	f           *os.File
	rootDir     string
	outwardName string
}

func (f *hdfsFile) Close() error {
	f.f.Close()
	return os.Remove(f.f.Name())
}

func (f *hdfsFile) Name() string {
	return f.f.Name()
}

func (f *hdfsFile) Read(p []byte) (int, error) {
	return f.f.Read(p)
}

func (f *hdfsFile) Write(p []byte) (int, error) {
	return f.f.Write(p)
}

func (f *hdfsFile) FileName() string {
	return f.outwardName
}

func New(cfg Config) (*t, error) {
	hdfs := &t{}
	hdfs.rootDir = cfg.rootDir
	hdfs.hdfsHost = cfg.hdfsHost
	hdfs.hdfsPort = cfg.hdfsPort
	hdfs.hdfsUser = cfg.hdfsUser
	return hdfs, nil
}

func (this *t) Open(name string, flags int) (domain.File, error) {
	// TODO(cwilliams): It's a security concern to not check if
	// name has any '..' in it
	f, err := this.getFile(name)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(f.Name(), os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	hdfsFile := &hdfsFile{}

	hdfsFile.f = file
	hdfsFile.rootDir = this.rootDir
	splitName := strings.Split(name, "/")
	hdfsFile.outwardName = splitName[len(splitName)-1]
	return hdfsFile, nil
}

// returns file on local filesystem. This file will have a "random-ish" name
// generated by the use of ioutil.TempFile
func (this *t) getFile(relativePath string) (*os.File, error) {
	fullPath := filepath.Join(this.rootDir, relativePath)

	url := fmt.Sprintf("http://%s@%s:%d/webhdfs/v1%s?op=OPEN&user.name=%s", this.hdfsUser, this.hdfsHost, this.hdfsPort, fullPath, this.hdfsUser)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}
