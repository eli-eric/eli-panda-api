package securityService

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"panda/apigateway/helpers"
	"panda/apigateway/middlewares"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/security-service/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type securityServiceMock struct {
	setUserEnabledFn func(userUID string, isEnabled bool) (string, error)
}

func (m *securityServiceMock) AuthenticateByUsernameAndPassword(username string, password string) (authUser models.UserAuthInfo, err error) {
	return models.UserAuthInfo{}, nil
}

func (m *securityServiceMock) GetUsersCodebook(facilityCode string) (result []codebookModels.Codebook, err error) {
	return nil, nil
}

func (m *securityServiceMock) GetUsersAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	return nil, nil
}

func (m *securityServiceMock) ChangeUserPassword(userName string, userUID string, passwords *models.ChangePasswordRequest) (err error) {
	return nil
}

func (m *securityServiceMock) SetUserEnabled(userUID string, isEnabled bool) (updatedUserUID string, err error) {
	if m.setUserEnabledFn != nil {
		return m.setUserEnabledFn(userUID, isEnabled)
	}
	return userUID, nil
}

func (m *securityServiceMock) GetEmployeesAutocompleteCodebook(searchText string, limit int, facilityCode string, filter *[]helpers.Filter, isAdmin bool) (result []codebookModels.Codebook, err error) {
	return nil, nil
}

func (m *securityServiceMock) GetProcurementersCodebook(facilityCode string) (result []codebookModels.Codebook, err error) {
	return nil, nil
}

func (m *securityServiceMock) GetTeamsAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	return nil, nil
}

func (m *securityServiceMock) GetContactPersonRolesAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	return nil, nil
}

func (m *securityServiceMock) GetUserByAzureIdToken(azureIdToken string, tenantId string) (authUser models.UserAuthInfo, err error) {
	return models.UserAuthInfo{}, nil
}

type userStatusValidatorMock struct {
	invalidated []string
	cache       []middlewares.UserStatusCacheEntry
}

func (m *userStatusValidatorMock) ValidateUserEnabled(userUID string) (isEnabled bool, err error) {
	return true, nil
}

func (m *userStatusValidatorMock) InvalidateUser(userUID string) {
	m.invalidated = append(m.invalidated, userUID)
}

func (m *userStatusValidatorMock) GetCacheEntries() []middlewares.UserStatusCacheEntry {
	return m.cache
}

func TestEnableUser_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/users/test-user-uid/enable", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/users/:userUID/enable")
	c.SetParamNames("userUID")
	c.SetParamValues("test-user-uid")

	validator := &userStatusValidatorMock{}
	h := NewSecurityHandlers(&securityServiceMock{}, validator)

	err := h.EnableUser()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"userUID":"test-user-uid","isEnabled":true}`, rec.Body.String())
	assert.Equal(t, []string{"test-user-uid"}, validator.invalidated)
}

func TestDisableUser_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/users/test-user-uid/disable", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/users/:userUID/disable")
	c.SetParamNames("userUID")
	c.SetParamValues("test-user-uid")

	validator := &userStatusValidatorMock{}
	h := NewSecurityHandlers(&securityServiceMock{}, validator)

	err := h.DisableUser()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"userUID":"test-user-uid","isEnabled":false}`, rec.Body.String())
	assert.Equal(t, []string{"test-user-uid"}, validator.invalidated)
}

func TestEnableUser_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/users/test-user-uid/enable", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/users/:userUID/enable")
	c.SetParamNames("userUID")
	c.SetParamValues("test-user-uid")

	h := NewSecurityHandlers(&securityServiceMock{
		setUserEnabledFn: func(userUID string, isEnabled bool) (string, error) {
			return "", errors.New("db write failed")
		},
	}, &userStatusValidatorMock{})

	err := h.EnableUser()(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
}

func TestGetUserStatusCache_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/cache", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expiresAt := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)
	validator := &userStatusValidatorMock{
		cache: []middlewares.UserStatusCacheEntry{
			{UserUID: "svc-user", IsEnabled: true, ExpiresAt: expiresAt},
		},
	}

	h := NewSecurityHandlers(&securityServiceMock{}, validator)

	err := h.GetUserStatusCache()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[{"userUID":"svc-user","isEnabled":true,"expiresAt":"2026-02-26T12:00:00Z"}]`, rec.Body.String())
}
