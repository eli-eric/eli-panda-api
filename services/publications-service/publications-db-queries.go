package publicationsservice

import (
	"fmt"
	"panda/apigateway/helpers"
	"strconv"
	"strings"
)

// ELIPANDA-503: server-side column filtering for publications.
//
// Filter ids are the canonical field names shared with the frontend filter
// sheet and the table columns. Every emitted clause is appended as a plain
// `AND (...)` predicate on the base `MATCH (n:Publication) ... WHERE` block —
// no MATCH/WITH — so the exact same predicate set runs inside both the list
// query and the count query, and can never change result cardinality.
// All values are passed as Cypher parameters; nothing is concatenated.

// publicationTextFilterFields are String properties matched with
// toLower(CONTAINS), per the ELIPANDA-503 field list.
var publicationTextFilterFields = []string{
	"title", "code", "doi", "allAuthors", "eliAuthors", "keywords",
	"longJournalTitle", "shortJournalTitle", "abstract", "citeAs",
	"wosNumber", "issn", "eissn", "eidScopus", "oecdFord", "note",
	"otherGrants", "webLink", "publisher", "publishPlace", "isbn",
	"bookTitle", "editionVolume", "proceedingsIsbn", "conferencePlace",
	"pages",
}

// publicationListFilterFields are String properties matched with IN $list
// (yearOfPublication is stored as a String, so the multiselect of years is a
// plain membership test).
var publicationListFilterFields = []string{
	"yearOfPublication", "eliPublication", "quartil", "quartilBasis",
	"language",
}

// publicationNumericRangeFields are Long/Double properties matched with an
// inclusive from/to range.
var publicationNumericRangeFields = []string{
	"impactFactor", "allAuthorsCount", "eliAuthorsCount", "pagesCount",
	"bookPagesCount", "volume", "issue",
}

// publicationDateRangeFields hold ISO-ish date strings (the form keeps them
// free-form, e.g. "2024", "2024-02", "2024-02-13"). The range is a
// lexicographic string comparison, which behaves as a real date range for any
// ISO-formatted value; a malformed value simply never matches. Keeping the
// comparison on the raw string avoids datetime() parse failures on the mixed
// formats the form has historically accepted.
var publicationDateRangeFields = []string{
	"dateOfPublication", "conferenceDate",
}

// publicationCodebookFilters maps the filter id of a codebook-relationship
// column to the relationship type and target label used to match it. Matching
// happens through a pattern comprehension predicate so a filtered query stays
// valid (and cardinality-neutral) in the count query too.
var publicationCodebookFilters = []struct {
	id      string
	relType string
	label   string
	alias   string
}{
	{"mediaTypeCb", "HAS_MEDIA_TYPE", "MediaType", "mediaTypeCb"},
	{"openAccessType", "HAS_OPEN_ACCESS_TYPE", "OpenAccessType", "openAccessType"},
	{"publishingCountry", "HAS_PUBLISHING_COUNTRY", "Country", "publishingCountry"},
	{"userCall", "HAS_USER_CALL", "UserCall", "userCall"},
	{"userExperimentCb", "HAS_USER_EXPERIMENT", "UserExperiment", "userExperimentCb"},
	{"experimentalSystemCb", "HAS_EXPERIMENTAL_SYSTEM", "ExperimentalSystem", "experimentalSystemCb"},
	{"publishFormatCb", "HAS_PUBLISH_FORMAT", "PublishFormat", "publishFormatCb"},
	{"conferenceScopeCb", "HAS_CONFERENCE_SCOPE", "ConferenceScope", "conferenceScopeCb"},
	// Relationship lists (values are node uids, not codebook objects only).
	{"grants", "HAS_GRANT", "Grant", "grant"},
	{"eliResearchers", "HAS_RESEARCHER", "Researcher", "researcher"},
}

// ApplyPublicationFilters appends every recognized ColumnFilter as an AND
// predicate on the query's base WHERE block. Filters combine with AND;
// filters with empty/blank values are ignored.
func ApplyPublicationFilters(query *helpers.DatabaseQuery, filtering *[]helpers.ColumnFilter) {
	if filtering == nil || len(*filtering) == 0 {
		return
	}

	for _, field := range publicationTextFilterFields {
		if value := helpers.GetFilterValueString(filtering, field); value != nil && *value != "" {
			param := filterParamName(field)
			query.Query += fmt.Sprintf(" AND toLower(n.%s) CONTAINS $%s", field, param)
			query.Parameters[param] = strings.ToLower(*value)
		}
	}

	for _, field := range publicationListFilterFields {
		if values := helpers.GetFilterValueListString(filtering, field); values != nil && len(*values) > 0 {
			param := filterParamName(field)
			query.Query += fmt.Sprintf(" AND n.%s IN $%s", field, param)
			query.Parameters[param] = *values
		}
	}

	for _, field := range publicationNumericRangeFields {
		if value := helpers.GetFilterValueRangeFloat64(filtering, field); value != nil {
			param := filterParamName(field)
			query.Query += fmt.Sprintf(
				" AND n.%s IS NOT NULL AND ($%sFrom IS NULL OR n.%s >= $%sFrom) AND ($%sTo IS NULL OR n.%s <= $%sTo)",
				field, param, field, param, param, field, param)
			query.Parameters[param+"From"] = value.Min
			query.Parameters[param+"To"] = value.Max
		}
	}

	for _, field := range publicationDateRangeFields {
		if from, to, ok := filterStringRange(filtering, field); ok {
			param := filterParamName(field)
			query.Query += fmt.Sprintf(
				" AND n.%s IS NOT NULL AND ($%sFrom IS NULL OR n.%s >= $%sFrom) AND ($%sTo IS NULL OR n.%s <= $%sTo)",
				field, param, field, param, param, field, param)
			query.Parameters[param+"From"] = from
			query.Parameters[param+"To"] = to
		}
	}

	for _, codebook := range publicationCodebookFilters {
		uids := filterUidList(filtering, codebook.id)
		if len(uids) == 0 {
			continue
		}
		param := filterParamName(codebook.id)
		// Pattern comprehension keeps cardinality at exactly one row per n.
		query.Query += fmt.Sprintf(
			" AND size([(n)-[:%s]->(%s:%s) WHERE %s.uid IN $%s | 1]) > 0",
			codebook.relType, codebook.alias, codebook.label, codebook.alias, param)
		query.Parameters[param] = uids
	}

	// Departments are denormalized into n.authorsDepartmentsArray as
	// "uid||name||count" entries (no HAS_DEPARTMENT relationship exists yet),
	// so a department filter matches the entry's uid up to its "||" separator
	// (a bare prefix would also match every uid that merely starts with it).
	if deptUids := filterUidList(filtering, "department"); len(deptUids) > 0 {
		query.Query += " AND ANY(x IN coalesce(n.authorsDepartmentsArray, []) WHERE ANY(d IN $filterDepartment WHERE x STARTS WITH d + '||'))"
		query.Parameters["filterDepartment"] = deptUids
	}
}

// filterParamName derives the Cypher parameter name for a filter id,
// e.g. "doi" -> "filterDoi", "mediaTypeCb" -> "filterMediaTypeCb".
func filterParamName(fieldID string) string {
	if fieldID == "" {
		return "filter"
	}
	return "filter" + strings.ToUpper(fieldID[:1]) + fieldID[1:]
}

// filterUidList extracts node uids from a filter value in any of the shapes
// the frontend sends: a single codebook object, a list of codebook objects
// (checkbox groups), a list of plain uid strings (combobox multiselect) or a
// single uid string.
func filterUidList(filtering *[]helpers.ColumnFilter, filterID string) []string {
	if filtering == nil {
		return nil
	}
	for _, f := range *filtering {
		if f.Id != filterID {
			continue
		}
		uids := uidListFromValue(f.Value)
		if len(uids) > 0 {
			return uids
		}
	}
	return nil
}

func uidListFromValue(value any) []string {
	switch typed := value.(type) {
	case map[string]interface{}:
		if uid, ok := typed["uid"].(string); ok && uid != "" {
			return []string{uid}
		}
	case []interface{}:
		var uids []string
		for _, item := range typed {
			uids = append(uids, uidListFromValue(item)...)
		}
		return uids
	case string:
		if typed != "" {
			return []string{typed}
		}
	}
	return nil
}

// filterStringRange reads a {min,max} range whose bounds are strings (dates)
// while tolerating numeric JSON bounds by stringifying them.
func filterStringRange(filtering *[]helpers.ColumnFilter, filterID string) (from any, to any, ok bool) {
	if filtering == nil {
		return nil, nil, false
	}
	for _, f := range *filtering {
		if f.Id != filterID {
			continue
		}
		value, isMap := f.Value.(map[string]interface{})
		if !isMap {
			continue
		}
		from = stringFromJSON(value["min"])
		to = stringFromJSON(value["max"])
		if from == nil && to == nil {
			continue
		}
		return from, to, true
	}
	return nil, nil, false
}

func stringFromJSON(value any) any {
	switch typed := value.(type) {
	case nil:
		return nil
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil
		}
		return strings.TrimSpace(typed)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	}
	return nil
}

// GetPublicationFilterOptionsQuery collects the distinct values and numeric
// bounds behind the publications filter sheet (ELIPANDA-503). Aggregation
// without grouping guarantees exactly one record even on an empty database.
func GetPublicationFilterOptionsQuery() (result helpers.DatabaseQuery) {
	result.Query = `
		MATCH (n:Publication) WHERE (n.deleted IS NULL OR n.deleted = false)
		RETURN {
			years: [x IN collect(DISTINCT n.yearOfPublication) WHERE x IS NOT NULL AND x <> '' | x],
			quartils: [x IN collect(DISTINCT n.quartil) WHERE x IS NOT NULL AND x <> '' | x],
			quartilBases: [x IN collect(DISTINCT n.quartilBasis) WHERE x IS NOT NULL AND x <> '' | x],
			languages: [x IN collect(DISTINCT n.language) WHERE x IS NOT NULL AND x <> '' | x],
			eliPublications: [x IN collect(DISTINCT n.eliPublication) WHERE x IS NOT NULL AND x <> '' | x],
			ranges: {
				impactFactor:    { min: min(n.impactFactor),    max: max(n.impactFactor) },
				allAuthorsCount: { min: min(n.allAuthorsCount), max: max(n.allAuthorsCount) },
				eliAuthorsCount: { min: min(n.eliAuthorsCount), max: max(n.eliAuthorsCount) },
				pagesCount:      { min: min(n.pagesCount),      max: max(n.pagesCount) },
				bookPagesCount:  { min: min(n.bookPagesCount),  max: max(n.bookPagesCount) },
				volume:          { min: min(n.volume),          max: max(n.volume) },
				issue:           { min: min(n.issue),           max: max(n.issue) }
			},
			dateBounds: {
				dateOfPublication: { min: min(n.dateOfPublication), max: max(n.dateOfPublication) },
				conferenceDate:    { min: min(n.conferenceDate),    max: max(n.conferenceDate) }
			}
		} AS result
	`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}
