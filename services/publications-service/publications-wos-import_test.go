package publicationsservice

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	codebookmodels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/publications-service/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type findExistingCall struct {
	doi                   string
	currentPublicationUID string
}

type findMediaTypeCall struct {
	code         string
	facilityCode string
}

type rememberResearcherIDCall struct {
	researcherUID string
	researcherID  string
	userUID       string
	makePrimary   bool
}

type stubWosImportRepository struct {
	existingPublication *models.WosExistingPublication
	findExistingErr     error
	findExistingCalls   []findExistingCall

	researchers          []wosResearcherRecord
	listResearchersErr   error
	listResearchersCalls int

	mediaType          *codebookmodels.Codebook
	findMediaTypeErr   error
	findMediaTypeCalls []findMediaTypeCall

	rememberedResearcherIDs   []string
	rememberedPrimaryID       string
	rememberResearcherIDErr   error
	rememberResearcherIDCalls []rememberResearcherIDCall
}

func (repo *stubWosImportRepository) findExistingPublication(
	doi,
	currentPublicationUID string,
) (*models.WosExistingPublication, error) {
	repo.findExistingCalls = append(repo.findExistingCalls, findExistingCall{
		doi:                   doi,
		currentPublicationUID: currentPublicationUID,
	})
	return repo.existingPublication, repo.findExistingErr
}

func (repo *stubWosImportRepository) listResearchersForWos() ([]wosResearcherRecord, error) {
	repo.listResearchersCalls++
	return repo.researchers, repo.listResearchersErr
}

func (repo *stubWosImportRepository) findMediaType(
	code,
	facilityCode string,
) (*codebookmodels.Codebook, error) {
	repo.findMediaTypeCalls = append(repo.findMediaTypeCalls, findMediaTypeCall{
		code:         code,
		facilityCode: facilityCode,
	})
	return repo.mediaType, repo.findMediaTypeErr
}

func (repo *stubWosImportRepository) rememberResearcherID(
	researcherUID,
	researcherID,
	userUID string,
	makePrimary bool,
) (rememberResearcherIDResult, error) {
	repo.rememberResearcherIDCalls = append(
		repo.rememberResearcherIDCalls,
		rememberResearcherIDCall{
			researcherUID: researcherUID,
			researcherID:  researcherID,
			userUID:       userUID,
			makePrimary:   makePrimary,
		},
	)
	if repo.rememberResearcherIDErr != nil {
		return rememberResearcherIDResult{}, repo.rememberResearcherIDErr
	}
	primary := repo.rememberedPrimaryID
	if primary == "" && len(repo.rememberedResearcherIDs) > 0 {
		primary = repo.rememberedResearcherIDs[0]
	}
	return rememberResearcherIDResult{
		researcherIDs: repo.rememberedResearcherIDs,
		primaryID:     primary,
	}, nil
}

func TestPreviewWosPublicationRejectsInvalidDOI(t *testing.T) {
	invalidDOIs := []string{
		"",
		"not-a-doi",
		"10.12/too-short-prefix",
		"https://doi.org/not-a-doi",
	}

	for _, doi := range invalidDOIs {
		t.Run(doi, func(t *testing.T) {
			service := &PublicationsService{wosRepository: &stubWosImportRepository{}}
			response := performPreviewRequest(t, service, "/v1/publications/wos-preview?doi="+doi)

			assertPublicationAPIError(
				t,
				response,
				http.StatusBadRequest,
				wosErrorInvalidDOI,
				false,
			)
		})
	}
}

func TestPreviewWosPublicationNormalizesCommonDOIForms(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{name: "doi prefix", raw: "  DOI:10.1000/ABC.Def  "},
		{name: "doi url", raw: "https://doi.org/10.1000/ABC.Def?source=clipboard"},
		{name: "legacy doi url and escaped slash", raw: "http://dx.doi.org/10.1000%2FABC.Def#section"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requestedDOI string
			upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				requestedDOI = strings.TrimPrefix(request.URL.Query().Get("q"), "DO=")
				writeWosResponse(t, response, `{
					"metadata":{"total":1,"page":1,"limit":50},
					"hits":[{
						"uid":"WOS:1",
						"title":"Normalized DOI",
						"identifiers":{"doi":"10.1000/abc.def"}
					}]
				}`)
			}))
			defer upstream.Close()

			service := newWosTestService(upstream.URL, &stubWosImportRepository{})
			preview, err := service.PreviewWosPublication(context.Background(), test.raw, "", "B")

			require.NoError(t, err)
			assert.Equal(t, "10.1000/abc.def", requestedDOI)
			assert.Equal(t, "10.1000/abc.def", preview.Doi)
			require.NotNil(t, preview.Values)
			assert.Equal(t, "10.1000/abc.def", valueOrEmpty(preview.Values.Doi))
		})
	}
}

func TestPreviewWosPublicationReturnsExistingPublicationBeforeCallingWos(t *testing.T) {
	repo := &stubWosImportRepository{existingPublication: &models.WosExistingPublication{
		Uid:   "publication-1",
		Code:  "ELI-2026-1",
		Title: "Existing laser paper",
		Doi:   "10.17077/etd.g638o927",
	}}
	upstreamCalls := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		upstreamCalls++
		response.WriteHeader(http.StatusInternalServerError)
	}))
	defer upstream.Close()

	service := newWosTestService(upstream.URL, repo)
	response := performPreviewRequest(
		t,
		service,
		"/v1/publications/wos-preview?doi=10.17077%2FETD.G638O927&currentPublicationUid=current-publication",
	)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, 0, upstreamCalls, "a duplicate check must not consume WoS quota")
	require.Len(t, repo.findExistingCalls, 1)
	assert.Equal(t, findExistingCall{
		doi:                   "10.17077/etd.g638o927",
		currentPublicationUID: "current-publication",
	}, repo.findExistingCalls[0])
	assert.Empty(t, repo.rememberResearcherIDCalls, "a preview must not persist researcher IDs")

	var preview models.WosPreviewResponse
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &preview))
	assert.Equal(t, "already-exists", preview.Status)
	assert.Equal(t, "10.17077/etd.g638o927", preview.Doi)
	assert.Equal(t, repo.existingPublication, preview.ExistingPublication)
}

func TestPreviewWosPublicationReportsMissingWosConfiguration(t *testing.T) {
	service := &PublicationsService{wosRepository: &stubWosImportRepository{}}
	response := performPreviewRequest(
		t,
		service,
		"/v1/publications/wos-preview?doi=10.17077%2Fetd.g638o927",
	)

	assertPublicationAPIError(
		t,
		response,
		http.StatusServiceUnavailable,
		wosErrorNotConfigured,
		false,
	)
}

func TestPreviewWosPublicationUsesOfficialStarterRequestParameters(t *testing.T) {
	const requestedDOI = "10.17077/etd.g638o927"

	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		assert.Equal(t, "/documents", request.URL.Path)
		assert.Equal(t, "WOS", request.URL.Query().Get("db"))
		assert.Equal(t, "DO="+requestedDOI, request.URL.Query().Get("q"))
		assert.Empty(t, request.URL.Query().Get("detail"), "full detail is the default and should be omitted")
		assert.Equal(t, "1", request.URL.Query().Get("page"))
		assert.Equal(t, "50", request.URL.Query().Get("limit"))
		assert.Equal(t, "test-key", request.Header.Get("X-ApiKey"))
		assert.Equal(t, "application/json", request.Header.Get("Accept"))

		writeWosResponse(t, response, `{
			"metadata":{"total":1,"page":1,"limit":50},
			"hits":[{
				"uid":"WOS:000123456789",
				"title":"Laser metadata",
				"identifiers":{"doi":"10.17077/ETD.G638O927"}
			}]
		}`)
	}))
	defer upstream.Close()

	service := newWosTestService(upstream.URL, &stubWosImportRepository{})
	response := performPreviewRequest(
		t,
		service,
		"/v1/publications/wos-preview?doi="+requestedDOI,
	)

	assert.Equal(t, http.StatusOK, response.Code)
	var preview models.WosPreviewResponse
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &preview))
	assert.Equal(t, "found", preview.Status)
	assert.Equal(t, requestedDOI, preview.Doi)
	require.NotNil(t, preview.Values)
	assert.Equal(t, requestedDOI, valueOrEmpty(preview.Values.Doi))
	assert.Equal(t, "Laser metadata", valueOrEmpty(preview.Values.Title))
}

func TestPreviewWosPublicationRequiresOneExactDOIMatch(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "no result",
			responseBody: `{
				"metadata":{"total":0,"page":1,"limit":50},
				"hits":[]
			}`,
			expectedStatus: http.StatusNotFound,
			expectedCode:   wosErrorRecordNotFound,
		},
		{
			name: "only a non-exact result",
			responseBody: `{
				"metadata":{"total":1,"page":1,"limit":50},
				"hits":[{"uid":"WOS:OTHER","identifiers":{"doi":"10.1000/other"}}]
			}`,
			expectedStatus: http.StatusNotFound,
			expectedCode:   wosErrorRecordNotFound,
		},
		{
			name: "more than one exact result",
			responseBody: `{
				"metadata":{"total":2,"page":1,"limit":50},
				"hits":[
					{"uid":"WOS:ONE","identifiers":{"doi":"10.1000/exact"}},
					{"uid":"WOS:TWO","identifiers":{"doi":"HTTPS://DOI.ORG/10.1000/EXACT"}}
				]
			}`,
			expectedStatus: http.StatusConflict,
			expectedCode:   wosErrorRecordAmbiguous,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				writeWosResponse(t, response, test.responseBody)
			}))
			defer upstream.Close()

			service := newWosTestService(upstream.URL, &stubWosImportRepository{})
			response := performPreviewRequest(
				t,
				service,
				"/v1/publications/wos-preview?doi=10.1000%2Fexact",
			)

			assertPublicationAPIError(
				t,
				response,
				test.expectedStatus,
				test.expectedCode,
				false,
			)
		})
	}
}

func TestPreviewWosPublicationMapsStarterFieldsAndAuthorMatches(t *testing.T) {
	mediaType := &codebookmodels.Codebook{
		UID:  "media-journal",
		Name: "Journal article",
		Code: "J",
	}
	repo := &stubWosImportRepository{
		mediaType: mediaType,
		researchers: []wosResearcherRecord{
			{
				UID:           "researcher-by-id",
				FirstName:     "Jane",
				LastName:      "Doe",
				ResearcherIDs: []string{" aab-1234-2020 ", "AAB-9999-2024"},
			},
			{UID: "researcher-by-name", FirstName: "John", LastName: "Smith"},
			{UID: "alex-one", FirstName: "Alex", LastName: "Brown"},
			{UID: "alex-two", FirstName: "Alex", LastName: "Brown"},
		},
	}

	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		writeWosResponse(t, response, `{
			"metadata":{"total":1,"page":1,"limit":50},
			"hits":[{
				"uid":"WOS:000987654321",
				"title":"Laser-driven metadata",
				"types":["Article"],
				"sourceTypes":["Journal"],
				"source":{
					"sourceTitle":"Journal of Laser Tests",
					"publishYear":2025,
					"publishMonth":"SEP",
					"volume":"42",
					"issue":"7",
					"articleNumber":"e0123",
					"pages":{"range":"101-109","begin":"101","end":"109","count":9}
				},
				"names":{"authors":[
					{"displayName":"Doe, Jane","wosStandard":"Doe J","researcherId":"AAB-1234-2020"},
					{"displayName":"Smith, John","wosStandard":"Smith J","researcherId":""},
					{"displayName":"Brown, Alex","wosStandard":"Brown A","researcherId":""},
					{"displayName":"No Match","wosStandard":"No Match","researcherId":""}
				]},
				"identifiers":{
					"doi":"10.1000/LASER.1",
					"issn":"1234-5678",
					"eissn":"8765-4321",
					"isbn":"978-1-23456-789-0"
				},
				"keywords":{"authorKeywords":[" lasers ","", "metadata"]},
				"links":{"record":"https://www.webofscience.com/wos/woscc/full-record/WOS:000987654321"}
			}]
		}`)
	}))
	defer upstream.Close()

	service := newWosTestService(upstream.URL, repo)
	response := performPreviewRequest(
		t,
		service,
		"/v1/publications/wos-preview?doi=10.1000%2Flaser.1",
	)

	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var preview models.WosPreviewResponse
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &preview))
	assert.Equal(t, "found", preview.Status)
	require.NotNil(t, preview.Values)
	assert.Equal(t, "10.1000/laser.1", valueOrEmpty(preview.Values.Doi))
	assert.Equal(t, "Laser-driven metadata", valueOrEmpty(preview.Values.Title))
	assert.Equal(t, "WOS:000987654321", valueOrEmpty(preview.Values.WosNumber))
	assert.Equal(t, "Journal of Laser Tests", valueOrEmpty(preview.Values.LongJournalTitle))
	assert.Equal(t, 42, intValueOrZero(preview.Values.Volume))
	assert.Equal(t, 7, intValueOrZero(preview.Values.Issue))
	assert.Equal(t, "101-109", valueOrEmpty(preview.Values.Pages))
	assert.Equal(t, 9, intValueOrZero(preview.Values.PagesCount))
	assert.Equal(t, "2025", valueOrEmpty(preview.Values.YearOfPublication))
	assert.Equal(t, "2025-09", valueOrEmpty(preview.Values.DateOfPublication))
	assert.Equal(t, "1234-5678", valueOrEmpty(preview.Values.Issn))
	assert.Equal(t, "8765-4321", valueOrEmpty(preview.Values.EIssn))
	assert.Equal(t, "978-1-23456-789-0", valueOrEmpty(preview.Values.Isbn))
	assert.Equal(
		t,
		"https://www.webofscience.com/wos/woscc/full-record/WOS:000987654321",
		valueOrEmpty(preview.Values.WebLink),
	)
	assert.Equal(t, "lasers; metadata", valueOrEmpty(preview.Values.Keywords))
	assert.Equal(t, "Doe, Jane; Smith, John; Brown, Alex; No Match", valueOrEmpty(preview.Values.AllAuthors))
	assert.Equal(t, 4, intValueOrZero(preview.Values.AllAuthorsCount))
	assert.Equal(t, mediaType, preview.Values.MediaTypeCb)
	assert.Equal(t, []findMediaTypeCall{{code: "J", facilityCode: "B"}}, repo.findMediaTypeCalls)

	require.Len(t, preview.Authors, 4)
	assert.Equal(t, "researcher-id", preview.Authors[0].Match.Kind)
	require.Len(t, preview.Authors[0].Match.Candidates, 1)
	assert.Equal(t, "researcher-by-id", preview.Authors[0].Match.Candidates[0].Uid)
	assert.Equal(t, "name", preview.Authors[1].Match.Kind)
	require.Len(t, preview.Authors[1].Match.Candidates, 1)
	assert.Equal(t, "researcher-by-name", preview.Authors[1].Match.Candidates[0].Uid)
	assert.Equal(t, "ambiguous", preview.Authors[2].Match.Kind)
	assert.ElementsMatch(t, []string{"alex-one", "alex-two"}, researcherUIDs(preview.Authors[2].Match.Candidates))
	assert.Equal(t, "none", preview.Authors[3].Match.Kind)
	assert.Empty(t, preview.Authors[3].Match.Candidates)
	assert.NotContains(t, preview.MissingImportableFields, "doi")
	assert.NotContains(t, preview.MissingImportableFields, "mediaTypeCb")
	assert.Contains(t, preview.UnavailableFields, "abstract")
	assert.Contains(t, preview.UnavailableFields, "openAccessType")
	assert.Contains(t, preview.UnavailableFields, "publishingCountry")
	assert.Contains(t, preview.UnavailableFields, "oecdFord")
	assert.Empty(t, repo.rememberResearcherIDCalls, "matching in a preview must stay read-only")
}

func TestPreviewWosPublicationMarksDuplicateResearcherIDAsAmbiguous(t *testing.T) {
	repo := &stubWosImportRepository{researchers: []wosResearcherRecord{
		{UID: "first", FirstName: "First", LastName: "Person", ResearcherIDs: []string{"AAA-1111-2022"}},
		{UID: "second", FirstName: "Second", LastName: "Person", ResearcherIDs: []string{"AAA-1111-2022"}},
	}}
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		writeWosResponse(t, response, `{
			"metadata":{"total":1,"page":1,"limit":50},
			"hits":[{
				"uid":"WOS:ONE",
				"title":"Ambiguous author",
				"identifiers":{"doi":"10.1000/ambiguous-author"},
				"names":{"authors":[{
					"displayName":"Someone Else",
					"wosStandard":"Else S",
					"researcherId":"AAA-1111-2022"
				}]}
			}]
		}`)
	}))
	defer upstream.Close()

	preview, err := newWosTestService(upstream.URL, repo).PreviewWosPublication(
		context.Background(),
		"10.1000/ambiguous-author",
		"",
		"B",
	)

	require.NoError(t, err)
	require.Len(t, preview.Authors, 1)
	assert.Equal(t, "ambiguous", preview.Authors[0].Match.Kind)
	assert.ElementsMatch(t, []string{"first", "second"}, researcherUIDs(preview.Authors[0].Match.Candidates))
}

func TestPreviewWosPublicationReportsMissingImportableFields(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		writeWosResponse(t, response, `{
			"metadata":{"total":1,"page":1,"limit":50},
			"hits":[{"uid":"WOS:SPARSE","identifiers":{"doi":"10.1000/sparse"}}]
		}`)
	}))
	defer upstream.Close()

	preview, err := newWosTestService(upstream.URL, &stubWosImportRepository{}).PreviewWosPublication(
		context.Background(),
		"10.1000/sparse",
		"",
		"B",
	)

	require.NoError(t, err)
	assert.NotContains(t, preview.MissingImportableFields, "doi")
	for _, field := range []string{
		"title",
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
	} {
		assert.Contains(t, preview.MissingImportableFields, field)
	}
	assert.NotNil(t, preview.Authors)
	assert.NotNil(t, preview.UnavailableFields)
}

func TestPreviewWosPublicationMapsUpstreamFailuresToStableErrors(t *testing.T) {
	tests := []struct {
		name           string
		upstreamStatus int
		body           string
		retryAfter     string
		expectedStatus int
		expectedCode   string
		expectedRetry  bool
		expectedHeader string
	}{
		{
			name:           "unauthorized key",
			upstreamStatus: http.StatusUnauthorized,
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   wosErrorAuthenticationFailed,
			expectedRetry:  false,
		},
		{
			name:           "rate limited",
			upstreamStatus: http.StatusTooManyRequests,
			retryAfter:     "120",
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   wosErrorRateLimited,
			expectedRetry:  true,
			expectedHeader: "120",
		},
		{
			name:           "server failure",
			upstreamStatus: http.StatusInternalServerError,
			expectedStatus: http.StatusBadGateway,
			expectedCode:   wosErrorUpstream,
			expectedRetry:  true,
		},
		{
			name:           "malformed success response",
			upstreamStatus: http.StatusOK,
			body:           `{this-is-not-json`,
			expectedStatus: http.StatusBadGateway,
			expectedCode:   wosErrorUpstream,
			expectedRetry:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				if test.retryAfter != "" {
					response.Header().Set("Retry-After", test.retryAfter)
				}
				response.WriteHeader(test.upstreamStatus)
				_, _ = response.Write([]byte(test.body))
			}))
			defer upstream.Close()

			service := newWosTestService(upstream.URL, &stubWosImportRepository{})
			response := performPreviewRequest(
				t,
				service,
				"/v1/publications/wos-preview?doi=10.1000%2Fupstream",
			)

			assertPublicationAPIError(
				t,
				response,
				test.expectedStatus,
				test.expectedCode,
				test.expectedRetry,
			)
			assert.Equal(t, test.expectedHeader, response.Header().Get("Retry-After"))
		})
	}
}

func TestPreviewWosPublicationDropsUnsafeRetryAfterHeader(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Retry-After", "not-a-valid-retry-value")
		response.WriteHeader(http.StatusTooManyRequests)
	}))
	defer upstream.Close()

	response := performPreviewRequest(
		t,
		newWosTestService(upstream.URL, &stubWosImportRepository{}),
		"/v1/publications/wos-preview?doi=10.1000%2Frate-limit",
	)

	assertPublicationAPIError(
		t,
		response,
		http.StatusServiceUnavailable,
		wosErrorRateLimited,
		true,
	)
	assert.Empty(t, response.Header().Get("Retry-After"))
}

func TestPreviewWosPublicationReportsUpstreamTimeout(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		<-request.Context().Done()
	}))
	defer upstream.Close()

	service := newWosTestService(upstream.URL, &stubWosImportRepository{})
	service.wosHTTPClient = &http.Client{Timeout: 20 * time.Millisecond}
	response := performPreviewRequest(
		t,
		service,
		"/v1/publications/wos-preview?doi=10.1000%2Ftimeout",
	)

	assertPublicationAPIError(
		t,
		response,
		http.StatusGatewayTimeout,
		wosErrorUpstreamTimeout,
		true,
	)
}

func TestRememberResearcherIDNormalizesAndReturnsAllStoredIDs(t *testing.T) {
	repo := &stubWosImportRepository{
		rememberedResearcherIDs: []string{"AAA-1111-2020", "AAB-1234-2020"},
	}
	service := &PublicationsService{wosRepository: repo}
	response := performRememberResearcherIDRequest(
		t,
		service,
		"researcher-1",
		"user-1",
		`{"researcherId":"  aab-1234-2020  "}`,
	)

	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	assert.Equal(t, []rememberResearcherIDCall{{
		researcherUID: "researcher-1",
		researcherID:  "AAB-1234-2020",
		userUID:       "user-1",
	}}, repo.rememberResearcherIDCalls)
	var result models.ResearcherIDsResponse
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
	assert.Equal(t, repo.rememberedResearcherIDs, result.ResearcherIDs)
}

func TestRememberResearcherIDRejectsInvalidInputWithoutWriting(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "malformed json", body: `{`},
		{name: "missing researcher id", body: `{}`},
		{name: "invalid researcher id", body: `{"researcherId":"not-an-id"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &stubWosImportRepository{}
			response := performRememberResearcherIDRequest(
				t,
				&PublicationsService{wosRepository: repo},
				"researcher-1",
				"user-1",
				test.body,
			)

			assertPublicationAPIError(
				t,
				response,
				http.StatusBadRequest,
				wosErrorInvalidResearcherID,
				false,
			)
			assert.Empty(t, repo.rememberResearcherIDCalls)
		})
	}
}

func TestRememberResearcherIDMapsRepositoryOutcomes(t *testing.T) {
	tests := []struct {
		name           string
		repositoryErr  error
		expectedStatus int
		expectedCode   string
		retryable      bool
	}{
		{
			name:           "researcher missing",
			repositoryErr:  errResearcherMissing,
			expectedStatus: http.StatusNotFound,
			expectedCode:   wosErrorResearcherNotFound,
		},
		{
			name:           "researcher id owned by another researcher",
			repositoryErr:  errResearcherIDOwned,
			expectedStatus: http.StatusConflict,
			expectedCode:   wosErrorResearcherIDConflict,
		},
		{
			name:           "unexpected repository failure",
			repositoryErr:  errors.New("database offline"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   wosErrorInternal,
			retryable:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &stubWosImportRepository{rememberResearcherIDErr: test.repositoryErr}
			response := performRememberResearcherIDRequest(
				t,
				&PublicationsService{wosRepository: repo},
				"researcher-1",
				"user-1",
				`{"researcherId":"AAB-1234-2020"}`,
			)

			assertPublicationAPIError(
				t,
				response,
				test.expectedStatus,
				test.expectedCode,
				test.retryable,
			)
			require.Len(t, repo.rememberResearcherIDCalls, 1)
		})
	}
}

func newWosTestService(
	upstreamURL string,
	repository wosImportRepository,
) *PublicationsService {
	return &PublicationsService{
		wosStarterApiUrl: upstreamURL,
		wosStarterApiKey: "test-key",
		wosHTTPClient:    &http.Client{Timeout: time.Second},
		wosRepository:    repository,
	}
}

func performPreviewRequest(
	t *testing.T,
	service *PublicationsService,
	target string,
) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, target, nil)
	response := httptest.NewRecorder()
	ctx := e.NewContext(request, response)
	ctx.Set("facilityCode", "B")

	err := NewPublicationsHandlers(service).PreviewWosPublication()(ctx)
	require.NoError(t, err)
	return response
}

func performRememberResearcherIDRequest(
	t *testing.T,
	service *PublicationsService,
	researcherUID,
	userUID,
	body string,
) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodPut,
		"/v1/researcher/"+researcherUID+"/researcher-ids",
		strings.NewReader(body),
	)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	ctx := e.NewContext(request, response)
	ctx.SetParamNames("uid")
	ctx.SetParamValues(researcherUID)
	ctx.Set("userUID", userUID)

	err := NewPublicationsHandlers(service).RememberResearcherID()(ctx)
	require.NoError(t, err)
	return response
}

func assertPublicationAPIError(
	t *testing.T,
	response *httptest.ResponseRecorder,
	expectedStatus int,
	expectedCode string,
	expectedRetryable bool,
) {
	t.Helper()

	assert.Equal(t, expectedStatus, response.Code, response.Body.String())
	var apiError models.PublicationAPIError
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &apiError))
	assert.Equal(t, expectedCode, apiError.Code)
	assert.NotEmpty(t, apiError.Message)
	assert.Equal(t, expectedRetryable, apiError.Retryable)
}

func writeWosResponse(t *testing.T, response http.ResponseWriter, body string) {
	t.Helper()
	response.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	_, err := response.Write([]byte(body))
	assert.NoError(t, err)
}

func researcherUIDs(researchers []models.WosResearcherCandidate) []string {
	result := make([]string, 0, len(researchers))
	for _, researcher := range researchers {
		result = append(result, researcher.Uid)
	}
	return result
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func intValueOrZero(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
