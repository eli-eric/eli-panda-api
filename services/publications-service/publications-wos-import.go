package publicationsservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	codebookmodels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/publications-service/models"
)

const (
	wosErrorInvalidDOI                 = "INVALID_DOI"
	wosErrorNotConfigured              = "WOS_NOT_CONFIGURED"
	wosErrorRecordNotFound             = "WOS_RECORD_NOT_FOUND"
	wosErrorRecordAmbiguous            = "WOS_RECORD_AMBIGUOUS"
	wosErrorAuthenticationFailed       = "WOS_AUTHENTICATION_FAILED"
	wosErrorRateLimited                = "WOS_RATE_LIMITED"
	wosErrorUpstream                   = "WOS_UPSTREAM_ERROR"
	wosErrorUpstreamTimeout            = "WOS_UPSTREAM_TIMEOUT"
	wosErrorInvalidResearcherID        = "INVALID_RESEARCHER_ID"
	wosErrorResearcherNotFound         = "RESEARCHER_NOT_FOUND"
	wosErrorResearcherIDConflict       = "RESEARCHER_ID_CONFLICT"
	wosErrorInternal                   = "INTERNAL_ERROR"
	wosRequestTimeout                  = 10 * time.Second
	wosResponseLimit             int64 = 5 << 20
)

var (
	doiPattern           = regexp.MustCompile(`(?i)^10\.\d{4,9}/\S+$`)
	researcherIDPattern  = regexp.MustCompile(`^[A-Z]{1,3}-\d{4}-\d{4}$`)
	errResearcherMissing = errors.New("researcher not found")
	errResearcherIDOwned = errors.New("researcher ID belongs to another researcher")
)

var wosImportableFields = []string{
	"doi",
	"title",
	"wosNumber",
	"longJournalTitle",
	"volume",
	"issue",
	"pages",
	"pagesCount",
	"yearOfPublication",
	"dateOfPublication",
	"issn",
	"eissn",
	"isbn",
	"webLink",
	"keywords",
	"allAuthors",
	"allAuthorsCount",
	"mediaTypeCb",
}

var wosUnavailableFields = []string{
	"abstract",
	"openAccessType",
	"publishingCountry",
	"oecdFord",
	"citeAs",
	"impactFactor",
	"quartilBasis",
	"quartil",
	"authorsDepartments",
	"grants",
	"code",
	"userCall",
	"userExperimentCb",
	"experimentalSystemCb",
	"note",
}

type wosResearcherRecord struct {
	UID                 string   `json:"uid"`
	FirstName           string   `json:"firstName"`
	LastName            string   `json:"lastName"`
	CurrentResearcherID string   `json:"currentResearcherId"`
	ResearcherIDs       []string `json:"researcherIds"`
}

type wosImportRepository interface {
	findExistingPublication(doi, currentPublicationUID string) (*models.WosExistingPublication, error)
	listResearchersForWos() ([]wosResearcherRecord, error)
	findMediaType(code, facilityCode string) (*codebookmodels.Codebook, error)
	rememberResearcherID(researcherUID, researcherID, userUID string, makePrimary bool) (rememberResearcherIDResult, error)
}

type publicationAPIError struct {
	StatusCode int
	Code       string
	Message    string
	Retryable  bool
	RetryAfter string
}

func (err *publicationAPIError) Error() string {
	return fmt.Sprintf("%s: %s", err.Code, err.Message)
}

func (err *publicationAPIError) response() models.PublicationAPIError {
	return models.PublicationAPIError{
		Code:      err.Code,
		Message:   err.Message,
		Retryable: err.Retryable,
	}
}

func normalizeDOI(raw string) (string, bool) {
	value := strings.TrimSpace(raw)
	lower := strings.ToLower(value)

	for _, prefix := range []string{
		"https://dx.doi.org/",
		"http://dx.doi.org/",
		"https://doi.org/",
		"http://doi.org/",
	} {
		if strings.HasPrefix(lower, prefix) {
			value = value[len(prefix):]
			if parsed, err := url.QueryUnescape(value); err == nil {
				value = parsed
			}
			if index := strings.IndexAny(value, "?#"); index >= 0 {
				value = value[:index]
			}
			break
		}
	}

	value = strings.TrimSpace(value)
	if len(value) >= 4 && strings.EqualFold(value[:4], "doi:") {
		value = strings.TrimSpace(value[4:])
	}

	value = strings.ToLower(strings.TrimSpace(value))
	if !doiPattern.MatchString(value) {
		return "", false
	}

	return value, true
}

func normalizeResearcherID(raw string) (string, bool) {
	value := strings.ToUpper(strings.TrimSpace(raw))
	return value, researcherIDPattern.MatchString(value)
}

// researcherIDYear reports the issue year encoded in a ResearcherID's suffix.
// Clarivate mints IDs as LETTERS-NNNN-YYYY, so the vintage travels with the
// value and needs no separate bookkeeping. Anything that does not parse is
// reported as unknown rather than as year zero, so an ID of unknown vintage
// never silently outranks — or is outranked by — one we can actually read.
func researcherIDYear(researcherID string) (int, bool) {
	value, valid := normalizeResearcherID(researcherID)
	if !valid {
		return 0, false
	}

	year, err := strconv.Atoi(value[len(value)-4:])
	if err != nil {
		return 0, false
	}

	return year, true
}

// newestResearcherID returns the ID with the highest issue year. It reports
// false when the answer is not defensible on its own: no readable candidate,
// or a tie between the newest ones. Those go to a human instead of being
// resolved by list order, because the wrong choice is exported to the
// government and nothing downstream would flag it.
func newestResearcherID(researcherIDs []string) (string, bool) {
	newest := ""
	newestYear := 0
	tied := false

	for _, candidate := range researcherIDs {
		year, ok := researcherIDYear(candidate)
		if !ok {
			continue
		}

		normalized, _ := normalizeResearcherID(candidate)
		switch {
		case newest == "" || year > newestYear:
			newest, newestYear, tied = normalized, year, false
		case year == newestYear && normalized != newest:
			tied = true
		}
	}

	if newest == "" || tied {
		return "", false
	}

	return newest, true
}

func (svc *PublicationsService) PreviewWosPublication(
	ctx context.Context,
	rawDOI,
	currentPublicationUID,
	facilityCode string,
) (models.WosPreviewResponse, error) {
	doi, valid := normalizeDOI(rawDOI)
	if !valid {
		return models.WosPreviewResponse{}, newPublicationAPIError(
			http.StatusBadRequest,
			wosErrorInvalidDOI,
			"Enter a valid DOI.",
			false,
		)
	}

	repository := svc.importRepository()
	if repository != nil {
		existingPublication, err := repository.findExistingPublication(doi, currentPublicationUID)
		if err != nil {
			return models.WosPreviewResponse{}, fmt.Errorf("find publication by DOI: %w", err)
		}
		if existingPublication != nil {
			return models.WosPreviewResponse{
				Status:              "already-exists",
				Doi:                 doi,
				ExistingPublication: existingPublication,
			}, nil
		}
	}

	_, hit, err := svc.fetchExactWosRecord(ctx, doi)
	if err != nil {
		return models.WosPreviewResponse{}, err
	}

	researchers := make([]wosResearcherRecord, 0)
	if repository != nil {
		researchers, err = repository.listResearchersForWos()
		if err != nil {
			return models.WosPreviewResponse{}, fmt.Errorf("list researchers for WoS matching: %w", err)
		}
	}

	values := mapWosImportValues(hit, doi)
	if repository != nil {
		mediaTypeCode := suggestMediaTypeCode(hit.WosTypes, hit.WosSourceTypes)
		if mediaTypeCode != "" {
			values.MediaTypeCb, err = repository.findMediaType(mediaTypeCode, facilityCode)
			if err != nil {
				return models.WosPreviewResponse{}, fmt.Errorf("find media type suggestion: %w", err)
			}
		}
	}

	return models.WosPreviewResponse{
		Status:                  "found",
		Doi:                     doi,
		Values:                  &values,
		Authors:                 matchWosAuthors(hit.WosNames.WosAuthors, researchers),
		MissingImportableFields: missingWosImportableFields(values),
		UnavailableFields:       append([]string(nil), wosUnavailableFields...),
	}, nil
}

func (svc *PublicationsService) RememberResearcherID(
	researcherUID,
	rawResearcherID,
	userUID string,
	makePrimary bool,
) (models.ResearcherIDsResponse, error) {
	researcherID, valid := normalizeResearcherID(rawResearcherID)
	if !valid {
		return models.ResearcherIDsResponse{}, newPublicationAPIError(
			http.StatusBadRequest,
			wosErrorInvalidResearcherID,
			"Enter a valid Web of Science ResearcherID.",
			false,
		)
	}

	repository := svc.importRepository()
	if repository == nil {
		return models.ResearcherIDsResponse{}, errors.New("WoS repository is unavailable")
	}

	remembered, err := repository.rememberResearcherID(researcherUID, researcherID, userUID, makePrimary)
	if errors.Is(err, errResearcherMissing) {
		return models.ResearcherIDsResponse{}, newPublicationAPIError(
			http.StatusNotFound,
			wosErrorResearcherNotFound,
			"The PANDA researcher could not be found.",
			false,
		)
	}
	if errors.Is(err, errResearcherIDOwned) {
		return models.ResearcherIDsResponse{}, newPublicationAPIError(
			http.StatusConflict,
			wosErrorResearcherIDConflict,
			"This ResearcherID is already linked to another PANDA researcher.",
			false,
		)
	}
	if err != nil {
		return models.ResearcherIDsResponse{}, fmt.Errorf("remember researcher ID: %w", err)
	}

	response := models.ResearcherIDsResponse{
		ResearcherIDs: remembered.researcherIDs,
		PrimaryID:     remembered.primaryID,
	}

	// Only worth raising when a newer ID exists and is not already current;
	// an ambiguous newest is left for a human rather than guessed at here.
	if newest, ok := newestResearcherID(remembered.researcherIDs); ok && newest != remembered.primaryID {
		response.SuggestedPrimary = newest
	}

	return response, nil
}

func (svc *PublicationsService) importRepository() wosImportRepository {
	if svc.wosRepository != nil {
		return svc.wosRepository
	}
	if svc.neo4jDriver == nil {
		return nil
	}
	return &neo4jWosImportRepository{driver: svc.neo4jDriver}
}

func (svc *PublicationsService) fetchExactWosRecord(
	ctx context.Context,
	doi string,
) (models.WosAPIResponse, models.WosHit, error) {
	response, err := svc.fetchWosResponse(ctx, doi)
	if err != nil {
		return models.WosAPIResponse{}, models.WosHit{}, err
	}

	exactMatches := make([]models.WosHit, 0, 1)
	for _, hit := range response.WosHits {
		hitDOI, valid := normalizeDOI(hit.WosIdentifiers.WosDOI)
		if valid && hitDOI == doi {
			exactMatches = append(exactMatches, hit)
		}
	}

	if len(exactMatches) == 0 {
		return models.WosAPIResponse{}, models.WosHit{}, newPublicationAPIError(
			http.StatusNotFound,
			wosErrorRecordNotFound,
			"This DOI was not found in Web of Science.",
			false,
		)
	}
	if len(exactMatches) > 1 {
		return models.WosAPIResponse{}, models.WosHit{}, newPublicationAPIError(
			http.StatusConflict,
			wosErrorRecordAmbiguous,
			"Web of Science returned more than one exact record for this DOI.",
			false,
		)
	}

	return response, exactMatches[0], nil
}

func (svc *PublicationsService) fetchWosResponse(
	ctx context.Context,
	doi string,
) (models.WosAPIResponse, error) {
	if strings.TrimSpace(svc.wosStarterApiUrl) == "" || strings.TrimSpace(svc.wosStarterApiKey) == "" {
		return models.WosAPIResponse{}, newPublicationAPIError(
			http.StatusServiceUnavailable,
			wosErrorNotConfigured,
			"Web of Science lookup is not configured.",
			false,
		)
	}

	endpoint, err := url.Parse(strings.TrimRight(svc.wosStarterApiUrl, "/") + "/documents")
	if err != nil || endpoint.Scheme == "" || endpoint.Host == "" {
		return models.WosAPIResponse{}, newPublicationAPIError(
			http.StatusBadGateway,
			wosErrorUpstream,
			"Web of Science is currently unavailable.",
			true,
		)
	}
	query := endpoint.Query()
	query.Set("db", "WOS")
	query.Set("q", "DO="+doi)
	query.Set("limit", "50")
	query.Set("page", "1")
	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return models.WosAPIResponse{}, fmt.Errorf("create WoS request: %w", err)
	}
	request.Header.Set("X-ApiKey", svc.wosStarterApiKey)
	request.Header.Set("Accept", "application/json")

	client := svc.wosHTTPClient
	if client == nil {
		client = newWosHTTPClient()
	}
	upstreamResponse, err := client.Do(request)
	if err != nil {
		var networkError net.Error
		if errors.Is(ctx.Err(), context.DeadlineExceeded) ||
			errors.As(err, &networkError) && networkError.Timeout() {
			return models.WosAPIResponse{}, newPublicationAPIError(
				http.StatusGatewayTimeout,
				wosErrorUpstreamTimeout,
				"Web of Science took too long to respond.",
				true,
			)
		}
		return models.WosAPIResponse{}, newPublicationAPIError(
			http.StatusBadGateway,
			wosErrorUpstream,
			"Web of Science is currently unavailable.",
			true,
		)
	}
	defer upstreamResponse.Body.Close()

	if upstreamResponse.StatusCode != http.StatusOK {
		return models.WosAPIResponse{}, mapWosStatusError(
			upstreamResponse.StatusCode,
			upstreamResponse.Header.Get("Retry-After"),
		)
	}

	var result models.WosAPIResponse
	decoder := json.NewDecoder(io.LimitReader(upstreamResponse.Body, wosResponseLimit))
	if err := decoder.Decode(&result); err != nil {
		return models.WosAPIResponse{}, newPublicationAPIError(
			http.StatusBadGateway,
			wosErrorUpstream,
			"Web of Science returned an invalid response.",
			true,
		)
	}

	return result, nil
}

func newWosHTTPClient() *http.Client {
	return &http.Client{
		Timeout: wosRequestTimeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func mapWosStatusError(status int, retryAfter string) *publicationAPIError {
	switch status {
	case http.StatusUnauthorized, http.StatusForbidden:
		return newPublicationAPIError(
			http.StatusServiceUnavailable,
			wosErrorAuthenticationFailed,
			"PANDA could not authenticate with Web of Science.",
			false,
		)
	case http.StatusTooManyRequests:
		err := newPublicationAPIError(
			http.StatusServiceUnavailable,
			wosErrorRateLimited,
			"The Web of Science lookup limit has been reached.",
			true,
		)
		err.RetryAfter = safeRetryAfter(retryAfter)
		return err
	default:
		return newPublicationAPIError(
			http.StatusBadGateway,
			wosErrorUpstream,
			"Web of Science is currently unavailable.",
			status >= http.StatusInternalServerError,
		)
	}
}

func safeRetryAfter(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if _, err := strconv.ParseUint(value, 10, 32); err == nil {
		return value
	}
	if _, err := http.ParseTime(value); err == nil {
		return value
	}
	return ""
}

func newPublicationAPIError(status int, code, message string, retryable bool) *publicationAPIError {
	return &publicationAPIError{
		StatusCode: status,
		Code:       code,
		Message:    message,
		Retryable:  retryable,
	}
}

func mapWosImportValues(hit models.WosHit, doi string) models.WosImportValues {
	values := models.WosImportValues{
		Doi:              stringValue(doi),
		Title:            stringValue(hit.WosTitle),
		WosNumber:        stringValue(hit.WosUID),
		LongJournalTitle: stringValue(hit.WosSource.WosSourceTitle),
		Volume:           positiveIntValue(hit.WosSource.WosVolume),
		Issue:            positiveIntValue(hit.WosSource.WosIssue),
		PagesCount:       positiveIntPointer(hit.WosSource.WosPages.WosCount),
		Issn:             stringValue(hit.WosIdentifiers.WosISSN),
		EIssn:            stringValue(hit.WosIdentifiers.WosEISSN),
		WebLink:          stringValue(hit.WosLinks.WosRecord),
	}

	pages := strings.TrimSpace(hit.WosSource.WosPages.WosRange)
	if pages == "" {
		pages = strings.TrimSpace(hit.WosSource.WosArticleNumber)
	}
	if pages == "" {
		begin := strings.TrimSpace(hit.WosSource.WosPages.WosBegin)
		end := strings.TrimSpace(hit.WosSource.WosPages.WosEnd)
		switch {
		case begin != "" && end != "" && begin != end:
			pages = begin + "-" + end
		case begin != "":
			pages = begin
		case end != "":
			pages = end
		}
	}
	values.Pages = stringValue(pages)

	if hit.WosSource.WosPublishYear > 0 {
		year := strconv.Itoa(hit.WosSource.WosPublishYear)
		values.YearOfPublication = &year
		if month := wosMonthNumber(hit.WosSource.WosPublishMonth); month != "" {
			date := year + "-" + month
			values.DateOfPublication = &date
		}
	}

	isbn := hit.WosIdentifiers.WosISBN
	if strings.TrimSpace(isbn) == "" {
		isbn = hit.WosIdentifiers.WosEISBN
	}
	values.Isbn = stringValue(isbn)

	keywords := nonEmptyStrings(hit.WosKeywords.WosAuthorKeywords)
	if len(keywords) > 0 {
		joined := strings.Join(keywords, "; ")
		values.Keywords = &joined
	}

	authorNames := make([]string, 0, len(hit.WosNames.WosAuthors))
	for _, author := range hit.WosNames.WosAuthors {
		name := strings.TrimSpace(author.WosDisplayName)
		if name == "" {
			name = strings.TrimSpace(author.WosStandard)
		}
		if name != "" {
			authorNames = append(authorNames, name)
		}
	}
	if len(authorNames) > 0 {
		joined := strings.Join(authorNames, "; ")
		count := len(authorNames)
		values.AllAuthors = &joined
		values.AllAuthorsCount = &count
	}

	return values
}

func missingWosImportableFields(values models.WosImportValues) []string {
	missing := make([]string, 0)
	for _, field := range wosImportableFields {
		present := false
		switch field {
		case "doi":
			present = values.Doi != nil
		case "title":
			present = values.Title != nil
		case "wosNumber":
			present = values.WosNumber != nil
		case "longJournalTitle":
			present = values.LongJournalTitle != nil
		case "volume":
			present = values.Volume != nil
		case "issue":
			present = values.Issue != nil
		case "pages":
			present = values.Pages != nil
		case "pagesCount":
			present = values.PagesCount != nil
		case "yearOfPublication":
			present = values.YearOfPublication != nil
		case "dateOfPublication":
			present = values.DateOfPublication != nil
		case "issn":
			present = values.Issn != nil
		case "eissn":
			present = values.EIssn != nil
		case "isbn":
			present = values.Isbn != nil
		case "webLink":
			present = values.WebLink != nil
		case "keywords":
			present = values.Keywords != nil
		case "allAuthors":
			present = values.AllAuthors != nil
		case "allAuthorsCount":
			present = values.AllAuthorsCount != nil
		case "mediaTypeCb":
			present = values.MediaTypeCb != nil
		}
		if !present {
			missing = append(missing, field)
		}
	}
	return missing
}

func matchWosAuthors(
	authors []models.WosAuthor,
	researchers []wosResearcherRecord,
) []models.WosImportAuthor {
	result := make([]models.WosImportAuthor, 0, len(authors))
	for index, author := range authors {
		match := models.WosAuthorMatch{Kind: "none", Candidates: make([]models.WosResearcherCandidate, 0)}
		if researcherID, valid := normalizeResearcherID(author.WosResearcherID); valid {
			match.Candidates = matchingResearchersByID(researchers, researcherID)
			switch len(match.Candidates) {
			case 1:
				match.Kind = "researcher-id"
			case 0:
			default:
				match.Kind = "ambiguous"
			}
		}

		if len(match.Candidates) == 0 {
			match.Candidates = matchingResearchersByName(researchers, author)
			switch len(match.Candidates) {
			case 1:
				match.Kind = "name"
			case 0:
				match.Kind = "none"
			default:
				match.Kind = "ambiguous"
			}
		}

		displayName := strings.TrimSpace(author.WosDisplayName)
		if displayName == "" {
			displayName = strings.TrimSpace(author.WosStandard)
		}
		result = append(result, models.WosImportAuthor{
			SourceIndex:  index,
			DisplayName:  displayName,
			WosStandard:  strings.TrimSpace(author.WosStandard),
			ResearcherID: strings.TrimSpace(author.WosResearcherID),
			Match:        match,
		})
	}
	return result
}

func matchingResearchersByID(researchers []wosResearcherRecord, researcherID string) []models.WosResearcherCandidate {
	result := make([]models.WosResearcherCandidate, 0)
	for _, researcher := range researchers {
		for _, candidateID := range researcher.ResearcherIDs {
			normalized, valid := normalizeResearcherID(candidateID)
			if valid && normalized == researcherID {
				result = appendUniqueResearcher(result, researcher)
				break
			}
		}
	}
	return result
}

func matchingResearchersByName(
	researchers []wosResearcherRecord,
	author models.WosAuthor,
) []models.WosResearcherCandidate {
	authorNames := nameVariants(author.WosDisplayName, author.WosStandard)
	if len(authorNames) == 0 {
		return []models.WosResearcherCandidate{}
	}

	result := make([]models.WosResearcherCandidate, 0)
	for _, researcher := range researchers {
		firstLast := normalizePersonName(researcher.FirstName + " " + researcher.LastName)
		lastFirst := normalizePersonName(researcher.LastName + " " + researcher.FirstName)
		if _, matches := authorNames[firstLast]; matches && firstLast != "" {
			result = appendUniqueResearcher(result, researcher)
			continue
		}
		if _, matches := authorNames[lastFirst]; matches && lastFirst != "" {
			result = appendUniqueResearcher(result, researcher)
		}
	}
	return result
}

func nameVariants(values ...string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if normalized := normalizePersonName(value); normalized != "" {
			result[normalized] = struct{}{}
		}
		if comma := strings.Index(value, ","); comma >= 0 {
			reordered := strings.TrimSpace(value[comma+1:]) + " " + strings.TrimSpace(value[:comma])
			if normalized := normalizePersonName(reordered); normalized != "" {
				result[normalized] = struct{}{}
			}
		}
	}
	return result
}

func normalizePersonName(value string) string {
	var builder strings.Builder
	lastWasSpace := true
	for _, character := range strings.ToLower(strings.TrimSpace(value)) {
		if unicode.IsLetter(character) || unicode.IsDigit(character) {
			builder.WriteRune(character)
			lastWasSpace = false
			continue
		}
		if !lastWasSpace {
			builder.WriteByte(' ')
			lastWasSpace = true
		}
	}
	return strings.TrimSpace(builder.String())
}

func appendUniqueResearcher(
	result []models.WosResearcherCandidate,
	researcher wosResearcherRecord,
) []models.WosResearcherCandidate {
	for _, existing := range result {
		if existing.Uid == researcher.UID {
			return result
		}
	}
	return append(result, models.WosResearcherCandidate{
		Uid:                 researcher.UID,
		FirstName:           researcher.FirstName,
		LastName:            researcher.LastName,
		CurrentResearcherID: researcher.CurrentResearcherID,
	})
}

func suggestMediaTypeCode(types, sourceTypes []string) string {
	values := append(append([]string(nil), types...), sourceTypes...)
	normalized := strings.ToLower(strings.Join(values, " "))
	switch {
	case strings.Contains(normalized, "book chapter"):
		return "C"
	case strings.Contains(normalized, "proceedings"), strings.Contains(normalized, "conference paper"):
		return "D"
	case strings.Contains(normalized, "article"), strings.Contains(normalized, "review"):
		return "J"
	default:
		return ""
	}
}

func wosMonthNumber(raw string) string {
	months := map[string]string{
		"JAN": "01", "FEB": "02", "MAR": "03", "APR": "04",
		"MAY": "05", "JUN": "06", "JUL": "07", "AUG": "08",
		"SEP": "09", "OCT": "10", "NOV": "11", "DEC": "12",
	}
	value := strings.ToUpper(strings.TrimSpace(raw))
	if len(value) >= 3 {
		return months[value[:3]]
	}
	return ""
}

func stringValue(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func positiveIntValue(value string) *int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return nil
	}
	return &parsed
}

func positiveIntPointer(value int) *int {
	if value <= 0 {
		return nil
	}
	return &value
}

func nonEmptyStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
