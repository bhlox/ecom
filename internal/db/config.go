package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionDBstring string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionDBstring)
	if err != nil {
		return nil, err
	}
	errs := db.Ping()
	if errs != nil {
		return nil, errs
	}
	log.Println("DB 🎯: Successfully connected!")
	return db, nil
}
