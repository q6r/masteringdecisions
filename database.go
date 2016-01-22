package main

import (
	"database/sql"
	"log"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

// InitDatabase initalizes the postgres database
// builds the schema of the table
// forigen key restriction is not handled in here but they're
// handled in each objects Save and Destroy methods
func InitDatabase() *gorp.DbMap {
	db, err := sql.Open("postgres",
		"user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(Person{}, "person").SetKeys(true, "person_id")
	dbmap.AddTableWithName(Decision{}, "decision").SetKeys(true, "decision_id")
	dbmap.AddTableWithName(Ballot{}, "ballot").SetKeys(true, "ballot_id")
	dbmap.AddTableWithName(Criterion{}, "criterion").SetKeys(true, "criterion_id")
	dbmap.AddTableWithName(Vote{}, "vote")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln(err)
	}
	return dbmap
}
