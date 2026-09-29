package models

// PublicationReporting is an editor-reviewed, publication-time reporting snapshot.
// Omission on a legacy save preserves it; submitting it replaces the whole snapshot.
type PublicationReporting struct {
	Classification         string                  `json:"classification" enums:"own-user,own-other,coauthorship,unclassified"`
	Reviewed               bool                    `json:"reviewed"`
	DocumentType           string                  `json:"documentType" enums:"article,proceedings,book-chapter,other,unknown"`
	JournalRankingStatus   string                  `json:"journalRankingStatus,omitempty" enums:"unknown,unranked"`
	DepartmentUIDs         []string                `json:"departmentUids"`
	Authors                []ReportingAuthor       `json:"authors"`
	UserCallUIDs           []string                `json:"userCallUids"`
	UserExperimentUIDs     []string                `json:"userExperimentUids"`
	ExperimentalSystemUIDs []string                `json:"experimentalSystemUids"`
	JournalMetrics         []JournalMetricSnapshot `json:"journalMetrics"`
	ReviewedAt             string                  `json:"reviewedAt,omitempty"`
	ReviewedBy             string                  `json:"reviewedBy,omitempty"`
}

type ReportingDepartment struct {
	DepartmentUID string `json:"departmentUid"`
}
type ReportingAuthor struct {
	ResearcherUID   string   `json:"researcherUid"`
	DepartmentUIDs  []string `json:"departmentUids"`
	IsFirstAuthor   *bool    `json:"isFirstAuthor"`
	IsCorresponding *bool    `json:"isCorresponding"`
}
type JournalMetricSnapshot struct {
	Source       string   `json:"source"`
	JournalID    string   `json:"journalId,omitempty"`
	Year         int      `json:"year"`
	Category     string   `json:"category"`
	Quartile     string   `json:"quartile"`
	Percentile   *float64 `json:"percentile"`
	ImpactFactor *float64 `json:"impactFactor"`
}

type ReportingCount struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}
type ReportingQualityCount struct {
	Quality string `json:"quality"`
	Count   int    `json:"count"`
}
type DepartmentReportingRow struct {
	UID              string `json:"uid"`
	Name             string `json:"name"`
	Q10Percent       int    `json:"q10Percent"`
	Q10To25          int    `json:"q10To25"`
	Q1Unsplit        int    `json:"q1Unsplit"`
	Q2               int    `json:"q2"`
	Q3               int    `json:"q3"`
	Q4               int    `json:"q4"`
	Proceedings      int    `json:"proceedings"`
	BookChapters     int    `json:"bookChapters"`
	Other            int    `json:"other"`
	Unranked         int    `json:"unranked"`
	Unknown          int    `json:"unknown"`
	TotalOwn         int    `json:"totalOwn"`
	CoAuthorship     int    `json:"coAuthorship"`
	Total            int    `json:"total"`
	UserPublications int    `json:"userPublications"`
}
type ReportingAuthorStats struct {
	ResearcherUID             string   `json:"researcherUid"`
	Name                      string   `json:"name"`
	DepartmentUIDs            []string `json:"departmentUids"`
	TotalAuthorships          int      `json:"totalAuthorships"`
	FirstAuthorCount          int      `json:"firstAuthorCount"`
	CorrespondingCount        int      `json:"correspondingCount"`
	UnknownFirstAuthorCount   int      `json:"unknownFirstAuthorCount"`
	UnknownCorrespondingCount int      `json:"unknownCorrespondingCount"`
}
type ReportingFraction struct {
	Q3Q4Count   int      `json:"q3q4Count"`
	RankedCount int      `json:"rankedCount"`
	Percent     *float64 `json:"percent"`
}
type ReportingTrend struct {
	Year int               `json:"year"`
	Own  ReportingFraction `json:"own"`
	User ReportingFraction `json:"user"`
}
type ExecutiveSummary struct {
	Year                         int                      `json:"year"`
	StartYear                    int                      `json:"startYear"`
	EndYear                      int                      `json:"endYear"`
	GeneratedAt                  string                   `json:"generatedAt"`
	PolicyVersion                string                   `json:"policyVersion"`
	TotalPublications            int                      `json:"totalPublications"`
	TotalOwnPublications         int                      `json:"totalOwnPublications"`
	TotalUserPublications        int                      `json:"totalUserPublications"`
	OtherPublications            int                      `json:"otherPublications"`
	CoauthorshipPublications     int                      `json:"coauthorshipPublications"`
	UnclassifiedPublications     int                      `json:"unclassifiedPublications"`
	PendingReviewPublications    int                      `json:"pendingReviewPublications"`
	DepartmentMatrix             []DepartmentReportingRow `json:"departmentMatrix"`
	OwnQuality                   []ReportingQualityCount  `json:"ownQuality"`
	UserQuality                  []ReportingQualityCount  `json:"userQuality"`
	UserPublicationsByCall       []ReportingCount         `json:"userPublicationsByCall"`
	UserPublicationsByDepartment []ReportingCount         `json:"userPublicationsByDepartment"`
	SystemBreakdown              []ReportingCount         `json:"systemBreakdown"`
	JournalFrequencies           []ReportingCount         `json:"journalFrequencies"`
	TopPublishingAuthors         []ReportingAuthorStats   `json:"topPublishingAuthors"`
	Q3Q4HistoricalTrend          []ReportingTrend         `json:"q3q4HistoricalTrend"`
}
