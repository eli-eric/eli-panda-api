package systemsService

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"panda/apigateway/helpers"
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

func TestParseSystemGraphQueryOptions_FilteredMode(t *testing.T) {
	c := newGraphQueryContext("search=power&systemLevels=technology_unit,key_systems&systemType=Cooling&relationshipTypes=is_powered_by,has_subsystem")

	options, err := parseSystemGraphQueryOptions(c)

	assert.NoError(t, err)
	assert.Equal(t, "power", options.Search)
	assert.Equal(t, []string{"TECHNOLOGY_UNIT", "KEY_SYSTEMS"}, options.SystemLevels)
	assert.Equal(t, "Cooling", options.SystemType)
	assert.Equal(t, []string{"IS_POWERED_BY", "HAS_SUBSYSTEM"}, options.RelationshipTypes)
	assert.Nil(t, options.LimitPerRelationshipType)
	assert.Equal(t, "", options.RelationshipType)
}

func TestParseSystemGraphQueryOptions_RelationshipTypeWithLimitPerTypeIsInvalid(t *testing.T) {
	c := newGraphQueryContext("relationshipType=IS_POWERED_BY&offset=20&limit=10&limitPerRelationshipType=20")

	_, err := parseSystemGraphQueryOptions(c)

	assert.EqualError(t, err, "relationshipType cannot be combined with limitPerRelationshipType")
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

func TestParseSystemGraphQueryOptions_LoadMoreWithFiltersIsValid(t *testing.T) {
	c := newGraphQueryContext("relationshipType=IS_POWERED_BY&search=abc")

	options, err := parseSystemGraphQueryOptions(c)

	assert.NoError(t, err)
	assert.Equal(t, "IS_POWERED_BY", options.RelationshipType)
	assert.Equal(t, "abc", options.Search)
}

func TestParseSystemGraphQueryOptions_FilterWithLimitPerTypeIsValid(t *testing.T) {
	c := newGraphQueryContext("search=abc&limitPerRelationshipType=20")

	options, err := parseSystemGraphQueryOptions(c)

	assert.NoError(t, err)
	assert.Equal(t, "abc", options.Search)
	if assert.NotNil(t, options.LimitPerRelationshipType) {
		assert.Equal(t, 20, *options.LimitPerRelationshipType)
	}
}

func TestParseSystemGraphQueryOptions_InvalidSystemLevels(t *testing.T) {
	c := newGraphQueryContext("systemLevels=INVALID")

	_, err := parseSystemGraphQueryOptions(c)

	assert.EqualError(t, err, "invalid systemLevels")
}

func TestParseSystemGraphQueryOptions_TrashSystemLevelIsAllowed(t *testing.T) {
	c := newGraphQueryContext("systemLevels=TECHNOLOGY_UNIT,TRASH")

	options, err := parseSystemGraphQueryOptions(c)

	assert.NoError(t, err)
	assert.Equal(t, []string{"TECHNOLOGY_UNIT", "TRASH"}, options.SystemLevels)
}

func TestParseSystemGraphQueryOptions_InvalidRelationshipTypes(t *testing.T) {
	c := newGraphQueryContext("relationshipTypes=INVALID")

	_, err := parseSystemGraphQueryOptions(c)

	assert.EqualError(t, err, "invalid relationshipTypes")
}

func TestSystemGraphValidationErrorMessage_GenericFallback(t *testing.T) {
	assert.Equal(t, "invalid graph query params", systemGraphValidationErrorMessage(nil))
	assert.Equal(t, "invalid graph query params", systemGraphValidationErrorMessage(helpers.ERR_INVALID_INPUT))
}

func TestSystemGraphValidationErrorMessage_SpecificWrappedMessage(t *testing.T) {
	err := fmt.Errorf("invalid systemType: %w", helpers.ERR_INVALID_INPUT)

	assert.Equal(t, "invalid systemType", systemGraphValidationErrorMessage(err))
}
