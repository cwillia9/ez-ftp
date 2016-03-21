package hdfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetFile(t *testing.T) {
	cfg := Config{}
	cfg.rootDir = "/user/cwilliams"
	cfg.hdfsHost = "hdedge001.marketdb.den01.pop"
	cfg.hdfsPort = 14000
	cfg.hdfsUser = "root"

	fs, err := New(cfg)

	f, err := fs.Open("data.csv", os.O_RDONLY)
	defer f.Close()

	if err != nil {
		t.Error(err)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Error(err)
	}

	if len(data) != 169 {
		t.Errorf("Expected data size: 169, actual data size %d", len(data))
	}

	if f.FileName() != "data.csv" {
		t.Error("Expected filename 'data.csv', got ", f.FileName())
	}

	fmt.Println("Filename:", f.Name())
}
