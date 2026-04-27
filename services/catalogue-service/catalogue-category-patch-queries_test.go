package catalogueService

import (
	"encoding/json"
	"testing"

	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	codebookModels "panda/apigateway/services/codebook-service/models"

	"github.com/stretchr/testify/assert"
)

func newOriginalCategory() *models.CatalogueCategory {
	return &models.CatalogueCategory{
		UID:        "cat-1",
		Name:       "Old Cat",
		Code:       "OC",
		SystemType: &codebookModels.Codebook{UID: "st-old", Name: "Old Type"},
	}
}

func parseCategoryChanges(t *testing.T, q helpers.DatabaseQuery) []map[string]interface{} {
	t.Helper()
	raw, ok := q.Parameters["changes"].(string)
	assert.True(t, ok, "changes parameter must be a string")
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(raw), &parsed))
	return parsed
}

func TestPatchCatalogueCategoryQuery_NameOnly(t *testing.T) {
	name := "New Cat"
	q := PatchCatalogueCategoryQuery("cat-1", &models.PatchCatalogueCategoryFields{Name: &name}, newOriginalCategory(), "user-1")

	assert.Equal(t, "uid", q.ReturnAlias)
	assert.Contains(t, q.Query, "MATCH(category:CatalogueCategory{uid: $uid})")
	assert.Contains(t, q.Query, "SET category.name = $name")
	assert.NotContains(t, q.Query, "SET category.code")
	assert.NotContains(t, q.Query, "HAS_SYSTEM_TYPE")
	assert.Contains(t, q.Query, `action: $action`)
	assert.Equal(t, "UPDATE", q.Parameters["action"])

	changes := parseCategoryChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "name", changes[0]["field"])
	assert.Equal(t, "Old Cat", changes[0]["oldValue"])
	assert.Equal(t, "New Cat", changes[0]["newValue"])
}

func TestPatchCatalogueCategoryQuery_SystemTypeChange_MatchesBeforeDelete(t *testing.T) {
	q := PatchCatalogueCategoryQuery("cat-1", &models.PatchCatalogueCategoryFields{
		SystemType: &models.Optional[codebookModels.Codebook]{Value: &codebookModels.Codebook{UID: "st-new", Name: "New Type"}},
	}, newOriginalCategory(), "user-1")

	idxMatchNew := indexOfStr(q.Query, "MATCH(newSystemType:SystemType{uid: $systemTypeUid})")
	idxDeleteOld := indexOfStr(q.Query, "DELETE oldST")
	assert.Greater(t, idxMatchNew, -1)
	assert.Greater(t, idxDeleteOld, idxMatchNew,
		"new systemType must be MATCHed before deleting the old HAS_SYSTEM_TYPE relationship")

	changes := parseCategoryChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Equal(t, "systemType", changes[0]["field"])
	assert.Equal(t, "codebook", changes[0]["type"])
}

func TestPatchCatalogueCategoryQuery_SystemTypeClear(t *testing.T) {
	q := PatchCatalogueCategoryQuery("cat-1", &models.PatchCatalogueCategoryFields{
		SystemType: &models.Optional[codebookModels.Codebook]{Value: nil},
	}, newOriginalCategory(), "user-1")

	assert.Contains(t, q.Query, "DELETE oldST")
	assert.NotContains(t, q.Query, "MERGE(category)-[:HAS_SYSTEM_TYPE]")

	changes := parseCategoryChanges(t, q)
	assert.Len(t, changes, 1)
	assert.Nil(t, changes[0]["newValue"])
}

func TestPatchCatalogueCategoryQuery_IdempotentName(t *testing.T) {
	same := "Old Cat"
	q := PatchCatalogueCategoryQuery("cat-1", &models.PatchCatalogueCategoryFields{Name: &same}, newOriginalCategory(), "user-1")

	assert.Contains(t, q.Query, "SET category.name = $name")
	assert.Equal(t, "[]", q.Parameters["changes"], "no-op PATCH still writes audit row with empty changes")
}

func indexOfStr(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}
