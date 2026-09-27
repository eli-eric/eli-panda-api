package publicationsservice

import (
	"testing"

	"panda/apigateway/services/testsetup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ELIPANDA-501: researchers can own several WoS ResearcherIDs; an editor who
// confirms a match can teach PANDA the id (PATCH /v1/researcher/:uid/
// researcher-ids) so the next import matches on the first tier.

const researcherIDsTestMarker = "ELI501TEST"

// HistoryLogQuery links the changed node to an existing (:User) via
// WAS_UPDATED_BY; a fixture user makes the audit observable.
const researcherIDsUserUID = "eli501-test-user"

func createResearcherIDsFixture(t *testing.T, uid string) {
	t.Helper()

	_, err := testsetup.TestSession.Run(`
		MERGE (u:User {uid: $userUid})
		CREATE (r:Researcher {
			uid: $uid, firstName: 'Jane', lastName: 'Doe',
			researcherId: 'AAA-1111-2022', researcherIds: ['AAA-1111-2022'],
			deleted: false, testMarker: $marker })`,
		map[string]interface{}{"uid": uid, "marker": researcherIDsTestMarker, "userUid": researcherIDsUserUID})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = testsetup.TestSession.Run(
			"MATCH (r:Researcher {uid: $uid}) DETACH DELETE r",
			map[string]interface{}{"uid": uid})
	})
}

func TestAppendResearcherIDsMergesDeduplicatesAndPersists(t *testing.T) {
	uid := "eli501-append-" + t.Name()
	createResearcherIDsFixture(t, uid)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// Lowercase and padded ids normalize; the already-known id deduplicates;
	// a malformed id is skipped by the service (the handler rejects earlier).
	result, err := service.AppendResearcherIDs(
		uid, []string{"bbb-2222-2023", " AAA-1111-2022 ", "not-an-id", "CCC-3333-2021"}, researcherIDsUserUID)
	require.NoError(t, err)
	assert.Equal(t, []string{"AAA-1111-2022", "BBB-2222-2023", "CCC-3333-2021"}, result.ResearcherIDs)

	record, err := testsetup.TestSession.Run(
		"MATCH (r:Researcher {uid: $uid}) RETURN r.researcherIds AS ids",
		map[string]interface{}{"uid": uid})
	require.NoError(t, err)
	row, err := record.Single()
	require.NoError(t, err)
	assert.Equal(t, []interface{}{"AAA-1111-2022", "BBB-2222-2023", "CCC-3333-2021"}, row.Values[0])

	// The change is audited like every other researcher write.
	audited, err := testsetup.TestSession.Run(
		"MATCH (r:Researcher {uid: $uid})-[h:WAS_UPDATED_BY]->(:User {uid: $userUid}) RETURN count(h) AS n",
		map[string]interface{}{"uid": uid, "userUid": researcherIDsUserUID})
	require.NoError(t, err)
	row, err = audited.Single()
	require.NoError(t, err)
	assert.Greater(t, row.Values[0].(int64), int64(0))
}

func TestAppendResearcherIDsUnknownResearcherFails(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	_, err := service.AppendResearcherIDs("eli501-does-not-exist", []string{"AAA-1111-2022"}, "test-user")
	assert.Error(t, err)
}

func TestListResearchersForWosMergesPrimaryAndLearnedIDs(t *testing.T) {
	uid := "eli501-list-" + t.Name()
	_, err := testsetup.TestSession.Run(`
		CREATE (r:Researcher {
			uid: $uid, firstName: 'Ada', lastName: 'Lovelace',
			researcherId: 'AAA-1111-2022',
			researcherIds: ['aaa-1111-2022', 'DDD-4444-2019'],
			deleted: false, testMarker: $marker })`,
		map[string]interface{}{"uid": uid, "marker": researcherIDsTestMarker})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = testsetup.TestSession.Run(
			"MATCH (r:Researcher {uid: $uid}) DETACH DELETE r",
			map[string]interface{}{"uid": uid})
	})

	repository := &neo4jWosImportRepository{driver: &testsetup.TestDriver}
	researchers, err := repository.listResearchersForWos()
	require.NoError(t, err)

	var mine wosResearcherRecord
	for _, researcher := range researchers {
		if researcher.UID == uid {
			mine = researcher
		}
	}
	require.NotEmpty(t, mine.UID, "the fixture researcher must be listed for matching")
	// Primary plus learned ids, all uppercased, duplicates kept out is not
	// required (matching only checks membership), but every id must be present.
	assert.Contains(t, mine.ResearcherIDs, "AAA-1111-2022")
	assert.Contains(t, mine.ResearcherIDs, "DDD-4444-2019")
}
