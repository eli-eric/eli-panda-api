package systemsService

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func directOnlyContext(rawQuery string) echo.Context {
	e := echo.New()
	target := "/v1/system/parent-uid/leaves"
	if rawQuery != "" {
		target += "?" + rawQuery
	}
	return e.NewContext(httptest.NewRequest(http.MethodGet, target, nil), httptest.NewRecorder())
}

func TestParseDirectOnlyParam_AbsentDefaultsToFalse(t *testing.T) {
	// Absent must stay the full-subtree behaviour — this is the pre-existing contract.
	got, err := parseDirectOnlyParam(directOnlyContext(""))

	assert.NoError(t, err)
	assert.False(t, got)
}

func TestParseDirectOnlyParam_EmptyValueDefaultsToFalse(t *testing.T) {
	// `?directOnly=` is what a client sends when it serialises an unset flag; it must
	// not be an error.
	got, err := parseDirectOnlyParam(directOnlyContext("directOnly="))

	assert.NoError(t, err)
	assert.False(t, got)
}

func TestParseDirectOnlyParam_AcceptsBoolEncodings(t *testing.T) {
	truthy := []string{"true", "TRUE", "True", "1", "t"}
	for _, raw := range truthy {
		got, err := parseDirectOnlyParam(directOnlyContext("directOnly=" + raw))
		assert.NoError(t, err, raw)
		assert.True(t, got, raw)
	}

	falsy := []string{"false", "FALSE", "0", "f"}
	for _, raw := range falsy {
		got, err := parseDirectOnlyParam(directOnlyContext("directOnly=" + raw))
		assert.NoError(t, err, raw)
		assert.False(t, got, raw)
	}
}

func TestParseDirectOnlyParam_RejectsGarbage(t *testing.T) {
	// Silently falling back to false would hand the caller the whole subtree while
	// they believe they asked for direct children only.
	for _, raw := range []string{"yes", "on", "2", "null"} {
		_, err := parseDirectOnlyParam(directOnlyContext("directOnly=" + raw))
		assert.Error(t, err, raw)
	}
}
