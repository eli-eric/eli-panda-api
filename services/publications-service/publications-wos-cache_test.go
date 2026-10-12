package publicationsservice

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"panda/apigateway/services/publications-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ELIPANDA-501: the Clarivate Starter quota is per key per day, so repeat
// lookups of the same DOI must be served from the in-process cache, while
// misses stay uncached (the record may have appeared upstream since).

func wosSingleHitBody(uid, doi, title string) string {
	return `{
		"metadata":{"total":1,"page":1,"limit":50},
		"hits":[{
			"uid":"` + uid + `",
			"title":"` + title + `",
			"identifiers":{"doi":"` + doi + `"}
		}]
	}`
}

func TestWosLookupServesRepeatDOIFromCache(t *testing.T) {
	upstreamCalls := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		upstreamCalls++
		doi := requestedWosDOI(request)
		writeWosResponse(t, response, wosSingleHitBody("WOS:CACHED:"+doi, doi, "Cached record"))
	}))
	defer upstream.Close()

	service := newWosTestService(upstream.URL, &stubWosImportRepository{})

	first, err := service.PreviewWosPublication(context.Background(), "10.1000/cached.1", "", "B")
	require.NoError(t, err)
	second, err := service.PreviewWosPublication(context.Background(), "https://doi.org/10.1000/CACHED.1", "", "B")
	require.NoError(t, err)

	assert.Equal(t, 1, upstreamCalls, "a repeat lookup inside the cache window must not consume quota")
	require.NotNil(t, first.Values)
	require.NotNil(t, second.Values)
	assert.Equal(t, "Cached record", valueOrEmpty(second.Values.Title))
	assert.Equal(t, "10.1000/cached.1", second.Doi)

	// A different DOI still goes upstream.
	_, err = service.PreviewWosPublication(context.Background(), "10.1000/other.2", "", "B")
	require.NoError(t, err)
	assert.Equal(t, 2, upstreamCalls)
}

func TestWosLookupDoesNotCacheMissesOrAmbiguity(t *testing.T) {
	upstreamCalls := 0
	mode := "empty"
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		upstreamCalls++
		switch mode {
		case "empty":
			writeWosResponse(t, response, `{"metadata":{"total":0,"page":1,"limit":50},"hits":[]}`)
		case "ambiguous":
			writeWosResponse(t, response, `{
				"metadata":{"total":2,"page":1,"limit":50},
				"hits":[
					{"uid":"WOS:A","title":"Twin A","identifiers":{"doi":"10.1000/twin"}},
					{"uid":"WOS:B","title":"Twin B","identifiers":{"doi":"10.1000/twin"}}
				]
			}`)
		}
	}))
	defer upstream.Close()

	service := newWosTestService(upstream.URL, &stubWosImportRepository{})

	_, err := service.PreviewWosPublication(context.Background(), "10.1000/twin", "", "B")
	require.Error(t, err)
	_, err = service.PreviewWosPublication(context.Background(), "10.1000/twin", "", "B")
	require.Error(t, err)
	assert.Equal(t, 2, upstreamCalls, "a not-found must never be cached")

	mode = "ambiguous"
	_, err = service.PreviewWosPublication(context.Background(), "10.1000/twin", "", "B")
	require.Error(t, err)
	assert.Equal(t, 3, upstreamCalls, "an ambiguous upstream answer must never be cached")
}

func TestWosLookupCacheExpiresAndStaysBounded(t *testing.T) {
	now := time.Date(2026, 9, 27, 12, 0, 0, 0, time.UTC)
	cache := newWosLookupCache(time.Hour, 2)
	cache.now = func() time.Time { return now }

	cache.store("a", models.WosHit{WosUID: "WOS:A"})
	// Distinct instants: entries expiring together would tie, and map
	// iteration order would then pick the evicted one at random.
	now = now.Add(time.Second)
	cache.store("b", models.WosHit{WosUID: "WOS:B"})

	now = now.Add(time.Minute)
	cache.store("c", models.WosHit{WosUID: "WOS:C"})
	// Cap is 2: the soonest-expiring entry (a, stored first) is evicted.
	assert.Equal(t, 2, len(cache.entries))
	_, found := cache.get("a")
	assert.False(t, found)
	cachedC, found := cache.get("c")
	require.True(t, found)
	assert.Equal(t, "WOS:C", cachedC.WosUID)

	now = now.Add(2 * time.Hour)
	_, found = cache.get("b")
	assert.False(t, found, "entries must expire after the TTL")
}

func TestWosLookupCacheIsNilSafe(t *testing.T) {
	// Services built without the constructor (tests, future wiring) must not
	// panic on cache access; they simply never hit.
	var cache *wosLookupCache
	_, found := cache.get("10.1000/x")
	assert.False(t, found)
	assert.NotPanics(t, func() { cache.store("10.1000/x", models.WosHit{}) })
}

func TestPreviewWosPublicationReportsFieldWarnings(t *testing.T) {
	repo := &stubWosImportRepository{}
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		writeWosResponse(t, response, `{
			"metadata":{"total":1,"page":1,"limit":50},
			"hits":[{
				"uid":"WOS:WARNINGS",
				"title":"Warning fixture",
				"types":["Article"],
				"sourceTypes":["Journal"],
				"source":{
					"sourceTitle":"Journal of Odd Volumes",
					"publishYear":2025,
					"publishMonth":"SEP",
					"volume":"42R",
					"issue":"7"
				},
				"identifiers":{"doi":"10.1000/warnings.1"}
			}]
		}`)
	}))
	defer upstream.Close()

	preview, err := newWosTestService(upstream.URL, repo).PreviewWosPublication(
		context.Background(), "10.1000/warnings.1", "", "B",
	)
	require.NoError(t, err)
	require.NotNil(t, preview.Values)

	// Non-numeric volume stays unset and is explained instead of guessed.
	assert.Nil(t, preview.Values.Volume)
	assert.Equal(t, 7, intValueOrZero(preview.Values.Issue))

	byCode := map[string]models.WosImportWarning{}
	for _, warning := range preview.Warnings {
		byCode[warning.Code] = warning
	}
	require.Contains(t, byCode, "VOLUME_NOT_NUMERIC")
	assert.Equal(t, "volume", byCode["VOLUME_NOT_NUMERIC"].Field)
	assert.Equal(t, "42R", byCode["VOLUME_NOT_NUMERIC"].Raw)
	assert.NotEmpty(t, byCode["VOLUME_NOT_NUMERIC"].Message)

	// WoS only knows the month; the missing day is a warning, not a fake date.
	require.Contains(t, byCode, "DATE_DAY_MISSING")
	assert.Equal(t, "dateOfPublication", byCode["DATE_DAY_MISSING"].Field)
	assert.Equal(t, "2025-09", valueOrEmpty(preview.Values.DateOfPublication))
	assert.NotContains(t, byCode, "ISSUE_NOT_NUMERIC")
}

func TestPreviewWosPublicationKeepsIsbnForBookChapter(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		writeWosResponse(t, response, `{
			"metadata":{"total":1,"page":1,"limit":50},
			"hits":[{
				"uid":"WOS:BOOK",
				"title":"A book chapter",
				"types":["Book Chapter"],
				"sourceTypes":["Book"],
				"identifiers":{"doi":"10.1000/book.1","isbn":"978-3-16-148410-0"}
			}]
		}`)
	}))
	defer upstream.Close()

	preview, err := newWosTestService(upstream.URL, &stubWosImportRepository{}).PreviewWosPublication(
		context.Background(), "10.1000/book.1", "", "B",
	)
	require.NoError(t, err)
	require.NotNil(t, preview.Values)
	assert.Equal(t, "978-3-16-148410-0", valueOrEmpty(preview.Values.Isbn))
	assert.NotContains(t, preview.UnavailableFields, "isbn")
}

func TestPreviewWosPublicationJSONCarriesWarningAndConfidenceFields(t *testing.T) {
	// FE (ELIPANDA-502) reads warnings[].code and authors[].match.confidence
	// straight off the wire; lock the JSON spelling.
	repo := &stubWosImportRepository{researchers: []wosResearcherRecord{
		{UID: "known", FirstName: "Jane", LastName: "Doe", ResearcherIDs: []string{"AAB-1234-2020"}},
	}}
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		writeWosResponse(t, response, `{
			"metadata":{"total":1,"page":1,"limit":50},
			"hits":[{
				"uid":"WOS:CONTRACT",
				"title":"Contract fixture",
				"types":["Article"],
				"source":{"publishYear":2025},
				"identifiers":{"doi":"10.1000/contract.1"},
				"names":{"authors":[
					{"displayName":"Doe, Jane","researcherId":"AAB-1234-2020"}
				]}
			}]
		}`)
	}))
	defer upstream.Close()

	response := performPreviewRequest(t, newWosTestService(upstream.URL, repo),
		"/v1/publications/wos/lookup?doi=10.1000%2Fcontract.1")
	require.Equal(t, http.StatusOK, response.Code)

	var wire struct {
		Warnings []struct {
			Code    string `json:"code"`
			Field   string `json:"field"`
			Raw     string `json:"raw"`
			Message string `json:"message"`
		} `json:"warnings"`
		Authors []struct {
			Match struct {
				Kind              string `json:"kind"`
				Confidence        string `json:"confidence"`
				KnownResearcherId bool   `json:"knownResearcherId"`
			} `json:"match"`
		} `json:"authors"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &wire))

	require.Len(t, wire.Authors, 1)
	assert.Equal(t, "researcher-id", wire.Authors[0].Match.Kind)
	assert.Equal(t, "EXACT_ID", wire.Authors[0].Match.Confidence)
	assert.True(t, wire.Authors[0].Match.KnownResearcherId)

	// Article: 2025 with no month yields the date warning on the wire.
	require.NotEmpty(t, wire.Warnings)
	assert.Equal(t, "DATE_DAY_MISSING", wire.Warnings[0].Code)
	assert.NotEmpty(t, wire.Warnings[0].Message)
}
