package systemsService

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSystemGraphLinksByUidAndTypeQuery_HasStablePagingOrder(t *testing.T) {
	query := GetSystemGraphLinksByUidAndTypeQuery("sys-1", "B", "IS_POWERED_BY", 20, 10)

	assert.Equal(t, "links", query.ReturnAlias)
	assert.Contains(t, query.Query, "ORDER BY link.relId ASC")
	assert.Contains(t, query.Query, "SKIP $offset")
	assert.Contains(t, query.Query, "LIMIT $limit")
	assert.Equal(t, 20, query.Parameters["offset"])
	assert.Equal(t, 10, query.Parameters["limit"])
}

func TestGetSystemGraphLinksByUidAndTypeCountQuery(t *testing.T) {
	query := GetSystemGraphLinksByUidAndTypeCountQuery("sys-1", "B", "IS_POWERED_BY")

	assert.Equal(t, "totalCount", query.ReturnAlias)
	assert.Contains(t, query.Query, "count(DISTINCT relationId) as totalCount")
}

func TestGetSystemGraphNodesByUidsQuery(t *testing.T) {
	query := GetSystemGraphNodesByUidsQuery([]string{"sys-1", "sys-2"}, "B")

	assert.Equal(t, "nodes", query.ReturnAlias)
	assert.Contains(t, query.Query, "WHERE sys.uid IN $uids")
	assert.Contains(t, query.Query, "ORDER BY nodes.uid ASC")
}

func TestGetSystemGraphFilteredLinksByUidQuery(t *testing.T) {
	query := GetSystemGraphFilteredLinksByUidQuery("sys-1", "B", []string{"HAS_SUBSYSTEM", "IS_POWERED_BY"}, "pump", []string{"TECHNOLOGY_UNIT"}, "Cooling")

	assert.Equal(t, "links", query.ReturnAlias)
	assert.Contains(t, query.Query, "uid: toString(id(r))")
	assert.Contains(t, query.Query, "ORDER BY link.relId ASC")
	assert.Contains(t, query.Query, "other.systemLevel IN $systemLevels")
	assert.NotContains(t, query.Query, "center.systemLevel IN $systemLevels")
	assert.Equal(t, true, query.Parameters["hasSearch"])
	assert.Equal(t, true, query.Parameters["hasSystemLevels"])
	assert.Equal(t, true, query.Parameters["hasSystemType"])
}

func TestGetSystemGraphLinksByUidAndTypeFilteredQuery_HasStablePagingOrder(t *testing.T) {
	query := GetSystemGraphLinksByUidAndTypeFilteredQuery("sys-1", "B", "IS_POWERED_BY", "pump", []string{"TECHNOLOGY_UNIT"}, "Cooling", 20, 10)

	assert.Equal(t, "links", query.ReturnAlias)
	assert.Contains(t, query.Query, "ORDER BY link.relId ASC")
	assert.Contains(t, query.Query, "SKIP $offset")
	assert.Contains(t, query.Query, "LIMIT $limit")
	assert.Contains(t, query.Query, "other.systemLevel IN $systemLevels")
	assert.NotContains(t, query.Query, "center.systemLevel IN $systemLevels")
	assert.Equal(t, 20, query.Parameters["offset"])
	assert.Equal(t, 10, query.Parameters["limit"])
	assert.Equal(t, true, query.Parameters["hasSearch"])
}

func TestGetSystemGraphLinksByUidAndTypeFilteredCountQuery(t *testing.T) {
	query := GetSystemGraphLinksByUidAndTypeFilteredCountQuery("sys-1", "B", "IS_POWERED_BY", "pump", []string{"TECHNOLOGY_UNIT"}, "Cooling")

	assert.Equal(t, "totalCount", query.ReturnAlias)
	assert.Contains(t, query.Query, "count(DISTINCT relationId) as totalCount")
	assert.Equal(t, true, query.Parameters["hasSystemLevels"])
	assert.Equal(t, true, query.Parameters["hasSystemType"])
}

func TestSystemTypeNameExistsInFacilityQuery(t *testing.T) {
	query := SystemTypeNameExistsInFacilityQuery("Cooling", "B")

	assert.Equal(t, "result", query.ReturnAlias)
	assert.Contains(t, query.Query, "toLower(st.name) = toLower($systemType)")
}
