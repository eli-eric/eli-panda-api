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
