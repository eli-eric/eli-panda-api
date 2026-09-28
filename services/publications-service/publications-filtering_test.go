package publicationsservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"
	"panda/apigateway/services/testsetup"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ELIPANDA-503: server-side column filtering and sorting for publications.
//
// Every fixture publication carries the ELI503TEST marker plus a unique
// substring inside its title ("eli503test"), so each test isolates itself by
// ANDing a title filter onto the filter under test. That keeps the
// expectations independent of any other rows living in the shared test graph.

const filterTestMarker = "ELI503TEST"

type filterFixture struct {
	pubA, pubB, pubC  string
	mediaJ, mediaC    string
	researcher, grant string
	departmentUid     string
}

func createFilterFixtures(t *testing.T) filterFixture {
	t.Helper()

	f := filterFixture{
		pubA:          "eli503-a-" + t.Name(),
		pubB:          "eli503-b-" + t.Name(),
		pubC:          "eli503-c-" + t.Name(),
		mediaJ:        "eli503-media-j-" + t.Name(),
		mediaC:        "eli503-media-c-" + t.Name(),
		researcher:    "eli503-res-" + t.Name(),
		grant:         "eli503-grant-" + t.Name(),
		departmentUid: "eli503-dept-" + t.Name(),
	}

	_, err := testsetup.TestSession.Run(`
		CREATE (a:Publication {
			testMarker: $marker, uid: $pubA, code: 'ELI-2024-001',
			title: 'ELI503TEST Quantum Optics Advances', doi: '10.1000/qoa.1',
			yearOfPublication: '2024', quartil: 'Q1', quartilBasis: 'JCR',
			language: 'eng', eliPublication: 'YES', impactFactor: 5.0,
			allAuthorsCount: 12, eliAuthorsCount: 3, pagesCount: 12, volume: 6, issue: 1,
			dateOfPublication: '2024-02-13',
			authorsDepartmentsArray: [$deptUid + '||Optics||3', 'other-dept||Other||1'],
			deleted: false, wosNumber: 'ELI503TEST-WOS-A' })
		CREATE (b:Publication {
			testMarker: $marker, uid: $pubB, code: 'ELI-2023-002',
			title: 'ELI503TEST Plasma Physics Review', doi: '10.1000/ppr.2',
			yearOfPublication: '2023', quartil: 'Q2', quartilBasis: 'JCR',
			language: 'hun', eliPublication: 'YES', impactFactor: 2.5,
			allAuthorsCount: 5, eliAuthorsCount: 1, pagesCount: 30, volume: 12, issue: 2,
			dateOfPublication: '2023', deleted: false, wosNumber: 'ELI503TEST-WOS-B' })
		CREATE (c:Publication {
			testMarker: $marker, uid: $pubC, code: 'ELI-2024-003',
			title: 'ELI503TEST Quantum Simulation Notes', doi: '10.1000/qsn.3',
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
		"marker":     filterTestMarker,
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
	return f
}

func cleanupFilterFixtures(t *testing.T) {
	t.Helper()
	_, err := testsetup.TestSession.Run(
		"MATCH (n) WHERE n.testMarker = $marker DETACH DELETE n",
		map[string]interface{}{"marker": filterTestMarker})
	require.NoError(t, err)
}

// isolationFilter restricts every query to this test file's fixtures. The
// marker lives on wosNumber so tests can still filter their own target column
// (duplicate filter ids would otherwise shadow each other).
func isolationFilter() []helpers.ColumnFilter {
	return []helpers.ColumnFilter{{Id: "wosNumber", Value: "eli503test"}}
}

func withIsolation(extra ...helpers.ColumnFilter) *[]helpers.ColumnFilter {
	all := append(isolationFilter(), extra...)
	return &all
}

func filterFixtureUids(publications []models.Publication) []string {
	uids := make([]string, 0, len(publications))
	for _, p := range publications {
		uids = append(uids, p.Uid)
	}
	return uids
}

func TestTextFilterMatchesCaseInsensitiveSubstring(t *testing.T) {
	createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, totalCount, err := service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "title", Value: "QUANTUM"}))

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"eli503-a-" + t.Name(), "eli503-c-" + t.Name()}, filterFixtureUids(result))
	assert.Equal(t, int64(2), totalCount)
}

func TestListFiltersUseMembership(t *testing.T) {
	createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, _, err := service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "yearOfPublication", Value: []any{"2024"}}))
	require.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{"eli503-a-" + t.Name(), "eli503-c-" + t.Name()},
		filterFixtureUids(result))

	result, _, err = service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "quartil", Value: []any{"Q1", "Q3"}}))
	require.NoError(t, err)
	assert.Equal(t, []string{"eli503-a-" + t.Name()}, filterFixtureUids(result))
}

func TestNumericRangeFilterIsInclusiveAndNullSafe(t *testing.T) {
	createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// Only-open upper bound.
	result, _, err := service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "impactFactor", Value: map[string]any{"min": 3.0}}))
	require.NoError(t, err)
	assert.Equal(t, []string{"eli503-a-" + t.Name()}, filterFixtureUids(result))

	// An empty range still applies: the null-impactFactor publication C drops
	// out, the two valued ones stay.
	result, _, err = service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "impactFactor", Value: map[string]any{}}))
	require.NoError(t, err)
	assert.Len(t, result, 2)

	result, _, err = service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "impactFactor", Value: map[string]any{"min": 2.5, "max": 2.5}}))
	require.NoError(t, err)
	assert.Equal(t, []string{"eli503-b-" + t.Name()}, filterFixtureUids(result))
}

func TestDateRangeFilterComparesIsoStrings(t *testing.T) {
	createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// '2023' sorts below '2024-01-01'; the null-date publication is excluded.
	result, _, err := service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "dateOfPublication", Value: map[string]any{"min": "2024-01-01"}}))
	require.NoError(t, err)
	assert.Equal(t, []string{"eli503-a-" + t.Name()}, filterFixtureUids(result))
	result, _, err = service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "dateOfPublication", Value: map[string]any{"min": "2023", "max": "2024-02-12"}}))
	require.NoError(t, err)
	assert.Equal(t, []string{"eli503-b-" + t.Name()}, filterFixtureUids(result))
}

func TestCodebookRelationshipFilterAcceptsObjectAndList(t *testing.T) {
	f := createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, _, err := service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "mediaTypeCb", Value: map[string]any{"uid": f.mediaJ, "name": "J - Peer-reviewed article"}}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))

	// Checkbox-group shape: a list of codebook objects.
	result, _, err = service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "mediaTypeCb", Value: []any{
			map[string]any{"uid": f.mediaC, "name": "C - Book chapter"},
		}}))
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{f.pubB, f.pubC}, filterFixtureUids(result))
}

func TestRelationshipListFiltersGrantAndResearcher(t *testing.T) {
	f := createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, _, err := service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "eliResearchers", Value: f.researcher}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))

	result, _, err = service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "grants", Value: f.grant}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))
}

func TestDepartmentFilterMatchesDenormalizedArrayPrefix(t *testing.T) {
	f := createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, _, err := service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "department", Value: f.departmentUid}))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))

	// A uid that is only a substring (not a prefix) of an entry matches nothing.
	result, _, err = service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "department", Value: "ptics"}))
	require.NoError(t, err)
	assert.Empty(t, result)

	// Nor does a uid that is a strict prefix of another department's uid: the
	// match ends at the "||" separator, so "…-dept-" is not "…-dept-<name>".
	result, _, err = service.GetPublications("", 1, 100, nil,
		withIsolation(helpers.ColumnFilter{Id: "department", Value: "eli503-dept-"}))
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestFiltersCombineWithAndAndCountMatchesData(t *testing.T) {
	f := createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	result, totalCount, err := service.GetPublications("", 1, 100, nil, withIsolation(
		helpers.ColumnFilter{Id: "yearOfPublication", Value: []any{"2024"}},
		helpers.ColumnFilter{Id: "eliPublication", Value: []any{"YES"}},
	))
	require.NoError(t, err)
	assert.Equal(t, []string{f.pubA}, filterFixtureUids(result))
	assert.Equal(t, int64(len(result)), totalCount)
}

func TestEmptyFilterValuesAreIgnored(t *testing.T) {
	createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	_, unfilteredCount, err := service.GetPublications("", 1, 100, nil, withIsolation())
	require.NoError(t, err)

	// Blank text and empty lists must not restrict the result set.
	result, totalCount, err := service.GetPublications("", 1, 100, nil, withIsolation(
		helpers.ColumnFilter{Id: "doi", Value: "   "},
		helpers.ColumnFilter{Id: "quartil", Value: []any{}},
	))
	require.NoError(t, err)
	assert.Equal(t, unfilteredCount, totalCount)
	assert.Len(t, result, int(unfilteredCount))
}

func TestSortingByRelationshipColumnAndUnknownIdsIgnored(t *testing.T) {
	f := createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// An unknown sort id (and a collection column) must be dropped, not pasted
	// into a property path; the valid relationship sort still applies.
	sorting := []helpers.Sorting{
		{ID: "notAColumnAtAll", DESC: false},
		{ID: "eliResearchers", DESC: true},
		{ID: "mediaTypeCb", DESC: true},
	}
	result, _, err := service.GetPublications("", 1, 100, &sorting, withIsolation())
	require.NoError(t, err)
	require.NotEmpty(t, result)
	// DESC on media type name: "J - ..." sorts above "C - ...".
	assert.Equal(t, f.pubA, result[0].Uid)

	sorting = []helpers.Sorting{{ID: "impactFactor", DESC: true}}
	result, _, err = service.GetPublications("", 1, 100, &sorting, withIsolation())
	require.NoError(t, err)
	assert.Equal(t, f.pubA, result[0].Uid)
}

func TestGetPublicationFilterOptionsAggregatesValuesAndBounds(t *testing.T) {
	createFilterFixtures(t)
	defer cleanupFilterFixtures(t)
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
