package catalogueService

import (
	"encoding/json"
	"testing"
	"time"

	"panda/apigateway/config"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/testsetup"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type patchFixture struct {
	itemUID     string
	categoryUID string
	category2   string
	supplierUID string
	supplier2   string
	userUID     string
	propAUID    string
	propBUID    string
	propCUID    string
	grpUID      string
}

func seedPatchFixture(t *testing.T) patchFixture {
	t.Helper()

	f := patchFixture{
		itemUID:     "pt-item-" + uuid.NewString(),
		categoryUID: "pt-cat-" + uuid.NewString(),
		category2:   "pt-cat2-" + uuid.NewString(),
		supplierUID: "pt-supp-" + uuid.NewString(),
		supplier2:   "pt-supp2-" + uuid.NewString(),
		userUID:     "pt-user-" + uuid.NewString(),
		propAUID:    "pt-propA-" + uuid.NewString(),
		propBUID:    "pt-propB-" + uuid.NewString(),
		propCUID:    "pt-propC-" + uuid.NewString(),
		grpUID:      "pt-grp-" + uuid.NewString(),
	}

	_, err := testsetup.TestSession.Run(`
		MERGE (typeStr:CatalogueCategoryPropertyType {code: 'text'}) ON CREATE SET typeStr.uid = randomUUID(), typeStr.name = 'Text'
		MERGE (typeNum:CatalogueCategoryPropertyType {code: 'number'}) ON CREATE SET typeNum.uid = randomUUID(), typeNum.name = 'Number'
		CREATE (cat:CatalogueCategory {uid: $categoryUID, name: 'PatchTestCat', code: 'PTC'})
		CREATE (cat2:CatalogueCategory {uid: $category2, name: 'PatchTestCat2', code: 'PTC2'})
		CREATE (supp:Supplier {uid: $supplierUID, name: 'PatchSupplier'})
		CREATE (supp2:Supplier {uid: $supplier2, name: 'PatchSupplier2'})
		CREATE (u:User {uid: $userUID, lastName: 'Tester', firstName: 'Patch'})
		CREATE (propA:CatalogueCategoryProperty {uid: $propAUID, name: 'Voltage'})
		CREATE (propB:CatalogueCategoryProperty {uid: $propBUID, name: 'Weight'})
		CREATE (propC:CatalogueCategoryProperty {uid: $propCUID, name: 'Note'})
		CREATE (propA)-[:IS_PROPERTY_TYPE]->(typeStr)
		CREATE (propB)-[:IS_PROPERTY_TYPE]->(typeNum)
		CREATE (propC)-[:IS_PROPERTY_TYPE]->(typeStr)
		CREATE (grp:CatalogueCategoryPropertyGroup {uid: $grpUID, name: 'General'})
		CREATE (cat)-[:HAS_GROUP]->(grp)
		CREATE (grp)-[:CONTAINS_PROPERTY]->(propA)
		CREATE (grp)-[:CONTAINS_PROPERTY]->(propB)
		CREATE (grp)-[:CONTAINS_PROPERTY]->(propC)
		CREATE (item:CatalogueItem {
			uid: $itemUID,
			name: 'Original Item',
			catalogueNumber: 'CN-ORIG',
			description: 'original description',
			lastUpdateTime: datetime()
		})
		CREATE (item)-[:BELONGS_TO_CATEGORY]->(cat)
		CREATE (item)-[:HAS_SUPPLIER]->(supp)
		CREATE (item)-[:HAS_CATALOGUE_PROPERTY {value: '12'}]->(propA)
		CREATE (item)-[:HAS_CATALOGUE_PROPERTY {value: '50'}]->(propB)
		CREATE (item)-[:HAS_CATALOGUE_PROPERTY {value: 'keep'}]->(propC)
		CREATE (item)-[:WAS_UPDATED_BY{at: datetime(), action: 'INSERT'}]->(u)
	`, map[string]interface{}{
		"itemUID":     f.itemUID,
		"categoryUID": f.categoryUID,
		"category2":   f.category2,
		"supplierUID": f.supplierUID,
		"supplier2":   f.supplier2,
		"userUID":     f.userUID,
		"propAUID":    f.propAUID,
		"propBUID":    f.propBUID,
		"propCUID":    f.propCUID,
		"grpUID":      f.grpUID,
	})
	assert.NoError(t, err)

	return f
}

func cleanupPatchFixture(f patchFixture) {
	testsetup.TestSession.Run(`
		MATCH (n) WHERE n.uid IN [
			$itemUID, $categoryUID, $category2, $supplierUID, $supplier2, $userUID,
			$propAUID, $propBUID, $propCUID, $grpUID
		] DETACH DELETE n
	`, map[string]interface{}{
		"itemUID":     f.itemUID,
		"categoryUID": f.categoryUID,
		"category2":   f.category2,
		"supplierUID": f.supplierUID,
		"supplier2":   f.supplier2,
		"userUID":     f.userUID,
		"propAUID":    f.propAUID,
		"propBUID":    f.propBUID,
		"propCUID":    f.propCUID,
		"grpUID":      f.grpUID,
	})
}

func newPatchSvc() *CatalogueService {
	return &CatalogueService{neo4jDriver: &testsetup.TestDriver, jwtSecret: config.Config{}.JwtSecret}
}

func readLatestChanges(t *testing.T, itemUID string) (string, string) {
	t.Helper()
	res, err := testsetup.TestSession.Run(`
		MATCH (i:CatalogueItem{uid: $uid})-[r:WAS_UPDATED_BY]->()
		WHERE r.action = 'PATCH'
		RETURN r.changes as changes, r.action as action
		ORDER BY r.at DESC LIMIT 1
	`, map[string]interface{}{"uid": itemUID})
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

func TestPatchCatalogueItem_UpdatesNameOnly(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, err := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)
	assert.NoError(t, err)

	newName := "Patched Name"
	fields := &models.PatchCatalogueItemFields{
		Name:           &newName,
		LastUpdateTime: original.LastUpdateTime,
	}

	updated, err := svc.PatchCatalogueItem(f.itemUID, fields, f.userUID)
	assert.NoError(t, err)
	assert.Equal(t, "Patched Name", updated.Name)
	assert.Equal(t, "CN-ORIG", updated.CatalogueNumber, "catalogueNumber must not change")
	assert.NotNil(t, updated.Description)
	assert.Equal(t, "original description", *updated.Description)
	assert.Len(t, updated.Details, 3, "all 3 details must remain intact")
	assert.NotNil(t, updated.Supplier)
	assert.Equal(t, f.supplierUID, updated.Supplier.UID, "supplier must not change")

	changes, action := readLatestChanges(t, f.itemUID)
	assert.Equal(t, "PATCH", action)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 1)
	assert.Equal(t, "name", parsed[0]["field"])
}

func TestPatchCatalogueItem_UpdatesOneDetailOnly(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)

	details := []models.CatalogueItemDetail{{
		Property: models.CatalogueCategoryProperty{UID: f.propAUID, Name: "Voltage", Type: models.CatalogueCategoryPropertyType{Code: "text"}},
		Value:    "24",
	}}
	fields := &models.PatchCatalogueItemFields{
		Details:        &details,
		LastUpdateTime: original.LastUpdateTime,
	}

	updated, err := svc.PatchCatalogueItem(f.itemUID, fields, f.userUID)
	assert.NoError(t, err)
	assert.Len(t, updated.Details, 3, "all 3 details must remain")

	byUID := map[string]string{}
	for _, d := range updated.Details {
		if v, ok := d.Value.(string); ok {
			byUID[d.Property.UID] = v
		}
	}
	assert.Equal(t, "24", byUID[f.propAUID], "Voltage updated")
	assert.Equal(t, "50", byUID[f.propBUID], "Weight untouched")
	assert.Equal(t, "keep", byUID[f.propCUID], "Note untouched")

	changes, _ := readLatestChanges(t, f.itemUID)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 1)
	assert.Equal(t, "Voltage", parsed[0]["field"])
	assert.Equal(t, "12", parsed[0]["oldValue"])
	assert.Equal(t, "24", parsed[0]["newValue"])
}

func TestPatchCatalogueItem_StaleTimestamp_Returns409(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)
	stale := original.LastUpdateTime.AddDate(-1, 0, 0)

	newName := "X"
	_, err := svc.PatchCatalogueItem(f.itemUID, &models.PatchCatalogueItemFields{
		Name:           &newName,
		LastUpdateTime: stale,
	}, f.userUID)
	assert.ErrorIs(t, err, helpers.ERR_CONFLICT)
}

func TestPatchCatalogueItem_NullDescription_Clears(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)

	fields := &models.PatchCatalogueItemFields{
		Description:    &models.Optional[string]{Value: nil},
		LastUpdateTime: original.LastUpdateTime,
	}
	updated, err := svc.PatchCatalogueItem(f.itemUID, fields, f.userUID)
	assert.NoError(t, err)
	assert.Nil(t, updated.Description, "description should be cleared to null")

	changes, _ := readLatestChanges(t, f.itemUID)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 1)
	assert.Equal(t, "description", parsed[0]["field"])
	assert.Nil(t, parsed[0]["newValue"])
}

func TestPatchCatalogueItem_IdempotentCall_EmptyChanges(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)

	sameName := original.Name
	fields := &models.PatchCatalogueItemFields{
		Name:           &sameName,
		LastUpdateTime: original.LastUpdateTime,
	}
	_, err := svc.PatchCatalogueItem(f.itemUID, fields, f.userUID)
	assert.NoError(t, err)

	changes, action := readLatestChanges(t, f.itemUID)
	assert.Equal(t, "PATCH", action)
	assert.Equal(t, "[]", changes, "idempotent PATCH should record empty changes array")
}

func TestPatchCatalogueItem_SupplierRelReplacement(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)

	fields := &models.PatchCatalogueItemFields{
		Supplier:       &models.Optional[codebookModels.Codebook]{Value: &codebookModels.Codebook{UID: f.supplier2, Name: "PatchSupplier2"}},
		LastUpdateTime: original.LastUpdateTime,
	}
	updated, err := svc.PatchCatalogueItem(f.itemUID, fields, f.userUID)
	assert.NoError(t, err)
	assert.NotNil(t, updated.Supplier)
	assert.Equal(t, f.supplier2, updated.Supplier.UID)

	changes, _ := readLatestChanges(t, f.itemUID)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 1)
	assert.Equal(t, "supplier", parsed[0]["field"])
	assert.Equal(t, "codebook", parsed[0]["type"])
}

func TestPatchCatalogueItem_CategoryReplacement(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)

	fields := &models.PatchCatalogueItemFields{
		Category:       &codebookModels.Codebook{UID: f.category2, Name: "PatchTestCat2"},
		LastUpdateTime: original.LastUpdateTime,
	}
	updated, err := svc.PatchCatalogueItem(f.itemUID, fields, f.userUID)
	assert.NoError(t, err)
	assert.Equal(t, f.category2, updated.Category.UID)

	changes, _ := readLatestChanges(t, f.itemUID)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 1)
	assert.Equal(t, "category", parsed[0]["field"])
	assert.Equal(t, "codebook", parsed[0]["type"])
}

func TestParsePatchCatalogueItemPayload_RequiresLastUpdateTime(t *testing.T) {
	raw := map[string]json.RawMessage{"name": json.RawMessage(`"X"`)}
	_, err := parsePatchCatalogueItemPayload(raw)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lastUpdateTime")
}

func TestParsePatchCatalogueItemPayload_NullDescription(t *testing.T) {
	raw := map[string]json.RawMessage{
		"lastUpdateTime": json.RawMessage(`"2026-04-21T10:00:00Z"`),
		"description":    json.RawMessage(`null`),
	}
	fields, err := parsePatchCatalogueItemPayload(raw)
	assert.NoError(t, err)
	assert.NotNil(t, fields.Description)
	assert.Nil(t, fields.Description.Value, "explicit null maps to Optional with nil Value")
}

func TestParsePatchCatalogueItemPayload_AbsentDescription(t *testing.T) {
	raw := map[string]json.RawMessage{
		"lastUpdateTime": json.RawMessage(`"2026-04-21T10:00:00Z"`),
	}
	fields, err := parsePatchCatalogueItemPayload(raw)
	assert.NoError(t, err)
	assert.Nil(t, fields.Description, "absent key maps to nil Optional pointer")
}

func TestParsePatchCatalogueItemPayload_CategoryEmptyUID_Rejected(t *testing.T) {
	raw := map[string]json.RawMessage{
		"lastUpdateTime": json.RawMessage(`"2026-04-21T10:00:00Z"`),
		"category":       json.RawMessage(`{}`),
	}
	_, err := parsePatchCatalogueItemPayload(raw)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category.uid")
}

func TestParsePatchCatalogueItemPayload_SupplierEmptyUID_Rejected(t *testing.T) {
	raw := map[string]json.RawMessage{
		"lastUpdateTime": json.RawMessage(`"2026-04-21T10:00:00Z"`),
		"supplier":       json.RawMessage(`{"name":"X"}`),
	}
	_, err := parsePatchCatalogueItemPayload(raw)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "supplier.uid")
}

func TestPatchCatalogueItem_UnknownUid_ReturnsNotFound(t *testing.T) {
	svc := newPatchSvc()

	newName := "X"
	_, err := svc.PatchCatalogueItem("does-not-exist-"+uuid.NewString(), &models.PatchCatalogueItemFields{
		Name:           &newName,
		LastUpdateTime: time.Now(),
	}, "anything")
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND)
}

func TestPatchCatalogueItem_UnknownSupplier_ReturnsValidationError(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)
	_, err := svc.PatchCatalogueItem(f.itemUID, &models.PatchCatalogueItemFields{
		Supplier:       &models.Optional[codebookModels.Codebook]{Value: &codebookModels.Codebook{UID: "nonexistent-" + uuid.NewString()}},
		LastUpdateTime: original.LastUpdateTime,
	}, f.userUID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "supplier not found")

	// verify old supplier relationship still exists
	reloaded, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)
	assert.NotNil(t, reloaded.Supplier, "original HAS_SUPPLIER must not be deleted when new supplier is invalid")
	assert.Equal(t, f.supplierUID, reloaded.Supplier.UID)
}

func TestPatchCatalogueItem_UnknownCategory_ReturnsValidationError(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)
	_, err := svc.PatchCatalogueItem(f.itemUID, &models.PatchCatalogueItemFields{
		Category:       &codebookModels.Codebook{UID: "nonexistent-" + uuid.NewString()},
		LastUpdateTime: original.LastUpdateTime,
	}, f.userUID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category not found")

	reloaded, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)
	assert.Equal(t, f.categoryUID, reloaded.Category.UID, "original category must not be deleted when new category is invalid")
}

func TestPatchCatalogueItem_UnknownPropertyUID_ReturnsValidationError(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)
	details := []models.CatalogueItemDetail{{
		Property: models.CatalogueCategoryProperty{UID: "not-in-category-" + uuid.NewString()},
		Value:    "X",
	}}
	_, err := svc.PatchCatalogueItem(f.itemUID, &models.PatchCatalogueItemFields{
		Details:        &details,
		LastUpdateTime: original.LastUpdateTime,
	}, f.userUID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "property")
}

func TestPatchCatalogueItem_TOCTOU_CaughtByCypherLock(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)

	// Simulate a concurrent writer advancing lastUpdateTime between our read and write.
	_, err := testsetup.TestSession.Run(
		`MATCH(i:CatalogueItem{uid: $uid}) SET i.lastUpdateTime = datetime() RETURN i`,
		map[string]interface{}{"uid": f.itemUID},
	)
	assert.NoError(t, err)

	// Now attempt PATCH with the original (stale) timestamp. The Go-level check would
	// normally catch this, but we call with the original ts to force the Cypher-level
	// check to be the one that rejects.
	newName := "Raced"
	_, err = svc.PatchCatalogueItem(f.itemUID, &models.PatchCatalogueItemFields{
		Name:           &newName,
		LastUpdateTime: original.LastUpdateTime,
	}, f.userUID)
	assert.ErrorIs(t, err, helpers.ERR_CONFLICT)
}

func TestParsePatchCatalogueItemPayload_NullCategory_Rejected(t *testing.T) {
	raw := map[string]json.RawMessage{
		"lastUpdateTime": json.RawMessage(`"2026-04-21T10:00:00Z"`),
		"category":       json.RawMessage(`null`),
	}
	_, err := parsePatchCatalogueItemPayload(raw)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category cannot be null")
}

func TestPatchCatalogueItemQuery_CypherLockRejectsStaleTimestamp(t *testing.T) {
	// This test bypasses the service-layer Go check by calling the Cypher directly,
	// proving the WHERE item.lastUpdateTime.epochMillis = ... guard actually works.
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)

	// Advance the stored timestamp so it no longer matches `original`.
	_, err := testsetup.TestSession.Run(
		`MATCH(i:CatalogueItem{uid: $uid}) SET i.lastUpdateTime = datetime() + duration({seconds: 1}) RETURN i`,
		map[string]interface{}{"uid": f.itemUID},
	)
	assert.NoError(t, err)

	// Build and execute the PATCH query with the original (now stale) timestamp.
	newName := "ForcedViaCypher"
	query := PatchCatalogueItemQuery(f.itemUID, &models.PatchCatalogueItemFields{
		Name:           &newName,
		LastUpdateTime: original.LastUpdateTime,
	}, &original, f.userUID)

	session, _ := helpers.NewNeo4jSession(testsetup.TestDriver)
	_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	assert.ErrorIs(t, err, helpers.ERR_NO_ROWS, "Cypher MATCH should return zero rows for stale lastUpdateTime")

	// verify the item name was NOT written
	reloaded, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)
	assert.Equal(t, "Original Item", reloaded.Name, "stale PATCH must not reach SET item.name")
}

func TestPatchCatalogueItem_CombinedCategoryAndDetails(t *testing.T) {
	f := seedPatchFixture(t)
	defer cleanupPatchFixture(f)
	svc := newPatchSvc()

	// Seed a property on the second category so we can include it in the combined PATCH.
	newPropUID := "pt-newcat-prop-" + uuid.NewString()
	newGrpUID := "pt-newcat-grp-" + uuid.NewString()
	_, err := testsetup.TestSession.Run(`
		MATCH (cat2:CatalogueCategory{uid: $cat2})
		MATCH (typeStr:CatalogueCategoryPropertyType{code: 'text'})
		CREATE (p:CatalogueCategoryProperty {uid: $propUID, name: 'Wavelength'})
		CREATE (p)-[:IS_PROPERTY_TYPE]->(typeStr)
		CREATE (grp:CatalogueCategoryPropertyGroup {uid: $grpUID, name: 'NewGrp'})
		CREATE (cat2)-[:HAS_GROUP]->(grp)
		CREATE (grp)-[:CONTAINS_PROPERTY]->(p)
	`, map[string]interface{}{"cat2": f.category2, "propUID": newPropUID, "grpUID": newGrpUID})
	assert.NoError(t, err)
	defer testsetup.TestSession.Run(`MATCH (n) WHERE n.uid IN [$p, $g] DETACH DELETE n`,
		map[string]interface{}{"p": newPropUID, "g": newGrpUID})

	original, _ := svc.GetCatalogueItemWithDetailsByUid(f.itemUID)

	// PATCH swaps category AND sets a detail belonging to the new category.
	details := []models.CatalogueItemDetail{{
		Property: models.CatalogueCategoryProperty{UID: newPropUID},
		Value:    "650nm",
	}}
	updated, err := svc.PatchCatalogueItem(f.itemUID, &models.PatchCatalogueItemFields{
		Category:       &codebookModels.Codebook{UID: f.category2, Name: "PatchTestCat2"},
		Details:        &details,
		LastUpdateTime: original.LastUpdateTime,
	}, f.userUID)
	assert.NoError(t, err, "combined category + new-category details should be accepted")
	assert.Equal(t, f.category2, updated.Category.UID)

	// verify the new detail value is persisted
	found := false
	for _, d := range updated.Details {
		if d.Property.UID == newPropUID {
			assert.Equal(t, "650nm", d.Value)
			found = true
			break
		}
	}
	assert.True(t, found, "new-category property should appear in the updated item")
}

func TestParsePatchCatalogueItemPayload_NullDescription_WhitespaceTolerant(t *testing.T) {
	raw := map[string]json.RawMessage{
		"lastUpdateTime": json.RawMessage(`"2026-04-21T10:00:00Z"`),
		"description":    json.RawMessage("  null\n"),
	}
	fields, err := parsePatchCatalogueItemPayload(raw)
	assert.NoError(t, err)
	assert.NotNil(t, fields.Description)
	assert.Nil(t, fields.Description.Value)
}
