package publicationsservice

import (
	"testing"
	"time"

	"panda/apigateway/services/publications-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const analyticsYear = 2025

func percentileOf(value float64) *float64 { return &value }

func rankedRow(uid, classification, quartile string, percentile *float64) reportingRow {
	return reportingRow{
		UID: uid, Year: analyticsYear, Classification: classification,
		Reviewed: true, DocumentType: "article",
		JournalMetrics: []models.JournalMetricSnapshot{
			{Source: "JCR", Year: analyticsYear, Category: "Physics", Quartile: quartile, Percentile: percentile},
		},
	}
}

func TestQualityBandSplitsQ1AtTheNinetiethPercentile(t *testing.T) {
	// The boundary is inclusive: exactly 90 is top-10%.
	assert.Equal(t, bandQ10Percent, rankedRow("a", reportingClassificationOwnUser, "Q1", percentileOf(90)).qualityBand())
	assert.Equal(t, bandQ10To25, rankedRow("b", reportingClassificationOwnUser, "Q1", percentileOf(89.99)).qualityBand())
	assert.Equal(t, bandQ10Percent, rankedRow("c", reportingClassificationOwnUser, "Q1", percentileOf(99.1)).qualityBand())

	// Without a percentile the split cannot be made, so Q1 stays unsplit rather
	// than being guessed into either half.
	assert.Equal(t, bandQ1Unsplit, rankedRow("d", reportingClassificationOwnUser, "Q1", nil).qualityBand())
}

func TestQualityBandKeepsMissingEvidenceDistinctFromUnranked(t *testing.T) {
	missing := reportingRow{UID: "a", Year: analyticsYear, DocumentType: "article"}
	assert.Equal(t, bandUnknown, missing.qualityBand())

	unranked := reportingRow{UID: "b", Year: analyticsYear, DocumentType: "article", JournalRankingStatus: reportingRankingStatusUnranked}
	assert.Equal(t, bandUnranked, unranked.qualityBand())

	// A metric from another year is not evidence about this paper.
	staleOnly := reportingRow{
		UID: "c", Year: analyticsYear, DocumentType: "article",
		JournalMetrics: []models.JournalMetricSnapshot{{Year: analyticsYear - 1, Quartile: "Q1", Percentile: percentileOf(95)}},
	}
	assert.Equal(t, bandUnknown, staleOnly.qualityBand())

	proceedings := reportingRow{UID: "d", Year: analyticsYear, DocumentType: "proceedings"}
	assert.Equal(t, bandProceedings, proceedings.qualityBand())
	chapter := reportingRow{UID: "e", Year: analyticsYear, DocumentType: "book-chapter"}
	assert.Equal(t, bandBookChapters, chapter.qualityBand())
}

func TestQualityBandUsesTheHighestPercentileForThePublicationYear(t *testing.T) {
	row := reportingRow{
		UID: "a", Year: analyticsYear, DocumentType: "article",
		JournalMetrics: []models.JournalMetricSnapshot{
			{Year: analyticsYear, Category: "Optics", Quartile: "Q2", Percentile: percentileOf(70)},
			{Year: analyticsYear, Category: "Physics", Quartile: "Q1", Percentile: percentileOf(94)},
			{Year: analyticsYear, Category: "Applied", Quartile: "Q1", Percentile: nil},
			{Year: analyticsYear + 1, Category: "Physics", Quartile: "Q4", Percentile: percentileOf(10)},
		},
	}
	assert.Equal(t, bandQ10Percent, row.qualityBand())
}

// The agreed acceptance fixture: 78 ranked own papers of which 22 are Q3 or Q4.
func TestQ3Q4TrendMatchesTheAgreedFixture(t *testing.T) {
	rows := make([]reportingRow, 0, 90)
	for i := 0; i < 22; i++ {
		rows = append(rows, rankedRow("q3q4", reportingClassificationOwnOther, "Q3", percentileOf(30)))
	}
	for i := 0; i < 56; i++ {
		rows = append(rows, rankedRow("upper", reportingClassificationOwnOther, "Q2", percentileOf(60)))
	}
	// Papers with no usable evidence sit outside both the numerator and the
	// denominator instead of quietly deflating the percentage.
	for i := 0; i < 12; i++ {
		rows = append(rows, reportingRow{Year: analyticsYear, Classification: reportingClassificationOwnOther, DocumentType: "proceedings"})
	}

	summary := buildExecutiveSummary(rows, analyticsYear, analyticsYear, analyticsYear, time.Now())
	require.Len(t, summary.Q3Q4HistoricalTrend, 1)

	own := summary.Q3Q4HistoricalTrend[0].Own
	assert.Equal(t, 22, own.Q3Q4Count)
	assert.Equal(t, 78, own.RankedCount)
	require.NotNil(t, own.Percent)
	assert.InDelta(t, 28.21, *own.Percent, 0.01)
}

func TestQ3Q4TrendReportsNoPercentForAnEmptyDenominator(t *testing.T) {
	rows := []reportingRow{
		{Year: analyticsYear, Classification: reportingClassificationOwnUser, DocumentType: "proceedings"},
		{Year: analyticsYear, Classification: reportingClassificationOwnUser, DocumentType: "article"},
	}
	summary := buildExecutiveSummary(rows, analyticsYear, analyticsYear, analyticsYear, time.Now())
	require.Len(t, summary.Q3Q4HistoricalTrend, 1)

	user := summary.Q3Q4HistoricalTrend[0].User
	assert.Equal(t, 0, user.RankedCount)
	assert.Equal(t, 0, user.Q3Q4Count)
	// Nil, not zero: no ranked papers means the rate is unknown, not 0%.
	assert.Nil(t, user.Percent)
}

func TestTrendCoversEveryRequestedYearEvenWhenEmpty(t *testing.T) {
	rows := []reportingRow{rankedRow("a", reportingClassificationOwnUser, "Q3", percentileOf(20))}
	summary := buildExecutiveSummary(rows, analyticsYear, analyticsYear-6, analyticsYear, time.Now())

	require.Len(t, summary.Q3Q4HistoricalTrend, 7)
	assert.Equal(t, analyticsYear-6, summary.Q3Q4HistoricalTrend[0].Year)
	assert.Equal(t, analyticsYear, summary.Q3Q4HistoricalTrend[6].Year)
	assert.Nil(t, summary.Q3Q4HistoricalTrend[0].Own.Percent)
}

func TestDepartmentCreditsDoNotInflateInstitutionalTotals(t *testing.T) {
	row := rankedRow("multi", reportingClassificationOwnUser, "Q1", percentileOf(95))
	row.Departments = []reportingReference{
		{UID: "d86", Name: "D86"}, {UID: "d88", Name: "D88"}, {UID: "d91", Name: "D91"},
		{UID: "d86", Name: "D86"}, // a duplicate must not earn a second credit
	}

	summary := buildExecutiveSummary([]reportingRow{row}, analyticsYear, analyticsYear, analyticsYear, time.Now())

	assert.Equal(t, 1, summary.TotalPublications, "one paper is one publication")
	assert.Equal(t, 1, summary.TotalOwnPublications)
	assert.Equal(t, 1, summary.TotalUserPublications)

	require.Len(t, summary.DepartmentMatrix, 3)
	credits := 0
	for _, department := range summary.DepartmentMatrix {
		credits += department.Total
		assert.Equal(t, 1, department.Q10Percent)
		assert.Equal(t, 1, department.TotalOwn)
		assert.Equal(t, 1, department.UserPublications)
	}
	assert.Equal(t, 3, credits, "the paper is credited once per department")
}

// The document's call chart accounts for 20 of 23 user papers. The three
// unexplained papers must be shown as unlinked, never assigned to a call.
func TestUnlinkedUserPapersAreSurfacedRatherThanRedistributed(t *testing.T) {
	rows := make([]reportingRow, 0, 23)
	for i := 0; i < 20; i++ {
		row := rankedRow("linked", reportingClassificationOwnUser, "Q2", percentileOf(60))
		row.UserCalls = []reportingReference{{UID: "call-1", Name: "Call 1"}}
		rows = append(rows, row)
	}
	for i := 0; i < 3; i++ {
		rows = append(rows, rankedRow("unlinked", reportingClassificationOwnUser, "Q2", percentileOf(60)))
	}

	summary := buildExecutiveSummary(rows, analyticsYear, analyticsYear, analyticsYear, time.Now())
	assert.Equal(t, 23, summary.TotalUserPublications)

	require.Len(t, summary.UserPublicationsByCall, 2)
	assert.Equal(t, "call-1", summary.UserPublicationsByCall[0].UID)
	assert.Equal(t, 20, summary.UserPublicationsByCall[0].Count)
	assert.Equal(t, unlinkedUID, summary.UserPublicationsByCall[1].UID)
	assert.Equal(t, 3, summary.UserPublicationsByCall[1].Count)
}

func TestClassificationTotalsKeepUserPapersASubsetOfOwn(t *testing.T) {
	rows := []reportingRow{
		rankedRow("a", reportingClassificationOwnUser, "Q1", percentileOf(95)),
		rankedRow("b", reportingClassificationOwnOther, "Q2", percentileOf(60)),
		rankedRow("c", reportingClassificationCoauthorship, "Q3", percentileOf(30)),
		{UID: "d", Year: analyticsYear, Classification: reportingClassificationUnclassified, DocumentType: "unknown"},
	}
	summary := buildExecutiveSummary(rows, analyticsYear, analyticsYear, analyticsYear, time.Now())

	assert.Equal(t, 4, summary.TotalPublications)
	assert.Equal(t, 2, summary.TotalOwnPublications, "own is own-user plus own-other")
	assert.Equal(t, 1, summary.TotalUserPublications)
	assert.Equal(t, 1, summary.OtherPublications)
	assert.Equal(t, 1, summary.CoauthorshipPublications)
	assert.Equal(t, 1, summary.UnclassifiedPublications)
	// Unclassified is reported beside the totals, never folded into them.
	assert.Equal(t, summary.TotalOwnPublications+summary.CoauthorshipPublications+summary.UnclassifiedPublications,
		summary.TotalPublications)
	assert.Equal(t, 1, summary.PendingReviewPublications)
}

func TestAuthorStatsKeepUnknownRolesVisible(t *testing.T) {
	yes := true
	rows := []reportingRow{
		{UID: "a", Year: analyticsYear, Classification: reportingClassificationOwnUser, DocumentType: "article",
			Authors: []reportingRowAuthor{
				{ResearcherUID: "r1", Name: "Nováková Jana", IsFirstAuthor: &yes, IsCorresponding: &yes},
				{ResearcherUID: "r2", Name: "Svoboda Petr"},
			}},
		{UID: "b", Year: analyticsYear, Classification: reportingClassificationOwnOther, DocumentType: "article",
			Authors: []reportingRowAuthor{{ResearcherUID: "r1", Name: "Nováková Jana"}}},
	}
	summary := buildExecutiveSummary(rows, analyticsYear, analyticsYear, analyticsYear, time.Now())

	require.Len(t, summary.TopPublishingAuthors, 2)
	first := summary.TopPublishingAuthors[0]
	assert.Equal(t, "r1", first.ResearcherUID)
	assert.Equal(t, 2, first.TotalAuthorships)
	assert.Equal(t, 1, first.FirstAuthorCount)
	assert.Equal(t, 1, first.CorrespondingCount)
	// A low confirmed count must be distinguishable from an unreviewed backlog.
	assert.Equal(t, 1, first.UnknownFirstAuthorCount)
	assert.Equal(t, 1, first.UnknownCorrespondingCount)

	second := summary.TopPublishingAuthors[1]
	assert.Equal(t, 0, second.FirstAuthorCount)
	assert.Equal(t, 1, second.UnknownFirstAuthorCount)
}

func TestSummaryIgnoresRowsOutsideTheReportingYear(t *testing.T) {
	rows := []reportingRow{
		rankedRow("current", reportingClassificationOwnUser, "Q1", percentileOf(95)),
		func() reportingRow {
			row := rankedRow("previous", reportingClassificationOwnUser, "Q4", percentileOf(10))
			row.Year = analyticsYear - 1
			row.JournalMetrics[0].Year = analyticsYear - 1
			return row
		}(),
	}
	summary := buildExecutiveSummary(rows, analyticsYear, analyticsYear-1, analyticsYear, time.Now())

	assert.Equal(t, 1, summary.TotalPublications, "headline totals cover the reporting year only")
	// The trend still spans the whole window.
	require.Len(t, summary.Q3Q4HistoricalTrend, 2)
	require.NotNil(t, summary.Q3Q4HistoricalTrend[0].Own.Percent)
	assert.InDelta(t, 100, *summary.Q3Q4HistoricalTrend[0].Own.Percent, 0.001)
	require.NotNil(t, summary.Q3Q4HistoricalTrend[1].Own.Percent)
	assert.InDelta(t, 0, *summary.Q3Q4HistoricalTrend[1].Own.Percent, 0.001)
}

func TestResolveReportWindowDefaultsToTheLatestCompletedYear(t *testing.T) {
	now := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	year, startYear, endYear := resolveReportWindow(0, 0, 0, now)
	assert.Equal(t, 2025, year)
	assert.Equal(t, 2025, endYear)
	assert.Equal(t, 2019, startYear, "the trend covers the year plus the preceding six")

	// An inverted range is clamped rather than returning nothing.
	_, startYear, endYear = resolveReportWindow(2025, 2030, 2025, now)
	assert.Equal(t, endYear, startYear)
}

func TestResolveReportWindowCapsAClientSuppliedSpan(t *testing.T) {
	now := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	// Every year in the window becomes an entry in the response, so an
	// unbounded span is a payload-inflation lever rather than a useful request.
	_, startYear, endYear := resolveReportWindow(2025, 1000, 9999, now)
	assert.Equal(t, 9999, endYear)
	assert.Equal(t, 9999-maxReportWindowYears+1, startYear)
	assert.LessOrEqual(t, endYear-startYear+1, maxReportWindowYears)

	// A normal request is untouched.
	_, startYear, endYear = resolveReportWindow(2025, 2019, 2025, now)
	assert.Equal(t, 2019, startYear)
	assert.Equal(t, 2025, endYear)
}
