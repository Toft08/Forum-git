package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver; the blank import ensures the driver is included "_" is important!!!
)

// InitDB initializes the SQLite database and returns a database connection object
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Println("Failed to initialize the database")
		log.Fatal(err)

	}
	return db
}
