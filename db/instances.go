package db

import (
	"context"
	"panda/apigateway/ioutils"

	"github.com/rs/zerolog/log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateNeo4jMainInstanceOrPanics(userName string, password string, host string, port string, schema string) *neo4j.DriverWithContext {
	ctx := context.Background()
	neo4jDriver, err := neo4j.NewDriverWithContext(
		schema+host+":"+port,
		neo4j.BasicAuth(userName, password, ""),
	)

	// Check error in DB driver instantiation
	if err != nil {
		ioutils.PanicOnError(err)
	}

	// Verify Connectivity
	err = neo4jDriver.VerifyConnectivity(ctx)

	// If connectivity fails, handle the error
	if err != nil {
		ioutils.PanicOnError(err)
	}

	log.Info().Msg("Neo4j security database connection established successfully.")

	return &neo4jDriver
}
