package publicationsservice

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"panda/apigateway/services/publications-service/models"
)

// PublicationEnrichmentConfig keeps provider credentials on the server. Empty
// URLs select official endpoints; an empty UnpaywallEmail disables that source.
type PublicationEnrichmentConfig struct {
	CrossrefURL    string
	CrossrefEmail  string
	UnpaywallURL   string
	UnpaywallEmail string
	HTTPClient     *http.Client
}

type PublicationEnrichmentService struct {
	publications *PublicationsService
	config       PublicationEnrichmentConfig
	client       *http.Client
}

type enrichmentSource struct {
	status        models.EnrichmentSourceStatus
	values        models.WosImportValues
	authors       []models.WosAuthor
	mediaTypeCode string
	datePrecision string
	openAccess    *models.EnrichmentOpenAccess
}

// NewPublicationEnrichmentService composes on the live publications service so
// the WoS credentials, HTTP client and Neo4j repository stay in one place.
func NewPublicationEnrichmentService(publications *PublicationsService, cfg PublicationEnrichmentConfig) *PublicationEnrichmentService {
	if cfg.CrossrefURL == "" {
		cfg.CrossrefURL = "https://api.crossref.org/works"
	}
	if cfg.UnpaywallURL == "" {
		cfg.UnpaywallURL = "https://api.unpaywall.org/v2"
	}
	client := newWosHTTPClient()
	if cfg.HTTPClient != nil {
		// Copy so a caller-supplied client keeps its transport but cannot
		// relax the timeout or redirect policy these lookups rely on.
		copyClient := *cfg.HTTPClient
		client = &copyClient
		client.Timeout = wosRequestTimeout
		client.CheckRedirect = newWosHTTPClient().CheckRedirect
	}
	return &PublicationEnrichmentService{publications: publications, config: cfg, client: client}
}

// Preview is read-only. Bibliographic suggestions prefer WoS when available,
// then Crossref; publication dates prefer Crossref's explicit publication date
// (published, print, online) as a single year/date pair. All disagreements remain
// visible and applying suggestions never persists a publication.
func (svc *PublicationEnrichmentService) Preview(ctx context.Context, rawDOI, currentPublicationUID, facilityCode string) (models.EnrichmentPreviewResponse, error) {
	doi, valid := normalizeDOI(rawDOI)
	if !valid || len(doi) > 2048 || strings.ContainsAny(doi, "\x00\r\n") {
		return models.EnrichmentPreviewResponse{}, newPublicationAPIError(http.StatusBadRequest, wosErrorInvalidDOI, "Enter a valid DOI.", false)
	}
	result := models.EnrichmentPreviewResponse{
		Doi: doi, Status: "not-found", Authors: []models.WosImportAuthor{},
		MissingImportableFields: []string{}, UnavailableFields: append([]string{}, wosUnavailableFields...),
		Sources: []models.EnrichmentSourceStatus{}, Provenance: map[string]models.EnrichmentOrigin{},
		Conflicts: []models.EnrichmentConflict{}, AuthorRolesStatus: "unknown", AffiliationStatus: "unknown",
	}
	repository := svc.publications.importRepository()
	if repository != nil {
		existing, err := repository.findExistingPublication(doi, currentPublicationUID)
		if err != nil {
			return result, fmt.Errorf("find existing enrichment publication: %w", err)
		}
		if existing != nil {
			result.Status, result.ExistingPublication = "already-exists", existing
			for _, provider := range []string{"crossref", "wos-starter", "unpaywall"} {
				result.Sources = append(result.Sources, models.EnrichmentSourceStatus{Provider: provider, Status: "skipped"})
			}
			return result, nil
		}
	}
	ctx, cancel := context.WithTimeout(ctx, wosRequestTimeout)
	defer cancel()
	sources := make([]enrichmentSource, 3)
	var group sync.WaitGroup
	for index, fetch := range []func(context.Context, string) enrichmentSource{svc.fetchCrossref, svc.fetchEnrichmentWos, svc.fetchUnpaywall} {
		group.Add(1)
		go func(index int, fetch func(context.Context, string) enrichmentSource) {
			defer group.Done()
			sources[index] = fetch(ctx, doi)
		}(index, fetch)
	}
	group.Wait()
	for _, source := range sources {
		result.Sources = append(result.Sources, source.status)
	}
	values := models.WosImportValues{Doi: &doi}
	var authorSource, mediaSource *enrichmentSource
	// Stable precedence also keeps the selected author order and researcher IDs
	// together, instead of merging unrelated names across provider records.
	for _, index := range []int{1, 0} {
		source := &sources[index]
		if source.status.Status != "ok" {
			continue
		}
		result.Status = "found"
		mergeEnrichmentValues(&values, source, &result)
		if authorSource == nil && len(source.authors) > 0 {
			authorSource = source
		}
		if mediaSource == nil && source.mediaTypeCode != "" {
			mediaSource = source
		}
	}
	// A year and date always come from the same source. Crossref dates can carry
	// day precision absent from Starter; DOI registration dates are never used.
	for _, index := range []int{0, 1} {
		source := &sources[index]
		if source.status.Status != "ok" || source.values.YearOfPublication == nil {
			continue
		}
		if values.YearOfPublication == nil {
			values.YearOfPublication, values.DateOfPublication = source.values.YearOfPublication, source.values.DateOfPublication
			result.PublicationDatePrecision = source.datePrecision
			result.Provenance["yearOfPublication"] = enrichmentOrigin(source)
			if values.DateOfPublication != nil {
				result.Provenance["dateOfPublication"] = enrichmentOrigin(source)
			}
		} else {
			recordEnrichmentConflict("yearOfPublication", values.YearOfPublication, source.values.YearOfPublication, source, &result)
			recordEnrichmentConflict("dateOfPublication", values.DateOfPublication, source.values.DateOfPublication, source, &result)
		}
	}
	if sources[2].status.Status == "ok" {
		result.OpenAccess = sources[2].openAccess
		result.Provenance["openAccess"] = enrichmentOrigin(&sources[2])
	}
	if result.Status != "found" {
		for _, source := range sources[:2] {
			if source.status.Status == "error" || source.status.Status == "ambiguous" {
				result.Status = "unavailable"
			}
		}
		result.MissingImportableFields = missingWosImportableFields(values)
		return result, nil
	}
	if authorSource != nil {
		researchers := []wosResearcherRecord{}
		if repository != nil {
			var err error
			researchers, err = repository.listResearchersForWos()
			if err != nil {
				return result, fmt.Errorf("list researchers for enrichment: %w", err)
			}
		}
		result.Authors = matchWosAuthors(authorSource.authors, researchers)
		result.Provenance["authors"] = enrichmentOrigin(authorSource)
	}
	if repository != nil && mediaSource != nil {
		var err error
		values.MediaTypeCb, err = repository.findMediaType(mediaSource.mediaTypeCode, facilityCode)
		if err != nil {
			return result, fmt.Errorf("resolve enrichment media type: %w", err)
		}
		if values.MediaTypeCb != nil {
			result.Provenance["mediaTypeCb"] = enrichmentOrigin(mediaSource)
		}
	}
	result.Values = &values
	result.MissingImportableFields = missingWosImportableFields(values)
	return result, nil
}

func enrichmentOrigin(source *enrichmentSource) models.EnrichmentOrigin {
	return models.EnrichmentOrigin{Provider: source.status.Provider, RetrievedAt: source.status.RetrievedAt}
}

// WosImportValues is a fixed pointer-only form contract. Walking its JSON tags
// makes every populated importable field receive provenance and conflict checks.
func mergeEnrichmentValues(values *models.WosImportValues, source *enrichmentSource, result *models.EnrichmentPreviewResponse) {
	destination, incoming := reflect.ValueOf(values).Elem(), reflect.ValueOf(source.values)
	for index := 0; index < incoming.NumField(); index++ {
		field := strings.Split(incoming.Type().Field(index).Tag.Get("json"), ",")[0]
		if field == "doi" || field == "yearOfPublication" || field == "dateOfPublication" || incoming.Field(index).IsNil() {
			continue
		}
		if destination.Field(index).IsNil() {
			destination.Field(index).Set(incoming.Field(index))
			result.Provenance[field] = enrichmentOrigin(source)
		} else {
			recordEnrichmentConflict(field, destination.Field(index).Interface(), incoming.Field(index).Interface(), source, result)
		}
	}
}

func recordEnrichmentConflict(field string, selected, alternative interface{}, source *enrichmentSource, result *models.EnrichmentPreviewResponse) {
	left, right := reflect.ValueOf(selected), reflect.ValueOf(alternative)
	if !left.IsValid() || !right.IsValid() || left.IsNil() || right.IsNil() {
		return
	}
	leftValue, rightValue := left.Elem().Interface(), right.Elem().Interface()
	if reflect.DeepEqual(leftValue, rightValue) {
		return
	}
	if a, ok := leftValue.(string); ok {
		if b, ok := rightValue.(string); ok && strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b)) {
			return
		}
	}
	result.Conflicts = append(result.Conflicts, models.EnrichmentConflict{
		Field: field, SelectedProvider: result.Provenance[field].Provider, SelectedValue: leftValue,
		AlternativeProvider: source.status.Provider, AlternativeValue: rightValue,
	})
}

func (svc *PublicationEnrichmentService) fetchEnrichmentWos(ctx context.Context, doi string) enrichmentSource {
	source := enrichmentSource{status: models.EnrichmentSourceStatus{Provider: "wos-starter"}}
	_, hit, err := svc.publications.fetchExactWosRecord(ctx, doi)
	if err != nil {
		source.status.Status, source.status.Code, source.status.Message = "error", "WOS_UPSTREAM_ERROR", "Web of Science is currently unavailable."
		var apiErr *publicationAPIError
		if errors.As(err, &apiErr) {
			source.status.Code, source.status.Message = apiErr.Code, apiErr.Message
			source.status.Retryable, source.status.RetryAfter = apiErr.Retryable, apiErr.RetryAfter
			switch apiErr.Code {
			case wosErrorNotConfigured:
				source.status.Status = "not-configured"
			case wosErrorRecordNotFound:
				source.status.Status = "not-found"
			case wosErrorRecordAmbiguous:
				source.status.Status = "ambiguous"
			}
		}
		return source
	}
	source.status.Status, source.status.RetrievedAt = "ok", time.Now().UTC().Format(time.RFC3339)
	// Provider-level mapping warnings are intentionally dropped here: the
	// enrichment response surfaces its own per-provider status, and a degraded
	// WoS field still has Crossref/Unpaywall candidates to fall back on.
	source.values, _ = mapWosImportValues(hit, doi)
	if hit.WosSource.WosPublishYear < 1 || hit.WosSource.WosPublishYear > 9999 {
		source.values.YearOfPublication, source.values.DateOfPublication = nil, nil
	} else {
		source.datePrecision = "year"
		if source.values.DateOfPublication != nil {
			source.datePrecision = "month"
		}
	}
	source.authors = hit.WosNames.WosAuthors
	source.mediaTypeCode = suggestMediaTypeCode(hit.WosTypes, hit.WosSourceTypes)
	return source
}
