package publicationsservice

import (
	"log"
	"os"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
)

var testDriver neo4j.Driver
var testSession neo4j.Session

func TestMain(m *testing.M) {
	var err error
	testDriver, err = neo4j.NewDriver("bolt://localhost:7682", neo4j.BasicAuth("neo4j", "xxxxx", ""))
	if err != nil {
		log.Fatal("Failed to connect to test Neo4j:", err)
	}

	testSession = testDriver.NewSession(neo4j.SessionConfig{})

	code := m.Run()

	// Cleanup
	testDriver.Close()
	os.Exit(code)
}

func TestGetPublicationByUid(t *testing.T) {
	service := NewPublicationsService(&testDriver, "", "")

	// Insert test data
	_, err := testSession.Run(`CREATE (p:Publication {uid: "test123"}) RETURN p`, nil)
	assert.NoError(t, err)

	// Run the actual test
	result, err := service.GetPublicationByUid("test123")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "test123", result.Uid)
}
