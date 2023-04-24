package db

import (
	"log"
	"panda/apigateway/ioutils"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func CreateNeo4jMainInstanceOrPanics(userName string, password string, host string, port string, schema string) *neo4j.Driver {
	neo4jDriver, err := neo4j.NewDriver(
		schema+host+":"+port,
		neo4j.BasicAuth(userName, password, ""),
	)

	// Check error in DB driver instantiation
	if err != nil {
		ioutils.PanicOnError(err)
	}

	// Verify Connectivity
	err = neo4jDriver.VerifyConnectivity()

	// If connectivity fails, handle the error
	if err != nil {
		ioutils.PanicOnError(err)
	}

	log.Println("Neo4j security database connection established successfully.")

	return &neo4jDriver
}
