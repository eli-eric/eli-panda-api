package publicationsservice

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"panda/apigateway/services/testsetup"
)

func TestGetPublicationByUid(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// Insert test data
	_, err := testsetup.TestSession.Run(`CREATE (p:Publication {uid: "test123"}) RETURN p`, nil)
	assert.NoError(t, err)

	// Run the actual test
	result, err := service.GetPublicationByUid("test123")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "test123", result.Uid)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (p:Publication {uid: "test123"}) DETACH DELETE p`, nil)
	assert.NoError(t, err)
}

func TestGetPublicationByUidNotFound(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// Run the actual test
	_, err := service.GetPublicationByUid("test123")

	// Assertions
	assert.Error(t, err)

	_, cleanupErr := testsetup.TestSession.Run(`MATCH (p:Publication {uid: "test123"}) DETACH DELETE p`, nil)
	assert.NoError(t, cleanupErr)
}
