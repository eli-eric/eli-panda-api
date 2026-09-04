package publicationsservice

import (
	"errors"
	"strings"
	"testing"

	"panda/apigateway/services/testsetup"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWosRepositoryFindExistingPublicationNormalizesStoredDOI(t *testing.T) {
	uid := "wos-publication-" + uuid.NewString()
	doi := "10.12345/repository-test-" + uuid.NewString()
	storedDOI := "HTTPS://DOI.ORG/" + strings.ToUpper(doi)
	repo := &neo4jWosImportRepository{driver: &testsetup.TestDriver}

	_, err := testsetup.TestSession.Run(`
		CREATE (:Publication {
			uid: $uid, doi: $doi, code: "ELI-TEST", title: "Repository test", updatedAt: datetime()
		})`, map[string]interface{}{"uid": uid, "doi": storedDOI})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = testsetup.TestSession.Run(
			`MATCH (publication:Publication {uid: $uid}) DETACH DELETE publication`,
			map[string]interface{}{"uid": uid},
		)
	})

	publication, err := repo.findExistingPublication(doi, "")
	require.NoError(t, err)
	if assert.NotNil(t, publication) {
		assert.Equal(t, uid, publication.Uid)
		assert.Equal(t, "ELI-TEST", publication.Code)
	}

	publication, err = repo.findExistingPublication(doi, uid)
	require.NoError(t, err)
	assert.Nil(t, publication)
}

func TestWosRepositoryRememberResearcherIDAppendsAndRejectsConflicts(t *testing.T) {
	targetUID := "wos-researcher-" + uuid.NewString()
	otherUID := "wos-researcher-" + uuid.NewString()
	userUID := "wos-user-" + uuid.NewString()
	repo := &neo4jWosImportRepository{driver: &testsetup.TestDriver}

	_, err := testsetup.TestSession.Run(`
		CREATE (:User {uid: $userUID})
		CREATE (:Researcher {
			uid: $targetUID, researcherId: "a-1111-2020", researcherIds: ["B-2222-2021"]
		})
		CREATE (:Researcher {
			uid: $otherUID, researcherIds: ["D-4444-2024"]
		})`, map[string]interface{}{
		"targetUID": targetUID,
		"otherUID":  otherUID,
		"userUID":   userUID,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = testsetup.TestSession.Run(
			`MATCH (node) WHERE node.uid IN $uids DETACH DELETE node`,
			map[string]interface{}{"uids": []string{targetUID, otherUID, userUID}},
		)
	})

	researcherIDs, err := repo.rememberResearcherID(targetUID, "C-3333-2023", userUID)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"A-1111-2020", "B-2222-2021", "C-3333-2023"}, researcherIDs)

	researcherIDs, err = repo.rememberResearcherID(targetUID, "C-3333-2023", userUID)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"A-1111-2020", "B-2222-2021", "C-3333-2023"}, researcherIDs)

	_, err = repo.rememberResearcherID(targetUID, "D-4444-2024", userUID)
	assert.ErrorIs(t, err, errResearcherIDOwned)

	_, err = repo.rememberResearcherID("missing-"+uuid.NewString(), "E-5555-2025", userUID)
	assert.True(t, errors.Is(err, errResearcherMissing), "expected missing researcher error, got %v", err)
}
