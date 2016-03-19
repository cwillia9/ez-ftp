package datastore

import (
	"database/sql"
	"log"

	// Set up to only use sqlite
	_ "github.com/mattn/go-sqlite3"
)

// DB is a global db handle
var db *sql.DB

var createPathMappingsTable = "CREATE TABLE IF NOT EXISTS `pathmappings` ( " +
	"`randstring` varchar(64) PRIMARY KEY, " +
	"`rootdir` varchar(512), " +
	"`path` VARCHAR(64) NULL, " +
	"`created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP " +
	");"

var createAuthorizedUsersTable = "CREATE TABLE IF NOT EXISTS `users` ( " +
	"`user` varchar(64) PRIMARY KEY, " +
	"`pass_hash` varchar(64), " +
	"`created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP " +
	");"

// InitDB instantiates and then validates new db connection
func InitDB(sqlitePath string) error {
	var err error
	db, err = sql.Open("sqlite3", sqlitePath)
	if err != nil {
		return err
	}
	_, err = db.Exec(createPathMappingsTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(createAuthorizedUsersTable)
	if err != nil {
		return err
	}

	return nil
}

// InsertFile new mapping
func InsertFile(id, root, path string) error {
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

// SelectFile path from id
func SelectFile(id string) (string, string, error) {
	var rootdir, path string
	err := db.QueryRow(
		"SELECT rootdir, path FROM pathmappings where randstring = ?;", id).Scan(&rootdir, &path)
	if err != nil {
		return "", "", err
	}
	return rootdir, path, nil
}

// SelectUser from db
func SelectUser(user string) (string, error) {
	var passhash string
	err := db.QueryRow("SELECT pass_hash FROM users where user = ?;", user).Scan(&passhash)
	if err != nil {
		return "", err
	}
	return passhash, nil
}
