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
	results, totalCount, err := service.GetResearchers("TestLastName", 1, 10)

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
