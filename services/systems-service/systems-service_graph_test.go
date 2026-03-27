package systemsService

import (
	"errors"
	"sort"
	"testing"

	"panda/apigateway/helpers"
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

func TestBuildSystemGraphRelationshipTotalMap(t *testing.T) {
	relationshipTypes := []string{"HAS_SUBSYSTEM", "IS_POWERED_BY", "IS_COOLED_BY"}
	counts := []systemGraphRelationshipCount{
		{Relationship: "HAS_SUBSYSTEM", Total: 5},
		{Relationship: "IS_POWERED_BY", Total: 2},
	}

	result := buildSystemGraphRelationshipTotalMap(relationshipTypes, counts)

	assert.Equal(t, int64(5), result["HAS_SUBSYSTEM"])
	assert.Equal(t, int64(2), result["IS_POWERED_BY"])
	assert.Equal(t, int64(0), result["IS_COOLED_BY"])
}

func TestCalculateSystemGraphHiddenLinksTotal(t *testing.T) {
	assert.Equal(t, int64(90), calculateSystemGraphHiddenLinksTotal(100, 10))
	assert.Equal(t, int64(0), calculateSystemGraphHiddenLinksTotal(10, 10))
	assert.Equal(t, int64(0), calculateSystemGraphHiddenLinksTotal(5, 10))
}

func TestInvalidSystemGraphInput(t *testing.T) {
	err := invalidSystemGraphInput("invalid systemType")

	assert.True(t, errors.Is(err, helpers.ERR_INVALID_INPUT))
	assert.Equal(t, "invalid systemType: INVALID_INPUT", err.Error())
}

func TestSortHierarchyRoots_PrioritizesRootsWithChildren(t *testing.T) {
	technologyUnit := "TECHNOLOGY_UNIT"
	keySystems := "KEY_SYSTEMS"
	trash := "TRASH"

	roots := []string{"root-b", "root-parent", "root-a"}
	nodesByUID := map[string]*models.SystemHierarchyNode{
		"root-a": {
			UID:         "root-a",
			Name:        "Alpha Root",
			SystemLevel: &trash,
		},
		"root-b": {
			UID:         "root-b",
			Name:        "Beta Root",
			SystemLevel: &keySystems,
		},
		"root-parent": {
			UID:         "root-parent",
			Name:        "Parent Root",
			SystemLevel: &technologyUnit,
		},
	}
	childrenByUID := map[string]map[string]struct{}{
		"root-parent": {
			"child-1": {},
		},
	}

	sortHierarchyRoots(roots, nodesByUID, childrenByUID)

	assert.Equal(t, []string{"root-parent", "root-b", "root-a"}, roots)
}

func TestCompareHierarchyNodes_SortsBySystemLevelOrderThenLevelThenName(t *testing.T) {
	technologyUnit := "TECHNOLOGY_UNIT"
	keySystems := "KEY_SYSTEMS"
	subsystems := "SUBSYSTEMS_AND_PARTS"
	trash := "TRASH"

	nodes := []*models.SystemHierarchyNode{
		{UID: "trash", Name: "Zulu", SystemLevel: &trash},
		{UID: "subsystems", Name: "Alpha", SystemLevel: &subsystems},
		{UID: "key", Name: "Zulu", SystemLevel: &keySystems},
		{UID: "technology", Name: "Zulu", SystemLevel: &technologyUnit},
	}

	sort.SliceStable(nodes, func(i, j int) bool {
		return compareHierarchyNodes(nodes[i], nodes[j])
	})

	assert.Equal(t, []string{"technology", "key", "subsystems", "trash"}, []string{nodes[0].UID, nodes[1].UID, nodes[2].UID, nodes[3].UID})
}
