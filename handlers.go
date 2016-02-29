package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cwillia9/ez-ftp/datastore"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET requests accepted on dl", http.StatusMethodNotAllowed)
		return
	}

	splt := strings.Split(r.URL.Path, "/")
	uuid := splt[len(splt)-1]

	fname, err := datastore.Select(uuid)
	if err != nil {
		// TODO(cwilliams): Doesn't exist or did we get a db err?
		fmt.Fprintf(w, "Record doesn't exist for uuid: "+uuid)
		return
	}
	log.Println("Serving file: " + fname)
	_, file := filepath.Split(fname)
	w.Header().Set("Content-Disposition", "attachment; filename="+file)
	http.ServeFile(w, r, fname)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "PUT" {
		http.Error(w, "Only POST requests accepted on ul", http.StatusMethodNotAllowed)
		return
	}
	file, handler, err := r.FormFile("uploadfile")
	log.Println(r.Form)
	if err != nil {
		log.Println(err)
		http.Error(w, "Expected uploadfile", http.StatusExpectationFailed)
		return
	}
	defer file.Close()

	//TODO(cwilliams): We want this to fail if the file already exists and let the
	//user know
	newfile := path.Join(cfg.RootDir, handler.Filename)
	f, err := os.OpenFile(newfile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			log.Printf("Tried creating file that already exists: " + newfile)
			http.Error(w, "File already exists", http.StatusConflict)
			return
		}
		log.Println(err)
		http.Error(w, "Oops....", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil && err != io.EOF {
		log.Println(err)
		http.Error(w, "Try again soon", http.StatusInternalServerError)
		return
	}
	log.Println("Successfully wrote file to: " + path.Join(cfg.RootDir, handler.Filename))

	randID := randomString(30)
	err = datastore.Insert(randID, cfg.RootDir, handler.Filename)
	// TODO(cwilliams): Did we fail because the entry is already there?
	if err != nil {
		log.Println(err)
		http.Error(w, "Try again soon", http.StatusInternalServerError)
	}
	log.Printf("Successfully stored new file %s/%s\n", cfg.RootDir, handler.Filename)
	fmt.Fprint(w, "new uuid: "+randID)
}

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
