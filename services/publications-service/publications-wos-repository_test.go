package publicationsservice

import (
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
