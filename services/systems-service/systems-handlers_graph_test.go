package systemsService

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newGraphQueryContext(rawQuery string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/system/test-uid/graph?"+rawQuery, nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func TestParseSystemGraphQueryOptions_LegacyMode(t *testing.T) {
	c := newGraphQueryContext("")

	options, err := parseSystemGraphQueryOptions(c)

	assert.NoError(t, err)
	assert.Nil(t, options.LimitPerRelationshipType)
	assert.False(t, options.IncludeRelationshipStats)
	assert.Equal(t, "", options.RelationshipType)
}

func TestParseSystemGraphQueryOptions_InitialMode(t *testing.T) {
	c := newGraphQueryContext("limitPerRelationshipType=20&includeRelationshipStats=true")

	options, err := parseSystemGraphQueryOptions(c)

	assert.NoError(t, err)
	if assert.NotNil(t, options.LimitPerRelationshipType) {
		assert.Equal(t, 20, *options.LimitPerRelationshipType)
	}
	assert.True(t, options.IncludeRelationshipStats)
	assert.Equal(t, "", options.RelationshipType)
}

func TestParseSystemGraphQueryOptions_LoadMorePriorityOverInitial(t *testing.T) {
	c := newGraphQueryContext("relationshipType=IS_POWERED_BY&offset=20&limit=10&limitPerRelationshipType=20")

	options, err := parseSystemGraphQueryOptions(c)

	assert.NoError(t, err)
	assert.Equal(t, "IS_POWERED_BY", options.RelationshipType)
	assert.Equal(t, 20, options.Offset)
	assert.Equal(t, 10, options.Limit)
	assert.Nil(t, options.LimitPerRelationshipType)
}

func TestParseSystemGraphQueryOptions_InvalidOffset(t *testing.T) {
	c := newGraphQueryContext("relationshipType=IS_POWERED_BY&offset=-1")

	_, err := parseSystemGraphQueryOptions(c)

	assert.EqualError(t, err, "invalid offset")
}

func TestParseSystemGraphQueryOptions_OffsetWithoutType(t *testing.T) {
	c := newGraphQueryContext("offset=10")

	_, err := parseSystemGraphQueryOptions(c)

	assert.EqualError(t, err, "offset and limit require relationshipType")
}

func TestParseSystemGraphQueryOptions_InvalidIncludeStats(t *testing.T) {
	c := newGraphQueryContext("includeRelationshipStats=abc")

	_, err := parseSystemGraphQueryOptions(c)

	assert.EqualError(t, err, "invalid includeRelationshipStats")
}
