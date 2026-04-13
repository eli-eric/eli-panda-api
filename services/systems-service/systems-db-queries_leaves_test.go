package systemsService

import (
	"testing"

	"panda/apigateway/helpers"

	"github.com/stretchr/testify/assert"
)

func TestGetSystemLeavesByParentUIDQuery_NoFilters(t *testing.T) {
	query := GetSystemLeavesByParentUIDQuery("parent-uid", "FAC", "", nil, nil, nil)

	assert.Equal(t, "systems", query.ReturnAlias)
	assert.Contains(t, query.Query, "OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone)")
	assert.Contains(t, query.Query, "OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st)")
	assert.Contains(t, query.Query, "OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc)")
	assert.Contains(t, query.Query, "OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsible)")
	assert.Contains(t, query.Query, "OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp)")
	assert.Contains(t, query.Query, "OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)")
	assert.Contains(t, query.Query, "SKIP $skip LIMIT $limit")
	assert.Equal(t, "parent-uid", query.Parameters["parentUID"])
	assert.Equal(t, "FAC", query.Parameters["facilityCode"])
}

func TestGetSystemLeavesByParentUIDQuery_SystemNameFilter(t *testing.T) {
	filters := []helpers.ColumnFilter{{Id: "name", Value: "Pump"}}
	query := GetSystemLeavesByParentUIDQuery("parent-uid", "FAC", "", nil, nil, &filters)

	assert.Contains(t, query.Query, "toLower(sys.name) CONTAINS $filterName")
	assert.Equal(t, "pump", query.Parameters["filterName"])
}

func TestGetSystemLeavesByParentUIDQuery_ZoneFilter(t *testing.T) {
	filters := []helpers.ColumnFilter{{Id: "zone", Value: map[string]interface{}{"uid": "zone-1", "name": "Zone A"}}}
	query := GetSystemLeavesByParentUIDQuery("parent-uid", "FAC", "", nil, nil, &filters)

	assert.Contains(t, query.Query, "MATCH (sys)-[:HAS_ZONE]->(zone) WHERE zone.uid = $filterZone")
	assert.NotContains(t, query.Query, "OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone)")
	assert.Equal(t, "zone-1", query.Parameters["filterZone"])
}

func TestGetSystemLeavesByParentUIDQuery_ItemUsageFilter(t *testing.T) {
	filters := []helpers.ColumnFilter{{Id: "itemUsage", Value: []interface{}{"usage-1", "usage-2"}}}
	query := GetSystemLeavesByParentUIDQuery("parent-uid", "FAC", "", nil, nil, &filters)

	assert.Contains(t, query.Query, "MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)")
	assert.Contains(t, query.Query, "MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage) WHERE itemUsage.uid IN $filterItemUsage")
	assert.Equal(t, &[]string{"usage-1", "usage-2"}, query.Parameters["filterItemUsage"])
}

func TestGetSystemLeavesByParentUIDQuery_PriceFilterAlone(t *testing.T) {
	min := 100.0
	max := 500.0
	filters := []helpers.ColumnFilter{{Id: "price", Value: map[string]interface{}{"min": min, "max": max}}}
	query := GetSystemLeavesByParentUIDQuery("parent-uid", "FAC", "", nil, nil, &filters)

	assert.Contains(t, query.Query, "MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)")
	assert.Contains(t, query.Query, "MATCH (physicalItem)<-[ol:HAS_ORDER_LINE]-(order)")
	assert.Contains(t, query.Query, "ol.price >= $filterPriceFrom")
}

func TestGetSystemLeavesByParentUIDQuery_OrderNameCaseInsensitive(t *testing.T) {
	filters := []helpers.ColumnFilter{{Id: "orderName", Value: "Main Order"}}
	query := GetSystemLeavesByParentUIDQuery("parent-uid", "FAC", "", nil, nil, &filters)

	assert.Contains(t, query.Query, "toLower(order.name) CONTAINS $filterOrderName")
	assert.Equal(t, "main order", query.Parameters["filterOrderName"])
}

func TestGetSystemLeavesByParentUIDQuery_DynamicPropFilterAlone(t *testing.T) {
	filters := []helpers.ColumnFilter{{Id: "prop-uid-1", Value: "test-value", Type: "text", PropType: "PHYSICAL_ITEM"}}
	query := GetSystemLeavesByParentUIDQuery("parent-uid", "FAC", "", nil, nil, &filters)

	assert.Contains(t, query.Query, "MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)")
	assert.Contains(t, query.Query, "MATCH(prop{uid: $propUID0})<-[pv]-(physicalItem)")
	assert.Contains(t, query.Query, "toLower(pv.value) contains $propFilterVal0")
}

func TestGetSystemLeavesByParentUIDQuery_ResponsibleVariableName(t *testing.T) {
	query := GetSystemLeavesByParentUIDQuery("parent-uid", "FAC", "", nil, nil, nil)

	assert.Contains(t, query.Query, "OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsible)")
	assert.NotContains(t, query.Query, "responsilbe")
}

func TestGetSystemLeavesByParentUIDQuery_Sorting(t *testing.T) {
	sorting := []helpers.Sorting{{ID: "importance", DESC: true}}
	query := GetSystemLeavesByParentUIDQuery("parent-uid", "FAC", "", nil, &sorting, nil)

	assert.Contains(t, query.Query, "ORDER BY toLower(systems.importance.name) DESC")
}

// Count query tests

func TestGetSystemLeavesByParentUIDCountQuery_NoFilters(t *testing.T) {
	query := GetSystemLeavesByParentUIDCountQuery("parent-uid", "FAC", "", nil)

	assert.Equal(t, "count", query.ReturnAlias)
	assert.Contains(t, query.Query, "RETURN COUNT(DISTINCT sys) as count")
	assert.NotContains(t, query.Query, "SKIP")
}

func TestGetSystemLeavesByParentUIDCountQuery_ItemUsageFilter(t *testing.T) {
	filters := []helpers.ColumnFilter{{Id: "itemUsage", Value: []interface{}{"usage-1"}}}
	query := GetSystemLeavesByParentUIDCountQuery("parent-uid", "FAC", "", &filters)

	assert.Contains(t, query.Query, "MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)")
	assert.Contains(t, query.Query, "MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage) WHERE itemUsage.uid IN $filterItemUsage")
}

func TestGetSystemLeavesByParentUIDCountQuery_PriceFilterAlone(t *testing.T) {
	min := 100.0
	filters := []helpers.ColumnFilter{{Id: "price", Value: map[string]interface{}{"min": min}}}
	query := GetSystemLeavesByParentUIDCountQuery("parent-uid", "FAC", "", &filters)

	assert.Contains(t, query.Query, "MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)")
	assert.Contains(t, query.Query, "ol.price >= $filterPriceFrom")
}

func TestGetSystemLeavesByParentUIDCountQuery_DynamicPropVariableScoping(t *testing.T) {
	filters := []helpers.ColumnFilter{{Id: "prop-uid-1", Value: "test", Type: "text", PropType: "PHYSICAL_ITEM"}}
	query := GetSystemLeavesByParentUIDCountQuery("parent-uid", "FAC", "", &filters)

	assert.Contains(t, query.Query, "WITH sys, physicalItem, catalogueItem MATCH(prop{uid:")
	assert.NotContains(t, query.Query, "WITH sys MATCH(prop{uid:")
}

func TestGetSystemLeavesByParentUIDCountQuery_CombinedPriceAndCatalogueName(t *testing.T) {
	min := 50.0
	filters := []helpers.ColumnFilter{
		{Id: "price", Value: map[string]interface{}{"min": min}},
		{Id: "catalogueName", Value: "sensor"},
	}
	query := GetSystemLeavesByParentUIDCountQuery("parent-uid", "FAC", "", &filters)

	// Both physicalItem and catalogueItem must stay in scope
	assert.Contains(t, query.Query, "WITH sys, physicalItem, catalogueItem MATCH (physicalItem)<-[ol:HAS_ORDER_LINE]-(order)")
	assert.Contains(t, query.Query, "WITH sys, physicalItem, catalogueItem WHERE toLower(catalogueItem.name) CONTAINS $filterCatalogueName")
}

func TestGetSystemLeavesByParentUIDCountQuery_SystemLevelFilterBeforeCategory(t *testing.T) {
	filters := []helpers.ColumnFilter{
		{Id: "systemLevel", Value: []interface{}{"KEY_SYSTEMS"}},
		{Id: "category", Value: map[string]interface{}{"uid": "cat-1", "name": "Cat"}},
	}
	query := GetSystemLeavesByParentUIDCountQuery("parent-uid", "FAC", "", &filters)

	// System-level filter should not drop category vars
	assert.Contains(t, query.Query, "sys.systemLevel IN $filterSystemLevel")
	assert.Contains(t, query.Query, "CatalogueCategory{uid:$filterCatalogueCategory}")
}
