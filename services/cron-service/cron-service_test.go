package cronservice

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"panda/apigateway/services/cron-service/models"

	"github.com/stretchr/testify/assert"
)

func newOKBaseServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "secret", r.Header.Get("x-api-key"))
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server
}

func TestGetEmployeesFromApi_ParsesRows(t *testing.T) {
	server := newOKBaseServer(t, http.StatusOK, `{"radky":[{"OJ_KOD":"1234"},{"OJ_KOD":"5678"}]}`)

	data, err := GetEmployeesFromApi(server.URL, "secret")

	assert.NoError(t, err)
	assert.Len(t, data.Rows, 2)
	assert.Equal(t, "1234", data.Rows[0].OJ_KOD)
}

func TestGetEmployeesFromApi_NonSuccessStatusIsError(t *testing.T) {
	server := newOKBaseServer(t, http.StatusInternalServerError, `oops`)

	_, err := GetEmployeesFromApi(server.URL, "secret")

	assert.ErrorContains(t, err, "status 500")
}

func TestGetEmployeesFromApi_InvalidJSONIsError(t *testing.T) {
	server := newOKBaseServer(t, http.StatusOK, `not json`)

	_, err := GetEmployeesFromApi(server.URL, "secret")

	assert.ErrorContains(t, err, "decode OKBase response")
}

func TestGetEmployeesFromApi_UnreachableIsError(t *testing.T) {
	server := newOKBaseServer(t, http.StatusOK, `{}`)
	server.Close()

	_, err := GetEmployeesFromApi(server.URL, "secret")

	assert.ErrorContains(t, err, "call OKBase")
}

func TestRunEmployeesSync_RejectsOverlappingRun(t *testing.T) {
	employeesSyncLock.Lock()
	defer employeesSyncLock.Unlock()

	_, err := RunEmployeesSync(nil, nil, employeesSyncTriggerManual)

	assert.True(t, errors.Is(err, ErrEmployeesSyncRunning))
}

func TestEmployeesSyncJobOutcome(t *testing.T) {
	status, message := employeesSyncJobOutcome(models.SyncEmployeesResult{SavedCount: 3}, nil)
	assert.Equal(t, models.CRON_JOB_STATUS_SUCCESS, status)
	assert.Equal(t, "3 employees saved, 0 errors", message)

	status, _ = employeesSyncJobOutcome(models.SyncEmployeesResult{SavedCount: 2, ErrorCount: 1}, nil)
	assert.Equal(t, models.CRON_JOB_STATUS_ERROR, status)

	status, message = employeesSyncJobOutcome(models.SyncEmployeesResult{}, errors.New("call OKBase: timeout"))
	assert.Equal(t, models.CRON_JOB_STATUS_ERROR, status)
	assert.Equal(t, "call OKBase: timeout", message)
}
