package main

import (
	"crypto/hmac"
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

	"github.com/cwillia9/ez-ftp/authentication"
	"github.com/cwillia9/ez-ftp/datastore"
	"github.com/cwillia9/ez-ftp/domain"
)

/*
This whole file needs to be refactored to not have all of the individual handlers have so much
logic. This effectively makes them untestable
*/

// TODO(cwilliams): Make this better. A proper HMAC implementation should use a signing string using
// something like the following
// StringToSign = HTTP-Verb + "\n" +
// 	Content-MD5 + "\n" +
// 	Content-Type + "\n" +
// 	Date + "\n" +
// 	CanonicalizedAmzHeaders +
// 	CanonicalizedResource;
// Signature = Base64( HMAC-SHA1( YourSecretAccessKeyID, UTF-8-Encoding-Of( StringToSign ) ) );
// Authorization = "AWS" + " " + AWSAccessKeyId + ":" + Signature;
//
// For now the authentication will just use a public key and shared secret
func hmacAuthentication(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		auth, ok := r.Header["Authorization"]
		if ok == false {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			log.Println("upload rejected: no Authorization given")
			return
		}

		authstring := auth[0]

		// auth should look like <public_key>:<signature>
		split := strings.Split(authstring, ":")
		if len(split) != 2 {
			http.Error(w, "Malformed Authorization", http.StatusUnauthorized)
			log.Println("upload rejected: malformed Authorization. Authorization: " + authstring)
			return
		}

		key, actualEncoding := split[0], split[1]

		passhash, err := datastore.SelectUser(key)
		if err != nil {
			http.Error(w, "No match found for public key "+key, http.StatusUnauthorized)
			log.Println("upload rejected: user not found. user: " + key)
			return
		}

		expectation := authentication.ComputeHmac1(key, passhash)

		if match := hmac.Equal([]byte(actualEncoding), []byte(expectation)); match == false {
			http.Error(w, "Authorization didn't match", http.StatusUnauthorized)
			log.Println("upload rejected: supplied mac encoding did not match expected for user " + key)
			return
		}

		fn(w, r)

	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request, system domain.FileSystem) {
	if r.Method != "GET" {
		http.Error(w, "Only GET requests accepted on dl", http.StatusMethodNotAllowed)
		return
	}

	splt := strings.Split(r.URL.Path, "/")
	uuid := splt[len(splt)-1]

	rootDir, path, err := datastore.SelectFile(uuid)
	if err != nil {
		// TODO(cwilliams): Doesn't exist or did we get a db err?
		fmt.Fprintf(w, "Record doesn't exist for uuid: "+uuid)
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}

	if rootDir != cfg.RootDir {
		log.Println("Attempting to download from from wrong root dir. Attempted", rootDir, "but running", cfg.RootDir)
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}

	f, err := system.Open(path, os.O_RDONLY)
	if err != nil {
		log.Println("Error opening file", err)
		http.Error(w, "Error retreiving file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	log.Println("Serving file: " + path)
	_, file := filepath.Split(path)
	w.Header().Set("Content-Disposition", "attachment; filename="+file)
	http.ServeFile(w, r, f.Name())
}

func uploadHandler(w http.ResponseWriter, r *http.Request, fs domain.FileSystem) {
	// Must be using a POST/PUT method
	if r.Method != "POST" && r.Method != "PUT" {
		http.Error(w, "Only POST requests accepted on /ul/", http.StatusMethodNotAllowed)
		return
	}

	// We expect the file to be called 'uploadfile'
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		http.Error(w, "Expected uploadfile", http.StatusExpectationFailed)
		return
	}
	defer file.Close()

	// TODO(cwilliams): Eventually we want to expose
	// some kind of admin api to view directory structure

	// Note: the call to FormFile above would already have parsed the form
	desiredPaths, ok := r.MultipartForm.Value["path"]
	if ok != true {
		log.Println("No path specified")
		http.Error(w, "expected path", http.StatusExpectationFailed)
		return
	}

	// We only support a single path
	desiredPath := desiredPaths[0]
	fmt.Println("desiredpath", desiredPath)

	newfile := path.Join(desiredPath, handler.Filename)

	if err = os.MkdirAll(path.Join(cfg.RootDir, desiredPath), 0777); err != nil {
		http.Error(w, "path failure", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// O_EXCL ensures that if the file already exists we will not overwrite it
	f, err := fs.Open(newfile, os.O_WRONLY|os.O_CREATE|os.O_EXCL)
	if err != nil {
		if os.IsExist(err) {
			// TODO(cwilliams): Add 'overwrite' flag functionality
			log.Printf("Tried creating file that already exists: " + newfile)
			http.Error(w, "File already exists", http.StatusConflict)
			return
		}
		log.Println(err)
		// This sucks....why couldn't we open a file?
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

	randID := randomString(32)
	err = datastore.InsertFile(randID, cfg.RootDir, path.Join(desiredPath, handler.Filename))
	// TODO(cwilliams): Did we fail because the entry is already there?
	// TODO(cwilliams): Did we fail beause the random string already exists? If so we should try again
	if err != nil {
		log.Println(err)
		http.Error(w, "Try again soon", http.StatusInternalServerError)
	}
	log.Printf("Successfully stored new file %s/%s\n", cfg.RootDir, handler.Filename)
	fmt.Fprint(w, "new uuid: "+randID)
}

func makeFsHandler(fn func(http.ResponseWriter, *http.Request, domain.FileSystem), fs domain.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, fs)
	}
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
