package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}