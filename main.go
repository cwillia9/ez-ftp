package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cwillia9/ez-ftp/config"
	"github.com/cwillia9/ez-ftp/datastore"
	"github.com/cwillia9/ez-ftp/localfs"
)

const (
	defaultTCPAddr    = "localhost:9999"
	defaultSqlitePath = "/usr/local/var/ez-ftp/sqlite.db"
)

var (
	cfg     *config.T
	rootDir string
)

func init() {
	cfg = &config.T{}

	flag.StringVar(&cfg.TCPAddr, "tcpAddr", defaultTCPAddr,
		"TCP address that the HTTP API should listen on")
	flag.StringVar(&cfg.RootDir, "rootDir", "",
		"(REQUIRED) Root directory to serve ftp files from")
	flag.StringVar(&cfg.SqlitePath, "sqlitePath", defaultSqlitePath,
		"Path to sqlite db")

	flag.Parse()

	if cfg.RootDir == "" {
		fmt.Println("rootDir option required")
		os.Exit(1)
	}

}

func main() {
	err := datastore.InitDB(cfg.SqlitePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully init'd database")

	system, err := localfs.New(cfg)
	if err != nil {
		log.Println("Failed instantiating file system")
	}

	http.HandleFunc("/dl/", downloadHandler)
	http.HandleFunc("/ul/", hmacAuthentication(uploadHandler))
	log.Fatal(http.ListenAndServe(":9999", nil))
}
