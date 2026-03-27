package zoneservice

import (
	"bytes"
	"panda/apigateway/services/zone-service/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"panda/apigateway/services/testsetup"
)

const testFacilityCode = "TEST_FACILITY"

func setupTestFacility(t *testing.T) {
	_, err := testsetup.TestSession.Run(
		`MERGE (f:Facility {code: $code}) RETURN f`,
		map[string]interface{}{"code": testFacilityCode},
	)
	assert.NoError(t, err)
}

func cleanupZone(uid string) {
	testsetup.TestSession.Run(`MATCH (z:Zone {uid: $uid}) DETACH DELETE z`, map[string]interface{}{"uid": uid})
}

func cleanupUser(uid string) {
	testsetup.TestSession.Run(`MATCH (u:User {uid: $uid}) DETACH DELETE u`, map[string]interface{}{"uid": uid})
}

func createTestUser(t *testing.T, uid string) {
	_, err := testsetup.TestSession.Run(
		`MERGE (u:User {uid: $uid}) RETURN u`,
		map[string]interface{}{"uid": uid},
	)
	assert.NoError(t, err)
}

func TestCreateRootZone(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	req := &models.ZoneCreateRequest{
		Name:  "Test Zone",
		Code:  "TZ-" + uuid.New().String()[:8],
		Notes: "test notes content",
	}

	result, err := service.CreateZone(testFacilityCode, userUID, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, result.UID)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Code, result.Code)
	assert.Equal(t, "test notes content", result.Notes)

	// verify in DB
	dbResult, err := service.GetZoneByUID(result.UID, testFacilityCode)
	assert.NoError(t, err)
	assert.Equal(t, req.Name, dbResult.Name)
	assert.Equal(t, "test notes content", dbResult.Notes)
	assert.Nil(t, dbResult.ParentZone)

	cleanupZone(result.UID)
	cleanupUser(userUID)
}

func TestCreateSubZone(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	parentReq := &models.ZoneCreateRequest{
		Name: "Parent Zone",
		Code: "PZ-" + uuid.New().String()[:8],
	}
	parent, err := service.CreateZone(testFacilityCode, userUID, parentReq)
	assert.NoError(t, err)

	subReq := &models.ZoneCreateRequest{
		Name:      "Sub Zone",
		Code:      "SZ-" + uuid.New().String()[:8],
		ParentUID: &parent.UID,
	}
	sub, err := service.CreateZone(testFacilityCode, userUID, subReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, sub.UID)

	// verify parent info
	dbSub, err := service.GetZoneByUID(sub.UID, testFacilityCode)
	assert.NoError(t, err)
	assert.NotNil(t, dbSub.ParentZone)
	assert.Equal(t, parent.UID, dbSub.ParentZone.UID)

	cleanupZone(sub.UID)
	cleanupZone(parent.UID)
	cleanupUser(userUID)
}

func TestCreateSubSubZone_Rejected(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	root, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Root", Code: "R-" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)

	sub, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Sub", Code: "S-" + uuid.New().String()[:8], ParentUID: &root.UID,
	})
	assert.NoError(t, err)

	// try create sub-sub — should fail
	_, err = service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "SubSub", Code: "SS-" + uuid.New().String()[:8], ParentUID: &sub.UID,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max 2 levels")

	cleanupZone(sub.UID)
	cleanupZone(root.UID)
	cleanupUser(userUID)
}

func TestCreateZone_DuplicateCode(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	code := "DUP-" + uuid.New().String()[:8]

	zone1, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Zone 1", Code: code,
	})
	assert.NoError(t, err)

	_, err = service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Zone 2", Code: code,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	cleanupZone(zone1.UID)
	cleanupUser(userUID)
}

func TestGetAllZones(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	root, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "List Root", Code: "LR-" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)

	sub, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "List Sub", Code: "LS-" + uuid.New().String()[:8], ParentUID: &root.UID,
	})
	assert.NoError(t, err)

	zones, _, err := service.GetAllZones(testFacilityCode, "", 1, 100, nil)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(zones), 2)

	found := false
	for _, z := range zones {
		if z.UID == sub.UID {
			found = true
			assert.NotNil(t, z.ParentZone)
			assert.Equal(t, root.UID, z.ParentZone.UID)
			break
		}
	}
	assert.True(t, found, "sub-zone should be in list")

	cleanupZone(sub.UID)
	cleanupZone(root.UID)
	cleanupUser(userUID)
}

func TestGetZoneByUID(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	zone, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Get Zone", Code: "GZ-" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)

	result, err := service.GetZoneByUID(zone.UID, testFacilityCode)
	assert.NoError(t, err)
	assert.Equal(t, zone.UID, result.UID)
	assert.Equal(t, "Get Zone", result.Name)

	cleanupZone(zone.UID)
	cleanupUser(userUID)
}

func TestGetZoneByUID_NotFound(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)

	_, err := service.GetZoneByUID("non-existent-uid", testFacilityCode)
	assert.Error(t, err)
}

func TestUpdateZone(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	zone, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Original", Code: "OR-" + uuid.New().String()[:8], Notes: "original notes",
	})
	assert.NoError(t, err)

	newCode := "UP-" + uuid.New().String()[:8]
	updatedNotes := "updated notes"
	updated, err := service.UpdateZone(zone.UID, testFacilityCode, userUID, &models.ZoneUpdateRequest{
		Name: "Updated", Code: newCode, Notes: &updatedNotes,
	})
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.Name)
	assert.Equal(t, newCode, updated.Code)
	assert.Equal(t, "updated notes", updated.Notes)

	cleanupZone(zone.UID)
	cleanupUser(userUID)
}

func TestUpdateZone_NotFound(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	_, err := service.UpdateZone("non-existent-uid", testFacilityCode, userUID, &models.ZoneUpdateRequest{
		Name: "Nope", Code: "NP",
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)

	cleanupUser(userUID)
}

func TestUpdateZone_SelfParent_Rejected(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	zone, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Self", Code: "SP-" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)

	_, err = service.UpdateZone(zone.UID, testFacilityCode, userUID, &models.ZoneUpdateRequest{
		Name: "Self", Code: zone.Code, ParentUID: &zone.UID,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be its own parent")

	cleanupZone(zone.UID)
	cleanupUser(userUID)
}

func TestUpdateZone_ReassignParent(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	rootA, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Root A", Code: "RA-" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)
	rootB, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Root B", Code: "RB-" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)
	sub, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Sub", Code: "RS-" + uuid.New().String()[:8], ParentUID: &rootA.UID,
	})
	assert.NoError(t, err)

	// reassign sub from rootA to rootB
	updated, err := service.UpdateZone(sub.UID, testFacilityCode, userUID, &models.ZoneUpdateRequest{
		Name: "Sub", Code: sub.Code, ParentUID: &rootB.UID,
	})
	assert.NoError(t, err)
	assert.NotNil(t, updated.ParentZone)
	assert.Equal(t, rootB.UID, updated.ParentZone.UID)

	cleanupZone(sub.UID)
	cleanupZone(rootA.UID)
	cleanupZone(rootB.UID)
	cleanupUser(userUID)
}

func TestDeleteZone(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	zone, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "To Delete", Code: "TD-" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)

	err = service.DeleteZone(zone.UID, testFacilityCode, userUID)
	assert.NoError(t, err)

	// verify soft deleted
	verifyResult, _ := testsetup.TestSession.Run(
		`MATCH (z:Zone {uid: $uid}) RETURN z.deleted as deleted`,
		map[string]interface{}{"uid": zone.UID},
	)
	if verifyResult.Next() {
		deleted, _ := verifyResult.Record().Get("deleted")
		assert.Equal(t, true, deleted)
	}

	cleanupZone(zone.UID)
	cleanupUser(userUID)
}

func TestDeleteZone_WithSubzones_Rejected(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	root, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Root", Code: "DR-" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)
	sub, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Sub", Code: "DS-" + uuid.New().String()[:8], ParentUID: &root.UID,
	})
	assert.NoError(t, err)

	err = service.DeleteZone(root.UID, testFacilityCode, userUID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrConflictSub)

	cleanupZone(sub.UID)
	cleanupZone(root.UID)
	cleanupUser(userUID)
}

func TestDeleteZone_WithSystemRef_Rejected(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	zone, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Zone With Sys", Code: "ZS-" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)

	// create system referencing existing zone (MATCH, not CREATE duplicate)
	sysUID := "test-sys-" + uuid.New().String()
	_, err = testsetup.TestSession.Run(
		`MATCH (z:Zone {uid: $zoneUID})
		 CREATE (s:System {uid: $sysUID, deleted: false})-[:HAS_ZONE]->(z)`,
		map[string]interface{}{"sysUID": sysUID, "zoneUID": zone.UID},
	)
	assert.NoError(t, err)

	err = service.DeleteZone(zone.UID, testFacilityCode, userUID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrConflictSys)

	testsetup.TestSession.Run(`MATCH (s:System {uid: $uid}) DETACH DELETE s`, map[string]interface{}{"uid": sysUID})
	cleanupZone(zone.UID)
	cleanupUser(userUID)
}

func TestImportZones_NewZones(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	code1 := "IMP1-" + uuid.New().String()[:8]
	code2 := "IMP2-" + uuid.New().String()[:8]

	csvData := "name,code\nImport Zone 1," + code1 + "\nImport Zone 2," + code2 + "\n"
	file := bytes.NewReader([]byte(csvData))

	result, err := service.ImportZones(testFacilityCode, userUID, file)
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Created)
	assert.Equal(t, 0, result.Skipped)

	testsetup.TestSession.Run(
		`MATCH (z:Zone) WHERE z.code IN [$c1, $c2] DETACH DELETE z`,
		map[string]interface{}{"c1": code1, "c2": code2},
	)
	cleanupUser(userUID)
}

func TestImportZones_SkipExisting(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	code := "SKIP-" + uuid.New().String()[:8]
	zone, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Existing", Code: code,
	})
	assert.NoError(t, err)

	csvData := "name,code\nExisting," + code + "\n"
	file := bytes.NewReader([]byte(csvData))

	result, err := service.ImportZones(testFacilityCode, userUID, file)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.Created)
	assert.Equal(t, 1, result.Skipped)

	cleanupZone(zone.UID)
	cleanupUser(userUID)
}

func TestImportZones_WithParentCode(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	parentCode := "PAR-" + uuid.New().String()[:8]
	parent, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Parent", Code: parentCode,
	})
	assert.NoError(t, err)

	childCode := "CHD-" + uuid.New().String()[:8]
	csvData := "name,code,parentCode\nChild Zone," + childCode + "," + parentCode + "\n"
	file := bytes.NewReader([]byte(csvData))

	result, err := service.ImportZones(testFacilityCode, userUID, file)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Created)

	// verify child has parent
	zones, _, err := service.GetAllZones(testFacilityCode, "", 1, 100, nil)
	assert.NoError(t, err)
	for _, z := range zones {
		if z.Code == childCode {
			assert.NotNil(t, z.ParentZone)
			assert.Equal(t, parent.UID, z.ParentZone.UID)
		}
	}

	testsetup.TestSession.Run(
		`MATCH (z:Zone) WHERE z.code IN [$c1, $c2] DETACH DELETE z`,
		map[string]interface{}{"c1": parentCode, "c2": childCode},
	)
	cleanupUser(userUID)
}

func TestGetAllZones_DefaultSort_RootZonesFirst(t *testing.T) {
	setupTestFacility(t)
	service := NewZoneService(&testsetup.TestDriver)
	userUID := "test-user-" + uuid.New().String()
	createTestUser(t, userUID)

	// create two root zones and two subzones with predictable codes
	// using ZZZ prefix to avoid collision with real data
	rootA, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Root A", Code: "ZZZA",
	})
	assert.NoError(t, err)
	rootB, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Root B", Code: "ZZZB",
	})
	assert.NoError(t, err)
	subA1, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Sub A1", Code: "ZZZA1", ParentUID: &rootA.UID,
	})
	assert.NoError(t, err)
	subB1, err := service.CreateZone(testFacilityCode, userUID, &models.ZoneCreateRequest{
		Name: "Sub B1", Code: "ZZZB1", ParentUID: &rootB.UID,
	})
	assert.NoError(t, err)

	// fetch with default sort (nil sorting)
	zones, _, err := service.GetAllZones(testFacilityCode, "ZZZ", 1, 100, nil)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(zones))

	// root zones (ParentZone == nil) must appear before subzones
	lastRootIdx := -1
	firstSubIdx := len(zones)
	for i, z := range zones {
		if z.ParentZone == nil {
			lastRootIdx = i
		} else if i < firstSubIdx {
			firstSubIdx = i
		}
	}
	assert.Greater(t, firstSubIdx, lastRootIdx, "all root zones should appear before any subzone")

	// within root zones, should be sorted by code ASC
	assert.Nil(t, zones[0].ParentZone)
	assert.Nil(t, zones[1].ParentZone)
	assert.True(t, zones[0].Code <= zones[1].Code, "root zones should be sorted by code ASC")

	cleanupZone(subA1.UID)
	cleanupZone(subB1.UID)
	cleanupZone(rootA.UID)
	cleanupZone(rootB.UID)
	cleanupUser(userUID)
}
