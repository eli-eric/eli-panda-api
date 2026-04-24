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

// ===== Group CRUD integration tests =====

func seedCategoryWithGroups(t *testing.T) (f categoryPatchFixture, groupAUID, groupBUID string) {
	t.Helper()
	f = seedCategoryPatchFixture(t)
	groupAUID = "cp-grpA-" + uuid.NewString()
	groupBUID = "cp-grpB-" + uuid.NewString()
	_, err := testsetup.TestSession.Run(`
		MATCH(c:CatalogueCategory{uid: $cat})
		CREATE (gA:CatalogueCategoryPropertyGroup {uid: $gA, name: 'Group A'})
		CREATE (gB:CatalogueCategoryPropertyGroup {uid: $gB, name: 'Group B'})
		CREATE (c)-[:HAS_GROUP]->(gA)
		CREATE (c)-[:HAS_GROUP]->(gB)
	`, map[string]interface{}{"cat": f.categoryUID, "gA": groupAUID, "gB": groupBUID})
	assert.NoError(t, err)
	return f, groupAUID, groupBUID
}

func cleanupGroups(ids ...string) {
	for _, id := range ids {
		testsetup.TestSession.Run(`MATCH (n) WHERE n.uid = $uid DETACH DELETE n`, map[string]interface{}{"uid": id})
	}
}

func TestCreateCatalogueCategoryGroup_AutoAssignsOrder(t *testing.T) {
	f := seedCategoryPatchFixture(t)
	defer cleanupCategoryPatchFixture(f)
	svc := newPatchSvc()

	// First group — no siblings, should get order=10
	g1, err := svc.CreateCatalogueCategoryGroup(f.categoryUID, &models.CreateCatalogueCategoryGroupFields{Name: "First"}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(g1.UID)
	assert.NotEmpty(t, g1.UID)
	assert.Equal(t, "First", g1.Name)
	assert.NotNil(t, g1.Order)
	assert.Equal(t, 10, *g1.Order)

	// Second group — max(siblings)+10 = 20
	g2, err := svc.CreateCatalogueCategoryGroup(f.categoryUID, &models.CreateCatalogueCategoryGroupFields{Name: "Second"}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(g2.UID)
	assert.Equal(t, 20, *g2.Order)

	// Third group — explicit order=15 overrides auto
	fifteen := 15
	g3, err := svc.CreateCatalogueCategoryGroup(f.categoryUID, &models.CreateCatalogueCategoryGroupFields{Name: "Third", Order: &fifteen}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(g3.UID)
	assert.Equal(t, 15, *g3.Order)
}

func TestCreateCatalogueCategoryGroup_UnknownCategory_ReturnsNotFound(t *testing.T) {
	svc := newPatchSvc()
	_, err := svc.CreateCatalogueCategoryGroup("missing-"+uuid.NewString(),
		&models.CreateCatalogueCategoryGroupFields{Name: "X"}, "user")
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND)
}

func TestPatchCatalogueCategoryGroup_RenameOnly(t *testing.T) {
	f, gA, gB := seedCategoryWithGroups(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	newName := "Renamed A"
	updated, err := svc.PatchCatalogueCategoryGroup(f.categoryUID, gA,
		&models.PatchCatalogueCategoryGroupFields{Name: &newName}, f.userUID)
	assert.NoError(t, err)
	assert.Equal(t, "Renamed A", updated.Name)

	changes, action := readLatestCategoryChanges(t, f.categoryUID)
	assert.Equal(t, "UPDATE", action)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 1)
	assert.Equal(t, "group."+gA+".name", parsed[0]["field"])
}

func TestPatchCatalogueCategoryGroup_WrongCategory_ReturnsNotFound(t *testing.T) {
	f, gA, gB := seedCategoryWithGroups(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)

	// Seed an unrelated category — gA doesn't belong to it.
	otherCat := "cp-other-" + uuid.NewString()
	_, err := testsetup.TestSession.Run(`CREATE (c:CatalogueCategory{uid: $uid, name: 'Other', code: 'OT'})`,
		map[string]interface{}{"uid": otherCat})
	assert.NoError(t, err)
	defer cleanupGroups(otherCat)

	svc := newPatchSvc()
	newName := "Tampered"
	_, err = svc.PatchCatalogueCategoryGroup(otherCat, gA,
		&models.PatchCatalogueCategoryGroupFields{Name: &newName}, f.userUID)
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND)
}

func TestPatchCatalogueCategoryGroup_LazyOrderSeed(t *testing.T) {
	f, gA, gB := seedCategoryWithGroups(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	// Both seeded groups have NULL order. PATCH gA.order=25 — service must lazy-seed
	// gB too (otherwise ORDER BY would place gB at the NULL-end 2147483647 sentinel
	// while gA jumps to 25 and the UI sees "gA first, everything else after").
	twentyFive := 25
	_, err := svc.PatchCatalogueCategoryGroup(f.categoryUID, gA,
		&models.PatchCatalogueCategoryGroupFields{Order: &twentyFive}, f.userUID)
	assert.NoError(t, err)

	// Use a raw Cypher read (the legacy GetCatalogueCategoryWithDetailsByUid filters
	// out empty groups) to verify both seeded groups have order values.
	res, qerr := testsetup.TestSession.Run(`
		MATCH(c:CatalogueCategory{uid: $uid})-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup)
		RETURN g.uid as uid, g.order as order
	`, map[string]interface{}{"uid": f.categoryUID})
	assert.NoError(t, qerr)

	orders := map[string]int64{}
	for res.Next() {
		uidRaw, _ := res.Record().Get("uid")
		ordRaw, _ := res.Record().Get("order")
		uidStr, _ := uidRaw.(string)
		if ord, ok := ordRaw.(int64); ok {
			orders[uidStr] = ord
		}
	}
	assert.Contains(t, orders, gA, "gA must have a non-nil order")
	assert.Contains(t, orders, gB, "lazy seed must assign order to gB even though PATCH only targets gA")
	assert.Equal(t, int64(25), orders[gA])
}

func TestDeleteCatalogueCategoryGroup_EmptyGroup_Succeeds(t *testing.T) {
	f, gA, gB := seedCategoryWithGroups(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gB) // gA should be gone after the DELETE
	svc := newPatchSvc()

	err := svc.DeleteCatalogueCategoryGroup(f.categoryUID, gA, f.userUID)
	assert.NoError(t, err)

	reloaded, _ := svc.GetCatalogueCategoryWithDetailsByUid(f.categoryUID)
	for _, g := range reloaded.Groups {
		assert.NotEqual(t, gA, g.UID, "gA must be gone")
	}
}

// ===== Property CRUD integration tests =====

// seedCategoryWithGroupsAndTypes adds a text-type node so create-property tests have a
// valid type.uid to reference. Returns both group UIDs as before plus the type UID.
func seedCategoryWithGroupsAndTypes(t *testing.T) (f categoryPatchFixture, gA, gB, typeUID string) {
	f, gA, gB = seedCategoryWithGroups(t)
	res, err := testsetup.TestSession.Run(
		`MERGE (t:CatalogueCategoryPropertyType {code: 'text'}) ON CREATE SET t.uid = randomUUID(), t.name = 'Text' RETURN t.uid as uid`,
		nil,
	)
	assert.NoError(t, err)
	if res.Next() {
		u, _ := res.Record().Get("uid")
		typeUID, _ = u.(string)
	}
	return f, gA, gB, typeUID
}

func TestCreateCatalogueCategoryProperty_HappyPath_AutoOrder(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	p, err := svc.CreateCatalogueCategoryProperty(f.categoryUID, gA, &models.CreateCatalogueCategoryPropertyFields{
		Name: "Voltage",
		Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(p.UID)
	assert.NotEmpty(t, p.UID)
	assert.Equal(t, "Voltage", p.Name)
	assert.NotNil(t, p.Order)
	assert.Equal(t, 10, *p.Order, "first property in empty group gets order=10")
	assert.Equal(t, typeUID, p.Type.UID)
}

func TestCreateCatalogueCategoryProperty_UnknownType_ReturnsValidationError(t *testing.T) {
	f, gA, gB := seedCategoryWithGroups(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	_, err := svc.CreateCatalogueCategoryProperty(f.categoryUID, gA, &models.CreateCatalogueCategoryPropertyFields{
		Name: "X",
		Type: models.CatalogueCategoryPropertyType{UID: "missing-" + uuid.NewString()},
	}, f.userUID)
	assert.ErrorIs(t, err, ErrPatchValidation)
	assert.Contains(t, err.Error(), "property type not found")
}

func TestPatchCatalogueCategoryProperty_RenameAndDefault(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	original, err := svc.CreateCatalogueCategoryProperty(f.categoryUID, gA, &models.CreateCatalogueCategoryPropertyFields{
		Name: "Voltage", Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(original.UID)

	newName := "Peak Voltage"
	newDefault := "230"
	updated, err := svc.PatchCatalogueCategoryProperty(f.categoryUID, original.UID, &models.PatchCatalogueCategoryPropertyFields{
		Name:         &newName,
		DefaultValue: &models.Optional[string]{Value: &newDefault},
	}, f.userUID)
	assert.NoError(t, err)
	assert.Equal(t, "Peak Voltage", updated.Name)
	assert.Equal(t, "230", updated.DefaultValue)

	changes, action := readLatestCategoryChanges(t, f.categoryUID)
	assert.Equal(t, "UPDATE", action)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 2, "expected 2 change entries (name + defaultValue)")
}

func TestPatchCatalogueCategoryProperty_Move_BetweenGroups(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	p, err := svc.CreateCatalogueCategoryProperty(f.categoryUID, gA, &models.CreateCatalogueCategoryPropertyFields{
		Name: "Moves", Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(p.UID)

	_, err = svc.PatchCatalogueCategoryProperty(f.categoryUID, p.UID, &models.PatchCatalogueCategoryPropertyFields{
		GroupUID: &gB,
	}, f.userUID)
	assert.NoError(t, err)

	// Verify CONTAINS_PROPERTY now comes from gB, not gA.
	res, qerr := testsetup.TestSession.Run(`
		MATCH(p:CatalogueCategoryProperty{uid: $pid})
		OPTIONAL MATCH(g:CatalogueCategoryPropertyGroup)-[:CONTAINS_PROPERTY]->(p)
		RETURN g.uid as gid
	`, map[string]interface{}{"pid": p.UID})
	assert.NoError(t, qerr)
	seenGroups := []string{}
	for res.Next() {
		if gid, _ := res.Record().Get("gid"); gid != nil {
			if s, ok := gid.(string); ok {
				seenGroups = append(seenGroups, s)
			}
		}
	}
	assert.Equal(t, []string{gB}, seenGroups, "property must be under gB only after move")
}

func TestPatchCatalogueCategoryProperty_PropertyFromDifferentCategory_Returns404(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)

	// Seed an unrelated category with its own group and property.
	otherCat := "cp-other-" + uuid.NewString()
	otherGroup := "cp-otherG-" + uuid.NewString()
	otherProp := "cp-otherP-" + uuid.NewString()
	_, err := testsetup.TestSession.Run(`
		MATCH(t:CatalogueCategoryPropertyType{uid: $tid})
		CREATE (c:CatalogueCategory{uid: $cid, name: 'Other', code: 'OT'})
		CREATE (g:CatalogueCategoryPropertyGroup{uid: $gid, name: 'Other G'})
		CREATE (c)-[:HAS_GROUP]->(g)
		CREATE (p:CatalogueCategoryProperty{uid: $pid, name: 'Other P'})
		CREATE (p)-[:IS_PROPERTY_TYPE]->(t)
		CREATE (g)-[:CONTAINS_PROPERTY]->(p)
	`, map[string]interface{}{"cid": otherCat, "gid": otherGroup, "pid": otherProp, "tid": typeUID})
	assert.NoError(t, err)
	defer cleanupGroups(otherCat, otherGroup, otherProp)

	svc := newPatchSvc()
	newName := "Tampered"
	_, err = svc.PatchCatalogueCategoryProperty(f.categoryUID, otherProp, &models.PatchCatalogueCategoryPropertyFields{
		Name: &newName,
	}, f.userUID)
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND, "flat URL with wrong category must 404")
}

func TestDeleteCatalogueCategoryProperty_Empty_Succeeds(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	p, err := svc.CreateCatalogueCategoryProperty(f.categoryUID, gA, &models.CreateCatalogueCategoryPropertyFields{
		Name: "Deletable", Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)

	err = svc.DeleteCatalogueCategoryProperty(f.categoryUID, p.UID, f.userUID)
	assert.NoError(t, err)

	_, err = svc.fetchCategoryProperty(f.categoryUID, p.UID)
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND)
}

func TestDeleteCatalogueCategoryProperty_BlocksWhenItemReferences(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	p, err := svc.CreateCatalogueCategoryProperty(f.categoryUID, gA, &models.CreateCatalogueCategoryPropertyFields{
		Name: "Referenced", Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(p.UID)

	itemUID := "cp-item-" + uuid.NewString()
	_, err = testsetup.TestSession.Run(`
		MATCH(c:CatalogueCategory{uid: $cid})
		MATCH(p:CatalogueCategoryProperty{uid: $pid})
		CREATE (i:CatalogueItem{uid: $iid, name: 'I', catalogueNumber: 'CN', lastUpdateTime: datetime()})
		CREATE (i)-[:BELONGS_TO_CATEGORY]->(c)
		CREATE (i)-[:HAS_CATALOGUE_PROPERTY {value: '42'}]->(p)
	`, map[string]interface{}{"cid": f.categoryUID, "pid": p.UID, "iid": itemUID})
	assert.NoError(t, err)
	defer cleanupGroups(itemUID)

	err = svc.DeleteCatalogueCategoryProperty(f.categoryUID, p.UID, f.userUID)
	assert.ErrorIs(t, err, helpers.ERR_DELETE_RELATED_ITEMS)

	_, err = svc.fetchCategoryProperty(f.categoryUID, p.UID)
	assert.NoError(t, err, "property must survive a rejected DELETE")
}

// ===== Physical property CRUD integration tests =====

func TestCreateCatalogueCategoryPhysicalProperty_HappyPath(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	p, err := svc.CreateCatalogueCategoryPhysicalProperty(f.categoryUID, &models.CreateCatalogueCategoryPhysicalPropertyFields{
		Name: "Weight",
		Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(p.UID)
	assert.Equal(t, "Weight", p.Name)
	assert.NotNil(t, p.Order)
	assert.Equal(t, 10, *p.Order)
}

func TestPatchCatalogueCategoryPhysicalProperty_RenameAndDefault(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	original, err := svc.CreateCatalogueCategoryPhysicalProperty(f.categoryUID, &models.CreateCatalogueCategoryPhysicalPropertyFields{
		Name: "Length", Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(original.UID)

	newName := "Total Length"
	newDefault := "2m"
	updated, err := svc.PatchCatalogueCategoryPhysicalProperty(f.categoryUID, original.UID, &models.PatchCatalogueCategoryPhysicalPropertyFields{
		Name:         &newName,
		DefaultValue: &models.Optional[string]{Value: &newDefault},
	}, f.userUID)
	assert.NoError(t, err)
	assert.Equal(t, "Total Length", updated.Name)
	assert.Equal(t, "2m", updated.DefaultValue)

	changes, action := readLatestCategoryChanges(t, f.categoryUID)
	assert.Equal(t, "UPDATE", action)
	var parsed []map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(changes), &parsed))
	assert.Len(t, parsed, 2)
	for _, c := range parsed {
		field, _ := c["field"].(string)
		assert.Contains(t, field, "physicalProperty.", "audit field must be namespaced under physicalProperty")
	}
}

func TestDeleteCatalogueCategoryPhysicalProperty_AlwaysSucceeds(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	p, err := svc.CreateCatalogueCategoryPhysicalProperty(f.categoryUID, &models.CreateCatalogueCategoryPhysicalPropertyFields{
		Name: "Width", Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)

	err = svc.DeleteCatalogueCategoryPhysicalProperty(f.categoryUID, p.UID, f.userUID)
	assert.NoError(t, err, "physical props aren't referenced by items, so DELETE always succeeds")

	_, err = svc.fetchCategoryPhysicalProperty(f.categoryUID, p.UID)
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND)
}

// ===== GET endpoints =====

func TestGetCatalogueCategoryGroup_HappyPath(t *testing.T) {
	f, gA, gB := seedCategoryWithGroups(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	got, err := svc.GetCatalogueCategoryGroup(f.categoryUID, gA)
	assert.NoError(t, err)
	assert.Equal(t, gA, got.UID)
	assert.Equal(t, "Group A", got.Name)
}

func TestGetCatalogueCategoryGroup_WrongCategory_Returns404(t *testing.T) {
	f, gA, gB := seedCategoryWithGroups(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)

	otherCat := "cp-other-" + uuid.NewString()
	_, err := testsetup.TestSession.Run(`CREATE (c:CatalogueCategory{uid: $uid, name: 'Other', code: 'OT'})`,
		map[string]interface{}{"uid": otherCat})
	assert.NoError(t, err)
	defer cleanupGroups(otherCat)

	svc := newPatchSvc()
	_, err = svc.GetCatalogueCategoryGroup(otherCat, gA)
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND)
}

func TestGetCatalogueCategoryProperty_HappyPath(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	created, err := svc.CreateCatalogueCategoryProperty(f.categoryUID, gA, &models.CreateCatalogueCategoryPropertyFields{
		Name: "Fetched", Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(created.UID)

	got, err := svc.GetCatalogueCategoryProperty(f.categoryUID, created.UID)
	assert.NoError(t, err)
	assert.Equal(t, created.UID, got.UID)
	assert.Equal(t, "Fetched", got.Name)
	assert.Equal(t, typeUID, got.Type.UID)
}

func TestGetCatalogueCategoryPhysicalProperty_HappyPath(t *testing.T) {
	f, gA, gB, typeUID := seedCategoryWithGroupsAndTypes(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	created, err := svc.CreateCatalogueCategoryPhysicalProperty(f.categoryUID, &models.CreateCatalogueCategoryPhysicalPropertyFields{
		Name: "Depth", Type: models.CatalogueCategoryPropertyType{UID: typeUID},
	}, f.userUID)
	assert.NoError(t, err)
	defer cleanupGroups(created.UID)

	got, err := svc.GetCatalogueCategoryPhysicalProperty(f.categoryUID, created.UID)
	assert.NoError(t, err)
	assert.Equal(t, created.UID, got.UID)
	assert.Equal(t, "Depth", got.Name)
}

func TestParsePatchCategoryPhysicalProperty_RejectsGroupUid(t *testing.T) {
	body := []byte(`{"name":"X","groupUid":"g1"}`)
	_, err := parsePatchCategoryPhysicalPropertyPayload(body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "groupUid is not valid on physical properties")
}

func TestDeleteCatalogueCategoryGroup_BlocksWhenPropertyHasItemValue(t *testing.T) {
	f, gA, gB := seedCategoryWithGroups(t)
	defer cleanupCategoryPatchFixture(f)
	defer cleanupGroups(gA, gB)
	svc := newPatchSvc()

	// Seed a property under gA and an item that references it.
	propUID := "cp-prop-" + uuid.NewString()
	itemUID := "cp-item-" + uuid.NewString()
	_, err := testsetup.TestSession.Run(`
		MATCH(g:CatalogueCategoryPropertyGroup{uid: $gid})
		MATCH(c:CatalogueCategory{uid: $cid})
		MERGE(typeStr:CatalogueCategoryPropertyType {code: 'text'}) ON CREATE SET typeStr.uid = randomUUID(), typeStr.name = 'Text'
		CREATE (p:CatalogueCategoryProperty {uid: $pid, name: 'P'})
		CREATE (p)-[:IS_PROPERTY_TYPE]->(typeStr)
		CREATE (g)-[:CONTAINS_PROPERTY]->(p)
		CREATE (i:CatalogueItem {uid: $iid, name: 'I', catalogueNumber: 'CN', lastUpdateTime: datetime()})
		CREATE (i)-[:BELONGS_TO_CATEGORY]->(c)
		CREATE (i)-[:HAS_CATALOGUE_PROPERTY {value: '5'}]->(p)
	`, map[string]interface{}{"gid": gA, "cid": f.categoryUID, "pid": propUID, "iid": itemUID})
	assert.NoError(t, err)
	defer cleanupGroups(propUID, itemUID)

	err = svc.DeleteCatalogueCategoryGroup(f.categoryUID, gA, f.userUID)
	assert.ErrorIs(t, err, helpers.ERR_DELETE_RELATED_ITEMS)

	// Group must still exist.
	reloaded, _ := svc.GetCatalogueCategoryWithDetailsByUid(f.categoryUID)
	var stillThere bool
	for _, g := range reloaded.Groups {
		if g.UID == gA {
			stillThere = true
		}
	}
	assert.True(t, stillThere, "group must survive a rejected DELETE")
}
