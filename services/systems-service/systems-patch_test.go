package systemsService

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func rawBody(t *testing.T, body string) map[string]json.RawMessage {
	t.Helper()
	raw := map[string]json.RawMessage{}
	require.NoError(t, json.Unmarshal([]byte(body), &raw))
	return raw
}

func TestParsePatchSystemPayload_AbsentNullAndValues(t *testing.T) {
	f, err := parsePatchSystemPayload(rawBody(t, `{
		"lastUpdateTime": "2026-01-01T10:00:00.123456789Z",
		"name": "N",
		"description": null,
		"zone": {"uid": "z1"},
		"location": null,
		"operators": [{"uid": "a"}, {"uid": "a"}, {"uid": "b"}],
		"minimalSpareParstCount": 3
	}`))
	require.NoError(t, err)
	assert.True(t, f.Token.Present)
	assert.Equal(t, 123456789, f.Token.Value.Nanosecond())
	assert.Equal(t, "N", *f.Name)
	require.NotNil(t, f.Description)
	assert.Nil(t, f.Description.Value, "explicit null = clear")
	assert.Nil(t, f.SystemCode, "absent = unchanged")
	assert.Equal(t, "z1", f.Zone.Value.UID)
	assert.Nil(t, f.Location.Value)
	assert.Equal(t, []models.UIDRef{{UID: "a"}, {UID: "b"}}, *f.Operators, "duplicates dropped")
	assert.Nil(t, f.MaintainedBy)
	assert.Equal(t, 3.0, *f.MinimalSpareParstCount.Value)
}

func TestParsePatchSystemPayload_NullToken(t *testing.T) {
	f, err := parsePatchSystemPayload(rawBody(t, `{"lastUpdateTime": null}`))
	require.NoError(t, err)
	assert.True(t, f.Token.Present)
	assert.Nil(t, f.Token.Value)

	f, err = parsePatchSystemPayload(rawBody(t, `{}`))
	require.NoError(t, err)
	assert.False(t, f.Token.Present)
}

func TestParsePatchSystemPayload_Rejects(t *testing.T) {
	cases := map[string]string{
		"unknown key":      `{"attribute": {"uid": "x"}}`,
		"null name":        `{"name": null}`,
		"ref without uid":  `{"zone": {}}`,
		"null list":        `{"operators": null}`,
		"list without uid": `{"maintainedBy": [{}]}`,
		"bad token":        `{"lastUpdateTime": "yesterday"}`,
		"wrong type":       `{"description": 5}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := parsePatchSystemPayload(rawBody(t, body))
			assert.Error(t, err)
		})
	}
}

func TestParsePatchPhysicalItemPayload(t *testing.T) {
	f, err := parsePatchPhysicalItemPayload(rawBody(t, `{"notes": null, "itemUsage": {"uid": "u"}}`))
	require.NoError(t, err)
	assert.Nil(t, f.Notes.Value)
	assert.Equal(t, "u", f.ItemUsage.Value.UID)
	assert.Nil(t, f.SerialNumber)

	_, err = parsePatchPhysicalItemPayload(rawBody(t, `{"name": "x"}`))
	assert.Error(t, err)
}

func TestFormatPropertyValue(t *testing.T) {
	s := func(v string) *string { return &v }
	assert.Equal(t, "N/A", formatPropertyValue(nil, "Text"))
	assert.Equal(t, "N/A", formatPropertyValue(s(""), "Text"))
	assert.Equal(t, "abc", formatPropertyValue(s("abc"), "Text"))
	assert.Equal(t, "1.5 - N/A", formatPropertyValue(s(`{"min":1.5}`), "Range"))
	assert.Equal(t, "N/A - N/A", formatPropertyValue(nil, "Range"))
	assert.Equal(t, "garbage", formatPropertyValue(s("garbage"), "Range"))
}

func TestResolveCatalogueProperties_LatestServiceWinsAndServiceOnly(t *testing.T) {
	s := func(v string) *string { return &v }
	older := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	newer := older.AddDate(0, 1, 0)
	group := []models.CodebookRef{{UID: "g", Name: "Group"}}
	raw := models.RawPhysicalItemProperties{
		Catalogue: []models.RawCatalogueProperty{
			{UID: "p1", Name: "Weight", Type: s("Text"), Value: s("10"), Groups: group},
			{UID: "p2", Name: "Same", Type: s("Text"), Value: s("x")},
		},
		ServiceItems: []models.RawServiceItemProperties{
			{Created: &older, Properties: []models.RawCatalogueProperty{{UID: "p1", Name: "Weight", Value: s("11"), Groups: group}}},
			{Created: &newer, Properties: []models.RawCatalogueProperty{
				{UID: "p1", Name: "Weight", Value: s("12"), Groups: group},
				{UID: "p2", Name: "Same", Value: s("x")},
				{UID: "p3", Name: "Extra", Value: s("y")},
			}},
		},
	}
	res := resolveCatalogueProperties(raw)

	require.Len(t, res.Groups, 2)
	assert.Equal(t, "no-group", res.Groups[0].Key, "General first")
	assert.True(t, res.HasOverriddenProperties)

	general := res.Groups[0].Properties
	require.Len(t, general, 2)
	assert.Equal(t, "Extra", general[0].Name)
	assert.Nil(t, general[0].Value, "service-only property has no catalogue value")
	assert.False(t, general[0].IsOverridden)
	assert.Equal(t, "Same", general[1].Name)
	assert.False(t, general[1].IsOverridden, "same value is not an override")

	weight := res.Groups[1].Properties[0]
	assert.Equal(t, "12", *weight.ServiceValue)
	assert.True(t, weight.IsOverridden)
}

func TestIsAllowedChildLevel(t *testing.T) {
	level := func(v string) *string { return &v }
	assert.True(t, isAllowedChildLevel(nil, "KEY_SYSTEMS"))
	assert.True(t, isAllowedChildLevel(level("TECHNOLOGY_UNIT"), "KEY_SYSTEMS"))
	assert.False(t, isAllowedChildLevel(level("SYSTEM_DOMAIN"), "KEY_SYSTEMS"))
	assert.False(t, isAllowedChildLevel(level("TRASH"), "TRASH"))
	assert.True(t, isAllowedChildLevel(level("LEGACY_VALUE"), "TRASH"), "unknown parent level is not blocked")
}

func TestPatchSystemQuery_TokenAndAudit(t *testing.T) {
	plan := newPatchPlan()
	plan.scalars["name"] = "N"
	plan.relations = append(plan.relations, relationPatch{relType: "HAS_ZONE", label: "Zone", uids: []string{"z"}})
	token := time.Date(2026, 1, 1, 10, 0, 0, 42, time.UTC)

	q := PatchSystemQuery("s1", "B", "u1", &token, plan, nil, "[]", true)
	assert.Contains(t, q.Query, "s.lastUpdateTime.nanosecond = $tokenNanosecond")
	assert.Equal(t, 42, q.Parameters["tokenNanosecond"])
	assert.Contains(t, q.Query, "s.name = $sys_name")
	assert.Contains(t, q.Query, "MATCH (n:Zone) WHERE n.uid IN $sys_rel0 MERGE (s)-[:HAS_ZONE]->(n)")
	assert.Contains(t, q.Query, `action: "UPDATE"`)
	assert.Equal(t, true, q.Parameters["outsideResponsibility"])

	legacy := PatchSystemQuery("s1", "B", "u1", nil, plan, nil, "[]", false)
	assert.Contains(t, legacy.Query, "s.lastUpdateTime IS NULL")
}

// patchServiceMock overrides only the PATCH-family methods exercised by handler tests.
type patchServiceMock struct {
	ISystemsService
	err error
}

func (m *patchServiceMock) PatchSystem(uid, facilityCode, userUID string, f *models.PatchSystemFields) (models.SystemDetail, error) {
	return models.SystemDetail{UID: uid}, m.err
}

func patchContext(body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/v1/system/s1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uid")
	c.SetParamValues("s1")
	c.Set("userUID", "u1")
	c.Set("facilityCode", "B")
	return c, rec
}

func httpCode(t *testing.T, err error, rec *httptest.ResponseRecorder) int {
	t.Helper()
	var he *echo.HTTPError
	if errors.As(err, &he) {
		return he.Code
	}
	require.NoError(t, err)
	return rec.Code
}

func TestPatchSystemHandler_ErrorMapping(t *testing.T) {
	cases := []struct {
		name string
		body string
		err  error
		code int
		in   string
	}{
		{"ok", `{"name":"x"}`, nil, http.StatusOK, `"uid":"s1"`},
		{"bad json", `{`, nil, http.StatusBadRequest, ""},
		{"unknown key", `{"foo":1}`, nil, http.StatusBadRequest, ""},
		{"validation", `{}`, &ValidationError{Msg: "bad"}, http.StatusBadRequest, ""},
		{"not found", `{}`, helpers.ERR_NOT_FOUND, http.StatusNotFound, ""},
		{"stale", `{}`, ErrStaleSystem, http.StatusConflict, `"code":"STALE"`},
		{"duplicate code", `{}`, helpers.ERR_DUPLICATE_SYSTEM_CODE, http.StatusConflict, `"code":"DUPLICATE_SYSTEM_CODE"`},
		{"internal", `{}`, errors.New("boom"), http.StatusInternalServerError, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := patchContext(tc.body)
			h := &SystemsHandlers{systemsService: &patchServiceMock{err: tc.err}}
			err := h.PatchSystem()(c)
			assert.Equal(t, tc.code, httpCode(t, err, rec))
			if tc.in != "" {
				assert.Contains(t, rec.Body.String(), tc.in)
			}
		})
	}
}
