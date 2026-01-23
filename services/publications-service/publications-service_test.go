package publicationsservice

import (
	"panda/apigateway/services/publications-service/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"panda/apigateway/services/testsetup"
)

func TestGetPublicationByUid(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// Insert test data
	_, err := testsetup.TestSession.Run(`CREATE (p:Publication {uid: "test123"}) RETURN p`, nil)
	assert.NoError(t, err)

	// Run the actual test
	result, err := service.GetPublicationByUid("test123")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "test123", result.Uid)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (p:Publication {uid: "test123"}) DETACH DELETE p`, nil)
	assert.NoError(t, err)
}

func TestGetPublicationByUidNotFound(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// Run the actual test
	_, err := service.GetPublicationByUid("test123")

	// Assertions
	assert.Error(t, err)

	_, cleanupErr := testsetup.TestSession.Run(`MATCH (p:Publication {uid: "test123"}) DETACH DELETE p`, nil)
	assert.NoError(t, cleanupErr)
}

// ==================== Researcher Tests ====================

func TestGetResearcherByUid(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "test-researcher-" + uuid.New().String()

	// Insert test data
	_, err := testsetup.TestSession.Run(
		`CREATE (r:Researcher {uid: $uid, firstName: "John", lastName: "Doe"}) RETURN r`,
		map[string]interface{}{"uid": testUid},
	)
	assert.NoError(t, err)

	// Run the actual test
	result, err := service.GetResearcherByUid(testUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, testUid, result.Uid)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (r:Researcher {uid: $uid}) DETACH DELETE r`, map[string]interface{}{"uid": testUid})
	assert.NoError(t, err)
}

func TestGetResearcherByUidNotFound(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "non-existent-researcher-" + uuid.New().String()

	// Run the actual test
	_, err := service.GetResearcherByUid(testUid)

	// Assertions
	assert.Error(t, err)
}

func TestCreateResearcher(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "test-researcher-create-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	orcid := "0000-0001-2345-6789"
	identificationNumber := "123456789"

	researcher := &models.Researcher{
		Uid:                  testUid,
		FirstName:            "Jane",
		LastName:             "Smith",
		ORCID:                &orcid,
		IdentificationNumber: &identificationNumber,
	}

	// Run the actual test
	result, err := service.CreateResearcher(researcher, userUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, testUid, result.Uid)
	assert.Equal(t, "Jane", result.FirstName)
	assert.Equal(t, "Smith", result.LastName)
	assert.Equal(t, &orcid, result.ORCID)

	// Verify in database
	dbResult, err := service.GetResearcherByUid(testUid)
	assert.NoError(t, err)
	assert.Equal(t, "Jane", dbResult.FirstName)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (r:Researcher {uid: $uid}) DETACH DELETE r`, map[string]interface{}{"uid": testUid})
	assert.NoError(t, err)
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": testUid})
}

func TestUpdateResearcher(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "test-researcher-update-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Insert test data
	_, err := testsetup.TestSession.Run(
		`CREATE (r:Researcher {uid: $uid, firstName: "Original", lastName: "Name"}) RETURN r`,
		map[string]interface{}{"uid": testUid},
	)
	assert.NoError(t, err)

	// Update researcher
	updatedResearcher := &models.Researcher{
		Uid:       testUid,
		FirstName: "Updated",
		LastName:  "Person",
	}

	result, err := service.UpdateResearcher(updatedResearcher, userUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.FirstName)
	assert.Equal(t, "Person", result.LastName)

	// Verify in database
	dbResult, err := service.GetResearcherByUid(testUid)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", dbResult.FirstName)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (r:Researcher {uid: $uid}) DETACH DELETE r`, map[string]interface{}{"uid": testUid})
	assert.NoError(t, err)
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": testUid})
}

func TestDeleteResearcher(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "test-researcher-delete-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Insert test data
	_, err := testsetup.TestSession.Run(
		`CREATE (r:Researcher {uid: $uid, firstName: "ToDelete", lastName: "Researcher"}) RETURN r`,
		map[string]interface{}{"uid": testUid},
	)
	assert.NoError(t, err)

	// Run the actual test (soft delete)
	err = service.DeleteResearcher(testUid, userUid)

	// Assertions
	assert.NoError(t, err)

	// Verify soft delete (deleted flag should be true)
	verifyResult, _ := testsetup.TestSession.Run(
		`MATCH (r:Researcher {uid: $uid}) RETURN r.deleted as deleted`,
		map[string]interface{}{"uid": testUid},
	)
	if verifyResult.Next() {
		deleted, _ := verifyResult.Record().Get("deleted")
		assert.Equal(t, true, deleted)
	}

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (r:Researcher {uid: $uid}) DETACH DELETE r`, map[string]interface{}{"uid": testUid})
	assert.NoError(t, err)
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": testUid})
}

func TestGetResearchers(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testPrefix := "test-researchers-list-" + uuid.New().String()

	// Insert test data
	for i := 1; i <= 3; i++ {
		_, err := testsetup.TestSession.Run(
			`CREATE (r:Researcher {uid: $uid, firstName: $firstName, lastName: "TestLastName"}) RETURN r`,
			map[string]interface{}{
				"uid":       testPrefix + "-" + string(rune('0'+i)),
				"firstName": "Researcher" + string(rune('0'+i)),
			},
		)
		assert.NoError(t, err)
	}

	// Run the actual test
	results, totalCount, err := service.GetResearchers("TestLastName", 1, 10, nil)

	// Assertions
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, totalCount, int64(3))
	assert.GreaterOrEqual(t, len(results), 3)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (r:Researcher) WHERE r.uid STARTS WITH $prefix DETACH DELETE r`, map[string]interface{}{"prefix": testPrefix})
	assert.NoError(t, err)
}

func TestCreateResearchers(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testPrefix := "test-bulk-create-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	researchers := []models.Researcher{
		{Uid: testPrefix + "-1", FirstName: "Bulk1", LastName: "Create1"},
		{Uid: testPrefix + "-2", FirstName: "Bulk2", LastName: "Create2"},
		{Uid: testPrefix + "-3", FirstName: "Bulk3", LastName: "Create3"},
	}

	// Run the actual test
	results, err := service.CreateResearchers(researchers, userUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 3, len(results))
	assert.Equal(t, "Bulk1", results[0].FirstName)
	assert.Equal(t, "Bulk2", results[1].FirstName)
	assert.Equal(t, "Bulk3", results[2].FirstName)

	// Verify in database
	for _, r := range researchers {
		dbResult, err := service.GetResearcherByUid(r.Uid)
		assert.NoError(t, err)
		assert.Equal(t, r.FirstName, dbResult.FirstName)
	}

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (r:Researcher) WHERE r.uid STARTS WITH $prefix DETACH DELETE r`, map[string]interface{}{"prefix": testPrefix})
	assert.NoError(t, err)
	_, _ = testsetup.TestSession.Run(`MATCH (h:History) WHERE h.objectUID STARTS WITH $prefix DETACH DELETE h`, map[string]interface{}{"prefix": testPrefix})
}

func TestGetResearcherWithCitizenship(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "test-researcher-citizenship-" + uuid.New().String()
	countryUid := "6626c1a0-a33b-4031-a12c-912e0ca2a650" // Czech Republic UID from migrations

	// Insert test data with citizenship relationship
	_, err := testsetup.TestSession.Run(
		`CREATE (r:Researcher {uid: $uid, firstName: "Czech", lastName: "Researcher"})
		 WITH r
		 MATCH (c:Country {uid: $countryUid})
		 CREATE (r)-[:HAS_CITIZENSHIP]->(c)
		 RETURN r`,
		map[string]interface{}{"uid": testUid, "countryUid": countryUid},
	)
	assert.NoError(t, err)

	// Run the actual test
	result, err := service.GetResearcherByUid(testUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, testUid, result.Uid)
	assert.Equal(t, "Czech", result.FirstName)
	if result.Citizenship != nil {
		assert.Equal(t, countryUid, result.Citizenship.UID)
	}

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (r:Researcher {uid: $uid}) DETACH DELETE r`, map[string]interface{}{"uid": testUid})
	assert.NoError(t, err)
}

// ==================== Publication-Researcher Connection Tests ====================

func TestCreatePublicationWithResearchers(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	pubUid := "test-pub-with-researchers-" + uuid.New().String()
	res1Uid := "test-res1-" + uuid.New().String()
	res2Uid := "test-res2-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Create researchers first
	_, err := testsetup.TestSession.Run(
		`CREATE (r1:Researcher {uid: $uid1, firstName: "John", lastName: "Doe"})
		 CREATE (r2:Researcher {uid: $uid2, firstName: "Jane", lastName: "Smith"})`,
		map[string]interface{}{"uid1": res1Uid, "uid2": res2Uid},
	)
	assert.NoError(t, err)

	// Create publication with researchers
	publication := &models.Publication{
		Uid:   pubUid,
		Title: "Test Publication",
		EliResearchers: []models.ResearcherRef{
			{Uid: res1Uid, FirstName: "John", LastName: "Doe"},
			{Uid: res2Uid, FirstName: "Jane", LastName: "Smith"},
		},
	}

	result, err := service.CreatePublication(publication, userUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, pubUid, result.Uid)
	assert.Equal(t, 2, len(result.EliResearchers))

	// Verify relationships in database
	dbResult, err := service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(dbResult.EliResearchers))

	// Clean up
	_, _ = testsetup.TestSession.Run(`MATCH (p:Publication {uid: $uid}) DETACH DELETE p`, map[string]interface{}{"uid": pubUid})
	_, _ = testsetup.TestSession.Run(`MATCH (r:Researcher) WHERE r.uid IN [$uid1, $uid2] DETACH DELETE r`, map[string]interface{}{"uid1": res1Uid, "uid2": res2Uid})
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": pubUid})
}

func TestGetPublicationWithResearchers(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	pubUid := "test-pub-get-researchers-" + uuid.New().String()
	resUid := "test-res-get-" + uuid.New().String()

	// Insert test data with HAS_RESEARCHER relationship
	_, err := testsetup.TestSession.Run(
		`CREATE (p:Publication {uid: $pubUid, title: "Test Publication"})
		 CREATE (r:Researcher {uid: $resUid, firstName: "Test", lastName: "Researcher"})
		 CREATE (p)-[:HAS_RESEARCHER]->(r)`,
		map[string]interface{}{"pubUid": pubUid, "resUid": resUid},
	)
	assert.NoError(t, err)

	// Run the actual test
	result, err := service.GetPublicationByUid(pubUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, pubUid, result.Uid)
	assert.Equal(t, 1, len(result.EliResearchers))
	assert.Equal(t, resUid, result.EliResearchers[0].Uid)
	assert.Equal(t, "Test", result.EliResearchers[0].FirstName)
	assert.Equal(t, "Researcher", result.EliResearchers[0].LastName)

	// Clean up
	_, _ = testsetup.TestSession.Run(`MATCH (p:Publication {uid: $uid}) DETACH DELETE p`, map[string]interface{}{"uid": pubUid})
	_, _ = testsetup.TestSession.Run(`MATCH (r:Researcher {uid: $uid}) DETACH DELETE r`, map[string]interface{}{"uid": resUid})
}

func TestUpdatePublicationAddResearchers(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	pubUid := "test-pub-add-researchers-" + uuid.New().String()
	resUid := "test-res-add-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Create publication without researchers
	_, err := testsetup.TestSession.Run(
		`CREATE (p:Publication {uid: $pubUid, title: "Test Publication"})
		 CREATE (r:Researcher {uid: $resUid, firstName: "New", lastName: "Researcher"})`,
		map[string]interface{}{"pubUid": pubUid, "resUid": resUid},
	)
	assert.NoError(t, err)

	// Update publication to add researchers
	publication := &models.Publication{
		Uid:   pubUid,
		Title: "Test Publication",
		EliResearchers: []models.ResearcherRef{
			{Uid: resUid, FirstName: "New", LastName: "Researcher"},
		},
	}

	_, err = service.UpdatePublication(publication, userUid)
	assert.NoError(t, err)

	// Verify relationships in database
	dbResult, err := service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dbResult.EliResearchers))
	assert.Equal(t, resUid, dbResult.EliResearchers[0].Uid)

	// Clean up
	_, _ = testsetup.TestSession.Run(`MATCH (p:Publication {uid: $uid}) DETACH DELETE p`, map[string]interface{}{"uid": pubUid})
	_, _ = testsetup.TestSession.Run(`MATCH (r:Researcher {uid: $uid}) DETACH DELETE r`, map[string]interface{}{"uid": resUid})
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": pubUid})
}

func TestUpdatePublicationRemoveResearchers(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	pubUid := "test-pub-remove-researchers-" + uuid.New().String()
	resUid := "test-res-remove-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Create publication with researcher
	_, err := testsetup.TestSession.Run(
		`CREATE (p:Publication {uid: $pubUid, title: "Test Publication"})
		 CREATE (r:Researcher {uid: $resUid, firstName: "ToRemove", lastName: "Researcher"})
		 CREATE (p)-[:HAS_RESEARCHER]->(r)`,
		map[string]interface{}{"pubUid": pubUid, "resUid": resUid},
	)
	assert.NoError(t, err)

	// Verify researcher is connected
	dbResult, err := service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dbResult.EliResearchers))

	// Update publication to remove all researchers
	publication := &models.Publication{
		Uid:            pubUid,
		Title:          "Test Publication",
		EliResearchers: []models.ResearcherRef{}, // empty list
	}

	_, err = service.UpdatePublication(publication, userUid)
	assert.NoError(t, err)

	// Verify no researchers connected
	dbResult, err = service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(dbResult.EliResearchers))

	// Verify researcher still exists (only relationship deleted)
	resResult, err := service.GetResearcherByUid(resUid)
	assert.NoError(t, err)
	assert.Equal(t, resUid, resResult.Uid)

	// Clean up
	_, _ = testsetup.TestSession.Run(`MATCH (p:Publication {uid: $uid}) DETACH DELETE p`, map[string]interface{}{"uid": pubUid})
	_, _ = testsetup.TestSession.Run(`MATCH (r:Researcher {uid: $uid}) DETACH DELETE r`, map[string]interface{}{"uid": resUid})
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": pubUid})
}

func TestUpdatePublicationReplaceResearchers(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	pubUid := "test-pub-replace-researchers-" + uuid.New().String()
	res1Uid := "test-res-old-" + uuid.New().String()
	res2Uid := "test-res-new-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Create publication with one researcher
	_, err := testsetup.TestSession.Run(
		`CREATE (p:Publication {uid: $pubUid, title: "Test Publication"})
		 CREATE (r1:Researcher {uid: $res1Uid, firstName: "Old", lastName: "Researcher"})
		 CREATE (r2:Researcher {uid: $res2Uid, firstName: "New", lastName: "Researcher"})
		 CREATE (p)-[:HAS_RESEARCHER]->(r1)`,
		map[string]interface{}{"pubUid": pubUid, "res1Uid": res1Uid, "res2Uid": res2Uid},
	)
	assert.NoError(t, err)

	// Verify old researcher is connected
	dbResult, err := service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dbResult.EliResearchers))
	assert.Equal(t, res1Uid, dbResult.EliResearchers[0].Uid)

	// Update publication to replace researcher
	publication := &models.Publication{
		Uid:   pubUid,
		Title: "Test Publication",
		EliResearchers: []models.ResearcherRef{
			{Uid: res2Uid, FirstName: "New", LastName: "Researcher"},
		},
	}

	_, err = service.UpdatePublication(publication, userUid)
	assert.NoError(t, err)

	// Verify new researcher connected, old one disconnected
	dbResult, err = service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dbResult.EliResearchers))
	assert.Equal(t, res2Uid, dbResult.EliResearchers[0].Uid)

	// Clean up
	_, _ = testsetup.TestSession.Run(`MATCH (p:Publication {uid: $uid}) DETACH DELETE p`, map[string]interface{}{"uid": pubUid})
	_, _ = testsetup.TestSession.Run(`MATCH (r:Researcher) WHERE r.uid IN [$uid1, $uid2] DETACH DELETE r`, map[string]interface{}{"uid1": res1Uid, "uid2": res2Uid})
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": pubUid})
}

// ==================== Grant Tests ====================

func TestGetGrantByUid(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "test-grant-" + uuid.New().String()

	// Insert test data
	_, err := testsetup.TestSession.Run(
		`CREATE (g:Grant {uid: $uid, code: "TEST-001", name: "Test Grant"}) RETURN g`,
		map[string]interface{}{"uid": testUid},
	)
	assert.NoError(t, err)

	// Run the actual test
	result, err := service.GetGrantByUid(testUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, testUid, result.Uid)
	assert.Equal(t, "TEST-001", result.Code)
	assert.Equal(t, "Test Grant", result.Name)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (g:Grant {uid: $uid}) DETACH DELETE g`, map[string]interface{}{"uid": testUid})
	assert.NoError(t, err)
}

func TestGetGrantByUidNotFound(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "non-existent-grant-" + uuid.New().String()

	// Run the actual test
	_, err := service.GetGrantByUid(testUid)

	// Assertions
	assert.Error(t, err)
}

func TestCreateGrant(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "test-grant-create-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()
	facilityCode := "B" // ELI facility

	grant := &models.Grant{
		Uid:  testUid,
		Code: "TEST-CREATE-001",
		Name: "Test Create Grant",
	}

	// Run the actual test
	result, err := service.CreateGrant(grant, userUid, facilityCode)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, testUid, result.Uid)
	assert.Equal(t, "TEST-CREATE-001", result.Code)
	assert.Equal(t, "Test Create Grant", result.Name)

	// Verify in database
	dbResult, err := service.GetGrantByUid(testUid)
	assert.NoError(t, err)
	assert.Equal(t, "TEST-CREATE-001", dbResult.Code)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (g:Grant {uid: $uid}) DETACH DELETE g`, map[string]interface{}{"uid": testUid})
	assert.NoError(t, err)
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": testUid})
}

func TestUpdateGrant(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "test-grant-update-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Insert test data
	_, err := testsetup.TestSession.Run(
		`CREATE (g:Grant {uid: $uid, code: "ORIG-001", name: "Original Grant"}) RETURN g`,
		map[string]interface{}{"uid": testUid},
	)
	assert.NoError(t, err)

	// Update grant
	updatedGrant := &models.Grant{
		Uid:  testUid,
		Code: "UPDATED-001",
		Name: "Updated Grant",
	}

	result, err := service.UpdateGrant(updatedGrant, userUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "UPDATED-001", result.Code)
	assert.Equal(t, "Updated Grant", result.Name)

	// Verify in database
	dbResult, err := service.GetGrantByUid(testUid)
	assert.NoError(t, err)
	assert.Equal(t, "UPDATED-001", dbResult.Code)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (g:Grant {uid: $uid}) DETACH DELETE g`, map[string]interface{}{"uid": testUid})
	assert.NoError(t, err)
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": testUid})
}

func TestDeleteGrant(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testUid := "test-grant-delete-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Insert test data
	_, err := testsetup.TestSession.Run(
		`CREATE (g:Grant {uid: $uid, code: "DELETE-001", name: "To Delete Grant"}) RETURN g`,
		map[string]interface{}{"uid": testUid},
	)
	assert.NoError(t, err)

	// Run the actual test (soft delete)
	err = service.DeleteGrant(testUid, userUid)

	// Assertions
	assert.NoError(t, err)

	// Verify soft delete (deleted flag should be true)
	verifyResult, _ := testsetup.TestSession.Run(
		`MATCH (g:Grant {uid: $uid}) RETURN g.deleted as deleted`,
		map[string]interface{}{"uid": testUid},
	)
	if verifyResult.Next() {
		deleted, _ := verifyResult.Record().Get("deleted")
		assert.Equal(t, true, deleted)
	}

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (g:Grant {uid: $uid}) DETACH DELETE g`, map[string]interface{}{"uid": testUid})
	assert.NoError(t, err)
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": testUid})
}

func TestGetGrants(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	testPrefix := "test-grants-list-" + uuid.New().String()
	facilityCode := "B"

	// Insert test data with facility relationship
	for i := 1; i <= 3; i++ {
		_, err := testsetup.TestSession.Run(
			`MATCH (f:Facility {code: $facilityCode})
			 CREATE (g:Grant {uid: $uid, code: $code, name: $name})
			 CREATE (g)-[:BELONGS_TO_FACILITY]->(f)
			 RETURN g`,
			map[string]interface{}{
				"uid":          testPrefix + "-" + string(rune('0'+i)),
				"code":         "LIST-00" + string(rune('0'+i)),
				"name":         "Test Grant " + string(rune('0'+i)),
				"facilityCode": facilityCode,
			},
		)
		assert.NoError(t, err)
	}

	// Run the actual test
	results, totalCount, err := service.GetGrants("Test Grant", 1, 10, facilityCode)

	// Assertions
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, totalCount, int64(3))
	assert.GreaterOrEqual(t, len(results), 3)

	// Clean up
	_, err = testsetup.TestSession.Run(`MATCH (g:Grant) WHERE g.uid STARTS WITH $prefix DETACH DELETE g`, map[string]interface{}{"prefix": testPrefix})
	assert.NoError(t, err)
}

// ==================== Publication-Grant Connection Tests ====================

func TestCreatePublicationWithGrants(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	pubUid := "test-pub-with-grants-" + uuid.New().String()
	grant1Uid := "test-grant1-" + uuid.New().String()
	grant2Uid := "test-grant2-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Create grants first
	_, err := testsetup.TestSession.Run(
		`CREATE (g1:Grant {uid: $uid1, code: "G1", name: "Grant One"})
		 CREATE (g2:Grant {uid: $uid2, code: "G2", name: "Grant Two"})`,
		map[string]interface{}{"uid1": grant1Uid, "uid2": grant2Uid},
	)
	assert.NoError(t, err)

	// Create publication with grants
	publication := &models.Publication{
		Uid:   pubUid,
		Title: "Test Publication with Grants",
		Grants: []models.GrantRef{
			{Uid: grant1Uid, Code: "G1", Name: "Grant One"},
			{Uid: grant2Uid, Code: "G2", Name: "Grant Two"},
		},
	}

	result, err := service.CreatePublication(publication, userUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, pubUid, result.Uid)
	assert.Equal(t, 2, len(result.Grants))

	// Verify relationships in database
	dbResult, err := service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(dbResult.Grants))

	// Clean up
	_, _ = testsetup.TestSession.Run(`MATCH (p:Publication {uid: $uid}) DETACH DELETE p`, map[string]interface{}{"uid": pubUid})
	_, _ = testsetup.TestSession.Run(`MATCH (g:Grant) WHERE g.uid IN [$uid1, $uid2] DETACH DELETE g`, map[string]interface{}{"uid1": grant1Uid, "uid2": grant2Uid})
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": pubUid})
}

func TestGetPublicationWithGrants(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	pubUid := "test-pub-get-grants-" + uuid.New().String()
	grantUid := "test-grant-get-" + uuid.New().String()

	// Insert test data with HAS_GRANT relationship
	_, err := testsetup.TestSession.Run(
		`CREATE (p:Publication {uid: $pubUid, title: "Test Publication"})
		 CREATE (g:Grant {uid: $grantUid, code: "TEST-G", name: "Test Grant"})
		 CREATE (p)-[:HAS_GRANT]->(g)`,
		map[string]interface{}{"pubUid": pubUid, "grantUid": grantUid},
	)
	assert.NoError(t, err)

	// Run the actual test
	result, err := service.GetPublicationByUid(pubUid)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, pubUid, result.Uid)
	assert.Equal(t, 1, len(result.Grants))
	assert.Equal(t, grantUid, result.Grants[0].Uid)
	assert.Equal(t, "TEST-G", result.Grants[0].Code)
	assert.Equal(t, "Test Grant", result.Grants[0].Name)

	// Clean up
	_, _ = testsetup.TestSession.Run(`MATCH (p:Publication {uid: $uid}) DETACH DELETE p`, map[string]interface{}{"uid": pubUid})
	_, _ = testsetup.TestSession.Run(`MATCH (g:Grant {uid: $uid}) DETACH DELETE g`, map[string]interface{}{"uid": grantUid})
}

func TestUpdatePublicationAddGrants(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	pubUid := "test-pub-add-grants-" + uuid.New().String()
	grantUid := "test-grant-add-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Create publication without grants
	_, err := testsetup.TestSession.Run(
		`CREATE (p:Publication {uid: $pubUid, title: "Test Publication"})
		 CREATE (g:Grant {uid: $grantUid, code: "NEW-G", name: "New Grant"})`,
		map[string]interface{}{"pubUid": pubUid, "grantUid": grantUid},
	)
	assert.NoError(t, err)

	// Update publication to add grants
	publication := &models.Publication{
		Uid:   pubUid,
		Title: "Test Publication",
		Grants: []models.GrantRef{
			{Uid: grantUid, Code: "NEW-G", Name: "New Grant"},
		},
	}

	_, err = service.UpdatePublication(publication, userUid)
	assert.NoError(t, err)

	// Verify relationships in database
	dbResult, err := service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dbResult.Grants))
	assert.Equal(t, grantUid, dbResult.Grants[0].Uid)

	// Clean up
	_, _ = testsetup.TestSession.Run(`MATCH (p:Publication {uid: $uid}) DETACH DELETE p`, map[string]interface{}{"uid": pubUid})
	_, _ = testsetup.TestSession.Run(`MATCH (g:Grant {uid: $uid}) DETACH DELETE g`, map[string]interface{}{"uid": grantUid})
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": pubUid})
}

func TestUpdatePublicationRemoveGrants(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	pubUid := "test-pub-remove-grants-" + uuid.New().String()
	grantUid := "test-grant-remove-" + uuid.New().String()
	userUid := "test-user-" + uuid.New().String()

	// Create publication with grant
	_, err := testsetup.TestSession.Run(
		`CREATE (p:Publication {uid: $pubUid, title: "Test Publication"})
		 CREATE (g:Grant {uid: $grantUid, code: "REMOVE-G", name: "To Remove Grant"})
		 CREATE (p)-[:HAS_GRANT]->(g)`,
		map[string]interface{}{"pubUid": pubUid, "grantUid": grantUid},
	)
	assert.NoError(t, err)

	// Verify grant is connected
	dbResult, err := service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dbResult.Grants))

	// Update publication to remove all grants
	publication := &models.Publication{
		Uid:    pubUid,
		Title:  "Test Publication",
		Grants: []models.GrantRef{}, // empty list
	}

	_, err = service.UpdatePublication(publication, userUid)
	assert.NoError(t, err)

	// Verify no grants connected
	dbResult, err = service.GetPublicationByUid(pubUid)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(dbResult.Grants))

	// Verify grant still exists (only relationship deleted)
	grantResult, err := service.GetGrantByUid(grantUid)
	assert.NoError(t, err)
	assert.Equal(t, grantUid, grantResult.Uid)

	// Clean up
	_, _ = testsetup.TestSession.Run(`MATCH (p:Publication {uid: $uid}) DETACH DELETE p`, map[string]interface{}{"uid": pubUid})
	_, _ = testsetup.TestSession.Run(`MATCH (g:Grant {uid: $uid}) DETACH DELETE g`, map[string]interface{}{"uid": grantUid})
	_, _ = testsetup.TestSession.Run(`MATCH (h:History {objectUID: $uid}) DETACH DELETE h`, map[string]interface{}{"uid": pubUid})
}
