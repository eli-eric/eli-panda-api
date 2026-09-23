package publicationsservice

import (
	"testing"
	"time"

	"panda/apigateway/services/publications-service/models"
	"panda/apigateway/services/testsetup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The shared test database is never cleaned automatically, so every fixture
// here is namespaced and removed by the test that created it.
const reportingTestPrefix = "reporting-test-"

func setupReportingFixtures(t *testing.T) {
	t.Helper()
	_, err := testsetup.TestSession.Run(`
		CREATE (:Department {uid: $d1, name: "Test Department 86"})
		CREATE (:Department {uid: $d2, name: "Test Department 88"})
		CREATE (:Researcher {uid: $r1, firstName: "Jana", lastName: "Nováková"})
		CREATE (:UserCall {uid: $c1, name: "Test Call 1"})
		CREATE (:ExperimentalSystem {uid: $s1, name: "Test ELIMAIA", code: "ELIMAIA"})
	`, map[string]interface{}{
		"d1": reportingTestPrefix + "dept-86",
		"d2": reportingTestPrefix + "dept-88",
		"r1": reportingTestPrefix + "researcher-1",
		"c1": reportingTestPrefix + "call-1",
		"s1": reportingTestPrefix + "system-1",
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err := testsetup.TestSession.Run(`
			MATCH (n)
			WHERE n.uid STARTS WITH $prefix
			OPTIONAL MATCH (n)-[:HAS_REPORTING]->(reporting:PublicationReporting)
			OPTIONAL MATCH (reporting)-[:HAS_REPORTING_AUTHOR]->(author:PublicationReportingAuthor)
			OPTIONAL MATCH (reporting)-[:HAS_JOURNAL_METRIC]->(metric:JournalMetric)
			DETACH DELETE n, reporting, author, metric
		`, map[string]interface{}{"prefix": reportingTestPrefix})
		require.NoError(t, err)
	})
}

func sampleReporting() *models.PublicationReporting {
	yes, no := true, false
	percentile := 92.5
	impact := 8.4
	return &models.PublicationReporting{
		Classification: reportingClassificationOwnUser,
		Reviewed:       true,
		DocumentType:   "article",
		DepartmentUIDs: []string{
			reportingTestPrefix + "dept-86",
			reportingTestPrefix + "dept-88",
		},
		Authors: []models.ReportingAuthor{{
			ResearcherUID:   reportingTestPrefix + "researcher-1",
			DepartmentUIDs:  []string{reportingTestPrefix + "dept-86"},
			IsFirstAuthor:   &yes,
			IsCorresponding: &no,
		}},
		UserCallUIDs:           []string{reportingTestPrefix + "call-1"},
		UserExperimentUIDs:     []string{},
		ExperimentalSystemUIDs: []string{reportingTestPrefix + "system-1"},
		JournalMetrics: []models.JournalMetricSnapshot{{
			Source: "JCR", JournalID: "0031-9007", Year: 2025,
			Category: "Physics, Multidisciplinary", Quartile: "Q1",
			Percentile: &percentile, ImpactFactor: &impact,
		}},
	}
}

func newReportingPublication(uid string) *models.Publication {
	return &models.Publication{
		Uid:               uid,
		Code:              "TEST-" + uid,
		Title:             "Reporting round-trip",
		YearOfPublication: "2025",
		LongJournalTitle:  "Physical Review Letters",
		EliPublication:    eliPublicationYes,
	}
}

func TestPublicationReportingRoundTrip(t *testing.T) {
	setupReportingFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	publication := newReportingPublication(reportingTestPrefix + "pub-roundtrip")
	publication.Reporting = sampleReporting()
	_, err := service.CreatePublication(publication, "test-user")
	require.NoError(t, err)

	stored, err := service.GetPublicationByUid(publication.Uid)
	require.NoError(t, err)
	require.NotNil(t, stored.Reporting)

	assert.Equal(t, reportingClassificationOwnUser, stored.Reporting.Classification)
	assert.True(t, stored.Reporting.Reviewed)
	assert.Equal(t, "article", stored.Reporting.DocumentType)

	// Both credited departments survive: the paper counts once in each.
	assert.ElementsMatch(t,
		[]string{reportingTestPrefix + "dept-86", reportingTestPrefix + "dept-88"},
		stored.Reporting.DepartmentUIDs)

	assert.Equal(t, []string{reportingTestPrefix + "call-1"}, stored.Reporting.UserCallUIDs)
	assert.Equal(t, []string{reportingTestPrefix + "system-1"}, stored.Reporting.ExperimentalSystemUIDs)
	assert.Empty(t, stored.Reporting.UserExperimentUIDs)

	require.Len(t, stored.Reporting.Authors, 1)
	author := stored.Reporting.Authors[0]
	assert.Equal(t, reportingTestPrefix+"researcher-1", author.ResearcherUID)
	assert.Equal(t, []string{reportingTestPrefix + "dept-86"}, author.DepartmentUIDs)
	require.NotNil(t, author.IsFirstAuthor)
	assert.True(t, *author.IsFirstAuthor)
	require.NotNil(t, author.IsCorresponding)
	assert.False(t, *author.IsCorresponding)

	require.Len(t, stored.Reporting.JournalMetrics, 1)
	metric := stored.Reporting.JournalMetrics[0]
	assert.Equal(t, 2025, metric.Year)
	assert.Equal(t, "Q1", metric.Quartile)
	require.NotNil(t, metric.Percentile)
	assert.InDelta(t, 92.5, *metric.Percentile, 0.001)

	// Reviewer identity is authored by the server, not accepted from the client.
	assert.Equal(t, "test-user", stored.Reporting.ReviewedBy)
	assert.NotEmpty(t, stored.Reporting.ReviewedAt)
}

func TestPublicationReportingUnknownAuthorRolesStayUnknown(t *testing.T) {
	setupReportingFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	publication := newReportingPublication(reportingTestPrefix + "pub-unknown-roles")
	publication.Reporting = &models.PublicationReporting{
		Classification: reportingClassificationUnclassified,
		DocumentType:   "unknown",
		Authors: []models.ReportingAuthor{{
			ResearcherUID: reportingTestPrefix + "researcher-1",
		}},
	}
	_, err := service.CreatePublication(publication, "test-user")
	require.NoError(t, err)

	stored, err := service.GetPublicationByUid(publication.Uid)
	require.NoError(t, err)
	require.NotNil(t, stored.Reporting)
	require.Len(t, stored.Reporting.Authors, 1)

	// An unstated role must come back unstated rather than as a false "no".
	assert.Nil(t, stored.Reporting.Authors[0].IsFirstAuthor)
	assert.Nil(t, stored.Reporting.Authors[0].IsCorresponding)

	// An unreviewed snapshot carries no reviewer.
	assert.False(t, stored.Reporting.Reviewed)
	assert.Empty(t, stored.Reporting.ReviewedBy)
	assert.Empty(t, stored.Reporting.ReviewedAt)
}

func TestUpdateWithoutReportingPreservesTheSnapshot(t *testing.T) {
	setupReportingFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	publication := newReportingPublication(reportingTestPrefix + "pub-preserve")
	publication.Reporting = sampleReporting()
	_, err := service.CreatePublication(publication, "test-user")
	require.NoError(t, err)

	// A client that knows nothing about reporting saves an unrelated edit.
	legacy := newReportingPublication(publication.Uid)
	legacy.Note = stringValue("Edited by a client that does not send reporting")
	_, err = service.UpdatePublication(legacy, "test-user")
	require.NoError(t, err)

	stored, err := service.GetPublicationByUid(publication.Uid)
	require.NoError(t, err)
	require.NotNil(t, stored.Reporting, "omitting reporting must not erase it")
	assert.Equal(t, reportingClassificationOwnUser, stored.Reporting.Classification)
	assert.Len(t, stored.Reporting.DepartmentUIDs, 2)
	assert.Len(t, stored.Reporting.JournalMetrics, 1)
	assert.True(t, stored.Reporting.Reviewed, "an unrelated edit keeps the review")
}

func TestUpdateChangingReportedDataInvalidatesTheReview(t *testing.T) {
	setupReportingFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	publication := newReportingPublication(reportingTestPrefix + "pub-invalidate")
	publication.Reporting = sampleReporting()
	_, err := service.CreatePublication(publication, "test-user")
	require.NoError(t, err)

	// The publication year is part of what was reviewed, so changing it without
	// resubmitting reporting must retract the confirmation.
	changed := newReportingPublication(publication.Uid)
	changed.YearOfPublication = "2024"
	_, err = service.UpdatePublication(changed, "test-user")
	require.NoError(t, err)

	stored, err := service.GetPublicationByUid(publication.Uid)
	require.NoError(t, err)
	require.NotNil(t, stored.Reporting)
	assert.False(t, stored.Reporting.Reviewed, "a changed publication year retracts the review")
	assert.Empty(t, stored.Reporting.ReviewedBy)
	// The editor's selections survive; only the confirmation is withdrawn.
	assert.Len(t, stored.Reporting.DepartmentUIDs, 2)
	assert.Equal(t, reportingClassificationOwnUser, stored.Reporting.Classification)
}

func TestUpdateWithExplicitReportingReplacesTheSnapshot(t *testing.T) {
	setupReportingFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	publication := newReportingPublication(reportingTestPrefix + "pub-replace")
	publication.Reporting = sampleReporting()
	_, err := service.CreatePublication(publication, "test-user")
	require.NoError(t, err)

	replacement := newReportingPublication(publication.Uid)
	replacement.Reporting = &models.PublicationReporting{
		Classification:       reportingClassificationCoauthorship,
		DocumentType:         "proceedings",
		JournalRankingStatus: reportingRankingStatusUnranked,
		DepartmentUIDs:       []string{reportingTestPrefix + "dept-88"},
	}
	_, err = service.UpdatePublication(replacement, "test-user")
	require.NoError(t, err)

	stored, err := service.GetPublicationByUid(publication.Uid)
	require.NoError(t, err)
	require.NotNil(t, stored.Reporting)
	assert.Equal(t, reportingClassificationCoauthorship, stored.Reporting.Classification)
	assert.Equal(t, "proceedings", stored.Reporting.DocumentType)
	assert.Equal(t, reportingRankingStatusUnranked, stored.Reporting.JournalRankingStatus)

	// Submitting reporting replaces the whole snapshot: the dropped department,
	// author, call, system and metric are gone rather than merged.
	require.Len(t, stored.Reporting.DepartmentUIDs, 1)
	assert.Equal(t, reportingTestPrefix+"dept-88", stored.Reporting.DepartmentUIDs[0])
	assert.Empty(t, stored.Reporting.Authors)
	assert.Empty(t, stored.Reporting.UserCallUIDs)
	assert.Empty(t, stored.Reporting.ExperimentalSystemUIDs)
	assert.Empty(t, stored.Reporting.JournalMetrics)

	// Nothing is orphaned behind the replaced snapshot.
	result, err := testsetup.TestSession.Run(`
		MATCH (n) WHERE n:PublicationReporting OR n:PublicationReportingAuthor OR n:JournalMetric
		WITH n WHERE NOT (n)<-[:HAS_REPORTING|HAS_REPORTING_AUTHOR|HAS_JOURNAL_METRIC]-()
		RETURN count(n) AS orphans`, nil)
	require.NoError(t, err)
	record, err := result.Single()
	require.NoError(t, err)
	assert.Equal(t, int64(0), record.Values[0])
}

func TestReportingRejectsValuesAnalyticsCannotInterpret(t *testing.T) {
	percentile := 140.0
	for name, reporting := range map[string]*models.PublicationReporting{
		"unknown classification": {Classification: "own", DocumentType: "article"},
		"unknown document type":  {Classification: reportingClassificationOwnUser, DocumentType: "preprint"},
		"unknown quartile": {Classification: reportingClassificationOwnUser, DocumentType: "article",
			JournalMetrics: []models.JournalMetricSnapshot{{Year: 2025, Category: "Physics", Quartile: "Q5"}}},
		"percentile out of range": {Classification: reportingClassificationOwnUser, DocumentType: "article",
			JournalMetrics: []models.JournalMetricSnapshot{{Year: 2025, Category: "Physics", Quartile: "Q1", Percentile: &percentile}}},
		"metric without category": {Classification: reportingClassificationOwnUser, DocumentType: "article",
			JournalMetrics: []models.JournalMetricSnapshot{{Year: 2025, Quartile: "Q1"}}},
		"author without researcher": {Classification: reportingClassificationOwnUser, DocumentType: "article",
			Authors: []models.ReportingAuthor{{}}},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Error(t, validatePublicationReporting(reporting))
		})
	}

	// A blank classification is an explicit "unclassified", not a rejection.
	blank := &models.PublicationReporting{}
	require.NoError(t, validatePublicationReporting(blank))
	assert.Equal(t, reportingClassificationUnclassified, blank.Classification)
	assert.Equal(t, "unknown", blank.DocumentType)
}

func TestExecutiveSummaryReadsSavedRecords(t *testing.T) {
	setupReportingFixtures(t)
	service := NewPublicationsService(&testsetup.TestDriver, "", "")

	// A reviewed user paper credited to two departments, in a top-10% journal.
	reviewed := newReportingPublication(reportingTestPrefix + "pub-summary-user")
	reviewed.Reporting = sampleReporting()
	_, err := service.CreatePublication(reviewed, "test-user")
	require.NoError(t, err)

	// A co-authored Q3 paper in the same year.
	coauthored := newReportingPublication(reportingTestPrefix + "pub-summary-coauthor")
	coauthored.LongJournalTitle = "Optics Express"
	coauthored.Reporting = &models.PublicationReporting{
		Classification: reportingClassificationCoauthorship,
		Reviewed:       true,
		DocumentType:   "article",
		DepartmentUIDs: []string{reportingTestPrefix + "dept-88"},
		JournalMetrics: []models.JournalMetricSnapshot{
			{Source: "JCR", Year: 2025, Category: "Optics", Quartile: "Q3"},
		},
	}
	_, err = service.CreatePublication(coauthored, "test-user")
	require.NoError(t, err)

	// A publication nobody has reviewed yet: part of the total and the backlog.
	unreviewed := newReportingPublication(reportingTestPrefix + "pub-summary-unreviewed")
	_, err = service.CreatePublication(unreviewed, "test-user")
	require.NoError(t, err)

	rows, err := service.listReportingRows(2025, 2025)
	require.NoError(t, err)

	// Other tests and seed data may share the database, so assert on our rows.
	ours := make([]reportingRow, 0, 3)
	for _, row := range rows {
		if len(row.UID) >= len(reportingTestPrefix) && row.UID[:len(reportingTestPrefix)] == reportingTestPrefix {
			ours = append(ours, row)
		}
	}
	require.Len(t, ours, 3)

	summary := buildExecutiveSummary(ours, 2025, 2025, 2025, time.Now())

	assert.Equal(t, 3, summary.TotalPublications)
	assert.Equal(t, 1, summary.TotalOwnPublications)
	assert.Equal(t, 1, summary.TotalUserPublications)
	assert.Equal(t, 1, summary.CoauthorshipPublications)
	assert.Equal(t, 1, summary.UnclassifiedPublications)
	assert.Equal(t, 1, summary.PendingReviewPublications)

	// D88 is credited by both reviewed papers, D86 only by the user paper, so
	// department credits (3) exceed the distinct total that carries them (2).
	byName := map[string]int{}
	for _, department := range summary.DepartmentMatrix {
		byName[department.Name] = department.Total
	}
	assert.Equal(t, 1, byName["Test Department 86"])
	assert.Equal(t, 2, byName["Test Department 88"])

	require.Len(t, summary.Q3Q4HistoricalTrend, 1)
	// The user cohort is the single top-10% paper: ranked, but not Q3/Q4.
	user := summary.Q3Q4HistoricalTrend[0].User
	assert.Equal(t, 1, user.RankedCount)
	assert.Equal(t, 0, user.Q3Q4Count)
	require.NotNil(t, user.Percent)
	assert.InDelta(t, 0, *user.Percent, 0.001)

	// The reviewed author round-trips with their confirmed roles intact.
	require.Len(t, summary.TopPublishingAuthors, 1)
	assert.Equal(t, reportingTestPrefix+"researcher-1", summary.TopPublishingAuthors[0].ResearcherUID)
	assert.Equal(t, "Nováková Jana", summary.TopPublishingAuthors[0].Name)
	assert.Equal(t, 1, summary.TopPublishingAuthors[0].FirstAuthorCount)

	journals := map[string]int{}
	for _, journal := range summary.JournalFrequencies {
		journals[journal.Name] = journal.Count
	}
	assert.Equal(t, 2, journals["Physical Review Letters"])
	assert.Equal(t, 1, journals["Optics Express"])
}
