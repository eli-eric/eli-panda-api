package publicationsservice

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"
	"panda/apigateway/services/testsetup"
)

// seedFilterablePublications creates a small, self-contained set of
// publications plus the codebook and researcher nodes they point at, so each
// filter can be asserted against a known answer. Every node carries the run's
// own suffix, so the fixtures never collide with the register's real data.
func seedFilterablePublications(t *testing.T) (suffix string, service IPublicationsService) {
	t.Helper()

	suffix = uuid.NewString()
	service = NewPublicationsService(&testsetup.TestDriver, "", "")

	_, err := testsetup.TestSession.Run(`
		CREATE (mt:MediaType {uid: "media-j-" + $s, name: "Journal article", code: "J"})
		CREATE (mtOther:MediaType {uid: "media-c-" + $s, name: "Book chapter", code: "C"})
		CREATE (country:Country {uid: "country-cz-" + $s, name: "Czechia", code: "CZ"})
		CREATE (res:Researcher {uid: "researcher-" + $s, firstName: "Ada", lastName: "Lovelace"})
		CREATE (grant:Grant {uid: "grant-" + $s, code: "G-1", name: "Grant One"})

		CREATE (a:Publication {
			uid: "pub-a-" + $s, code: "PUB-A-" + $s,
			title: "Laser beam diagnostics", yearOfPublication: "2024",
			eliPublication: "YES", impactFactor: 5.5, allAuthorsCount: 3,
			dateOfPublication: "2024-06", quartil: "Q1", language: "eng",
			authorsDepartmentsArray: ["dept-optics-" + $s + "||Optics||2"],
			deleted: false
		})
		CREATE (b:Publication {
			uid: "pub-b-" + $s, code: "PUB-B-" + $s,
			title: "Plasma acceleration study", yearOfPublication: "2023",
			eliPublication: "NO", impactFactor: 1.2, allAuthorsCount: 9,
			dateOfPublication: "2023-02", quartil: "Q3", language: "cze",
			authorsDepartmentsArray: ["dept-plasma-" + $s + "||Plasma||1"],
			deleted: false
		})
		CREATE (c:Publication {
			uid: "pub-c-" + $s, code: "PUB-C-" + $s,
			title: "Laser plasma interaction", yearOfPublication: "2025",
			eliPublication: "YES", impactFactor: 9.9, allAuthorsCount: 1,
			dateOfPublication: "2025-11", quartil: "Q1", language: "eng",
			deleted: false
		})

		CREATE (a)-[:HAS_MEDIA_TYPE]->(mt)
		CREATE (b)-[:HAS_MEDIA_TYPE]->(mtOther)
		CREATE (a)-[:HAS_PUBLISHING_COUNTRY]->(country)
		CREATE (a)-[:HAS_RESEARCHER]->(res)
		CREATE (a)-[:HAS_GRANT]->(grant)
		`, map[string]interface{}{"s": suffix})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = testsetup.TestSession.Run(
			`MATCH (n) WHERE n.uid ENDS WITH $s DETACH DELETE n`,
			map[string]interface{}{"s": suffix})
	})

	return suffix, service
}

// codesOf returns the seeded publications' codes, which identify them without
// depending on ordering.
func codesOf(publications []models.Publication) []string {
	codes := make([]string, 0, len(publications))
	for _, publication := range publications {
		codes = append(codes, publication.Code)
	}
	return codes
}

func TestApplyPublicationFiltersEachFilterInIsolation(t *testing.T) {
	suffix, service := seedFilterablePublications(t)

	tests := []struct {
		name      string
		filter    helpers.ColumnFilter
		wantCodes []string
	}{
		{
			name:      "text contains, case insensitive",
			filter:    helpers.ColumnFilter{Id: "title", Value: "LASER"},
			wantCodes: []string{"PUB-A-" + suffix, "PUB-C-" + suffix},
		},
		{
			name:      "list of years",
			filter:    helpers.ColumnFilter{Id: "yearOfPublication", Value: []interface{}{"2024", "2025"}},
			wantCodes: []string{"PUB-A-" + suffix, "PUB-C-" + suffix},
		},
		{
			name:      "eliPublication flag",
			filter:    helpers.ColumnFilter{Id: "eliPublication", Value: []interface{}{"NO"}},
			wantCodes: []string{"PUB-B-" + suffix},
		},
		{
			name:      "numeric range, both bounds",
			filter:    helpers.ColumnFilter{Id: "impactFactor", Value: map[string]interface{}{"min": 2.0, "max": 6.0}},
			wantCodes: []string{"PUB-A-" + suffix},
		},
		{
			name:      "numeric range, lower bound only",
			filter:    helpers.ColumnFilter{Id: "allAuthorsCount", Value: map[string]interface{}{"min": 3.0}},
			wantCodes: []string{"PUB-A-" + suffix, "PUB-B-" + suffix},
		},
		{
			name:      "date range over partial YYYY-MM values",
			filter:    helpers.ColumnFilter{Id: "dateOfPublication", Value: map[string]interface{}{"min": "2024-01", "max": "2024-12"}},
			wantCodes: []string{"PUB-A-" + suffix},
		},
		{
			name:      "codebook relationship by uid",
			filter:    helpers.ColumnFilter{Id: "mediaType", Value: map[string]interface{}{"uid": "media-j-" + suffix, "name": "Journal article"}},
			wantCodes: []string{"PUB-A-" + suffix},
		},
		{
			name: "codebook as a list of uids, from a checkbox group",
			filter: helpers.ColumnFilter{Id: "mediaType", Value: []interface{}{
				"media-j-" + suffix, "media-c-" + suffix,
			}},
			wantCodes: []string{"PUB-A-" + suffix, "PUB-B-" + suffix},
		},
		{
			name:      "researcher relationship list",
			filter:    helpers.ColumnFilter{Id: "eliResearchers", Value: []interface{}{"researcher-" + suffix}},
			wantCodes: []string{"PUB-A-" + suffix},
		},
		{
			name:      "grant relationship, single picked value",
			filter:    helpers.ColumnFilter{Id: "grant", Value: map[string]interface{}{"uid": "grant-" + suffix, "name": "Grant One"}},
			wantCodes: []string{"PUB-A-" + suffix},
		},
		{
			name:      "department over the denormalized array",
			filter:    helpers.ColumnFilter{Id: "department", Value: map[string]interface{}{"uid": "dept-optics-" + suffix, "name": "Optics"}},
			wantCodes: []string{"PUB-A-" + suffix},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filters := []helpers.ColumnFilter{{Id: "code", Value: suffix}, test.filter}
			result, totalCount, err := service.GetPublications("", 1, 100, nil, &filters)
			require.NoError(t, err)

			assert.ElementsMatch(t, test.wantCodes, codesOf(result))
			assert.Equal(t, int64(len(result)), totalCount,
				"totalCount must describe the same filtered set as the page")
		})
	}
}

func TestApplyPublicationFiltersCombineWithAnd(t *testing.T) {
	suffix, service := seedFilterablePublications(t)

	filters := []helpers.ColumnFilter{
		{Id: "code", Value: suffix},
		{Id: "title", Value: "laser"},
		{Id: "eliPublication", Value: []interface{}{"YES"}},
		{Id: "impactFactor", Value: map[string]interface{}{"min": 8.0}},
	}

	result, totalCount, err := service.GetPublications("", 1, 100, nil, &filters)
	require.NoError(t, err)

	// Only PUB-C is a laser paper, ELI-flagged, and above impact factor 8.
	assert.Contains(t, codesOf(result), "PUB-C-"+suffix)
	assert.NotContains(t, codesOf(result), "PUB-A-"+suffix)
	assert.NotContains(t, codesOf(result), "PUB-B-"+suffix)
	assert.Equal(t, int64(len(result)), totalCount)
}

func TestApplyPublicationFiltersIgnoresEmptyValues(t *testing.T) {
	suffix, service := seedFilterablePublications(t)

	baseline, baselineCount, err := service.GetPublications("", 1, 100, nil,
		&[]helpers.ColumnFilter{{Id: "code", Value: suffix}})
	require.NoError(t, err)

	withEmpties := []helpers.ColumnFilter{
		{Id: "code", Value: suffix},
		{Id: "title", Value: "   "},
		{Id: "yearOfPublication", Value: []interface{}{}},
		{Id: "impactFactor", Value: map[string]interface{}{}},
		{Id: "dateOfPublication", Value: map[string]interface{}{"min": ""}},
		{Id: "mediaType", Value: map[string]interface{}{"uid": ""}},
	}
	result, totalCount, err := service.GetPublications("", 1, 100, nil, &withEmpties)
	require.NoError(t, err)

	assert.ElementsMatch(t, codesOf(baseline), codesOf(result),
		"empty filter values must not narrow the result set")
	assert.Equal(t, baselineCount, totalCount)
}

func TestPublicationsTotalCountMatchesFilteredPage(t *testing.T) {
	suffix, service := seedFilterablePublications(t)

	filters := []helpers.ColumnFilter{{Id: "code", Value: suffix}}
	result, totalCount, err := service.GetPublications("", 1, 100, nil, &filters)
	require.NoError(t, err)

	assert.Len(t, result, 3)
	assert.Equal(t, int64(3), totalCount)
}

func TestPublicationsSortingByRelationshipColumn(t *testing.T) {
	suffix, service := seedFilterablePublications(t)

	// Only A and B carry a media type; sorting must run on the related node's
	// name rather than a property of the publication itself.
	filters := []helpers.ColumnFilter{{Id: "code", Value: suffix}}
	sorting := []helpers.Sorting{{ID: "mediaType", DESC: false}}

	result, _, err := service.GetPublications("", 1, 100, &sorting, &filters)
	require.NoError(t, err)
	require.Len(t, result, 3)

	// "Book chapter" sorts before "Journal article"; the row without a media
	// type sorts with the nulls and is not asserted on.
	positions := map[string]int{}
	for index, publication := range result {
		positions[publication.Code] = index
	}
	assert.Less(t, positions["PUB-B-"+suffix], positions["PUB-A-"+suffix],
		"Book chapter must sort before Journal article")
}

func TestMapPublicationSortField(t *testing.T) {
	assert.Equal(t, "n.title", mapPublicationSortField("title"))
	assert.Equal(t, "mediaTypeCb.name", mapPublicationSortField("mediaType"))
	assert.Equal(t, "experimentalSystemCb.name", mapPublicationSortField("experimentalSystem"))

	// Unknown and collected columns are not sortable and must not become a
	// property path.
	assert.Equal(t, "", mapPublicationSortField("eliResearchers"))
	assert.Equal(t, "", mapPublicationSortField("definitely not a column"))
	assert.Equal(t, "", mapPublicationSortField(""))
}

func TestGetPublicationsSortingClauseRejectsUnknownIds(t *testing.T) {
	assert.Equal(t, defaultPublicationsSortingClause, getPublicationsSortingClause(nil))
	assert.Equal(t, defaultPublicationsSortingClause,
		getPublicationsSortingClause(&[]helpers.Sorting{}))

	// An unknown id alone falls back rather than emitting "ORDER BY".
	assert.Equal(t, defaultPublicationsSortingClause,
		getPublicationsSortingClause(&[]helpers.Sorting{{ID: "nope", DESC: true}}))

	// A known id survives alongside an unknown one.
	clause := getPublicationsSortingClause(&[]helpers.Sorting{
		{ID: "nope"}, {ID: "title", DESC: true},
	})
	assert.Contains(t, clause, "n.title")
	assert.NotContains(t, clause, "nope")
}

func TestApplyPublicationFiltersParameterizesEveryValue(t *testing.T) {
	filters := []helpers.ColumnFilter{
		{Id: "title", Value: "'; MATCH (x) DETACH DELETE x //"},
		{Id: "mediaType", Value: map[string]interface{}{"uid": "' OR 1=1 //"}},
	}

	query := buildPublicationsCountQuery("", &filters)

	assert.NotContains(t, query.Query, "DETACH DELETE x")
	assert.NotContains(t, query.Query, "OR 1=1")
	assert.Equal(t, "'; match (x) detach delete x //", query.Parameters["filterTitle"])
	assert.Equal(t, "' OR 1=1 //", query.Parameters["filterMediaType"])
}

func TestApplyPublicationFiltersSurvivesMalformedCodebookValues(t *testing.T) {
	// helpers.GetFilterValueCodebook asserts on "name" unchecked and panics
	// when it is absent, so a hand-written columnFilter could take the endpoint
	// down. Building the query must stay a pure function of whatever arrives.
	malformed := []helpers.ColumnFilter{
		{Id: "mediaType", Value: map[string]interface{}{"uid": "media-uid"}},
		{Id: "openAccessType", Value: map[string]interface{}{"name": "no uid here"}},
		{Id: "publishingCountry", Value: "not an object at all"},
		{Id: "grant", Value: map[string]interface{}{"uid": 42}},
		{Id: "department", Value: map[string]interface{}{}},
	}

	require.NotPanics(t, func() {
		query := buildPublicationsCountQuery("", &malformed)

		// The one usable value still filters; the rest are ignored.
		assert.Equal(t, "media-uid", query.Parameters["filterMediaType"])
		assert.NotContains(t, query.Parameters, "filterOpenAccessType")
		assert.NotContains(t, query.Parameters, "filterPublishingCountry")
		assert.NotContains(t, query.Parameters, "filterGrant")
		assert.NotContains(t, query.Parameters, "filterDepartment")
	})
}
