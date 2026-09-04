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

	remembered, err := repo.rememberResearcherID(targetUID, "C-3333-2023", userUID, false)
	require.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{"A-1111-2020", "B-2222-2021", "C-3333-2023"}, remembered.researcherIDs)

	remembered, err = repo.rememberResearcherID(targetUID, "C-3333-2023", userUID, false)
	require.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{"A-1111-2020", "B-2222-2021", "C-3333-2023"}, remembered.researcherIDs)

	_, err = repo.rememberResearcherID(targetUID, "D-4444-2024", userUID, false)
	assert.ErrorIs(t, err, errResearcherIDOwned)

	_, err = repo.rememberResearcherID("missing-"+uuid.NewString(), "E-5555-2025", userUID, false)
	assert.True(t, errors.Is(err, errResearcherMissing), "expected missing researcher error, got %v", err)
}

func TestWosRepositoryRememberResearcherIDMaintainsCurrentID(t *testing.T) {
	userUID := "wos-user-" + uuid.NewString()
	repo := &neo4jWosImportRepository{driver: &testsetup.TestDriver}

	_, err := testsetup.TestSession.Run(
		`CREATE (:User {uid: $userUID})`, map[string]interface{}{"userUID": userUID})
	require.NoError(t, err)

	createdUIDs := []string{userUID}
	t.Cleanup(func() {
		_, _ = testsetup.TestSession.Run(
			`MATCH (node) WHERE node.uid IN $uids DETACH DELETE node`,
			map[string]interface{}{"uids": createdUIDs},
		)
	})

	// A ResearcherID may only ever belong to one researcher, so every case below
	// uses its own values. The ZZ prefix is unused by the real register.
	newResearcher := func(t *testing.T, currentID string) string {
		t.Helper()
		researcherUID := "wos-researcher-" + uuid.NewString()
		createdUIDs = append(createdUIDs, researcherUID)
		_, err := testsetup.TestSession.Run(`
			CREATE (r:Researcher {uid: $uid})
			SET r.researcherId = CASE WHEN $currentID = "" THEN NULL ELSE $currentID END`,
			map[string]interface{}{"uid": researcherUID, "currentID": currentID})
		require.NoError(t, err)
		return researcherUID
	}

	currentIDOf := func(t *testing.T, researcherUID string) string {
		t.Helper()
		result, err := testsetup.TestSession.Run(
			`MATCH (r:Researcher {uid: $uid}) RETURN coalesce(r.researcherId, "") AS id`,
			map[string]interface{}{"uid": researcherUID})
		require.NoError(t, err)
		record, err := result.Single()
		require.NoError(t, err)
		value, _ := record.Get("id")
		id, _ := value.(string)
		return id
	}

	t.Run("fills an empty current ID with the newest known", func(t *testing.T) {
		researcherUID := newResearcher(t, "")

		remembered, err := repo.rememberResearcherID(researcherUID, "ZZA-1000-2015", userUID, false)
		require.NoError(t, err)
		assert.Equal(t, "ZZA-1000-2015", remembered.primaryID)
		assert.Equal(t, "ZZA-1000-2015", currentIDOf(t, researcherUID))
	})

	t.Run("an older confirmed ID does not displace the current one", func(t *testing.T) {
		researcherUID := newResearcher(t, "ZZB-2000-2023")

		remembered, err := repo.rememberResearcherID(researcherUID, "ZZB-2001-2015", userUID, false)
		require.NoError(t, err)
		assert.Equal(t, "ZZB-2000-2023", remembered.primaryID)
		assert.ElementsMatch(t,
			[]string{"ZZB-2000-2023", "ZZB-2001-2015"}, remembered.researcherIDs)
		assert.Equal(t, "ZZB-2000-2023", currentIDOf(t, researcherUID))
	})

	t.Run("a newer confirmed ID waits until promotion is asked for", func(t *testing.T) {
		researcherUID := newResearcher(t, "ZZC-3000-2015")

		remembered, err := repo.rememberResearcherID(researcherUID, "ZZC-3001-2023", userUID, false)
		require.NoError(t, err)
		assert.Equal(t, "ZZC-3000-2015", remembered.primaryID,
			"a newer ID must not promote itself")

		remembered, err = repo.rememberResearcherID(researcherUID, "ZZC-3001-2023", userUID, true)
		require.NoError(t, err)
		assert.Equal(t, "ZZC-3001-2023", remembered.primaryID)
		assert.Equal(t, "ZZC-3001-2023", currentIDOf(t, researcherUID))
	})

	t.Run("the current ID is always a member of the array", func(t *testing.T) {
		researcherUID := newResearcher(t, "ZZD-4000-2014")

		remembered, err := repo.rememberResearcherID(researcherUID, "ZZD-4001-2020", userUID, false)
		require.NoError(t, err)
		assert.Contains(t, remembered.researcherIDs, remembered.primaryID)
	})
}
