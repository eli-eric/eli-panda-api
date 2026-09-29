package publicationsservice

import (
	"errors"
	"strings"
	"time"

	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/publications-service/models"

	"github.com/google/uuid"
)

// Reporting classification and document-type vocabularies. Anything outside them
// is rejected rather than stored, so analytics never has to guess what a value meant.
const (
	reportingClassificationOwnUser      = "own-user"
	reportingClassificationOwnOther     = "own-other"
	reportingClassificationCoauthorship = "coauthorship"
	reportingClassificationUnclassified = "unclassified"

	reportingRankingStatusUnknown  = "unknown"
	reportingRankingStatusUnranked = "unranked"
)

var reportingClassifications = map[string]bool{
	reportingClassificationOwnUser:      true,
	reportingClassificationOwnOther:     true,
	reportingClassificationCoauthorship: true,
	reportingClassificationUnclassified: true,
}

var reportingDocumentTypes = map[string]bool{
	"article": true, "proceedings": true, "book-chapter": true, "other": true, "unknown": true,
}

var reportingQuartiles = map[string]bool{"Q1": true, "Q2": true, "Q3": true, "Q4": true}

// reportRelevantChange reports whether a publication edit invalidates a saved
// review. A snapshot records what was true at review time, so once any of these
// change the earlier confirmation no longer describes the record.
func reportRelevantChange(old, updated models.Publication) bool {
	if old.YearOfPublication != updated.YearOfPublication ||
		old.Doi != updated.Doi ||
		old.LongJournalTitle != updated.LongJournalTitle ||
		old.EliPublication != updated.EliPublication ||
		!equalStringPointers(old.DateOfPublication, updated.DateOfPublication) ||
		!equalStringPointers(old.Quartil, updated.Quartil) ||
		!equalStringPointers(old.QuartilBasis, updated.QuartilBasis) ||
		!equalStringPointers(old.Issn, updated.Issn) ||
		!equalStringPointers(old.EIssn, updated.EIssn) {
		return true
	}
	if codebookUID(old.MediaTypeCb) != codebookUID(updated.MediaTypeCb) {
		return true
	}
	if !equalFloatPointers(old.ImpactFactor, updated.ImpactFactor) {
		return true
	}
	if !equalResearcherSets(old.EliResearchers, updated.EliResearchers) {
		return true
	}
	return !equalStringSets(departmentUIDsOf(old), departmentUIDsOf(updated))
}

// validatePublicationReporting rejects values the analytics queries could not
// interpret. Empty enum values are normalized to their explicit "unknown"
// members rather than left blank, so a missing value is never mistaken for a
// meaningful one.
func validatePublicationReporting(reporting *models.PublicationReporting) error {
	if reporting == nil {
		return nil
	}
	if strings.TrimSpace(reporting.Classification) == "" {
		reporting.Classification = reportingClassificationUnclassified
	}
	if !reportingClassifications[reporting.Classification] {
		return errors.New("reporting classification must be own-user, own-other, coauthorship or unclassified")
	}
	if strings.TrimSpace(reporting.DocumentType) == "" {
		reporting.DocumentType = "unknown"
	}
	if !reportingDocumentTypes[reporting.DocumentType] {
		return errors.New("reporting document type is not recognised")
	}
	if reporting.JournalRankingStatus != "" &&
		reporting.JournalRankingStatus != reportingRankingStatusUnknown &&
		reporting.JournalRankingStatus != reportingRankingStatusUnranked {
		return errors.New("journal ranking status must be unknown or unranked")
	}
	for _, metric := range reporting.JournalMetrics {
		if !reportingQuartiles[metric.Quartile] {
			return errors.New("journal metric quartile must be Q1, Q2, Q3 or Q4")
		}
		if metric.Year < 1900 || metric.Year > 9999 {
			return errors.New("journal metric year is out of range")
		}
		if strings.TrimSpace(metric.Category) == "" {
			return errors.New("journal metric category is required")
		}
		if metric.Percentile != nil && (*metric.Percentile < 0 || *metric.Percentile > 100) {
			return errors.New("journal metric percentile must be between 0 and 100")
		}
	}
	for _, author := range reporting.Authors {
		if strings.TrimSpace(author.ResearcherUID) == "" {
			return errors.New("reporting author requires a researcher")
		}
	}
	for _, departmentUID := range reporting.DepartmentUIDs {
		if strings.TrimSpace(departmentUID) == "" {
			return errors.New("reporting department requires a department")
		}
	}
	return nil
}

// getPublicationReporting returns nil when the publication has never been
// reviewed. Absence is a real state and must not be filled in with defaults.
func (svc *PublicationsService) getPublicationReporting(uid string) (*models.PublicationReporting, error) {
	session, err := helpers.NewNeo4jSession(*svc.neo4jDriver)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	query := helpers.DatabaseQuery{
		Query: `
			MATCH (p:Publication {uid: $uid})-[:HAS_REPORTING]->(reporting:PublicationReporting)
			RETURN {
				classification: coalesce(reporting.classification, "unclassified"),
				reviewed: coalesce(reporting.reviewed, false),
				documentType: coalesce(reporting.documentType, "unknown"),
				journalRankingStatus: reporting.journalRankingStatus,
				reviewedAt: reporting.reviewedAt,
				reviewedBy: reporting.reviewedBy,
				departmentUids: [(reporting)-[:REPORTED_DEPARTMENT]->(department:Department) | department.uid],
				userCallUids: [(reporting)-[:REPORTED_USER_CALL]->(call:UserCall) | call.uid],
				userExperimentUids: [(reporting)-[:REPORTED_USER_EXPERIMENT]->(experiment:UserExperiment) | experiment.uid],
				experimentalSystemUids: [(reporting)-[:REPORTED_EXPERIMENTAL_SYSTEM]->(system:ExperimentalSystem) | system.uid],
				authors: [(reporting)-[:HAS_REPORTING_AUTHOR]->(author:PublicationReportingAuthor) | {
					researcherUid: head([(author)-[:IS_RESEARCHER]->(researcher:Researcher) | researcher.uid]),
					departmentUids: [(author)-[:REPORTED_DEPARTMENT]->(authorDepartment:Department) | authorDepartment.uid],
					isFirstAuthor: author.isFirstAuthor,
					isCorresponding: author.isCorresponding
				}],
				journalMetrics: [(reporting)-[:HAS_JOURNAL_METRIC]->(metric:JournalMetric) | {
					source: coalesce(metric.source, "JCR"),
					journalId: metric.journalId,
					year: metric.year,
					category: metric.category,
					quartile: metric.quartile,
					percentile: metric.percentile,
					impactFactor: metric.impactFactor
				}]
			} AS reporting`,
		ReturnAlias: "reporting",
		Parameters:  map[string]interface{}{"uid": uid},
	}

	reporting, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.PublicationReporting](session, query)
	if errors.Is(err, helpers.ERR_NO_ROWS) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	normalizeReportingCollections(&reporting)
	return &reporting, nil
}

// publicationReportingQueries builds the whole snapshot write. Submitting
// reporting replaces it outright, so the previous subgraph is removed first;
// returning one slice lets the caller run all of it in a single transaction
// alongside the publication's own update.
func publicationReportingQueries(
	publicationUID string,
	reporting *models.PublicationReporting,
	userUID string,
	now time.Time,
) []helpers.DatabaseQuery {
	queries := []helpers.DatabaseQuery{deletePublicationReportingQuery(publicationUID)}

	reportingUID := uuid.NewString()
	// Reviewer identity and time are authored here rather than taken from the
	// request, so a client cannot claim someone else reviewed a record.
	reviewedAt, reviewedBy := interface{}(nil), interface{}(nil)
	if reporting.Reviewed {
		reviewedAt, reviewedBy = now.UTC().Format(time.RFC3339), userUID
	}

	queries = append(queries, helpers.DatabaseQuery{
		Query: `
			MATCH (p:Publication {uid: $publicationUid})
			CREATE (p)-[:HAS_REPORTING]->(reporting:PublicationReporting {
				uid: $uid,
				classification: $classification,
				reviewed: $reviewed,
				documentType: $documentType,
				journalRankingStatus: $journalRankingStatus,
				reviewedAt: $reviewedAt,
				reviewedBy: $reviewedBy
			})`,
		Parameters: map[string]interface{}{
			"publicationUid":       publicationUID,
			"uid":                  reportingUID,
			"classification":       reporting.Classification,
			"reviewed":             reporting.Reviewed,
			"documentType":         reporting.DocumentType,
			"journalRankingStatus": nullableString(reporting.JournalRankingStatus),
			"reviewedAt":           reviewedAt,
			"reviewedBy":           reviewedBy,
		},
	})

	for _, departmentUID := range reporting.DepartmentUIDs {
		queries = append(queries, linkReportingQuery(reportingUID, "Department", "REPORTED_DEPARTMENT", departmentUID))
	}
	for _, callUID := range reporting.UserCallUIDs {
		queries = append(queries, linkReportingQuery(reportingUID, "UserCall", "REPORTED_USER_CALL", callUID))
	}
	for _, experimentUID := range reporting.UserExperimentUIDs {
		queries = append(queries, linkReportingQuery(reportingUID, "UserExperiment", "REPORTED_USER_EXPERIMENT", experimentUID))
	}
	for _, systemUID := range reporting.ExperimentalSystemUIDs {
		queries = append(queries, linkReportingQuery(reportingUID, "ExperimentalSystem", "REPORTED_EXPERIMENTAL_SYSTEM", systemUID))
	}

	for _, author := range reporting.Authors {
		authorUID := uuid.NewString()
		queries = append(queries, helpers.DatabaseQuery{
			Query: `
				MATCH (reporting:PublicationReporting {uid: $reportingUid})
				MATCH (researcher:Researcher {uid: $researcherUid})
				WHERE researcher.deleted IS NULL OR researcher.deleted = false
				CREATE (reporting)-[:HAS_REPORTING_AUTHOR]->(author:PublicationReportingAuthor {
					uid: $uid,
					isFirstAuthor: $isFirstAuthor,
					isCorresponding: $isCorresponding
				})
				CREATE (author)-[:IS_RESEARCHER]->(researcher)`,
			Parameters: map[string]interface{}{
				"reportingUid":    reportingUID,
				"researcherUid":   author.ResearcherUID,
				"uid":             authorUID,
				"isFirstAuthor":   nullableBool(author.IsFirstAuthor),
				"isCorresponding": nullableBool(author.IsCorresponding),
			},
		})
		for _, departmentUID := range author.DepartmentUIDs {
			queries = append(queries, helpers.DatabaseQuery{
				Query: `
					MATCH (author:PublicationReportingAuthor {uid: $authorUid})
					MATCH (department:Department {uid: $departmentUid})
					MERGE (author)-[:REPORTED_DEPARTMENT]->(department)`,
				Parameters: map[string]interface{}{"authorUid": authorUID, "departmentUid": departmentUID},
			})
		}
	}

	for _, metric := range reporting.JournalMetrics {
		source := strings.TrimSpace(metric.Source)
		if source == "" {
			source = "JCR"
		}
		queries = append(queries, helpers.DatabaseQuery{
			Query: `
				MATCH (reporting:PublicationReporting {uid: $reportingUid})
				CREATE (reporting)-[:HAS_JOURNAL_METRIC]->(:JournalMetric {
					uid: $uid,
					source: $source,
					journalId: $journalId,
					year: $year,
					category: $category,
					quartile: $quartile,
					percentile: $percentile,
					impactFactor: $impactFactor
				})`,
			Parameters: map[string]interface{}{
				"reportingUid": reportingUID,
				"uid":          uuid.NewString(),
				"source":       source,
				"journalId":    nullableString(metric.JournalID),
				"year":         metric.Year,
				"category":     strings.TrimSpace(metric.Category),
				"quartile":     metric.Quartile,
				"percentile":   nullableFloat(metric.Percentile),
				"impactFactor": nullableFloat(metric.ImpactFactor),
			},
		})
	}

	return queries
}

func deletePublicationReportingQuery(publicationUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `
			MATCH (p:Publication {uid: $publicationUid})-[:HAS_REPORTING]->(reporting:PublicationReporting)
			OPTIONAL MATCH (reporting)-[:HAS_REPORTING_AUTHOR]->(author:PublicationReportingAuthor)
			OPTIONAL MATCH (reporting)-[:HAS_JOURNAL_METRIC]->(metric:JournalMetric)
			DETACH DELETE reporting, author, metric`,
		Parameters: map[string]interface{}{"publicationUid": publicationUID},
	}
}

// invalidatePublicationReviewQuery clears the review flag without touching the
// rest of the snapshot, so the editor's earlier selections survive and only the
// confirmation has to be repeated.
func invalidatePublicationReviewQuery(publicationUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `
			MATCH (p:Publication {uid: $publicationUid})-[:HAS_REPORTING]->(reporting:PublicationReporting)
			SET reporting.reviewed = false, reporting.reviewedAt = null, reporting.reviewedBy = null`,
		Parameters: map[string]interface{}{"publicationUid": publicationUID},
	}
}

func linkReportingQuery(reportingUID, label, relationship, targetUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `
			MATCH (reporting:PublicationReporting {uid: $reportingUid})
			MATCH (target:` + label + ` {uid: $targetUid})
			MERGE (reporting)-[:` + relationship + `]->(target)`,
		Parameters: map[string]interface{}{"reportingUid": reportingUID, "targetUid": targetUID},
	}
}

// normalizeReportingCollections replaces nil slices with empty ones so the JSON
// response always carries arrays, and drops author rows whose researcher has
// since been deleted rather than emitting an author with no identity.
func normalizeReportingCollections(reporting *models.PublicationReporting) {
	if reporting.DepartmentUIDs == nil {
		reporting.DepartmentUIDs = []string{}
	}
	if reporting.UserCallUIDs == nil {
		reporting.UserCallUIDs = []string{}
	}
	if reporting.UserExperimentUIDs == nil {
		reporting.UserExperimentUIDs = []string{}
	}
	if reporting.ExperimentalSystemUIDs == nil {
		reporting.ExperimentalSystemUIDs = []string{}
	}
	if reporting.JournalMetrics == nil {
		reporting.JournalMetrics = []models.JournalMetricSnapshot{}
	}
	authors := make([]models.ReportingAuthor, 0, len(reporting.Authors))
	for _, author := range reporting.Authors {
		if strings.TrimSpace(author.ResearcherUID) == "" {
			continue
		}
		if author.DepartmentUIDs == nil {
			author.DepartmentUIDs = []string{}
		}
		authors = append(authors, author)
	}
	reporting.Authors = authors
}

func nullableString(value string) interface{} {
	if trimmed := strings.TrimSpace(value); trimmed != "" {
		return trimmed
	}
	return nil
}

func nullableBool(value *bool) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

func nullableFloat(value *float64) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

func codebookUID(value *codebookModels.Codebook) string {
	if value == nil {
		return ""
	}
	return value.UID
}

func departmentUIDsOf(publication models.Publication) []string {
	result := make([]string, 0, len(publication.AuthorsDepartments))
	for _, entry := range publication.AuthorsDepartments {
		if entry.Department.UID != "" {
			result = append(result, entry.Department.UID)
		}
	}
	return result
}

func equalStringPointers(left, right *string) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func equalFloatPointers(left, right *float64) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func equalResearcherSets(left, right []models.ResearcherRef) bool {
	leftUIDs := make([]string, 0, len(left))
	for _, item := range left {
		leftUIDs = append(leftUIDs, item.Uid)
	}
	rightUIDs := make([]string, 0, len(right))
	for _, item := range right {
		rightUIDs = append(rightUIDs, item.Uid)
	}
	return equalStringSets(leftUIDs, rightUIDs)
}

func equalStringSets(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	counts := make(map[string]int, len(left))
	for _, value := range left {
		counts[value]++
	}
	for _, value := range right {
		counts[value]--
		if counts[value] < 0 {
			return false
		}
	}
	return true
}

// listReportingRows returns one flat row per publication in the window,
// including publications that have never been reviewed: they are part of the
// institutional total and of the unclassified backlog, and dropping them would
// quietly shrink every headline number.
func (svc *PublicationsService) listReportingRows(startYear, endYear int) ([]reportingRow, error) {
	session, err := helpers.NewNeo4jSession(*svc.neo4jDriver)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	query := helpers.DatabaseQuery{
		Query: `
			MATCH (p:Publication)
			WHERE (p.deleted IS NULL OR p.deleted = false)
			  AND p.yearOfPublication IS NOT NULL
			  AND toInteger(p.yearOfPublication) >= $startYear
			  AND toInteger(p.yearOfPublication) <= $endYear
			OPTIONAL MATCH (p)-[:HAS_REPORTING]->(reporting:PublicationReporting)
			RETURN {
				uid: p.uid,
				year: toInteger(p.yearOfPublication),
				journalTitle: coalesce(p.longJournalTitle, ""),
				classification: coalesce(reporting.classification, "unclassified"),
				reviewed: coalesce(reporting.reviewed, false),
				documentType: coalesce(reporting.documentType, "unknown"),
				journalRankingStatus: coalesce(reporting.journalRankingStatus, ""),
				departments: [(reporting)-[:REPORTED_DEPARTMENT]->(department:Department) |
					{uid: department.uid, name: coalesce(department.name, department.uid)}],
				userCalls: [(reporting)-[:REPORTED_USER_CALL]->(call:UserCall) |
					{uid: call.uid, name: coalesce(call.name, call.uid)}],
				experimentalSystems: [(reporting)-[:REPORTED_EXPERIMENTAL_SYSTEM]->(system:ExperimentalSystem) |
					{uid: system.uid, name: coalesce(system.name, system.uid)}],
				authors: [(reporting)-[:HAS_REPORTING_AUTHOR]->(author:PublicationReportingAuthor) | {
					researcherUid: head([(author)-[:IS_RESEARCHER]->(researcher:Researcher) | researcher.uid]),
					name: head([(author)-[:IS_RESEARCHER]->(researcher:Researcher) |
						trim(coalesce(researcher.lastName, "") + " " + coalesce(researcher.firstName, ""))]),
					departmentUids: [(author)-[:REPORTED_DEPARTMENT]->(department:Department) | department.uid],
					isFirstAuthor: author.isFirstAuthor,
					isCorresponding: author.isCorresponding
				}],
				journalMetrics: [(reporting)-[:HAS_JOURNAL_METRIC]->(metric:JournalMetric) | {
					source: coalesce(metric.source, "JCR"),
					journalId: metric.journalId,
					year: metric.year,
					category: metric.category,
					quartile: metric.quartile,
					percentile: metric.percentile,
					impactFactor: metric.impactFactor
				}]
			} AS row`,
		ReturnAlias: "row",
		Parameters:  map[string]interface{}{"startYear": startYear, "endYear": endYear},
	}

	rows, err := helpers.GetNeo4jArrayOfNodes[reportingRow](session, query)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []reportingRow{}
	}
	return rows, nil
}
