package publicationsservice

import (
	"sort"
	"strings"
	"time"

	"panda/apigateway/services/publications-service/models"
)

// reportingPolicyVersion identifies the counting rules behind a report. Change
// it whenever a rule below changes, so an exported report can be traced back to
// the policy that produced it.
const reportingPolicyVersion = "2026-09-reporting-v1"

// Quality bands. The first three are document types that cannot carry a journal
// quartile; the rest are derived from the confirmed JCR evidence.
const (
	bandQ10Percent   = "q10Percent"
	bandQ10To25      = "q10To25"
	bandQ1Unsplit    = "q1Unsplit"
	bandQ2           = "q2"
	bandQ3           = "q3"
	bandQ4           = "q4"
	bandProceedings  = "proceedings"
	bandBookChapters = "bookChapters"
	bandOther        = "other"
	bandUnranked     = "unranked"
	bandUnknown      = "unknown"
)

// q10PercentThreshold splits Q1. At or above the 90th percentile a paper is in
// the top 10% of its category; below it, it is in the 10-25% band.
const q10PercentThreshold = 90.0

type reportingReference struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type reportingRowAuthor struct {
	ResearcherUID   string   `json:"researcherUid"`
	Name            string   `json:"name"`
	DepartmentUIDs  []string `json:"departmentUids"`
	IsFirstAuthor   *bool    `json:"isFirstAuthor"`
	IsCorresponding *bool    `json:"isCorresponding"`
}

// reportingRow is one publication as the analytics rules see it. Rows are
// fetched flat and aggregated in Go, so the counting rules stay readable and
// testable without a database.
type reportingRow struct {
	UID                  string                         `json:"uid"`
	Year                 int                            `json:"year"`
	JournalTitle         string                         `json:"journalTitle"`
	Classification       string                         `json:"classification"`
	Reviewed             bool                           `json:"reviewed"`
	DocumentType         string                         `json:"documentType"`
	JournalRankingStatus string                         `json:"journalRankingStatus"`
	Departments          []reportingReference           `json:"departments"`
	UserCalls            []reportingReference           `json:"userCalls"`
	ExperimentalSystems  []reportingReference           `json:"experimentalSystems"`
	Authors              []reportingRowAuthor           `json:"authors"`
	JournalMetrics       []models.JournalMetricSnapshot `json:"journalMetrics"`
}

func (row reportingRow) isOwn() bool {
	return row.Classification == reportingClassificationOwnUser ||
		row.Classification == reportingClassificationOwnOther
}

// isUser reports a user publication. User papers are a subset of own papers, so
// own-user is the only classification that qualifies.
func (row reportingRow) isUser() bool {
	return row.Classification == reportingClassificationOwnUser
}

// qualityBand assigns exactly one band. Document types that cannot be ranked
// win outright; everything else is decided by the confirmed JCR evidence for
// the publication year, and the absence of evidence stays visible as "unknown"
// rather than being folded into a quartile.
func (row reportingRow) qualityBand() string {
	switch row.DocumentType {
	case "proceedings":
		return bandProceedings
	case "book-chapter":
		return bandBookChapters
	case "other":
		return bandOther
	}

	metric, found := row.bestMetricForPublicationYear()
	if !found {
		if row.JournalRankingStatus == reportingRankingStatusUnranked {
			return bandUnranked
		}
		return bandUnknown
	}

	switch metric.Quartile {
	case "Q1":
		if metric.Percentile == nil {
			// A Q1 without a percentile cannot be split; saying which half it
			// falls in would be an invention.
			return bandQ1Unsplit
		}
		if *metric.Percentile >= q10PercentThreshold {
			return bandQ10Percent
		}
		return bandQ10To25
	case "Q2":
		return bandQ2
	case "Q3":
		return bandQ3
	case "Q4":
		return bandQ4
	default:
		return bandUnknown
	}
}

// bestMetricForPublicationYear picks the highest category percentile recorded
// for the publication's own year. Metrics from other years are ignored: a
// journal's standing moves, and reporting a paper against a later year's
// ranking would misstate it. A metric without a percentile only wins when no
// ranked alternative exists for that year.
func (row reportingRow) bestMetricForPublicationYear() (models.JournalMetricSnapshot, bool) {
	var best models.JournalMetricSnapshot
	found := false
	for _, metric := range row.JournalMetrics {
		if metric.Year != row.Year {
			continue
		}
		if !found {
			best, found = metric, true
			continue
		}
		if metricRanksHigher(metric, best) {
			best = metric
		}
	}
	return best, found
}

func metricRanksHigher(candidate, incumbent models.JournalMetricSnapshot) bool {
	switch {
	case candidate.Percentile == nil:
		return false
	case incumbent.Percentile == nil:
		return true
	default:
		return *candidate.Percentile > *incumbent.Percentile
	}
}

// rankedBands are the bands that represent a known journal quartile. They are
// the denominator of every Q3+Q4 percentage.
var rankedBands = map[string]bool{
	bandQ10Percent: true, bandQ10To25: true, bandQ1Unsplit: true,
	bandQ2: true, bandQ3: true, bandQ4: true,
}

func buildExecutiveSummary(rows []reportingRow, year, startYear, endYear int, generatedAt time.Time) models.ExecutiveSummary {
	summary := models.ExecutiveSummary{
		Year: year, StartYear: startYear, EndYear: endYear,
		GeneratedAt:   generatedAt.UTC().Format(time.RFC3339),
		PolicyVersion: reportingPolicyVersion,
	}

	current := make([]reportingRow, 0, len(rows))
	for _, row := range rows {
		if row.Year == year {
			current = append(current, row)
		}
	}

	// Institutional totals count distinct publications. Department, call and
	// system totals below are overlapping credits and may sum higher; they are
	// never used to derive these.
	for _, row := range current {
		summary.TotalPublications++
		switch row.Classification {
		case reportingClassificationOwnUser:
			summary.TotalOwnPublications++
			summary.TotalUserPublications++
		case reportingClassificationOwnOther:
			summary.TotalOwnPublications++
			summary.OtherPublications++
		case reportingClassificationCoauthorship:
			summary.CoauthorshipPublications++
		default:
			summary.UnclassifiedPublications++
		}
		if !row.Reviewed {
			summary.PendingReviewPublications++
		}
	}

	summary.DepartmentMatrix = buildDepartmentMatrix(current)
	summary.OwnQuality = buildQualityCounts(current, reportingRow.isOwn)
	summary.UserQuality = buildQualityCounts(current, reportingRow.isUser)
	summary.UserPublicationsByCall = buildUserCallCounts(current)
	summary.UserPublicationsByDepartment = buildUserDepartmentCounts(current)
	summary.SystemBreakdown = buildSystemCounts(current)
	summary.JournalFrequencies = buildJournalFrequencies(current)
	summary.TopPublishingAuthors = buildAuthorStats(current)
	summary.Q3Q4HistoricalTrend = buildTrend(rows, startYear, endYear)
	return summary
}

// buildDepartmentMatrix credits each publication once in every credited
// department. Rows therefore may sum above the institutional total, which is
// the agreed counting rule and why both are reported separately.
func buildDepartmentMatrix(rows []reportingRow) []models.DepartmentReportingRow {
	byUID := map[string]*models.DepartmentReportingRow{}
	order := []string{}

	for _, row := range rows {
		for _, department := range distinctReferences(row.Departments) {
			entry, exists := byUID[department.UID]
			if !exists {
				entry = &models.DepartmentReportingRow{UID: department.UID, Name: department.Name}
				byUID[department.UID] = entry
				order = append(order, department.UID)
			}
			entry.Total++
			switch row.qualityBand() {
			case bandQ10Percent:
				entry.Q10Percent++
			case bandQ10To25:
				entry.Q10To25++
			case bandQ1Unsplit:
				entry.Q1Unsplit++
			case bandQ2:
				entry.Q2++
			case bandQ3:
				entry.Q3++
			case bandQ4:
				entry.Q4++
			case bandProceedings:
				entry.Proceedings++
			case bandBookChapters:
				entry.BookChapters++
			case bandOther:
				entry.Other++
			case bandUnranked:
				entry.Unranked++
			default:
				entry.Unknown++
			}
			if row.isOwn() {
				entry.TotalOwn++
			}
			if row.Classification == reportingClassificationCoauthorship {
				entry.CoAuthorship++
			}
			if row.isUser() {
				entry.UserPublications++
			}
		}
	}

	result := make([]models.DepartmentReportingRow, 0, len(order))
	for _, uid := range order {
		result = append(result, *byUID[uid])
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func buildQualityCounts(rows []reportingRow, include func(reportingRow) bool) []models.ReportingQualityCount {
	counts := map[string]int{}
	for _, row := range rows {
		if include(row) {
			counts[row.qualityBand()]++
		}
	}
	// Every band is reported, including the empty ones, so a reader can see that
	// a band is genuinely zero rather than missing from the response.
	bands := []string{
		bandQ10Percent, bandQ10To25, bandQ1Unsplit, bandQ2, bandQ3, bandQ4,
		bandProceedings, bandBookChapters, bandOther, bandUnranked, bandUnknown,
	}
	result := make([]models.ReportingQualityCount, 0, len(bands))
	for _, band := range bands {
		result = append(result, models.ReportingQualityCount{Quality: band, Count: counts[band]})
	}
	return result
}

// unlinkedUID labels the bucket for user publications with no confirmed link of
// a given kind. The count is surfaced rather than redistributed, so a chart can
// never silently account for papers it has no evidence about.
const unlinkedUID = "__unlinked__"

func buildUserCallCounts(rows []reportingRow) []models.ReportingCount {
	return buildLinkCounts(rows, func(row reportingRow) []reportingReference { return row.UserCalls }, "No confirmed call")
}

func buildSystemCounts(rows []reportingRow) []models.ReportingCount {
	return buildLinkCounts(rows, func(row reportingRow) []reportingReference { return row.ExperimentalSystems }, "No confirmed system")
}

func buildUserDepartmentCounts(rows []reportingRow) []models.ReportingCount {
	return buildLinkCounts(rows, func(row reportingRow) []reportingReference { return row.Departments }, "No credited department")
}

func buildLinkCounts(rows []reportingRow, links func(reportingRow) []reportingReference, unlinkedLabel string) []models.ReportingCount {
	counts := map[string]*models.ReportingCount{}
	order := []string{}
	unlinked := 0

	for _, row := range rows {
		if !row.isUser() {
			continue
		}
		references := distinctReferences(links(row))
		if len(references) == 0 {
			unlinked++
			continue
		}
		for _, reference := range references {
			entry, exists := counts[reference.UID]
			if !exists {
				entry = &models.ReportingCount{UID: reference.UID, Name: reference.Name}
				counts[reference.UID] = entry
				order = append(order, reference.UID)
			}
			entry.Count++
		}
	}

	result := make([]models.ReportingCount, 0, len(order)+1)
	for _, uid := range order {
		result = append(result, *counts[uid])
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Name < result[j].Name
	})
	if unlinked > 0 {
		result = append(result, models.ReportingCount{UID: unlinkedUID, Name: unlinkedLabel, Count: unlinked})
	}
	return result
}

func buildJournalFrequencies(rows []reportingRow) []models.ReportingCount {
	counts := map[string]*models.ReportingCount{}
	for _, row := range rows {
		title := strings.TrimSpace(row.JournalTitle)
		if title == "" {
			continue
		}
		key := strings.ToLower(title)
		entry, exists := counts[key]
		if !exists {
			entry = &models.ReportingCount{UID: key, Name: title}
			counts[key] = entry
		}
		entry.Count++
	}
	result := make([]models.ReportingCount, 0, len(counts))
	for _, entry := range counts {
		result = append(result, *entry)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Name < result[j].Name
	})
	return result
}

// buildAuthorStats counts authorships, not publications, and keeps the unknown
// role counts beside the confirmed ones so a low first-author count is
// distinguishable from an unreviewed backlog.
func buildAuthorStats(rows []reportingRow) []models.ReportingAuthorStats {
	stats := map[string]*models.ReportingAuthorStats{}
	order := []string{}

	for _, row := range rows {
		seen := map[string]bool{}
		for _, author := range row.Authors {
			if author.ResearcherUID == "" || seen[author.ResearcherUID] {
				continue
			}
			seen[author.ResearcherUID] = true

			entry, exists := stats[author.ResearcherUID]
			if !exists {
				entry = &models.ReportingAuthorStats{
					ResearcherUID:  author.ResearcherUID,
					Name:           author.Name,
					DepartmentUIDs: []string{},
				}
				stats[author.ResearcherUID] = entry
				order = append(order, author.ResearcherUID)
			}
			entry.TotalAuthorships++
			switch {
			case author.IsFirstAuthor == nil:
				entry.UnknownFirstAuthorCount++
			case *author.IsFirstAuthor:
				entry.FirstAuthorCount++
			}
			switch {
			case author.IsCorresponding == nil:
				entry.UnknownCorrespondingCount++
			case *author.IsCorresponding:
				entry.CorrespondingCount++
			}
			entry.DepartmentUIDs = appendDistinct(entry.DepartmentUIDs, author.DepartmentUIDs)
		}
	}

	result := make([]models.ReportingAuthorStats, 0, len(order))
	for _, uid := range order {
		result = append(result, *stats[uid])
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].TotalAuthorships != result[j].TotalAuthorships {
			return result[i].TotalAuthorships > result[j].TotalAuthorships
		}
		return result[i].Name < result[j].Name
	})
	return result
}

func buildTrend(rows []reportingRow, startYear, endYear int) []models.ReportingTrend {
	byYear := map[int][]reportingRow{}
	for _, row := range rows {
		byYear[row.Year] = append(byYear[row.Year], row)
	}
	trend := make([]models.ReportingTrend, 0, endYear-startYear+1)
	for year := startYear; year <= endYear; year++ {
		yearRows := byYear[year]
		trend = append(trend, models.ReportingTrend{
			Year: year,
			Own:  buildFraction(yearRows, reportingRow.isOwn),
			User: buildFraction(yearRows, reportingRow.isUser),
		})
	}
	return trend
}

// buildFraction divides Q3+Q4 by the papers with a known quartile. Proceedings,
// book chapters, confirmed unranked journals and papers still missing evidence
// are excluded from both sides, and an empty denominator yields a nil percent
// that the UI renders as N/A rather than as zero.
func buildFraction(rows []reportingRow, include func(reportingRow) bool) models.ReportingFraction {
	fraction := models.ReportingFraction{}
	for _, row := range rows {
		if !include(row) {
			continue
		}
		band := row.qualityBand()
		if !rankedBands[band] {
			continue
		}
		fraction.RankedCount++
		if band == bandQ3 || band == bandQ4 {
			fraction.Q3Q4Count++
		}
	}
	if fraction.RankedCount > 0 {
		percent := float64(fraction.Q3Q4Count) / float64(fraction.RankedCount) * 100
		fraction.Percent = &percent
	}
	return fraction
}

func distinctReferences(references []reportingReference) []reportingReference {
	result := make([]reportingReference, 0, len(references))
	seen := map[string]bool{}
	for _, reference := range references {
		if reference.UID == "" || seen[reference.UID] {
			continue
		}
		seen[reference.UID] = true
		result = append(result, reference)
	}
	return result
}

func appendDistinct(target []string, values []string) []string {
	for _, value := range values {
		found := false
		for _, existing := range target {
			if existing == value {
				found = true
				break
			}
		}
		if !found && value != "" {
			target = append(target, value)
		}
	}
	return target
}

// defaultReportYear is the most recent completed calendar year. The dashboard
// opens on it and advances on its own each January.
func defaultReportYear(now time.Time) int { return now.Year() - 1 }

// GetExecutiveSummary calculates the management report from saved records. A
// database failure is returned as an error: an empty or partial report would
// look like a real one.
func (svc *PublicationsService) GetExecutiveSummary(year, startYear, endYear int) (models.ExecutiveSummary, error) {
	rows, err := svc.listReportingRows(startYear, endYear)
	if err != nil {
		return models.ExecutiveSummary{}, err
	}
	return buildExecutiveSummary(rows, year, startYear, endYear, time.Now()), nil
}

// resolveReportWindow applies the agreed defaults and keeps the range sane: the
// latest completed calendar year, plus the six years before it.
func resolveReportWindow(year, startYear, endYear int, now time.Time) (int, int, int) {
	if year <= 0 {
		year = defaultReportYear(now)
	}
	if endYear <= 0 {
		endYear = year
	}
	if startYear <= 0 {
		startYear = endYear - 6
	}
	if startYear > endYear {
		startYear = endYear
	}
	return year, startYear, endYear
}
