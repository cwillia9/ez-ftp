package datastore

import (
	"database/sql"
	"log"
	pathlib "path"

	// Set up to only use sqlite
	_ "github.com/mattn/go-sqlite3"
)

// DB is a global db handle
var db *sql.DB

var createTable = "CREATE TABLE IF NOT EXISTS `pathmappings` ( " +
	"`randstring` varchar(64) PRIMARY KEY, " +
	"`rootdir` varchar(512), " +
	"`path` VARCHAR(64) NULL, " +
	"`created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP " +
	");"

// InitDB instantiates and then validates new db connection
func InitDB(sqlitePath string) error {
	var err error
	db, err = sql.Open("sqlite3", sqlitePath)
	if err != nil {
		return err
	}
	_, err = db.Exec(createTable)
	if err != nil {
		return err
	}

	return nil
}

// Insert new mapping
func Insert(id, root, path string) error {
	log.Printf("ID: %s root: %s path: %s\n", id, root, path)
	log.Println(db)
	stmt, err := db.Prepare("INSERT INTO pathmappings" +
		"(randstring, rootdir, path) " +
		"values(?,?,?);")
	if err != nil {
		return err
	}
	stmt.Exec(id, root, path)
	if err != nil {
		return err
	}
	return nil
}

// Select path from id
func Select(id string) (string, error) {
	var rootdir, path string
	err := db.QueryRow(
		"SELECT rootdir, path FROM pathmappings where randstring = ?;", id).Scan(&rootdir, &path)
	if err != nil {
		return "", err
	}
	return pathlib.Join(rootdir, path), nil
}
