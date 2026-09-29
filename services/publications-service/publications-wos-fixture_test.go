package publicationsservice

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	codebookmodels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/publications-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ELIPANDA-501/502 shared fixture: testdata/wos-lookup-raw.json is a Clarivate
// Starter /documents response in the recorded wire format and
// testdata/wos-lookup-preview.json is the exact lookup response PANDA derives
// from it. The frontend repository checks in the same preview and renders its
// dialog tests from it, so the two sides cannot drift without a red test.
//
// Regenerate after an intended contract change with
//
//	UPDATE_WOS_FIXTURE=1 go test ./services/publications-service/ -run TestWosSharedFixture
const (
	wosFixtureRawPath     = "testdata/wos-lookup-raw.json"
	wosFixturePreviewPath = "testdata/wos-lookup-preview.json"
)

func wosFixtureRepository() *stubWosImportRepository {
	return &stubWosImportRepository{
		existingPublication: &models.WosExistingPublication{
			Uid: "pub-eli-2024-017", Code: "ELI-2024-017",
			Title: "Anomalous absorption of relativistic laser pulses in near-critical plasmas",
			Doi:   "10.1103/physrevresearch.6.013126",
		},
		researchers: []wosResearcherRecord{
			{UID: "res-novak", FirstName: "Jan", LastName: "Novák", ResearcherIDs: []string{"AAB-1234-2020"}},
			{UID: "res-dvorak", FirstName: "Petr", LastName: "Dvořák"},
			{UID: "res-svoboda-1", FirstName: "Martin", LastName: "Svoboda"},
			{UID: "res-svoboda-2", FirstName: "Martin", LastName: "Svoboda"},
		},
		mediaType: &codebookmodels.Codebook{UID: "media-j", Name: "J - Peer-reviewed article", Code: "J"},
	}
}

func TestWosSharedFixture(t *testing.T) {
	raw, err := os.ReadFile(wosFixtureRawPath)
	require.NoError(t, err)

	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		writeWosResponse(t, response, string(raw))
	}))
	defer upstream.Close()

	preview, err := newWosTestService(upstream.URL, wosFixtureRepository()).PreviewWosPublication(
		context.Background(), "https://doi.org/10.1103/PhysRevResearch.6.013126", "", "B",
	)
	require.NoError(t, err)

	actual, err := json.MarshalIndent(preview, "", "  ")
	require.NoError(t, err)
	actual = append(actual, '\n')

	if os.Getenv("UPDATE_WOS_FIXTURE") == "1" {
		require.NoError(t, os.MkdirAll(filepath.Dir(wosFixturePreviewPath), 0o755))
		require.NoError(t, os.WriteFile(wosFixturePreviewPath, actual, 0o644))
	}

	expected, err := os.ReadFile(wosFixturePreviewPath)
	require.NoError(t, err)
	assert.JSONEq(t, string(expected), string(actual))

	// The fixture must keep exercising every dialog branch.
	assert.Equal(t, "already-exists", preview.Status)
	assert.Equal(t, "WOS:001164928200001", preview.WosUid)
	assert.NotEmpty(t, preview.RecordUrl)
	confidences := make([]string, 0, len(preview.Authors))
	for _, author := range preview.Authors {
		confidences = append(confidences, author.Match.Confidence)
	}
	assert.Equal(t, []string{"EXACT_ID", "NAME", "AMBIGUOUS", "NONE"}, confidences)
	codes := make([]string, 0, len(preview.Warnings))
	for _, warning := range preview.Warnings {
		codes = append(codes, warning.Code)
	}
	assert.ElementsMatch(t, []string{"ISSUE_NOT_NUMERIC", "DATE_DAY_MISSING"}, codes)
}

func TestWosImportValuesOnlyEmitImportableFields(t *testing.T) {
	// Every key the mapper can emit is a member of wosImportableFields — the
	// list the frontend freezes as WOS_IMPORTABLE_FIELDS.
	importable := make(map[string]struct{}, len(wosImportableFields))
	for _, field := range wosImportableFields {
		importable[field] = struct{}{}
	}

	valuesType := reflect.TypeOf(models.WosImportValues{})
	for index := range valuesType.NumField() {
		name := strings.Split(valuesType.Field(index).Tag.Get("json"), ",")[0]
		_, ok := importable[name]
		assert.Truef(t, ok, "WosImportValues.%s (json %q) is not in wosImportableFields",
			valuesType.Field(index).Name, name)
	}
	assert.Equal(t, valuesType.NumField(), len(wosImportableFields),
		"every importable field needs a WosImportValues slot")
}

func TestNormalizePersonNameFoldsDiacritics(t *testing.T) {
	assert.Equal(t, normalizePersonName("Dvorak Petr"), normalizePersonName("Dvořák, Petr"))
	assert.Equal(t, "sykora jiri", normalizePersonName("Sýkora, Jiří"))
}
