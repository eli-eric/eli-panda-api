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

func TestHasSystemGraphFilters(t *testing.T) {
	assert.False(t, hasSystemGraphFilters(models.SystemGraphQueryOptions{}))
	assert.True(t, hasSystemGraphFilters(models.SystemGraphQueryOptions{Search: "abc"}))
	assert.True(t, hasSystemGraphFilters(models.SystemGraphQueryOptions{SystemLevels: []string{"TECHNOLOGY_UNIT"}}))
	assert.True(t, hasSystemGraphFilters(models.SystemGraphQueryOptions{SystemType: "Cooling"}))
	assert.True(t, hasSystemGraphFilters(models.SystemGraphQueryOptions{RelationshipTypes: []string{"HAS_SUBSYSTEM"}}))
}

func TestGetSystemGraphRelationshipTypes(t *testing.T) {
	custom := getSystemGraphRelationshipTypes(models.SystemGraphQueryOptions{RelationshipTypes: []string{"HAS_SUBSYSTEM"}})
	defaultTypes := getSystemGraphRelationshipTypes(models.SystemGraphQueryOptions{})

	assert.Equal(t, []string{"HAS_SUBSYSTEM"}, custom)
	assert.Equal(t, allowedSystemGraphRelationshipTypes, defaultTypes)
}

func TestIsSystemGraphRelationshipTypeAllowedByFilter(t *testing.T) {
	assert.True(t, isSystemGraphRelationshipTypeAllowedByFilter("HAS_SUBSYSTEM", nil))
	assert.True(t, isSystemGraphRelationshipTypeAllowedByFilter("HAS_SUBSYSTEM", []string{"HAS_SUBSYSTEM", "IS_POWERED_BY"}))
	assert.False(t, isSystemGraphRelationshipTypeAllowedByFilter("IS_COOLED_BY", []string{"HAS_SUBSYSTEM"}))
}

func TestBuildSystemGraphRelationshipStats(t *testing.T) {
	links := []models.SystemGraphLink{
		{Relationship: "HAS_SUBSYSTEM"},
		{Relationship: "HAS_SUBSYSTEM"},
		{Relationship: "IS_POWERED_BY"},
	}

	stats := buildSystemGraphRelationshipStats(links)

	assert.Equal(t, int64(2), stats["HAS_SUBSYSTEM"].Total)
	assert.Equal(t, int64(2), stats["HAS_SUBSYSTEM"].Returned)
	assert.False(t, stats["HAS_SUBSYSTEM"].HasMore)
	assert.Equal(t, int64(1), stats["IS_POWERED_BY"].Total)
}
