package models

// EnrichmentPreviewRequest requests suggestions only; saving uses the existing
// publication create/update endpoints after the user reviews the preview.
type EnrichmentPreviewRequest struct {
	Doi                   string `json:"doi"`
	CurrentPublicationUID string `json:"currentPublicationUid,omitempty"`
}

type EnrichmentPreviewResponse struct {
	Status                   string                      `json:"status"`
	Doi                      string                      `json:"doi"`
	ExistingPublication      *WosExistingPublication     `json:"existingPublication,omitempty"`
	Values                   *WosImportValues            `json:"values,omitempty"`
	Authors                  []WosImportAuthor           `json:"authors"`
	MissingImportableFields  []string                    `json:"missingImportableFields"`
	UnavailableFields        []string                    `json:"unavailableFields"`
	Sources                  []EnrichmentSourceStatus    `json:"sources"`
	Provenance               map[string]EnrichmentOrigin `json:"provenance"`
	Conflicts                []EnrichmentConflict        `json:"conflicts"`
	PublicationDatePrecision string                      `json:"publicationDatePrecision,omitempty"`
	OpenAccess               *EnrichmentOpenAccess       `json:"openAccess,omitempty"`
	AuthorRolesStatus        string                      `json:"authorRolesStatus"`
	AffiliationStatus        string                      `json:"affiliationStatus"`
}

type EnrichmentSourceStatus struct {
	Provider    string `json:"provider"`
	Status      string `json:"status"`
	Code        string `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Retryable   bool   `json:"retryable"`
	RetryAfter  string `json:"retryAfter,omitempty"`
	RetrievedAt string `json:"retrievedAt,omitempty"`
}

type EnrichmentOrigin struct {
	Provider    string `json:"provider"`
	RetrievedAt string `json:"retrievedAt"`
}

type EnrichmentConflict struct {
	Field               string      `json:"field"`
	SelectedProvider    string      `json:"selectedProvider"`
	SelectedValue       interface{} `json:"selectedValue"`
	AlternativeProvider string      `json:"alternativeProvider"`
	AlternativeValue    interface{} `json:"alternativeValue"`
}

// Open-access metadata is informational until a user chooses the corresponding
// institutional codebook entry. Missing metadata must not be interpreted as closed.
type EnrichmentOpenAccess struct {
	Status string `json:"status"`
	URL    string `json:"url,omitempty"`
	PDFURL string `json:"pdfUrl,omitempty"`
}
