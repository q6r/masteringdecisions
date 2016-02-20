package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/astaxie/beego/config"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

// InitDatabase initalizes the postgres database
// builds the schema of the table
// forigen key restriction is not handled in here but they're
// handled in each objects Save and Destroy methods
func InitDatabase(conf config.Configer) *gorp.DbMap {
	dbsrc := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s",
		conf.String("database::user"),
		conf.String("database::name"),
		conf.String("database::sslmode"),
		conf.String("database::password"))

	db, err := sql.Open("postgres", dbsrc)
	if err != nil {
		log.Fatalf("Unable to connect to postgres : %#v\n", err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(Person{}, "person").SetKeys(true, "person_id")
	dbmap.AddTableWithName(Decision{}, "decision").SetKeys(true, "decision_id")
	dbmap.AddTableWithName(Ballot{}, "ballot").SetKeys(true, "ballot_id")
	dbmap.AddTableWithName(Alternative{}, "alternative").SetKeys(true, "alternative_id")
	dbmap.AddTableWithName(Criterion{}, "criterion").SetKeys(true, "criterion_id")
	dbmap.AddTableWithName(Vote{}, "vote")
	dbmap.AddTableWithName(Rating{}, "rating")

	if err = dbmap.CreateTablesIfNotExists(); err != nil {
		log.Fatalln(err)
	}

	// Always create an admin account
	_, err = dbmap.Exec("DELETE FROM person WHERE person_id=0")
	if err != nil {
		log.Fatalln(err)
	}

	hashed := HashPassword(conf.String("admin::password"))
	_, err = dbmap.Exec("INSERT INTO person VALUES(0,$1,$2,$3,$4)",
		conf.String("admin::email"),
		hashed,
		conf.String("admin::name_first"),
		conf.String("admin::name_last"))
	if err != nil {
		log.Fatalln(err)
	}

	return dbmap
}
