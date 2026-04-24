package catalogueService

import (
	"encoding/json"
	"testing"

	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/testsetup"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type categoryPatchFixture struct {
	categoryUID   string
	systemTypeA   string
	systemTypeB   string
	userUID       string
}

func seedCategoryPatchFixture(t *testing.T) categoryPatchFixture {
	t.Helper()

	f := categoryPatchFixture{
		categoryUID: "cp-cat-" + uuid.NewString(),
		systemTypeA: "cp-st-a-" + uuid.NewString(),
		systemTypeB: "cp-st-b-" + uuid.NewString(),
		userUID:     "cp-user-" + uuid.NewString(),
	}

	_, err := testsetup.TestSession.Run(`
		CREATE (cat:CatalogueCategory {uid: $categoryUID, name: 'CP Original', code: 'CPO'})
		CREATE (stA:SystemType {uid: $systemTypeA, name: 'Type A', code: 'TA'})
		CREATE (stB:SystemType {uid: $systemTypeB, name: 'Type B', code: 'TB'})
		CREATE (u:User {uid: $userUID, lastName: 'Tester', firstName: 'Cat'})
		CREATE (cat)-[:HAS_SYSTEM_TYPE]->(stA)
	`, map[string]interface{}{
		"categoryUID": f.categoryUID,
		"systemTypeA": f.systemTypeA,
		"systemTypeB": f.systemTypeB,
		"userUID":     f.userUID,
	})
	assert.NoError(t, err)
	return f
}

func cleanupCategoryPatchFixture(f categoryPatchFixture) {
	testsetup.TestSession.Run(`
		MATCH (n) WHERE n.uid IN [$categoryUID, $systemTypeA, $systemTypeB, $userUID]
		DETACH DELETE n
	`, map[string]interface{}{
		"categoryUID": f.categoryUID,
		"systemTypeA": f.systemTypeA,
		"systemTypeB": f.systemTypeB,
		"userUID":     f.userUID,
	})
}

func readLatestCategoryChanges(t *testing.T, categoryUID string) (string, string) {
	t.Helper()
	res, err := testsetup.TestSession.Run(`
		MATCH (c:CatalogueCategory{uid: $uid})-[r:WAS_UPDATED_BY]->()
		RETURN r.changes as changes, r.action as action
		ORDER BY r.at DESC LIMIT 1
	`, map[string]interface{}{"uid": categoryUID})
	assert.NoError(t, err)
	if res.Next() {
		rec := res.Record()
		changes, _ := rec.Get("changes")
		action, _ := rec.Get("action")
		s, _ := changes.(string)
		a, _ := action.(string)
		return s, a
	}
	return "", ""
}

func TestPatchCatalogueCategory_UpdatesNameOnly(t *testing.T) {
	f := seedCategoryPatchFixture(t)
	defer cleanupCategoryPatchFixture(f)
	svc := newPatchSvc()

	newName := "CP Patched"
	updated, err := svc.PatchCatalogueCategory(f.categoryUID, &models.PatchCatalogueCategoryFields{Name: &newName}, f.userUID)
	assert.NoError(t, err)
	assert.Equal(t, "CP Patched", updated.Name)
	assert.Equal(t, "CPO", updated.Code, "code must not change")
	assert.NotNil(t, updated.SystemType)
	assert.Equal(t, f.systemTypeA, updated.SystemType.UID, "systemType must stay intact")

	changes, action := readLatestCategoryChanges(t, f.categoryUID)
	assert.Equal(t, "UPDATE", action)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 1)
	assert.Equal(t, "name", parsed[0]["field"])
}

func TestPatchCatalogueCategory_UnknownUid_ReturnsNotFound(t *testing.T) {
	svc := newPatchSvc()
	newName := "X"
	_, err := svc.PatchCatalogueCategory("does-not-exist-"+uuid.NewString(),
		&models.PatchCatalogueCategoryFields{Name: &newName}, "anyone")
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND)
}

func TestPatchCatalogueCategory_UnknownSystemType_ReturnsValidationError(t *testing.T) {
	f := seedCategoryPatchFixture(t)
	defer cleanupCategoryPatchFixture(f)
	svc := newPatchSvc()

	_, err := svc.PatchCatalogueCategory(f.categoryUID, &models.PatchCatalogueCategoryFields{
		SystemType: &models.Optional[codebookModels.Codebook]{Value: &codebookModels.Codebook{UID: "missing-" + uuid.NewString()}},
	}, f.userUID)
	assert.ErrorIs(t, err, ErrPatchValidation)
	assert.Contains(t, err.Error(), "systemType not found")

	reloaded, _ := svc.GetCatalogueCategoryWithDetailsByUid(f.categoryUID)
	assert.Equal(t, f.systemTypeA, reloaded.SystemType.UID, "original systemType relationship must stay")
}

func TestPatchCatalogueCategory_SystemTypeSwap(t *testing.T) {
	f := seedCategoryPatchFixture(t)
	defer cleanupCategoryPatchFixture(f)
	svc := newPatchSvc()

	updated, err := svc.PatchCatalogueCategory(f.categoryUID, &models.PatchCatalogueCategoryFields{
		SystemType: &models.Optional[codebookModels.Codebook]{Value: &codebookModels.Codebook{UID: f.systemTypeB, Name: "Type B"}},
	}, f.userUID)
	assert.NoError(t, err)
	assert.Equal(t, f.systemTypeB, updated.SystemType.UID)

	changes, _ := readLatestCategoryChanges(t, f.categoryUID)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 1)
	assert.Equal(t, "systemType", parsed[0]["field"])
	assert.Equal(t, "codebook", parsed[0]["type"])
}

func TestPatchCatalogueCategory_SystemTypeClear(t *testing.T) {
	f := seedCategoryPatchFixture(t)
	defer cleanupCategoryPatchFixture(f)
	svc := newPatchSvc()

	updated, err := svc.PatchCatalogueCategory(f.categoryUID, &models.PatchCatalogueCategoryFields{
		SystemType: &models.Optional[codebookModels.Codebook]{Value: nil},
	}, f.userUID)
	assert.NoError(t, err)
	assert.Nil(t, updated.SystemType, "systemType relationship should be cleared")
}

func TestParsePatchCatalogueCategoryPayload_NullNameRejected(t *testing.T) {
	raw := map[string]json.RawMessage{"name": json.RawMessage(`null`)}
	_, err := parsePatchCatalogueCategoryPayload(raw)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be null")
}

func TestParsePatchCatalogueCategoryPayload_SystemTypeEmptyUIDRejected(t *testing.T) {
	raw := map[string]json.RawMessage{"systemType": json.RawMessage(`{"name":"x"}`)}
	_, err := parsePatchCatalogueCategoryPayload(raw)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "systemType.uid")
}
