package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	var err error
	d := &Database{}
	d.DB, err = sql.Open("postgres",
		"user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	return d
}
