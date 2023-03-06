package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func MigrateNeo4jMainInstance(userName string, password string, host string, port string) {
	m, err := migrate.New(
		"file://db/neo4j/migrations",
		"neo4j://"+userName+":"+password+"@"+host+":"+port+"?x-multi-statement=true")
	// if there is a db error log and shut down
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Applay migrations...")
	// if there is an error in migrations log and shut down, if its successful or there are no changes we can continue
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
	log.Println("Migrations OK")
}
