package systemsService

import (
	"testing"

	"panda/apigateway/services/systems-service/models"

	"github.com/stretchr/testify/assert"
)

func TestIsAllowedSystemGraphRelationshipType(t *testing.T) {
	assert.True(t, isAllowedSystemGraphRelationshipType("HAS_SUBSYSTEM"))
	assert.True(t, isAllowedSystemGraphRelationshipType("IS_COOLED_BY"))
	assert.False(t, isAllowedSystemGraphRelationshipType("INVALID_REL"))
}

func TestCollectSystemGraphNodeUIDs(t *testing.T) {
	links := []models.SystemGraphLink{
		{Source: "sys-2", Target: "sys-3"},
		{Source: "sys-1", Target: "sys-2"},
		{Source: "sys-3", Target: "sys-4"},
	}

	result := collectSystemGraphNodeUIDs("sys-1", links)

	assert.Equal(t, []string{"sys-1", "sys-2", "sys-3", "sys-4"}, result)
}
