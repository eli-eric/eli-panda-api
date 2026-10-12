package publicationsservice

import (
	"encoding/json"
	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"
	"panda/apigateway/services/testsetup"
	"testing"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ELIPANDA-503: server-side column filtering and sorting for publications.
//
// Every fixture run has a unique marker in its titles. Queries use the normal
// search parameter to isolate that run without occupying any column filter,
// so even title and wosNumber filters can be tested independently.

const filterTestMarker = "ELI503TEST"

type filterFixture struct {
	pubA, pubB, pubC  string
	mediaJ, mediaC    string
	researcher, grant string
	departmentUid     string
	marker            string
}

func createFilterFixtures(t *testing.T) filterFixture {
	t.Helper()

	marker := filterTestMarker + "-" + t.Name() + "-" + uuid.NewString()
	f := filterFixture{
		pubA:          "eli503-a-" + marker,
		pubB:          "eli503-b-" + marker,
		pubC:          "eli503-c-" + marker,
		mediaJ:        "eli503-media-j-" + marker,
		mediaC:        "eli503-media-c-" + marker,
		researcher:    "eli503-res-" + marker,
		grant:         "eli503-grant-" + marker,
		departmentUid: "eli503-dept-" + marker,
		marker:        marker,
	}
	t.Cleanup(func() { cleanupFilterFixtures(t, f) })
	session := testsetup.TestDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	created, err := session.Run(`
		CREATE (a:Publication {
			testMarker: $marker, uid: $pubA, code: 'ELI-2024-001',
			title: $marker + ' Quantum Optics Advances', doi: '10.1000/qoa.1',
			yearOfPublication: '2024', quartil: 'Q1', quartilBasis: 'JCR',
			language: 'eng', eliPublication: 'YES', impactFactor: 5.0,
			allAuthorsCount: 12, eliAuthorsCount: 3, pagesCount: 12, volume: 6, issue: 1,
			dateOfPublication: '2024-02-13',
			authorsDepartmentsArray: [$deptUid + '||Optics||3', 'other-dept||Other||1'],
			deleted: false, wosNumber: 'ELI503TEST-WOS-A' })
		CREATE (b:Publication {
			testMarker: $marker, uid: $pubB, code: 'ELI-2023-002',
			title: $marker + ' Plasma Physics Review', doi: '10.1000/ppr.2',
			yearOfPublication: '2023', quartil: 'Q2', quartilBasis: 'JCR',
			language: 'hun', eliPublication: 'YES', impactFactor: 2.5,
			allAuthorsCount: 5, eliAuthorsCount: 1, pagesCount: 30, volume: 12, issue: 2,
			dateOfPublication: '2023', deleted: false, wosNumber: 'ELI503TEST-WOS-B' })
		CREATE (c:Publication {
			testMarker: $marker, uid: $pubC, code: 'ELI-2024-003',
			title: $marker + ' Quantum Simulation Notes', doi: '10.1000/qsn.3',
			yearOfPublication: '2024', eliPublication: 'NO',
			allAuthorsCount: 2, eliAuthorsCount: 0, pagesCount: 4, volume: 1, issue: 1,
			deleted: false, wosNumber: 'ELI503TEST-WOS-C' })
		CREATE (j:MediaType {testMarker: $marker, uid: $mediaJ, name: 'J - Peer-reviewed article', code: 'J'})
		CREATE (mc:MediaType {testMarker: $marker, uid: $mediaC, name: 'C - Book chapter', code: 'C'})
		CREATE (r:Researcher {testMarker: $marker, uid: $researcher, firstName: 'Jan', lastName: 'Novak'})
		CREATE (g:Grant {testMarker: $marker, uid: $grant, code: 'CZ.2.01'})
		CREATE (a)-[:HAS_MEDIA_TYPE]->(j)
		CREATE (b)-[:HAS_MEDIA_TYPE]->(mc)
		CREATE (c)-[:HAS_MEDIA_TYPE]->(mc)
		CREATE (a)-[:HAS_RESEARCHER]->(r)
		CREATE (a)-[:HAS_GRANT]->(g)
	`, map[string]interface{}{
		"marker":     f.marker,
		"pubA":       f.pubA,
		"pubB":       f.pubB,
		"pubC":       f.pubC,
		"mediaJ":     f.mediaJ,
		"mediaC":     f.mediaC,
		"researcher": f.researcher,
		"grant":      f.grant,
		"deptUid":    f.departmentUid,
	})
	require.NoError(t, err)
	_, err = created.Consume()
	require.NoError(t, err)
	return f
}

func cleanupFilterFixtures(t *testing.T, f filterFixture) {
	t.Helper()
	session := testsetup.TestDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	deleted, err := session.Run(
		"MATCH (n) WHERE n.testMarker = $marker DETACH DELETE n",
		map[string]interface{}{"marker": f.marker})
	require.NoError(t, err)
	_, err = deleted.Consume()
	require.NoError(t, err)
}

func publicationFilters(extra ...helpers.ColumnFilter) *[]helpers.ColumnFilter {
	return &extra
}

func filterFixtureUids(publications []models.Publication) []string {
	uids := make([]string, 0, len(publications))
	for _, p := range publications {
		uids = append(uids, p.Uid)
	}
	return uids
}

func TestTextFilterMatchesCaseInsensitiveSubstring(t *testing.T) {
	f := createFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, totalCount, err := service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "title", Value: "QUANTUM"}))

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{f.pubA, f.pubC}, filterFixtureUids(result))
	assert.Equal(t, int64(2), totalCount)
}

func TestNormalizeFilterOptionsEmptyArrays(t *testing.T) {
	options := normalizeFilterOptions(models.PublicationFilterOptions{})
	encoded, err := json.Marshal(options)
	require.NoError(t, err)
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(encoded, &fields))
	for _, field := range []string{"years", "quartils", "quartilBases", "languages", "eliPublications"} {
		assert.JSONEq(t, "[]", string(fields[field]), field)
	}
	assert.JSONEq(t, "{}", string(fields["ranges"]))
	assert.JSONEq(t, "{}", string(fields["dateBounds"]))
}

func TestWosNumberFilterDoesNotShadowFixtureIsolation(t *testing.T) {
	f := createFilterFixtures(t)
	other := createFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, totalCount, err := service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "wosNumber", Value: "WOS-B"}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubB}, filterFixtureUids(result))
	assert.Equal(t, int64(1), totalCount)

	cleanupFilterFixtures(t, f)
	result, totalCount, err = service.GetPublications(other.marker, 1, 100, nil, nil)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{other.pubA, other.pubB, other.pubC}, filterFixtureUids(result))
	assert.Equal(t, int64(3), totalCount)
}

func TestListFiltersUseMembership(t *testing.T) {
	f := createFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, _, err := service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "yearOfPublication", Value: []any{"2024"}}))
	require.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{f.pubA, f.pubC},
		filterFixtureUids(result))

	result, _, err = service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "quartil", Value: []any{"Q1", "Q3"}}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))
}

func TestNumericRangeFilterIsInclusiveAndNullSafe(t *testing.T) {
	f := createFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// Only-open upper bound.
	result, _, err := service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "impactFactor", Value: map[string]any{"min": 3.0}}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))

	// An empty range is a no-op: both bounds open means "no constraint", so
	// the null-impactFactor publication C stays in — emitting a bare
	// IS NOT NULL here would silently drop every null-valued row.
	result, _, err = service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "impactFactor", Value: map[string]any{}}))
	require.NoError(t, err)
	assert.Len(t, result, 3)

	result, _, err = service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "impactFactor", Value: map[string]any{"min": 2.5, "max": 2.5}}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubB}, filterFixtureUids(result))
}

func TestApplyPublicationFiltersMalformedInputDoesNotPanic(t *testing.T) {
	// A hand-written columnFilter (anything but what the sheet emits) must be
	// ignored or coerced, never panic the request goroutine into a 500.
	filters := []helpers.ColumnFilter{
		{Id: "title", Value: 5},
		{Id: "code", Value: ""},
		{Id: "yearOfPublication", Value: "2024"},
		{Id: "quartil", Value: []interface{}{2024, nil, "Q1"}},
		{Id: "impactFactor", Value: map[string]any{"min": "3"}},
		{Id: "volume", Value: "not-a-range"},
		{Id: "dateOfPublication", Value: map[string]any{"min": 2024.0}},
		{Id: "mediaTypeCb", Value: 42},
	}
	query := helpers.DatabaseQuery{
		Query:      "MATCH (n:Publication) WHERE true",
		Parameters: map[string]interface{}{},
	}
	require.NotPanics(t, func() { ApplyPublicationFilters(&query, &filters) })

	// Wrong-typed values produce no predicate at all (a bare string for a list
	// field is malformed input, not a one-item list) ...
	assert.NotContains(t, query.Query, "filterTitle")
	assert.NotContains(t, query.Query, "filterCode")
	assert.NotContains(t, query.Query, "filterYearOfPublication")
	assert.NotContains(t, query.Query, "filterVolume")
	assert.NotContains(t, query.Query, "filterMediaTypeCb")
	// ... non-string items in a string list are dropped ...
	assert.Equal(t, []string{"Q1"}, query.Parameters["filterQuartil"])
	// ... numeric-string bounds are coerced to numbers ...
	assert.Equal(t, 3.0, *query.Parameters["filterImpactFactorFrom"].(*float64))
	// ... and numeric date bounds are stringified for the lexicographic range.
	assert.Equal(t, "2024", query.Parameters["filterDateOfPublicationFrom"])
}

func TestDateRangeFilterComparesIsoStrings(t *testing.T) {
	f := createFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// '2023' sorts below '2024-01-01'; the null-date publication is excluded.
	result, _, err := service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "dateOfPublication", Value: map[string]any{"min": "2024-01-01"}}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))
	result, _, err = service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "dateOfPublication", Value: map[string]any{"min": "2023", "max": "2024-02-12"}}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubB}, filterFixtureUids(result))
}

func TestCodebookRelationshipFilterAcceptsObjectAndList(t *testing.T) {
	f := createFilterFixtures(t)

	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, _, err := service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "mediaTypeCb", Value: map[string]any{"uid": f.mediaJ, "name": "J - Peer-reviewed article"}}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))

	// Checkbox-group shape: a list of codebook objects.
	result, _, err = service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "mediaTypeCb", Value: []any{
			map[string]any{"uid": f.mediaC, "name": "C - Book chapter"},
		}}))
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{f.pubB, f.pubC}, filterFixtureUids(result))
}

func TestRelationshipListFiltersGrantAndResearcher(t *testing.T) {
	f := createFilterFixtures(t)

	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, _, err := service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "eliResearchers", Value: f.researcher}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))

	result, _, err = service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "grants", Value: f.grant}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))
}

func TestDepartmentFilterMatchesDenormalizedArrayPrefix(t *testing.T) {
	f := createFilterFixtures(t)

	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, _, err := service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "department", Value: f.departmentUid}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))

	// A uid that is only a substring (not a prefix) of an entry matches nothing.
	result, _, err = service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "department", Value: "ptics"}))
	require.NoError(t, err)
	assert.Empty(t, result)

	// Nor does a uid that is a strict prefix of another department's uid: the
	// match ends at the "||" separator, so "…-dept-" is not "…-dept-<name>".
	result, _, err = service.GetPublications(f.marker, 1, 100, nil,
		publicationFilters(helpers.ColumnFilter{Id: "department", Value: "eli503-dept-"}))
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestFiltersCombineWithAndAndCountMatchesData(t *testing.T) {
	f := createFilterFixtures(t)

	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, totalCount, err := service.GetPublications(f.marker, 1, 100, nil, publicationFilters(
		helpers.ColumnFilter{Id: "yearOfPublication", Value: []any{"2024"}},
		helpers.ColumnFilter{Id: "eliPublication", Value: []any{"YES"}},
	))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))
	assert.Equal(t, int64(len(result)), totalCount)
}

func TestEmptyFilterValuesAreIgnored(t *testing.T) {
	f := createFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")
	_, unfilteredCount, err := service.GetPublications(f.marker, 1, 100, nil, publicationFilters())
	require.NoError(t, err)

	// Blank text and empty lists must not restrict the result set.
	result, totalCount, err := service.GetPublications(f.marker, 1, 100, nil, publicationFilters(
		helpers.ColumnFilter{Id: "doi", Value: "   "},
		helpers.ColumnFilter{Id: "quartil", Value: []any{}},
	))
	require.NoError(t, err)
	assert.Equal(t, unfilteredCount, totalCount)
	assert.Len(t, result, int(unfilteredCount))
}

func TestSortingByRelationshipColumnAndUnknownIdsIgnored(t *testing.T) {
	f := createFilterFixtures(t)

	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// An unknown sort id (and a collection column) must be dropped, not pasted
	// into a property path; the valid relationship sort still applies.
	sorting := []helpers.Sorting{
		{ID: "notAColumnAtAll", DESC: false},
		{ID: "eliResearchers", DESC: true},
		{ID: "mediaTypeCb", DESC: true},
	}
	result, _, err := service.GetPublications(f.marker, 1, 100, &sorting, publicationFilters())
	require.NoError(t, err)
	require.NotEmpty(t, result)
	// DESC on media type name: "J - ..." sorts above "C - ...".
	assert.Equal(t, f.pubA, result[0].Uid)

	sorting = []helpers.Sorting{{ID: "impactFactor", DESC: true}}
	result, _, err = service.GetPublications(f.marker, 1, 100, &sorting, publicationFilters())
	require.NoError(t, err)
	assert.Equal(t, f.pubA, result[0].Uid)
}

func TestGetPublicationFilterOptionsAggregatesValuesAndBounds(t *testing.T) {
	createFilterFixtures(t)

	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	options, err := service.GetPublicationFilterOptions()
	require.NoError(t, err)

	assert.Contains(t, options.Years, "2024")
	assert.Contains(t, options.Years, "2023")
	// Years come back newest first.
	require.GreaterOrEqual(t, len(options.Years), 2)
	assert.Less(t, indexOfValue(options.Years, "2024"), indexOfValue(options.Years, "2023"))

	assert.Contains(t, options.Quartils, "Q1")
	require.NotNil(t, options.Ranges["impactFactor"].Min)
	require.NotNil(t, options.Ranges["impactFactor"].Max)
	// The global bounds must cover the fixture values (other rows may widen them).
	assert.LessOrEqual(t, *options.Ranges["impactFactor"].Min, 2.5)
	assert.GreaterOrEqual(t, *options.Ranges["impactFactor"].Max, 5.0)
}

func indexOfValue(values []string, value string) int {
	for i, v := range values {
		if v == value {
			return i
		}
	}
	return -1
}
