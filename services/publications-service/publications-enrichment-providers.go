package publicationsservice

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"panda/apigateway/services/publications-service/models"
)

const (
	enrichmentProviderCrossref  = "crossref"
	enrichmentProviderUnpaywall = "unpaywall"
)

// crossrefWork is the subset of a Crossref `/works/{doi}` record PANDA maps onto
// the publication form. Every field stays optional: Crossref routinely omits
// pagination, ISSNs and subjects, and an absent value must remain absent rather
// than becoming an empty string in the form.
type crossrefWork struct {
	Message struct {
		Title          []string         `json:"title"`
		ContainerTitle []string         `json:"container-title"`
		Volume         string           `json:"volume"`
		Issue          string           `json:"issue"`
		Page           string           `json:"page"`
		Type           string           `json:"type"`
		URL            string           `json:"URL"`
		Subject        []string         `json:"subject"`
		ISBN           []string         `json:"ISBN"`
		ISSN           []string         `json:"ISSN"`
		ISSNType       []crossrefISSN   `json:"issn-type"`
		Author         []crossrefAuthor `json:"author"`
		// Publication dates, in the order of preference below. `created` is
		// deliberately absent: it records when Crossref first received the
		// metadata, not when the work was published.
		Published       crossrefDate `json:"published"`
		PublishedPrint  crossrefDate `json:"published-print"`
		PublishedOnline crossrefDate `json:"published-online"`
	} `json:"message"`
}

type crossrefISSN struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type crossrefAuthor struct {
	Given  string `json:"given"`
	Family string `json:"family"`
	Name   string `json:"name"`
	ORCID  string `json:"ORCID"`
}

// crossrefDate carries partial precision: date-parts may hold a year, a year and
// month, or a full date. The length is the precision and must not be padded.
type crossrefDate struct {
	DateParts [][]int `json:"date-parts"`
}

// unpaywallRecord is the subset of the Unpaywall v2 DOI response PANDA uses.
// Open access is informational until an editor picks the matching codebook entry.
type unpaywallRecord struct {
	OaStatus     string `json:"oa_status"`
	IsOa         bool   `json:"is_oa"`
	BestOaLocall *struct {
		URL       string `json:"url"`
		URLForPdf string `json:"url_for_pdf"`
	} `json:"best_oa_location"`
}

// fetchCrossref returns bibliographic suggestions for a DOI. It never persists
// anything and never fails the preview as a whole: an unreachable Crossref is
// reported as one unavailable source while the other providers still contribute.
func (svc *PublicationEnrichmentService) fetchCrossref(ctx context.Context, doi string) enrichmentSource {
	source := enrichmentSource{status: models.EnrichmentSourceStatus{Provider: enrichmentProviderCrossref}}

	endpoint, err := enrichmentEndpoint(svc.config.CrossrefURL, doi)
	if err != nil {
		source.status = enrichmentUpstreamStatus(enrichmentProviderCrossref, "Crossref is currently unavailable.")
		return source
	}
	// The mailto contact puts PANDA in Crossref's polite pool, which is rate
	// limited far more generously than anonymous traffic.
	if email := strings.TrimSpace(svc.config.CrossrefEmail); email != "" {
		query := endpoint.Query()
		query.Set("mailto", email)
		endpoint.RawQuery = query.Encode()
	}

	var work crossrefWork
	if status, ok := svc.decodeEnrichmentJSON(ctx, enrichmentProviderCrossref, endpoint, &work); !ok {
		source.status = status
		return source
	}

	source.status.Status = "ok"
	source.status.RetrievedAt = time.Now().UTC().Format(time.RFC3339)
	source.values = mapCrossrefValues(work, doi)
	source.values.YearOfPublication, source.values.DateOfPublication, source.datePrecision = crossrefPublicationDate(work)
	source.authors = mapCrossrefAuthors(work.Message.Author)
	// Crossref type names are hyphenated ("book-chapter"); the shared classifier
	// matches on spaced words.
	source.mediaTypeCode = suggestMediaTypeCode([]string{strings.ReplaceAll(work.Message.Type, "-", " ")}, nil)
	return source
}

// fetchUnpaywall resolves open-access state only. It contributes no bibliographic
// values, so a failure here never degrades the rest of the preview.
func (svc *PublicationEnrichmentService) fetchUnpaywall(ctx context.Context, doi string) enrichmentSource {
	source := enrichmentSource{status: models.EnrichmentSourceStatus{Provider: enrichmentProviderUnpaywall}}

	// Unpaywall requires an institutional contact on every request. Without one
	// the source stays unconfigured rather than silently reporting "not found".
	email := strings.TrimSpace(svc.config.UnpaywallEmail)
	if email == "" {
		source.status.Status = "not-configured"
		source.status.Code = wosErrorNotConfigured
		source.status.Message = "Unpaywall lookup is not configured."
		return source
	}

	endpoint, err := enrichmentEndpoint(svc.config.UnpaywallURL, doi)
	if err != nil {
		source.status = enrichmentUpstreamStatus(enrichmentProviderUnpaywall, "Unpaywall is currently unavailable.")
		return source
	}
	query := endpoint.Query()
	query.Set("email", email)
	endpoint.RawQuery = query.Encode()

	var record unpaywallRecord
	if status, ok := svc.decodeEnrichmentJSON(ctx, enrichmentProviderUnpaywall, endpoint, &record); !ok {
		source.status = status
		return source
	}

	source.status.Status = "ok"
	source.status.RetrievedAt = time.Now().UTC().Format(time.RFC3339)
	source.openAccess = mapUnpaywallOpenAccess(record)
	return source
}

// enrichmentEndpoint appends the DOI as a single path element. Assigning the
// decoded Path lets net/url escape it on String(), so a DOI's slashes and any
// reserved characters survive without being double-encoded.
func enrichmentEndpoint(base, doi string) (*url.URL, error) {
	endpoint, err := url.Parse(strings.TrimRight(base, "/"))
	if err != nil {
		return nil, err
	}
	if endpoint.Scheme == "" || endpoint.Host == "" {
		return nil, errors.New("enrichment endpoint is not absolute")
	}
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + "/" + doi
	return endpoint, nil
}

// decodeEnrichmentJSON performs the GET and decodes the body. The bool reports
// whether target was populated; when false the returned status explains why, in
// the per-provider vocabulary the preview response exposes to the client.
func (svc *PublicationEnrichmentService) decodeEnrichmentJSON(
	ctx context.Context,
	provider string,
	endpoint *url.URL,
	target interface{},
) (models.EnrichmentSourceStatus, bool) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return enrichmentUpstreamStatus(provider, enrichmentUnavailableMessage(provider)), false
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "PANDA/1.0 (https://panda.eli-laser.eu; publication metadata enrichment)")

	client := svc.client
	if client == nil {
		client = newWosHTTPClient()
	}
	response, err := client.Do(request)
	if err != nil {
		var networkError net.Error
		if errors.Is(ctx.Err(), context.DeadlineExceeded) ||
			errors.As(err, &networkError) && networkError.Timeout() {
			status := enrichmentUpstreamStatus(provider, enrichmentTimeoutMessage(provider))
			status.Code = wosErrorUpstreamTimeout
			return status, false
		}
		return enrichmentUpstreamStatus(provider, enrichmentUnavailableMessage(provider)), false
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return enrichmentStatusFromResponse(provider, response), false
	}

	decoder := json.NewDecoder(io.LimitReader(response.Body, wosResponseLimit))
	if err := decoder.Decode(target); err != nil {
		status := enrichmentUpstreamStatus(provider, enrichmentInvalidMessage(provider))
		return status, false
	}
	return models.EnrichmentSourceStatus{Provider: provider, Status: "ok"}, true
}

func enrichmentStatusFromResponse(provider string, response *http.Response) models.EnrichmentSourceStatus {
	status := models.EnrichmentSourceStatus{Provider: provider}
	switch response.StatusCode {
	case http.StatusNotFound, http.StatusGone:
		status.Status = "not-found"
		status.Code = wosErrorRecordNotFound
		status.Message = enrichmentNotFoundMessage(provider)
	case http.StatusUnauthorized, http.StatusForbidden:
		status.Status = "error"
		status.Code = wosErrorAuthenticationFailed
		status.Message = "PANDA could not authenticate with " + enrichmentProviderName(provider) + "."
	case http.StatusTooManyRequests:
		status.Status = "error"
		status.Code = wosErrorRateLimited
		status.Message = "The " + enrichmentProviderName(provider) + " lookup limit has been reached."
		status.Retryable = true
		status.RetryAfter = safeRetryAfter(response.Header.Get("Retry-After"))
	default:
		status.Status = "error"
		status.Code = wosErrorUpstream
		status.Message = enrichmentUnavailableMessage(provider)
		status.Retryable = response.StatusCode >= http.StatusInternalServerError
	}
	return status
}

func enrichmentUpstreamStatus(provider, message string) models.EnrichmentSourceStatus {
	return models.EnrichmentSourceStatus{
		Provider:  provider,
		Status:    "error",
		Code:      wosErrorUpstream,
		Message:   message,
		Retryable: true,
	}
}

func enrichmentProviderName(provider string) string {
	if provider == enrichmentProviderUnpaywall {
		return "Unpaywall"
	}
	return "Crossref"
}

func enrichmentUnavailableMessage(provider string) string {
	return enrichmentProviderName(provider) + " is currently unavailable."
}

func enrichmentTimeoutMessage(provider string) string {
	return enrichmentProviderName(provider) + " took too long to respond."
}

func enrichmentInvalidMessage(provider string) string {
	return enrichmentProviderName(provider) + " returned an invalid response."
}

func enrichmentNotFoundMessage(provider string) string {
	if provider == enrichmentProviderUnpaywall {
		return "Unpaywall has no open-access record for this DOI."
	}
	return "This DOI was not found in Crossref."
}

func mapCrossrefValues(work crossrefWork, doi string) models.WosImportValues {
	message := work.Message
	values := models.WosImportValues{
		Doi:              stringValue(doi),
		Title:            stringValue(firstNonEmpty(message.Title)),
		LongJournalTitle: stringValue(firstNonEmpty(message.ContainerTitle)),
		Volume:           positiveIntValue(message.Volume),
		Issue:            positiveIntValue(message.Issue),
		Pages:            stringValue(message.Page),
		Isbn:             stringValue(firstNonEmpty(message.ISBN)),
		WebLink:          stringValue(message.URL),
	}

	values.Issn, values.EIssn = crossrefISSNs(message.ISSN, message.ISSNType)

	if subjects := nonEmptyStrings(message.Subject); len(subjects) > 0 {
		joined := strings.Join(subjects, "; ")
		values.Keywords = &joined
	}

	names := make([]string, 0, len(message.Author))
	for _, author := range message.Author {
		if name := crossrefAuthorName(author); name != "" {
			names = append(names, name)
		}
	}
	if len(names) > 0 {
		joined := strings.Join(names, "; ")
		count := len(names)
		values.AllAuthors = &joined
		values.AllAuthorsCount = &count
	}

	return values
}

// crossrefISSNs prefers the typed list, which distinguishes print from
// electronic. The untyped ISSN array is only a fallback and its order carries no
// meaning, so a single value is reported as the print ISSN rather than guessed.
func crossrefISSNs(issns []string, typed []crossrefISSN) (print, electronic *string) {
	for _, entry := range typed {
		switch strings.ToLower(strings.TrimSpace(entry.Type)) {
		case "print":
			if print == nil {
				print = stringValue(entry.Value)
			}
		case "electronic":
			if electronic == nil {
				electronic = stringValue(entry.Value)
			}
		}
	}
	if print != nil || electronic != nil {
		return print, electronic
	}
	return stringValue(firstNonEmpty(issns)), nil
}

func crossrefAuthorName(author crossrefAuthor) string {
	given := strings.TrimSpace(author.Given)
	family := strings.TrimSpace(author.Family)
	switch {
	case given != "" && family != "":
		return given + " " + family
	case family != "":
		return family
	case given != "":
		return given
	default:
		// Organisational authors carry a single `name` instead of given/family.
		return strings.TrimSpace(author.Name)
	}
}

func mapCrossrefAuthors(authors []crossrefAuthor) []models.WosAuthor {
	result := make([]models.WosAuthor, 0, len(authors))
	for _, author := range authors {
		name := crossrefAuthorName(author)
		if name == "" {
			continue
		}
		result = append(result, models.WosAuthor{WosDisplayName: name})
	}
	return result
}

// crossrefPublicationDate returns the year, an optional date and the precision
// actually available. Precision is never padded: a work Crossref only dates to a
// year must not gain a fabricated January 1st.
func crossrefPublicationDate(work crossrefWork) (year, date *string, precision string) {
	message := work.Message
	for _, candidate := range []crossrefDate{message.Published, message.PublishedPrint, message.PublishedOnline} {
		if len(candidate.DateParts) == 0 {
			continue
		}
		parts := candidate.DateParts[0]
		if len(parts) == 0 || parts[0] < 1 || parts[0] > 9999 {
			continue
		}
		yearText := strconv.Itoa(parts[0])
		switch {
		case len(parts) >= 3 && validMonth(parts[1]) && validDay(parts[2]):
			full := yearText + "-" + twoDigits(parts[1]) + "-" + twoDigits(parts[2])
			return &yearText, &full, "day"
		case len(parts) >= 2 && validMonth(parts[1]):
			month := yearText + "-" + twoDigits(parts[1])
			return &yearText, &month, "month"
		default:
			return &yearText, nil, "year"
		}
	}
	return nil, nil, ""
}

func mapUnpaywallOpenAccess(record unpaywallRecord) *models.EnrichmentOpenAccess {
	status := strings.TrimSpace(record.OaStatus)
	if status == "" {
		// An absent oa_status is unknown, not closed.
		if !record.IsOa && record.BestOaLocall == nil {
			return nil
		}
		status = "unknown"
	}
	openAccess := &models.EnrichmentOpenAccess{Status: status}
	if record.BestOaLocall != nil {
		openAccess.URL = strings.TrimSpace(record.BestOaLocall.URL)
		openAccess.PDFURL = strings.TrimSpace(record.BestOaLocall.URLForPdf)
	}
	return openAccess
}

func firstNonEmpty(values []string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func validMonth(month int) bool { return month >= 1 && month <= 12 }

func validDay(day int) bool { return day >= 1 && day <= 31 }

func twoDigits(value int) string {
	if value < 10 {
		return "0" + strconv.Itoa(value)
	}
	return strconv.Itoa(value)
}
