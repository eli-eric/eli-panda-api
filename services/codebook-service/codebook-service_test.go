package codebookService

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"panda/apigateway/helpers"
	"panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/testsetup"
	"panda/apigateway/shared"
)

func TestUpdateCodebook_DuplicateCodeCaseInsensitive_ReturnsConflict(t *testing.T) {
	svc := &CodebookService{neo4jDriver: &testsetup.TestDriver}
	userUID := "test-user-" + uuid.New().String()
	targetUID := "test-unit-target-" + uuid.New().String()
	otherUID := "test-unit-other-" + uuid.New().String()

	createTestUser(t, userUID)
	createTestUnit(t, targetUID, "Target Unit", "UNIT-OLD")
	createTestUnit(t, otherUID, "Other Unit", "UNIT-ABC")

	_, err := svc.UpdateCodebook("UNIT", "", userUID, []string{shared.ROLE_CODEBOOKS_ADMIN}, &models.Codebook{
		UID:  targetUID,
		Name: "Target Unit Updated",
		Code: " unit-abc ",
	})

	assert.ErrorIs(t, err, helpers.ERR_CONFLICT)

	cleanupTestCodebookData(t, targetUID, otherUID, userUID)
}

func TestUpdateCodebook_SameUIDSameCode_AllowsUpdate(t *testing.T) {
	svc := &CodebookService{neo4jDriver: &testsetup.TestDriver}
	userUID := "test-user-" + uuid.New().String()
	targetUID := "test-unit-target-" + uuid.New().String()

	createTestUser(t, userUID)
	createTestUnit(t, targetUID, "Original Name", "UNIT-ABC")

	_, err := svc.UpdateCodebook("UNIT", "", userUID, []string{shared.ROLE_CODEBOOKS_ADMIN}, &models.Codebook{
		UID:  targetUID,
		Name: "Updated Name",
		Code: " UNIT-ABC ",
	})

	assert.NoError(t, err)

	name, code := getUnitNameAndCode(t, targetUID)
	assert.Equal(t, "Updated Name", name)
	assert.Equal(t, "UNIT-ABC", code)

	cleanupTestCodebookData(t, targetUID, userUID)
}

func TestUpdateCodebook_NameOnly_DoesNotChangeCode(t *testing.T) {
	svc := &CodebookService{neo4jDriver: &testsetup.TestDriver}
	userUID := "test-user-" + uuid.New().String()
	targetUID := "test-unit-target-" + uuid.New().String()

	createTestUser(t, userUID)
	createTestUnit(t, targetUID, "Original Name", "UNIT-CODE")

	_, err := svc.UpdateCodebook("UNIT", "", userUID, []string{shared.ROLE_CODEBOOKS_ADMIN}, &models.Codebook{
		UID:  targetUID,
		Name: "Name Only Update",
		Code: "",
	})

	assert.NoError(t, err)

	name, code := getUnitNameAndCode(t, targetUID)
	assert.Equal(t, "Name Only Update", name)
	assert.Equal(t, "UNIT-CODE", code)

	cleanupTestCodebookData(t, targetUID, userUID)
}

func createTestUser(t *testing.T, uid string) {
	t.Helper()
	_, err := testsetup.TestSession.Run(`CREATE (:User {uid: $uid})`, map[string]interface{}{"uid": uid})
	assert.NoError(t, err)
}

func createTestUnit(t *testing.T, uid string, name string, code string) {
	t.Helper()
	_, err := testsetup.TestSession.Run(`CREATE (:Unit {uid: $uid, name: $name, code: $code})`, map[string]interface{}{
		"uid":  uid,
		"name": name,
		"code": code,
	})
	assert.NoError(t, err)
}

func getUnitNameAndCode(t *testing.T, uid string) (name string, code string) {
	t.Helper()
	result, err := testsetup.TestSession.Run(`MATCH (u:Unit {uid: $uid}) RETURN u.name as name, u.code as code`, map[string]interface{}{"uid": uid})
	assert.NoError(t, err)
	assert.True(t, result.Next())

	rawName, _ := result.Record().Get("name")
	rawCode, _ := result.Record().Get("code")

	name, _ = rawName.(string)
	code, _ = rawCode.(string)

	return name, code
}

func cleanupTestCodebookData(t *testing.T, uids ...string) {
	t.Helper()
	for _, uid := range uids {
		_, err := testsetup.TestSession.Run(`MATCH (n {uid: $uid}) DETACH DELETE n`, map[string]interface{}{"uid": uid})
		assert.NoError(t, err)
	}
}
