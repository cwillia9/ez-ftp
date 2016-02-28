package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type FSMap struct {
	mu      sync.Mutex
	Mapping map[string]string
}

var fsMap FSMap

func init() {
	fsMap = FSMap{sync.Mutex{}, make(map[string]string)}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET requests accepted on dl", http.StatusMethodNotAllowed)
		return
	}

	splt := strings.Split(r.URL.Path, "/")
	uuid := splt[len(splt)-1]

	fsMap.mu.Lock()
	fname, ok := fsMap.Mapping[string(uuid)]
	log.Println(fsMap.Mapping)
	fsMap.mu.Unlock()
	if ok == false {
		fmt.Fprintf(w, "Record doesn't exist for uuid: "+uuid)
		return
	}
	fmt.Fprint(w, "Requested uuid: "+uuid+" associated file: "+fname)
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
	f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(f, file)
	if err != nil && err != io.EOF {
		log.Println(err)
		http.Error(w, "Try again soon", http.StatusInternalServerError)
		return
	}

	randID := randomString(30)
	fsMap.mu.Lock()
	fsMap.Mapping[randID] = handler.Filename
	log.Println(fsMap.Mapping)
	fsMap.mu.Unlock()
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

func main() {
	http.HandleFunc("/dl/", downloadHandler)
	http.HandleFunc("/ul/", uploadHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
