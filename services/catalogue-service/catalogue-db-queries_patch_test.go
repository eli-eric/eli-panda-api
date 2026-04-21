package catalogueService

import (
	"encoding/json"
	"strings"
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
	assert.Contains(t, q.Query, `action: "UPDATE"`)
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
	assert.Contains(t, q.Query, `action: "UPDATE"`)
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
	assert.Contains(t, q.Query, `action: "UPDATE"`)
	assert.Equal(t, "[]", q.Parameters["changes"])
	assert.Equal(t, "user-1", q.Parameters["userUID"])
	assert.Equal(t, "item-1", q.Parameters["uid"])
}

func TestPatchCatalogueItemQuery_ItemMatchGuardsLastUpdateTime(t *testing.T) {
	q := PatchCatalogueItemQuery("item-1", &models.PatchCatalogueItemFields{}, newOriginalItem(), "user-1")

	assert.Contains(t, q.Query, "WHERE item.lastUpdateTime.epochMillis = $lastUpdateTimeMillis",
		"MATCH should include lastUpdateTime precondition to defend against TOCTOU races")
	assert.Equal(t, newOriginalItem().LastUpdateTime.UnixMilli(), q.Parameters["lastUpdateTimeMillis"])
}

func TestPatchCatalogueItemQuery_SupplierChange_MatchesNewBeforeDelete(t *testing.T) {
	fields := &models.PatchCatalogueItemFields{
		Supplier: &models.Optional[codebookModels.Codebook]{Value: &codebookModels.Codebook{UID: "supp-new", Name: "New"}},
	}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	idxMatchNew := strings.Index(q.Query, "MATCH(newSup:Supplier{uid: $supplierUid})")
	idxDeleteOld := strings.Index(q.Query, "DELETE oldSup")
	assert.Greater(t, idxMatchNew, -1, "must MATCH new supplier")
	assert.Greater(t, idxDeleteOld, idxMatchNew,
		"must MATCH new supplier BEFORE deleting the old relationship")
}

func TestPatchCatalogueItemQuery_CategoryChange_MatchesNewBeforeDelete(t *testing.T) {
	fields := &models.PatchCatalogueItemFields{
		Category: &codebookModels.Codebook{UID: "cat-new", Name: "New Cat"},
	}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	idxMatchNew := strings.Index(q.Query, "MATCH(newCat:CatalogueCategory{uid: $categoryUid})")
	idxDeleteOld := strings.Index(q.Query, "DELETE oldCat")
	assert.Greater(t, idxMatchNew, -1)
	assert.Greater(t, idxDeleteOld, idxMatchNew,
		"must MATCH new category BEFORE deleting the old relationship")
}

func TestPatchCatalogueItemQuery_AuditUsesDBSourcedMetadata(t *testing.T) {
	// Client sends a spoofed Name and Type.Code; builder should ignore those
	// in favor of originalItem's canonical metadata.
	details := []models.CatalogueItemDetail{{
		Property: models.CatalogueCategoryProperty{
			UID:  "prop-B",
			Name: "SPOOFED_NAME",
			Type: models.CatalogueCategoryPropertyType{Code: "string"},
		},
		Value: "99",
	}}
	fields := &models.PatchCatalogueItemFields{Details: &details}
	q := PatchCatalogueItemQuery("item-1", fields, newOriginalItem(), "user-1")

	changes := parseChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "Weight", changes[0]["field"], "field should come from originalItem, not payload")
	assert.Equal(t, "number", changes[0]["type"], "type should come from originalItem's Type.Code, not payload")
}

func TestPatchCatalogueItemQuery_DetailFloatFormatEdgeCases(t *testing.T) {
	// Pin the encodeDetailValueForStorage behavior for tricky float formats so a future
	// change doesn't silently alter how numbers are stored / compared.
	cases := []struct {
		name     string
		oldValue interface{}
		newValue interface{}
		expected string
		changed  bool
	}{
		{"integer float matches integer string", "1", float64(1), "1", false},
		{"integer float with trailing .0", "1.0", float64(1), "1", true}, // drift: stored as "1.0", payload normalizes to "1"
		{"decimal preserves precision", "3.14", float64(3.14), "3.14", false},
		{"large int-valued float uses exp notation", "1e+20", float64(1e20), "1e+20", false},
		{"negative", "-2.5", float64(-2.5), "-2.5", false},
		{"bool true", "true", true, "true", false},
		{"bool false", "false", false, "false", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			orig := newOriginalItem()
			// override prop-B's stored value for this case
			for i := range orig.Details {
				if orig.Details[i].Property.UID == "prop-B" {
					orig.Details[i].Value = tc.oldValue
				}
			}

			details := []models.CatalogueItemDetail{{
				Property: models.CatalogueCategoryProperty{UID: "prop-B"},
				Value:    tc.newValue,
			}}
			q := PatchCatalogueItemQuery("item-1", &models.PatchCatalogueItemFields{Details: &details}, orig, "user-1")

			assert.Equal(t, tc.expected, q.Parameters["propValue0"], "stored form must match expected normalized output")

			changes := parseChanges(t, q)
			if tc.changed {
				assert.Len(t, changes, 1, "expected value drift to register as a change")
			} else {
				assert.Len(t, changes, 0, "semantically equal old/new should not register a change")
			}
		})
	}
}

func TestPatchCatalogueItemQuery_DetailNumberTypeCoercion(t *testing.T) {
	// originalItem has prop-B (Weight, type "number") stored as string "50".
	// Payload sends JSON number 50 (float64 in Go) → should NOT register as change.
	detailsSame := []models.CatalogueItemDetail{{
		Property: models.CatalogueCategoryProperty{UID: "prop-B"},
		Value:    float64(50),
	}}
	q := PatchCatalogueItemQuery("item-1", &models.PatchCatalogueItemFields{Details: &detailsSame}, newOriginalItem(), "user-1")
	changes := parseChanges(t, q)
	assert.Len(t, changes, 0, "float 50 normalizes to '50' and matches stored string '50'")

	// Changing the value should register.
	detailsDiff := []models.CatalogueItemDetail{{
		Property: models.CatalogueCategoryProperty{UID: "prop-B"},
		Value:    float64(75),
	}}
	q2 := PatchCatalogueItemQuery("item-1", &models.PatchCatalogueItemFields{Details: &detailsDiff}, newOriginalItem(), "user-1")
	changes2 := parseChanges(t, q2)
	assert.Len(t, changes2, 1)
	assert.Equal(t, "50", changes2[0]["oldValue"])
	assert.Equal(t, "75", changes2[0]["newValue"])
}
