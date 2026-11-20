package cronservice

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"panda/apigateway/services/cron-service/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// cronServiceMock is a hand-written ICronService stub; syncFn controls SyncEmployees.
type cronServiceMock struct {
	syncFn func() (models.SyncEmployeesResult, error)
}

func (m *cronServiceMock) LogCronJobStart(job *models.CronJob) (string, error) { return "", nil }

func (m *cronServiceMock) LogCronJobEnd(uid, jobStatus, jobMessage string) error { return nil }

func (m *cronServiceMock) GetCronJobHistory() ([]models.CronJobHistory, error) {
	return []models.CronJobHistory{}, nil
}

func (m *cronServiceMock) SyncEmployees() (models.SyncEmployeesResult, error) {
	return m.syncFn()
}

func runSyncEmployees(t *testing.T, syncFn func() (models.SyncEmployeesResult, error)) (*httptest.ResponseRecorder, error) {
	t.Helper()
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodPost, "/v1/cron/sync-employees", nil), rec)

	h := NewCronHandlers(&cronServiceMock{syncFn: syncFn})
	return rec, h.SyncEmployees()(c)
}

func TestSyncEmployees_ReturnsCounts(t *testing.T) {
	rec, err := runSyncEmployees(t, func() (models.SyncEmployeesResult, error) {
		return models.SyncEmployeesResult{SavedCount: 12, ErrorCount: 1}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"savedCount":12,"errorCount":1}`, rec.Body.String())
}

func TestSyncEmployees_AlreadyRunningIsConflict(t *testing.T) {
	_, err := runSyncEmployees(t, func() (models.SyncEmployeesResult, error) {
		return models.SyncEmployeesResult{}, ErrEmployeesSyncRunning
	})

	var httpErr *echo.HTTPError
	assert.True(t, errors.As(err, &httpErr))
	assert.Equal(t, http.StatusConflict, httpErr.Code)
}

func TestSyncEmployees_FetchFailureIsBadGateway(t *testing.T) {
	_, err := runSyncEmployees(t, func() (models.SyncEmployeesResult, error) {
		return models.SyncEmployeesResult{}, errors.New("OKBase returned status 500")
	})

	var httpErr *echo.HTTPError
	assert.True(t, errors.As(err, &httpErr))
	assert.Equal(t, http.StatusBadGateway, httpErr.Code)
}
