package systemsService

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// canEditServiceMock embeds ISystemsService (a nil interface) so it satisfies the large
// interface without implementing every method; tests override only what they exercise.
// Any un-overridden method called would panic — tests are written so that never happens.
type canEditServiceMock struct {
	ISystemsService
	canEditFn func(systemUID, userUID string) (models.CanEditSystemResult, error)
	itemFn    func(itemUID string) (string, error)
	calls     *[]string
}

func (m *canEditServiceMock) CanEditSystem(systemUID, userUID string) (models.CanEditSystemResult, error) {
	if m.calls != nil {
		*m.calls = append(*m.calls, systemUID)
	}
	return m.canEditFn(systemUID, userUID)
}

func (m *canEditServiceMock) GetSystemUIDByItemUID(itemUID string) (string, error) {
	return m.itemFn(itemUID)
}

func canEditContext(method, target string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userUID", "actor-uid")
	c.Set("facilityCode", "B")
	c.Set("userRoles", []string{"systems-edit"})
	return c, rec
}

func assertForbidden(t *testing.T, err error) {
	t.Helper()
	var he *echo.HTTPError
	if assert.True(t, errors.As(err, &he), "expected echo.HTTPError, got %v", err) {
		assert.Equal(t, http.StatusForbidden, he.Code)
	}
}

func TestGuardSystemEdit_Allowed(t *testing.T) {
	c, _ := canEditContext(http.MethodPut, "/v1/system/s1")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		canEditFn: func(s, u string) (models.CanEditSystemResult, error) {
			return models.CanEditSystemResult{Result: true}, nil
		},
	}}
	assert.NoError(t, h.guardSystemEdit(c, "s1"))
}

func TestGuardSystemEdit_Forbidden(t *testing.T) {
	c, _ := canEditContext(http.MethodPut, "/v1/system/s1")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		canEditFn: func(s, u string) (models.CanEditSystemResult, error) {
			return models.CanEditSystemResult{Result: false}, nil
		},
	}}
	assertForbidden(t, h.guardSystemEdit(c, "s1"))
}

func TestGuardSystemEdit_MissingSystemPasses(t *testing.T) {
	c, _ := canEditContext(http.MethodDelete, "/v1/system/missing")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		canEditFn: func(s, u string) (models.CanEditSystemResult, error) {
			return models.CanEditSystemResult{}, helpers.ERR_NO_ROWS
		},
	}}
	// non-existent system -> guard passes (handler's own flow 404s)
	assert.NoError(t, h.guardSystemEdit(c, "missing"))
}

func TestGuardSystemEdit_ForbidIfAnyForbidden(t *testing.T) {
	c, _ := canEditContext(http.MethodPost, "/v1/systems/move")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		canEditFn: func(s, u string) (models.CanEditSystemResult, error) {
			// allowed for s1, forbidden for the destination parent
			return models.CanEditSystemResult{Result: s == "s1"}, nil
		},
	}}
	assertForbidden(t, h.guardSystemEdit(c, "s1", "destParent"))
}

func TestGuardSystemEdit_SkipsEmptyUids(t *testing.T) {
	var calls []string
	c, _ := canEditContext(http.MethodPost, "/v1/systems/move")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		calls: &calls,
		canEditFn: func(s, u string) (models.CanEditSystemResult, error) {
			return models.CanEditSystemResult{Result: true}, nil
		},
	}}
	assert.NoError(t, h.guardSystemEdit(c, "", "s1", ""))
	assert.Equal(t, []string{"s1"}, calls) // only the non-empty uid triggers a check
}

func TestCanEditSystemHandler_OK(t *testing.T) {
	c, rec := canEditContext(http.MethodGet, "/v1/system/s1/can-edit")
	c.SetParamNames("uid")
	c.SetParamValues("s1")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		canEditFn: func(s, u string) (models.CanEditSystemResult, error) {
			return models.CanEditSystemResult{
				Result:       true,
				Responsibles: []models.SystemResponsible{{UID: "A", FirstName: "Ann", LastName: "Lee"}},
			}, nil
		},
	}}
	err := h.CanEditSystem()(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"result":true`)
	assert.Contains(t, rec.Body.String(), `"uid":"A"`)
}

func TestCanEditSystemHandler_NotFound(t *testing.T) {
	c, _ := canEditContext(http.MethodGet, "/v1/system/missing/can-edit")
	c.SetParamNames("uid")
	c.SetParamValues("missing")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		canEditFn: func(s, u string) (models.CanEditSystemResult, error) {
			return models.CanEditSystemResult{}, helpers.ERR_NO_ROWS
		},
	}}
	err := h.CanEditSystem()(c)
	var he *echo.HTTPError
	if assert.True(t, errors.As(err, &he)) {
		assert.Equal(t, http.StatusNotFound, he.Code)
	}
}

func TestGuardByItem_ForbiddenAfterResolve(t *testing.T) {
	c, _ := canEditContext(http.MethodPut, "/v1/physical-item/item1/properties")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		itemFn: func(itemUID string) (string, error) { return "sysOfItem", nil },
		canEditFn: func(s, u string) (models.CanEditSystemResult, error) {
			assert.Equal(t, "sysOfItem", s)
			return models.CanEditSystemResult{Result: false}, nil
		},
	}}
	assertForbidden(t, h.guardSystemEditByItem(c, "item1"))
}

func TestGuardByItem_MissingItemPasses(t *testing.T) {
	c, _ := canEditContext(http.MethodPut, "/v1/physical-item/item1/properties")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		itemFn: func(itemUID string) (string, error) { return "", helpers.ERR_NO_ROWS },
	}}
	assert.NoError(t, h.guardSystemEditByItem(c, "item1"))
}

func TestUpdateSystemHandler_Forbidden(t *testing.T) {
	// guard fails -> handler returns 403 before ever calling the (unimplemented) UpdateSystem service
	c, _ := canEditContext(http.MethodPut, "/v1/system/s1")
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.SetParamNames("uid")
	c.SetParamValues("s1")
	h := &SystemsHandlers{systemsService: &canEditServiceMock{
		canEditFn: func(s, u string) (models.CanEditSystemResult, error) {
			return models.CanEditSystemResult{Result: false}, nil
		},
	}}
	assertForbidden(t, h.UpdateSystem()(c))
}
