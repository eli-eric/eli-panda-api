package publicationsservice

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"panda/apigateway/services/publications-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newEnrichmentTestService points every provider at test doubles. Any provider
// left without a URL is simply unreachable, which is itself a case worth
// exercising.
func newEnrichmentTestService(
	crossrefURL, wosURL, unpaywallURL, unpaywallEmail string,
	repository wosImportRepository,
) *PublicationEnrichmentService {
	publications := &PublicationsService{
		wosStarterApiUrl: wosURL,
		wosStarterApiKey: "test-key",
		wosHTTPClient:    &http.Client{Timeout: 2 * time.Second},
		wosRepository:    repository,
	}
	return NewPublicationEnrichmentService(publications, PublicationEnrichmentConfig{
		CrossrefURL:    crossrefURL,
		CrossrefEmail:  "library@example.org",
		UnpaywallURL:   unpaywallURL,
		UnpaywallEmail: unpaywallEmail,
		HTTPClient:     &http.Client{Timeout: 2 * time.Second},
	})
}

func jsonServer(t *testing.T, status int, body string, headers map[string]string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server
}

func sourceByProvider(t *testing.T, sources []models.EnrichmentSourceStatus, provider string) models.EnrichmentSourceStatus {
	t.Helper()
	for _, source := range sources {
		if source.Provider == provider {
			return source
		}
	}
	t.Fatalf("no source reported for provider %q", provider)
	return models.EnrichmentSourceStatus{}
}

const crossrefFullBody = `{"message":{
	"title":["Relativistic electron beams from a laser wakefield"],
	"container-title":["Physical Review Letters"],
	"volume":"126","issue":"12","page":"123601","type":"journal-article",
	"URL":"https://doi.org/10.1103/physrevlett.126.123601",
	"subject":["General Physics and Astronomy","Optics"],
	"ISSN":["0031-9007","1079-7114"],
	"issn-type":[{"value":"0031-9007","type":"print"},{"value":"1079-7114","type":"electronic"}],
	"author":[
		{"given":"Jana","family":"Nováková","ORCID":"http://orcid.org/0000-0002-1825-0097"},
		{"given":"Petr","family":"Svoboda"}
	],
	"created":{"date-parts":[[2019,1,1]]},
	"published":{"date-parts":[[2025,3,17]]}
}}`

func TestFetchCrossrefMapsBibliographicMetadata(t *testing.T) {
	server := jsonServer(t, http.StatusOK, crossrefFullBody, nil)
	service := newEnrichmentTestService(server.URL, "", "", "", &stubWosImportRepository{})

	source := service.fetchCrossref(context.Background(), "10.1103/physrevlett.126.123601")

	require.Equal(t, "ok", source.status.Status)
	assert.Equal(t, enrichmentProviderCrossref, source.status.Provider)
	assert.NotEmpty(t, source.status.RetrievedAt)

	require.NotNil(t, source.values.Title)
	assert.Equal(t, "Relativistic electron beams from a laser wakefield", *source.values.Title)
	require.NotNil(t, source.values.LongJournalTitle)
	assert.Equal(t, "Physical Review Letters", *source.values.LongJournalTitle)
	require.NotNil(t, source.values.Volume)
	assert.Equal(t, 126, *source.values.Volume)
	require.NotNil(t, source.values.Pages)
	assert.Equal(t, "123601", *source.values.Pages)

	// The typed ISSN list distinguishes print from electronic; a bare list would
	// only have been guessed at by position.
	require.NotNil(t, source.values.Issn)
	assert.Equal(t, "0031-9007", *source.values.Issn)
	require.NotNil(t, source.values.EIssn)
	assert.Equal(t, "1079-7114", *source.values.EIssn)

	require.NotNil(t, source.values.AllAuthorsCount)
	assert.Equal(t, 2, *source.values.AllAuthorsCount)

	require.Len(t, source.authors, 2)
	assert.Equal(t, "Jana Nováková", source.authors[0].WosDisplayName)
	assert.Equal(t, "http://orcid.org/0000-0002-1825-0097", source.authors[0].Orcid)
	assert.Empty(t, source.authors[1].Orcid)

	// journal-article maps onto the same media-type code WoS resolves.
	assert.Equal(t, "J", source.mediaTypeCode)
}

func TestFetchCrossrefUsesPublicationDatesNotTheCreatedDate(t *testing.T) {
	// `created` is when Crossref received the metadata, six years before
	// publication here. Using it would file the paper under the wrong year.
	server := jsonServer(t, http.StatusOK, crossrefFullBody, nil)
	service := newEnrichmentTestService(server.URL, "", "", "", &stubWosImportRepository{})

	source := service.fetchCrossref(context.Background(), "10.1103/x")

	require.NotNil(t, source.values.YearOfPublication)
	assert.Equal(t, "2025", *source.values.YearOfPublication)
	require.NotNil(t, source.values.DateOfPublication)
	assert.Equal(t, "2025-03-17", *source.values.DateOfPublication)
	assert.Equal(t, "day", source.datePrecision)
}

func TestFetchCrossrefPreservesPartialDatePrecision(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		body      string
		year      string
		date      string
		precision string
	}{
		{
			name:      "year only",
			body:      `{"message":{"published":{"date-parts":[[2023]]}}}`,
			year:      "2023",
			precision: "year",
		},
		{
			name:      "year and month",
			body:      `{"message":{"published":{"date-parts":[[2023,7]]}}}`,
			year:      "2023",
			date:      "2023-07",
			precision: "month",
		},
		{
			name:      "falls back to print then online",
			body:      `{"message":{"published-print":{"date-parts":[[2022,11,2]]}}}`,
			year:      "2022",
			date:      "2022-11-02",
			precision: "day",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			server := jsonServer(t, http.StatusOK, testCase.body, nil)
			service := newEnrichmentTestService(server.URL, "", "", "", &stubWosImportRepository{})

			source := service.fetchCrossref(context.Background(), "10.1/x")

			require.NotNil(t, source.values.YearOfPublication)
			assert.Equal(t, testCase.year, *source.values.YearOfPublication)
			assert.Equal(t, testCase.precision, source.datePrecision)
			if testCase.date == "" {
				// A year-only record must not gain an invented January 1st.
				assert.Nil(t, source.values.DateOfPublication)
				return
			}
			require.NotNil(t, source.values.DateOfPublication)
			assert.Equal(t, testCase.date, *source.values.DateOfPublication)
		})
	}
}

func TestFetchCrossrefReportsUpstreamFailuresWithoutInventingData(t *testing.T) {
	repository := &stubWosImportRepository{}

	t.Run("unknown DOI", func(t *testing.T) {
		server := jsonServer(t, http.StatusNotFound, `{"message":"not found"}`, nil)
		source := newEnrichmentTestService(server.URL, "", "", "", repository).
			fetchCrossref(context.Background(), "10.1/missing")
		assert.Equal(t, "not-found", source.status.Status)
		assert.Nil(t, source.values.Title)
	})

	t.Run("rate limited", func(t *testing.T) {
		server := jsonServer(t, http.StatusTooManyRequests, `{}`, map[string]string{"Retry-After": "120"})
		source := newEnrichmentTestService(server.URL, "", "", "", repository).
			fetchCrossref(context.Background(), "10.1/busy")
		assert.Equal(t, "error", source.status.Status)
		assert.Equal(t, wosErrorRateLimited, source.status.Code)
		assert.True(t, source.status.Retryable)
		assert.Equal(t, "120", source.status.RetryAfter)
	})

	t.Run("malformed response", func(t *testing.T) {
		server := jsonServer(t, http.StatusOK, `{"message":`, nil)
		source := newEnrichmentTestService(server.URL, "", "", "", repository).
			fetchCrossref(context.Background(), "10.1/broken")
		assert.Equal(t, "error", source.status.Status)
		assert.Contains(t, source.status.Message, "invalid response")
	})

	t.Run("timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		}))
		defer server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		source := newEnrichmentTestService(server.URL, "", "", "", repository).fetchCrossref(ctx, "10.1/slow")
		assert.Equal(t, "error", source.status.Status)
		assert.True(t, source.status.Retryable)
	})
}

func TestFetchUnpaywallStaysDisabledWithoutAContact(t *testing.T) {
	// Unpaywall requires an email on every request, so an unconfigured contact
	// is reported as such rather than as a failed or empty lookup.
	service := newEnrichmentTestService("", "", "https://unpaywall.invalid", "", &stubWosImportRepository{})

	source := service.fetchUnpaywall(context.Background(), "10.1/x")

	assert.Equal(t, "not-configured", source.status.Status)
	assert.Equal(t, wosErrorNotConfigured, source.status.Code)
	assert.Nil(t, source.openAccess)
}

func TestFetchUnpaywallMapsOpenAccessState(t *testing.T) {
	body := `{"oa_status":"gold","is_oa":true,"best_oa_location":{"url":"https://example.org/article","url_for_pdf":"https://example.org/article.pdf"}}`
	server := jsonServer(t, http.StatusOK, body, nil)
	service := newEnrichmentTestService("", "", server.URL, "library@example.org", &stubWosImportRepository{})

	source := service.fetchUnpaywall(context.Background(), "10.1/x")

	require.Equal(t, "ok", source.status.Status)
	require.NotNil(t, source.openAccess)
	assert.Equal(t, "gold", source.openAccess.Status)
	assert.Equal(t, "https://example.org/article", source.openAccess.URL)
	assert.Equal(t, "https://example.org/article.pdf", source.openAccess.PDFURL)

	// Unpaywall contributes no bibliographic values.
	assert.Nil(t, source.values.Title)
	assert.Empty(t, source.authors)
}

func TestFetchUnpaywallKeepsAbsentMetadataUnknownRatherThanClosed(t *testing.T) {
	server := jsonServer(t, http.StatusOK, `{}`, nil)
	service := newEnrichmentTestService("", "", server.URL, "library@example.org", &stubWosImportRepository{})

	source := service.fetchUnpaywall(context.Background(), "10.1/x")

	require.Equal(t, "ok", source.status.Status)
	assert.Nil(t, source.openAccess, "no open-access record is not the same as closed access")
}

func TestPreviewSurvivesAProviderBeingDown(t *testing.T) {
	crossref := jsonServer(t, http.StatusOK, crossrefFullBody, nil)
	unpaywall := jsonServer(t, http.StatusServiceUnavailable, `{}`, nil)
	repository := &stubWosImportRepository{
		researchers: []wosResearcherRecord{
			{UID: "r-1", FirstName: "Jana", LastName: "Nováková", Orcids: []string{"0000-0002-1825-0097"}},
		},
	}
	// The WoS URL is empty, so that provider reports itself unconfigured.
	service := newEnrichmentTestService(crossref.URL, "", unpaywall.URL, "library@example.org", repository)

	result, err := service.Preview(context.Background(), "10.1103/PhysRevLett.126.123601", "", "B")
	require.NoError(t, err)

	// One healthy provider is enough to return a useful preview.
	assert.Equal(t, "found", result.Status)
	require.NotNil(t, result.Values)
	require.NotNil(t, result.Values.Title)

	// Every provider's outcome is reported so the editor can see what was and
	// was not consulted.
	require.Len(t, result.Sources, 3)
	assert.Equal(t, "ok", sourceByProvider(t, result.Sources, enrichmentProviderCrossref).Status)
	assert.Equal(t, "not-configured", sourceByProvider(t, result.Sources, "wos-starter").Status)
	assert.Equal(t, "error", sourceByProvider(t, result.Sources, enrichmentProviderUnpaywall).Status)

	// A failed Unpaywall leaves open access unknown rather than guessed.
	assert.Nil(t, result.OpenAccess)

	// Provenance names where each accepted value came from.
	assert.Equal(t, enrichmentProviderCrossref, result.Provenance["title"].Provider)

	// The ORCID carried by Crossref resolves the author to a stored researcher.
	require.Len(t, result.Authors, 2)
	assert.Equal(t, "orcid", result.Authors[0].Match.Kind)
	require.Len(t, result.Authors[0].Match.Candidates, 1)
	assert.Equal(t, "r-1", result.Authors[0].Match.Candidates[0].Uid)
	assert.Equal(t, "none", result.Authors[1].Match.Kind)
}

func TestPreviewNeverWrites(t *testing.T) {
	crossref := jsonServer(t, http.StatusOK, crossrefFullBody, nil)
	repository := &stubWosImportRepository{}
	service := newEnrichmentTestService(crossref.URL, "", "", "", repository)

	_, err := service.Preview(context.Background(), "10.1103/PhysRevLett.126.123601", "existing-uid", "B")
	require.NoError(t, err)

	// The repository exposes only reads; the preview must touch nothing else.
	require.Len(t, repository.findExistingCalls, 1)
	assert.Equal(t, "10.1103/physrevlett.126.123601", repository.findExistingCalls[0].doi)
	assert.Equal(t, "existing-uid", repository.findExistingCalls[0].currentPublicationUID)
	assert.Equal(t, 1, repository.listResearchersCalls)
}

func TestPreviewSkipsProvidersForAnAlreadyImportedDOI(t *testing.T) {
	crossrefCalls := 0
	crossref := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		crossrefCalls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(crossrefFullBody))
	}))
	defer crossref.Close()

	repository := &stubWosImportRepository{
		existingPublication: &models.WosExistingPublication{Uid: "pub-1", Code: "PUB-1", Title: "Already here"},
	}
	service := newEnrichmentTestService(crossref.URL, "", "", "", repository)

	result, err := service.Preview(context.Background(), "10.1103/PhysRevLett.126.123601", "", "B")
	require.NoError(t, err)

	assert.Equal(t, "already-exists", result.Status)
	require.NotNil(t, result.ExistingPublication)
	assert.Equal(t, "pub-1", result.ExistingPublication.Uid)
	// No provider quota is spent on a DOI PANDA already holds.
	assert.Equal(t, 0, crossrefCalls)
	require.Len(t, result.Sources, 3)
	for _, source := range result.Sources {
		assert.Equal(t, "skipped", source.Status)
	}
}

func TestPreviewRejectsAnInvalidDOI(t *testing.T) {
	service := newEnrichmentTestService("", "", "", "", &stubWosImportRepository{})

	_, err := service.Preview(context.Background(), "definitely not a doi", "", "B")

	require.Error(t, err)
	assert.Contains(t, err.Error(), wosErrorInvalidDOI)
}
