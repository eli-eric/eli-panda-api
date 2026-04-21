package catalogueService

import (
	"encoding/json"
	"testing"

	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	codebookModels "panda/apigateway/services/codebook-service/models"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T { return &v }

func newOriginalItem() *models.CatalogueItem {
	desc := "old desc"
	return &models.CatalogueItem{
		Uid:             "item-1",
		Name:            "Old Name",
		CatalogueNumber: "CN-001",
		Description:     &desc,
		Category:        codebookModels.Codebook{UID: "cat-old", Name: "Old Cat"},
		Supplier:        &codebookModels.Codebook{UID: "supp-old", Name: "Old Supp"},
		Details: []models.CatalogueItemDetail{
			{Property: models.CatalogueCategoryProperty{UID: "prop-A", Name: "Voltage", Type: models.CatalogueCategoryPropertyType{Code: "text"}}, Value: "12"},
			{Property: models.CatalogueCategoryProperty{UID: "prop-B", Name: "Weight", Type: models.CatalogueCategoryPropertyType{Code: "number"}}, Value: "50"},
			{Property: models.CatalogueCategoryProperty{UID: "prop-C", Name: "Note", Type: models.CatalogueCategoryPropertyType{Code: "text"}}, Value: "keep"},
		},
	}
}

func parseChanges(t *testing.T, q helpers.DatabaseQuery) []map[string]interface{} {
	raw, ok := q.Parameters["changes"].(string)
	assert.True(t, ok, "changes parameter must be a string")
	var parsed []map[string]interface{}
	err := json.Unmarshal([]byte(raw), &parsed)
	assert.NoError(t, err)
	return parsed
}

func TestPatchCatalogueItemQuery_NameOnly(t *testing.T) {
	fields := &models.PatchCatalogueItemFields{Name: ptr("New Name")}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	assert.Equal(t, "uid", q.ReturnAlias)
	assert.Contains(t, q.Query, "SET item.name = $name")
	assert.NotContains(t, q.Query, "SET item.description")
	assert.NotContains(t, q.Query, "SET item.catalogueNumber")
	assert.NotContains(t, q.Query, "HAS_CATALOGUE_PROPERTY")
	assert.NotContains(t, q.Query, "HAS_SUPPLIER")
	assert.NotContains(t, q.Query, "BELONGS_TO_CATEGORY")
	assert.Contains(t, q.Query, `action: "PATCH"`)
	assert.Contains(t, q.Query, "changes: $changes")
	assert.Equal(t, "New Name", q.Parameters["name"])

	changes := parseChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "name", changes[0]["field"])
	assert.Equal(t, "string", changes[0]["type"])
	assert.Equal(t, "Old Name", changes[0]["oldValue"])
	assert.Equal(t, "New Name", changes[0]["newValue"])
}

func TestPatchCatalogueItemQuery_DescriptionNull(t *testing.T) {
	fields := &models.PatchCatalogueItemFields{Description: &models.Optional[string]{Value: nil}}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	assert.Contains(t, q.Query, "SET item.description = $description")
	assert.Nil(t, q.Parameters["description"])

	changes := parseChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "description", changes[0]["field"])
	assert.Equal(t, "old desc", changes[0]["oldValue"])
	assert.Nil(t, changes[0]["newValue"])
}

func TestPatchCatalogueItemQuery_DetailsMergeNoDelete(t *testing.T) {
	details := []models.CatalogueItemDetail{
		{Property: models.CatalogueCategoryProperty{UID: "prop-A", Name: "Voltage", Type: models.CatalogueCategoryPropertyType{Code: "text"}}, Value: "24"},
	}
	fields := &models.PatchCatalogueItemFields{Details: &details}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	assert.Contains(t, q.Query, "MERGE(item)-[r_pp0:HAS_CATALOGUE_PROPERTY]->(pp0)")
	assert.Contains(t, q.Query, "SET r_pp0.value = $propValue0")
	assert.NotContains(t, q.Query, "DELETE r_propToDelete")
	assert.Equal(t, "prop-A", q.Parameters["propUID0"])
	assert.Equal(t, "24", q.Parameters["propValue0"])

	changes := parseChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "Voltage", changes[0]["field"])
	assert.Equal(t, "string", changes[0]["type"])
	assert.Equal(t, "12", changes[0]["oldValue"])
	assert.Equal(t, "24", changes[0]["newValue"])
}

func TestPatchCatalogueItemQuery_DetailsEmptyArray_NoOp(t *testing.T) {
	empty := []models.CatalogueItemDetail{}
	fields := &models.PatchCatalogueItemFields{Details: &empty}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	assert.NotContains(t, q.Query, "HAS_CATALOGUE_PROPERTY")
	changes := parseChanges(t, q)
	assert.Len(t, changes, 0)
}

func TestPatchCatalogueItemQuery_SupplierChange(t *testing.T) {
	fields := &models.PatchCatalogueItemFields{
		Supplier: &models.Optional[codebookModels.Codebook]{Value: &codebookModels.Codebook{UID: "supp-new", Name: "New Supp"}},
	}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	assert.Contains(t, q.Query, "OPTIONAL MATCH (item)-[oldSup:HAS_SUPPLIER]->()")
	assert.Contains(t, q.Query, "DELETE oldSup")
	assert.Contains(t, q.Query, "MATCH(newSup:Supplier{uid: $supplierUid})")
	assert.Contains(t, q.Query, "MERGE(item)-[:HAS_SUPPLIER]->(newSup)")
	assert.Equal(t, "supp-new", q.Parameters["supplierUid"])

	changes := parseChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "supplier", changes[0]["field"])
	assert.Equal(t, "codebook", changes[0]["type"])
}

func TestPatchCatalogueItemQuery_SupplierClear(t *testing.T) {
	fields := &models.PatchCatalogueItemFields{
		Supplier: &models.Optional[codebookModels.Codebook]{Value: nil},
	}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	assert.Contains(t, q.Query, "DELETE oldSup")
	assert.NotContains(t, q.Query, "MERGE(item)-[:HAS_SUPPLIER]")
	assert.Nil(t, q.Parameters["supplierUid"])

	changes := parseChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "supplier", changes[0]["field"])
	assert.Nil(t, changes[0]["newValue"])
}

func TestPatchCatalogueItemQuery_CategoryChange(t *testing.T) {
	fields := &models.PatchCatalogueItemFields{
		Category: &codebookModels.Codebook{UID: "cat-new", Name: "New Cat"},
	}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	assert.Contains(t, q.Query, "OPTIONAL MATCH (item)-[oldCat:BELONGS_TO_CATEGORY]->()")
	assert.Contains(t, q.Query, "DELETE oldCat")
	assert.Contains(t, q.Query, "MERGE(item)-[:BELONGS_TO_CATEGORY]->(newCat)")
	assert.Equal(t, "cat-new", q.Parameters["categoryUid"])

	changes := parseChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "category", changes[0]["field"])
	assert.Equal(t, "codebook", changes[0]["type"])
}

func TestPatchCatalogueItemQuery_NoChanges_EmptyChangesArray(t *testing.T) {
	fields := &models.PatchCatalogueItemFields{Name: ptr("Old Name")}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	assert.Contains(t, q.Query, "SET item.name = $name")
	assert.Contains(t, q.Query, `action: "PATCH"`)
	assert.Equal(t, "[]", q.Parameters["changes"])
}

func TestPatchCatalogueItemQuery_ChangesJsonShape(t *testing.T) {
	fields := &models.PatchCatalogueItemFields{Name: ptr("X"), CatalogueNumber: ptr("CN-002")}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	changes := parseChanges(t, q)
	assert.Len(t, changes, 2)
	for _, c := range changes {
		assert.Contains(t, c, "field")
		assert.Contains(t, c, "type")
		assert.Contains(t, c, "oldValue")
		assert.Contains(t, c, "newValue")
	}
}

func TestPatchCatalogueItemQuery_DetailRangeValue_SerializesAsJson(t *testing.T) {
	rangeVal := map[string]interface{}{"min": 1.0, "max": 10.0}
	details := []models.CatalogueItemDetail{{
		Property: models.CatalogueCategoryProperty{UID: "prop-range", Name: "Frequency", Type: models.CatalogueCategoryPropertyType{Code: "range"}},
		Value:    rangeVal,
	}}
	fields := &models.PatchCatalogueItemFields{Details: &details}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	raw, ok := q.Parameters["propValue0"].(string)
	assert.True(t, ok, "range value must be serialized to a JSON string")
	var parsed map[string]float64
	assert.NoError(t, json.Unmarshal([]byte(raw), &parsed))
	assert.Equal(t, 1.0, parsed["min"])
	assert.Equal(t, 10.0, parsed["max"])

	changes := parseChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "Frequency", changes[0]["field"])
	assert.Equal(t, "string", changes[0]["type"], "range maps to string ChangeType per spec")
	assert.Equal(t, raw, changes[0]["newValue"], "new value in change entry is the JSON string form")
}

func TestPatchCatalogueItemQuery_AlwaysSetsLastUpdateAndAuditRel(t *testing.T) {
	q := PatchCatalogueItemQuery("item-1", &models.PatchCatalogueItemFields{}, newOriginalItem(), "user-1")

	assert.Contains(t, q.Query, "SET item.lastUpdateTime = datetime()")
	assert.Contains(t, q.Query, "CREATE(item)-[:WAS_UPDATED_BY")
	assert.Contains(t, q.Query, `action: "PATCH"`)
	assert.Equal(t, "[]", q.Parameters["changes"])
	assert.Equal(t, "user-1", q.Parameters["userUID"])
	assert.Equal(t, "item-1", q.Parameters["uid"])
}
