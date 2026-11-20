package cronservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	"panda/apigateway/services/cron-service/models"
	"strings"
	"sync"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// what started the employees sync, stored in the cron job history message
const (
	employeesSyncTriggerScheduled = "scheduled"
	employeesSyncTriggerManual    = "manual"
)

// ErrEmployeesSyncRunning is returned when a sync (scheduled or manual) is already in progress
var ErrEmployeesSyncRunning = errors.New("employees sync is already running")

// shared by the scheduled and the manual run so they never overlap
var employeesSyncLock sync.Mutex

var okBaseClient = &http.Client{Timeout: 60 * time.Second}

type CronService struct {
	neo4jDriver *neo4j.Driver
	scheduler   *cron.Cron
	jwtSecret   string
	settings    *config.Config
}

type ICronService interface {
	LogCronJobStart(job *models.CronJob) (uid string, err error)
	LogCronJobEnd(uid, jobStatus, jobMessage string) (err error)
	GetCronJobHistory() (cronJobHistory []models.CronJobHistory, err error)
	SyncEmployees() (result models.SyncEmployeesResult, err error)
}

func NewCronService(settings *config.Config, driver *neo4j.Driver) ICronService {

	// Schedule a job for Sync eli-bm employees
	c := createScheduler(driver, settings)
	if c == nil {
		log.Error().Msgf("Cron service initialization failed: %s", "Cron service is nil")
	}

	return &CronService{neo4jDriver: driver, jwtSecret: settings.JwtSecret, scheduler: c, settings: settings}
}

func createScheduler(driver *neo4j.Driver, settings *config.Config) *cron.Cron {
	c := cron.New()
	if c == nil {
		return nil
	}

	_, err := c.AddFunc("0 2 * * *", func() {
		log.Info().Msgf("Cron service started: %v", "SyncEliBeamlinesEmployees")
		_, err := RunEmployeesSync(driver, settings, employeesSyncTriggerScheduled)

		if err != nil {
			log.Error().Msgf("Cron service SyncEliBeamlinesEmployees error: %v", err.Error())
		}
	})

	if err != nil {
		log.Error().Msgf("Cron service SyncEliBeamlinesEmployees initialization failed: %v", err.Error())
		return nil
	}

	_, err = c.AddFunc("@every 10s", func() {

		err := SetAllSystemTypes(driver, settings)

		// log only if there is an error
		if err != nil {
			log.Error().Msgf("Cron service SetAllSystemTypes error: %v", err.Error())
		}
	})

	if err != nil {
		log.Error().Msgf("Cron service SetAllSystemTypes initialization failed: %v", err.Error())
		return nil
	}

	c.Start()
	return c
}

func (svc *CronService) LogCronJobStart(job *models.CronJob) (uid string, err error) {
	return logCronJobStart(svc.neo4jDriver, job)
}

func (svc *CronService) LogCronJobEnd(uid, jobStatus, jobMessage string) (err error) {
	return logCronJobEnd(svc.neo4jDriver, uid, jobStatus, jobMessage)
}

func logCronJobStart(neo4jDriver *neo4j.Driver, job *models.CronJob) (uid string, err error) {
	session, _ := helpers.NewNeo4jSession(*neo4jDriver)
	defer session.Close()

	query := LogCronJobStartQuery(job.JobUID, job.JobType, job.JobName)
	uid, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return uid, err
}

func logCronJobEnd(neo4jDriver *neo4j.Driver, uid, jobStatus, jobMessage string) (err error) {
	session, _ := helpers.NewNeo4jSession(*neo4jDriver)
	defer session.Close()

	query := LogCronJobEndQuery(uid, jobStatus, jobMessage)
	_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return err
}

func (svc *CronService) GetCronJobHistory() (cronJobHistory []models.CronJobHistory, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	query := GetCronJobHistoryQuery()
	cronJobHistory, err = helpers.GetNeo4jArrayOfNodes[models.CronJobHistory](session, query)

	helpers.ProcessArrayResult(&cronJobHistory, err)

	return cronJobHistory, err
}

func (svc *CronService) SyncEmployees() (result models.SyncEmployeesResult, err error) {
	return RunEmployeesSync(svc.neo4jDriver, svc.settings, employeesSyncTriggerManual)
}

// RunEmployeesSync runs the HR sync guarded against overlapping runs and records it in the cron job history
func RunEmployeesSync(neo4jDriver *neo4j.Driver, settings *config.Config, trigger string) (result models.SyncEmployeesResult, err error) {
	if !employeesSyncLock.TryLock() {
		return result, ErrEmployeesSyncRunning
	}
	defer employeesSyncLock.Unlock()

	job := models.CRON_JOB_SYNC_EMPLOYEES_FROM_HR_ELI_BM
	historyUid, logErr := logCronJobStart(neo4jDriver, &job)
	if logErr != nil {
		log.Error().Msgf("SyncEliBeamlinesEmployees: failed to log job start: %v", logErr.Error())
	}

	savedCount, rowErrs, err := SyncEliBeamlinesEmployees(neo4jDriver, settings)

	for _, rowErr := range rowErrs {
		log.Error().Msgf("SyncEliBeamlinesEmployees error: %v", rowErr.Error())
	}

	result = models.SyncEmployeesResult{SavedCount: savedCount, ErrorCount: len(rowErrs)}

	if historyUid != "" {
		jobStatus, jobMessage := employeesSyncJobOutcome(result, err)
		jobMessage = fmt.Sprintf("%s run: %s", trigger, jobMessage)
		if logErr := logCronJobEnd(neo4jDriver, historyUid, jobStatus, jobMessage); logErr != nil {
			log.Error().Msgf("SyncEliBeamlinesEmployees: failed to log job end: %v", logErr.Error())
		}
	}

	return result, err
}

func employeesSyncJobOutcome(result models.SyncEmployeesResult, err error) (jobStatus, jobMessage string) {
	if err != nil {
		return models.CRON_JOB_STATUS_ERROR, err.Error()
	}

	jobMessage = fmt.Sprintf("%d employees saved, %d errors", result.SavedCount, result.ErrorCount)
	if result.ErrorCount > 0 {
		return models.CRON_JOB_STATUS_ERROR, jobMessage
	}

	return models.CRON_JOB_STATUS_SUCCESS, jobMessage
}

// SyncEliBeamlinesEmployees fetches employees from HR (OKBase) and merges them into the database.
// err is set when the employees could not be fetched at all, rowErrs holds per-employee write failures.
func SyncEliBeamlinesEmployees(neo4jDriver *neo4j.Driver, settings *config.Config) (succesCount int, rowErrs []error, err error) {
	url := settings.ApiIntegrationBeamlinesOKBaseUrl
	apiKey := settings.ApiIntegrationBeamlinesOKBaseApiKey

	employees, err := GetEmployeesFromApi(url, apiKey)
	if err != nil {
		return 0, nil, err
	}

	session, _ := helpers.NewNeo4jSession(*neo4jDriver)
	defer session.Close()

	for _, employee := range employees.Rows {
		// an Employee without a number cannot be matched on the next run, it would only pile up
		if strings.TrimSpace(employee.OSOBNI_CISLO) == "" {
			rowErrs = append(rowErrs, fmt.Errorf("employee %s %s has no OSOBNI_CISLO, skipped", employee.JMENO, employee.PRIJMENI))
			continue
		}

		if rowErr := syncEmployee(session, employee); rowErr != nil {
			rowErrs = append(rowErrs, rowErr)
		} else {
			succesCount++
		}
	}

	log.Info().Msgf("SyncEliBeamlinesEmployees finished: %v employees saved. Errors: %v", succesCount, len(rowErrs))

	return succesCount, rowErrs, nil
}

// syncEmployee merges one employee; duplicated Employee nodes (no unique constraint on employeeNumber)
// are left untouched and reported, so they are not updated and linked to the same User again
func syncEmployee(session neo4j.Session, employee models.SyncEliBeamlinesEmployee) error {
	count, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CountEmployeesByNumberQuery(employee.OSOBNI_CISLO))
	if err != nil {
		return err
	}

	if count > 1 {
		return fmt.Errorf("employeeNumber %s has %d Employee nodes, skipped", employee.OSOBNI_CISLO, count)
	}

	return helpers.WriteNeo4jAndReturnNothing(session, SyncEliBeamlinesEmployeeQuery(employee, "B"))
}

// use http post to get the employees from the API
func GetEmployeesFromApi(url, apiKey string) (data models.SyncEliBeamlinesEmployeesResponse, err error) {

	// JSON body
	body := []byte(`{
		"kod": "EXT_WS_LDAP",
		"uzivatelska": true,
		"export": false   
	}`)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return data, fmt.Errorf("create OKBase request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	resp, err := okBaseClient.Do(req)
	if err != nil {
		return data, fmt.Errorf("call OKBase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return data, fmt.Errorf("OKBase returned status %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("decode OKBase response: %w", err)
	}

	return data, nil
}

func SetAllSystemTypes(neo4jDriver *neo4j.Driver, settings *config.Config) (err error) {

	session, _ := helpers.NewNeo4jSession(*neo4jDriver)
	defer session.Close()

	err = helpers.WriteNeo4jAndReturnNothing(session, SetSystemTypesQuery())

	return err
}
