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
	calls     *[]string
}

func (m *canEditServiceMock) CanEditSystem(systemUID, userUID string) (models.CanEditSystemResult, error) {
	if m.calls != nil {
		*m.calls = append(*m.calls, systemUID)
	}
	return m.canEditFn(systemUID, userUID)
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
