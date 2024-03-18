package db

import (
	"database/sql"
	"log"
)

func Connect(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
