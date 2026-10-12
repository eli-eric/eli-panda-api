package publicationsservice

import (
	"fmt"
	"sync"
	"testing"

	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"
	"panda/apigateway/services/testsetup"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ELIPANDA-501: researchers can own several WoS ResearcherIDs; an editor who
// confirms a match can teach PANDA the id (PATCH /v1/researcher/:uid/
// researcher-ids) so the next import matches on the first tier.

const researcherIDsTestMarker = "ELI501TEST"

func createResearcherIDsFixture(t *testing.T, uid string) string {
	t.Helper()

	userUID := "eli501-test-user-" + uuid.NewString()
	marker := researcherIDsTestMarker + "-" + t.Name() + "-" + uuid.NewString()
	t.Cleanup(func() {
		session := testsetup.TestDriver.NewSession(neo4j.SessionConfig{})
		defer session.Close()
		deleted, err := session.Run(
			"MATCH (n) WHERE n.testMarker = $marker DETACH DELETE n",
			map[string]interface{}{"marker": marker})
		require.NoError(t, err)
		_, err = deleted.Consume()
		require.NoError(t, err)
	})
	session := testsetup.TestDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	created, err := session.Run(`
		CREATE (u:User {uid: $userUid, testMarker: $marker})
		CREATE (r:Researcher {
			uid: $uid, firstName: 'Jane', lastName: 'Doe',
			researcherId: 'AAA-1111-2022', researcherIds: ['AAA-1111-2022'],
			deleted: false, testMarker: $marker })`,
		map[string]interface{}{"uid": uid, "marker": marker, "userUid": userUID})
	require.NoError(t, err)
	_, err = created.Consume()
	require.NoError(t, err)
	return userUID
}

func TestAppendResearcherIDsMergesDeduplicatesAndPersists(t *testing.T) {
	uid := "eli501-append-" + t.Name() + "-" + uuid.NewString()
	researcherIDsUserUID := createResearcherIDsFixture(t, uid)
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

func storedResearcherIDs(t *testing.T, uid string) []string {
	t.Helper()
	record, err := testsetup.TestSession.Run(
		"MATCH (r:Researcher {uid: $uid}) RETURN r.researcherIds AS ids",
		map[string]interface{}{"uid": uid})
	require.NoError(t, err)
	row, err := record.Single()
	require.NoError(t, err)
	ids := []string{}
	for _, id := range row.Values[0].([]interface{}) {
		ids = append(ids, id.(string))
	}
	return ids
}

func TestAppendResearcherIDsKeepsConcurrentAppends(t *testing.T) {
	uid := "eli501-concurrent-" + t.Name() + "-" + uuid.NewString()
	researcherIDsUserUID := createResearcherIDsFixture(t, uid)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// Editors confirming different WoS authors of the same person at the same
	// time must all be remembered; a read-merge-write keeps only the last one.
	ids := make([]string, 20)
	for index := range ids {
		ids[index] = fmt.Sprintf("B%c-%04d-2020", 'A'+index, index)
	}
	var group sync.WaitGroup
	errs := make(chan error, len(ids))
	for _, id := range ids {
		group.Add(1)
		go func(id string) {
			defer group.Done()
			_, err := service.AppendResearcherIDs(uid, []string{id}, researcherIDsUserUID)
			errs <- err
		}(id)
	}
	group.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	assert.ElementsMatch(t, append([]string{"AAA-1111-2022"}, ids...), storedResearcherIDs(t, uid))
}

func TestAppendResearcherIDsKeepsExistingEntriesVerbatim(t *testing.T) {
	uid := "eli501-legacy-" + t.Name() + "-" + uuid.NewString()
	researcherIDsUserUID := createResearcherIDsFixture(t, uid)
	// Values that predate the ResearcherID format (or its casing) are the
	// editors' data; an unrelated append must neither drop nor rewrite them.
	_, err := testsetup.TestSession.Run(
		"MATCH (r:Researcher {uid: $uid}) SET r.researcherIds = ['aaa-1111-2022', 'legacy 42']",
		map[string]interface{}{"uid": uid})
	require.NoError(t, err)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// AAA-1111-2022 is already known under another casing, so only BBB is new.
	result, err := service.AppendResearcherIDs(uid, []string{"AAA-1111-2022", "bbb-2222-2023"}, researcherIDsUserUID)
	require.NoError(t, err)
	expected := []string{"aaa-1111-2022", "legacy 42", "BBB-2222-2023"}
	assert.Equal(t, expected, result.ResearcherIDs)
	assert.Equal(t, expected, storedResearcherIDs(t, uid))
}

func TestUpdateResearcherPreservesIDsLearnedAfterFormLoad(t *testing.T) {
	uid := "eli501-stale-put-" + t.Name() + "-" + uuid.NewString()
	researcherIDsUserUID := createResearcherIDsFixture(t, uid)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	staleForm, err := service.GetResearcherByUid(uid)
	require.NoError(t, err)
	_, err = service.AppendResearcherIDs(uid, []string{"BBB-2222-2023"}, researcherIDsUserUID)
	require.NoError(t, err)

	staleForm.FirstName = "Janet"
	staleForm.ResearcherIDs = nil
	updated, err := service.UpdateResearcher(&staleForm, researcherIDsUserUID)
	require.NoError(t, err)
	assert.Equal(t, "Janet", updated.FirstName)
	assert.Equal(t, []string{"AAA-1111-2022", "BBB-2222-2023"}, updated.ResearcherIDs)
	assert.Equal(t, updated.ResearcherIDs, storedResearcherIDs(t, uid))
}

func TestUpdateResearcherBodyCarriedIDsReplaceTheList(t *testing.T) {
	uid := "eli501-replace-put-" + t.Name() + "-" + uuid.NewString()
	researcherIDsUserUID := createResearcherIDsFixture(t, uid)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// The correction path: a body that carries researcherIds is authoritative,
	// so a wrong id learned via "remember this match" can be fixed from the UI.
	form, err := service.GetResearcherByUid(uid)
	require.NoError(t, err)
	form.ResearcherIDs = []string{"AAA-1111-2022", "EEE-5555-2021"}
	updated, err := service.UpdateResearcher(&form, researcherIDsUserUID)
	require.NoError(t, err)
	assert.Equal(t, []string{"AAA-1111-2022", "EEE-5555-2021"}, updated.ResearcherIDs)
	assert.Equal(t, updated.ResearcherIDs, storedResearcherIDs(t, uid))

	// An explicit empty list clears; a body that omits the field keeps it.
	form.ResearcherIDs = []string{}
	updated, err = service.UpdateResearcher(&form, researcherIDsUserUID)
	require.NoError(t, err)
	assert.Empty(t, storedResearcherIDs(t, uid))

	form.ResearcherIDs = nil
	_, err = service.UpdateResearcher(&form, researcherIDsUserUID)
	require.NoError(t, err)
	assert.Empty(t, storedResearcherIDs(t, uid))
}

func TestBuildResearcherUpdateQueryIDReplacement(t *testing.T) {
	stored := models.Researcher{Uid: "query-only", ResearcherIDs: []string{"AAA-1111-2022"}}
	for _, test := range []struct {
		name string
		ids  []string
	}{
		{"omitted", nil},
		{"unchanged-but-authoritative", []string{"AAA-1111-2022"}},
		{"clear", []string{}},
		{"replace", []string{"BBB-2222-2023"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			form := stored
			form.ResearcherIDs = test.ids
			query := buildResearcherUpdateQuery(form, stored)
			if test.ids == nil {
				assert.NotContains(t, query.Query, "SET n.researcherIds")
				assert.NotContains(t, query.Parameters, "researcherIds")
			} else {
				assert.Contains(t, query.Query, "SET n.researcherIds = $researcherIds")
				assert.Equal(t, test.ids, query.Parameters["researcherIds"])
			}
			assert.Equal(t, test.ids, form.ResearcherIDs)
		})
	}
}

func TestAppendResearcherIDsDeletedResearcherFails(t *testing.T) {
	uid := "eli501-deleted-" + t.Name() + "-" + uuid.NewString()
	researcherIDsUserUID := createResearcherIDsFixture(t, uid)
	_, err := testsetup.TestSession.Run(
		"MATCH (r:Researcher {uid: $uid}) SET r.deleted = true",
		map[string]interface{}{"uid": uid})
	require.NoError(t, err)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// The write matches only live nodes; reporting success on a deleted
	// researcher would hand the caller a 200 with unchanged data.
	_, err = service.AppendResearcherIDs(uid, []string{"BBB-2222-2023"}, researcherIDsUserUID)
	assert.ErrorIs(t, err, helpers.ERR_NO_ROWS)
}

func TestAppendResearcherIDsUnknownResearcherFails(t *testing.T) {
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	_, err := service.AppendResearcherIDs("eli501-does-not-exist", []string{"AAA-1111-2022"}, "test-user")
	assert.Error(t, err)
}

func TestListResearchersForWosMergesPrimaryAndLearnedIDs(t *testing.T) {
	uid := "eli501-list-" + t.Name() + "-" + uuid.NewString()
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
