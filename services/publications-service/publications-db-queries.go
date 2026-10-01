package publicationsservice

import (
	"fmt"
	"sort"
	"strings"

	"panda/apigateway/helpers"
)

// The filter and sort ids below are the column ids the frontend sends from
// publications.columns.tsx. Filters that have no column (book and conference
// fields, department) keep the property name.

// publicationTextFilters map a filter id to the node property matched with a
// case-insensitive CONTAINS.
var publicationTextFilters = map[string]string{
	"title":             "title",
	"code":              "code",
	"doi":               "doi",
	"allAuthors":        "allAuthors",
	"eliAuthors":        "eliAuthors",
	"keywords":          "keywords",
	"longJournalTitle":  "longJournalTitle",
	"shortJournalTitle": "shortJournalTitle",
	"abstract":          "abstract",
	"citeAs":            "citeAs",
	"wosNumber":         "wosNumber",
	"issn":              "issn",
	"eissn":             "eissn",
	"eidScopus":         "eidScopus",
	"oecdFord":          "oecdFord",
	"note":              "note",
	"otherGrants":       "otherGrants",
	"webLink":           "webLink",
	"publisher":         "publisher",
	"publishPlace":      "publishPlace",
	"isbn":              "isbn",
	"bookTitle":         "bookTitle",
	"editionVolume":     "editionVolume",
	"proceedingsIsbn":   "proceedingsIsbn",
	"conferencePlace":   "conferencePlace",
	"pages":             "pages",
}

// publicationListFilters map a filter id to a node property matched with IN.
var publicationListFilters = map[string]string{
	"yearOfPublication": "yearOfPublication",
	"eliPublication":    "eliPublication",
	"quartil":           "quartil",
	"quartilBasis":      "quartilBasis",
	"language":          "language",
}

// publicationNumericRangeFilters map a filter id to a numeric node property.
var publicationNumericRangeFilters = map[string]string{
	"impactFactor":    "impactFactor",
	"allAuthorsCount": "allAuthorsCount",
	"eliAuthorsCount": "eliAuthorsCount",
	"pagesCount":      "pagesCount",
	"bookPagesCount":  "bookPagesCount",
	"volume":          "volume",
	"issue":           "issue",
}

// publicationDateRangeFilters map a filter id to a date-ish property. Both are
// stored as strings, so they are compared lexically — correct for the
// zero-padded YYYY, YYYY-MM and YYYY-MM-DD forms the register holds.
var publicationDateRangeFilters = map[string]string{
	"dateOfPublication": "dateOfPublication",
	"conferenceDate":    "conferenceDate",
}

// publicationCodebookFilters map a filter id to the relationship reaching its
// codebook node.
var publicationCodebookFilters = map[string]string{
	"mediaType":          "HAS_MEDIA_TYPE",
	"openAccessType":     "HAS_OPEN_ACCESS_TYPE",
	"publishingCountry":  "HAS_PUBLISHING_COUNTRY",
	"userCall":           "HAS_USER_CALL",
	"userExperiment":     "HAS_USER_EXPERIMENT",
	"experimentalSystem": "HAS_EXPERIMENTAL_SYSTEM",
	"publishFormat":      "HAS_PUBLISH_FORMAT",
	"conferenceScope":    "HAS_CONFERENCE_SCOPE",
}

// publicationRelationListFilters map a filter id to the relationship and label
// of a to-many link filtered by a list of uids.
var publicationRelationListFilters = map[string]struct {
	relationship string
	label        string
}{
	"grant":          {relationship: "HAS_GRANT", label: "Grant"},
	"eliResearchers": {relationship: "HAS_RESEARCHER", label: "Researcher"},
}

// ApplyPublicationFilters appends one AND-ed predicate per active filter to a
// query that has already opened its WHERE clause. Every predicate is written so
// it holds on the bare (n:Publication) match, which is what lets the list query
// and the count query share it and stay in agreement.
//
// Empty values are skipped by the helpers, so an empty filter never narrows the
// result set. Values are always parameterized; only ids drawn from the maps
// above ever reach the query string.
func ApplyPublicationFilters(query *helpers.DatabaseQuery, filtering *[]helpers.ColumnFilter) {
	if filtering == nil || len(*filtering) == 0 {
		return
	}

	for _, id := range sortedKeys(publicationTextFilters) {
		if value := filterString(filtering, id); value != "" {
			parameter := filterParameterName(id)
			query.Query += fmt.Sprintf(" AND toLower(n.%s) CONTAINS $%s ", publicationTextFilters[id], parameter)
			query.Parameters[parameter] = strings.ToLower(value)
		}
	}

	for _, id := range sortedKeys(publicationListFilters) {
		if values := filterStringList(filtering, id); len(values) > 0 {
			parameter := filterParameterName(id)
			query.Query += fmt.Sprintf(" AND n.%s IN $%s ", publicationListFilters[id], parameter)
			query.Parameters[parameter] = values
		}
	}

	for _, id := range sortedKeys(publicationNumericRangeFilters) {
		applyRange(query, id, publicationNumericRangeFilters[id],
			helpers.GetFilterValueRangeFloat64(filtering, id))
	}

	for _, id := range sortedKeys(publicationDateRangeFilters) {
		applyDateRange(query, id, publicationDateRangeFilters[id],
			getFilterValueRangeString(filtering, id))
	}

	for _, id := range sortedKeys(publicationCodebookFilters) {
		// A checkbox group sends a list of uids, a combobox sends one codebook
		// object. Both mean "the publication points at one of these".
		if uids := filterStringList(filtering, id); len(uids) > 0 {
			parameter := filterParameterName(id)
			query.Query += fmt.Sprintf(
				" AND EXISTS { MATCH (n)-[:%s]->(c) WHERE c.uid IN $%s } ",
				publicationCodebookFilters[id], parameter)
			query.Parameters[parameter] = uids
			continue
		}

		if uid := filterCodebookUID(filtering, id); uid != "" {
			parameter := filterParameterName(id)
			query.Query += fmt.Sprintf(
				" AND EXISTS { MATCH (n)-[:%s]->(c) WHERE c.uid = $%s } ",
				publicationCodebookFilters[id], parameter)
			query.Parameters[parameter] = uid
		}
	}

	for _, id := range sortedRelationKeys(publicationRelationListFilters) {
		if values := filterStringList(filtering, id); len(values) > 0 {
			relation := publicationRelationListFilters[id]
			parameter := filterParameterName(id)
			query.Query += fmt.Sprintf(
				" AND EXISTS { MATCH (n)-[:%s]->(rel:%s) WHERE rel.uid IN $%s } ",
				relation.relationship, relation.label, parameter)
			query.Parameters[parameter] = values
			continue
		}

		// A single picked value arrives as a codebook rather than a list.
		if uid := filterCodebookUID(filtering, id); uid != "" {
			relation := publicationRelationListFilters[id]
			parameter := filterParameterName(id)
			query.Query += fmt.Sprintf(
				" AND EXISTS { MATCH (n)-[:%s]->(rel:%s) WHERE rel.uid = $%s } ",
				relation.relationship, relation.label, parameter)
			query.Parameters[parameter] = uid
		}
	}

	applyDepartmentFilter(query, filtering)
}

// applyDepartmentFilter matches the denormalized authorsDepartmentsArray, whose
// entries are "uid||name||count". A Publication has no HAS_DEPARTMENT
// relationship today; introducing one would let this become an EXISTS match
// like the others.
func applyDepartmentFilter(query *helpers.DatabaseQuery, filtering *[]helpers.ColumnFilter) {
	if uids := filterStringList(filtering, "department"); len(uids) > 0 {
		query.Query += ` AND ANY(entry IN coalesce(n.authorsDepartmentsArray, []) WHERE ANY(uid IN $filterDepartment WHERE entry STARTS WITH uid)) `
		query.Parameters["filterDepartment"] = uids
		return
	}

	uid := filterCodebookUID(filtering, "department")
	if uid == "" {
		uid = filterString(filtering, "department")
	}
	if uid == "" {
		return
	}

	query.Query += ` AND ANY(entry IN coalesce(n.authorsDepartmentsArray, []) WHERE entry STARTS WITH $filterDepartment) `
	query.Parameters["filterDepartment"] = uid
}

func applyRange(query *helpers.DatabaseQuery, id, property string, value *helpers.RangeFloat64Nullable) {
	if value == nil || (value.Min == nil && value.Max == nil) {
		return
	}

	parameter := filterParameterName(id)
	if value.Min != nil {
		query.Query += fmt.Sprintf(" AND n.%s >= $%sMin ", property, parameter)
		query.Parameters[parameter+"Min"] = *value.Min
	}
	if value.Max != nil {
		query.Query += fmt.Sprintf(" AND n.%s <= $%sMax ", property, parameter)
		query.Parameters[parameter+"Max"] = *value.Max
	}
}

func applyDateRange(query *helpers.DatabaseQuery, id, property string, value *rangeStringNullable) {
	if value == nil || (value.Min == nil && value.Max == nil) {
		return
	}

	parameter := filterParameterName(id)
	if value.Min != nil {
		query.Query += fmt.Sprintf(" AND n.%s >= $%sMin ", property, parameter)
		query.Parameters[parameter+"Min"] = *value.Min
	}
	if value.Max != nil {
		query.Query += fmt.Sprintf(" AND n.%s <= $%sMax ", property, parameter)
		query.Parameters[parameter+"Max"] = *value.Max
	}
}

type rangeStringNullable struct {
	Min *string
	Max *string
}

// getFilterValueRangeString reads a {min, max} filter whose bounds are strings.
// helpers has no string-range reader; dates are stored as strings here, and
// parsing them into floats to reuse GetFilterValueRangeFloat64 would lose the
// partial YYYY and YYYY-MM forms the register holds.
func getFilterValueRangeString(filters *[]helpers.ColumnFilter, filterID string) *rangeStringNullable {
	if filters == nil {
		return nil
	}

	for _, filter := range *filters {
		if filter.Id != filterID {
			continue
		}

		value, ok := filter.Value.(map[string]interface{})
		if !ok {
			return nil
		}

		result := rangeStringNullable{Min: boundString(value["min"]), Max: boundString(value["max"])}
		if result.Min == nil && result.Max == nil {
			return nil
		}
		return &result
	}

	return nil
}

// filterCodebookUID reads the uid of a codebook-shaped filter value.
//
// helpers.GetFilterValueCodebook asserts on both uid and name unchecked
// (helpers/database.go:940), so a codebook filter arriving without a name —
// which a hand-written columnFilter query string easily produces — panics the
// request instead of being ignored. Only the uid is needed here, so reading it
// defensively keeps a malformed query string from taking the endpoint down.
func filterCodebookUID(filters *[]helpers.ColumnFilter, filterID string) string {
	if filters == nil {
		return ""
	}

	for _, filter := range *filters {
		if filter.Id != filterID {
			continue
		}

		value, ok := filter.Value.(map[string]interface{})
		if !ok {
			return ""
		}

		uid, ok := value["uid"].(string)
		if !ok {
			return ""
		}
		return strings.TrimSpace(uid)
	}

	return ""
}

// filterStringList reads a list-of-strings filter value.
//
// helpers.GetFilterValueListString asserts the value is []interface{} and every
// element a string, both unchecked (helpers/database.go:948), so a single
// picked value arriving as an object panics the request. Anything that is not a
// usable string is skipped instead.
func filterStringList(filters *[]helpers.ColumnFilter, filterID string) []string {
	if filters == nil {
		return nil
	}

	for _, filter := range *filters {
		if filter.Id != filterID {
			continue
		}

		raw, ok := filter.Value.([]interface{})
		if !ok {
			return nil
		}

		values := make([]string, 0, len(raw))
		for _, item := range raw {
			text, ok := item.(string)
			if !ok || strings.TrimSpace(text) == "" {
				continue
			}
			values = append(values, strings.TrimSpace(text))
		}
		return values
	}

	return nil
}

// filterString reads a plain text filter value, ignoring anything that is not a
// non-empty string.
func filterString(filters *[]helpers.ColumnFilter, filterID string) string {
	if filters == nil {
		return ""
	}

	for _, filter := range *filters {
		if filter.Id != filterID {
			continue
		}

		text, ok := filter.Value.(string)
		if !ok {
			return ""
		}
		return strings.TrimSpace(text)
	}

	return ""
}

func boundString(value interface{}) *string {
	text, ok := value.(string)
	if !ok || strings.TrimSpace(text) == "" {
		return nil
	}
	trimmed := strings.TrimSpace(text)
	return &trimmed
}

// filterParameterName keeps generated parameter names stable and collision-free
// across the filter groups.
func filterParameterName(id string) string {
	return "filter" + strings.ToUpper(id[:1]) + id[1:]
}

// sortedKeys keeps the generated Cypher deterministic, which makes the queries
// readable in logs and the tests assertable.
func sortedKeys(source map[string]string) []string {
	keys := make([]string, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedRelationKeys(source map[string]struct {
	relationship string
	label        string
}) []string {
	keys := make([]string, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
