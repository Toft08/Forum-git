package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver; the blank import ensures the driver is included "_" is important!!!
)

// initDB initializes the SQLite database and returns a database connection object.
func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
		_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
	`)
	}
	return db
}
