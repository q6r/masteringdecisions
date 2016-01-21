package main

import (
	"database/sql"
	"log"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

func InitDatabase() *gorp.DbMap {
	db, err := sql.Open("postgres",
		"user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(Person{}, "person").SetKeys(true, "person_id")

	// TODO : need to remove decision autoincr ???
	dbmap.AddTableWithName(Decision{}, "decision").SetKeys(true, "decision_id")

	dbmap.AddTableWithName(Ballot{}, "ballot").SetKeys(true, "ballot_id")

	// TODO : need to remove criterion autoincr ???
	dbmap.AddTableWithName(Criterion{}, "criterion").SetKeys(true, "criterion_id")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln(err)
	}
	return dbmap
}
