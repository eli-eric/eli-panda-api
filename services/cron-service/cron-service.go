package cronservice

import (
	"bytes"
	"encoding/json"
	"net/http"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	"panda/apigateway/services/cron-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type CronService struct {
	neo4jDriver *neo4j.Driver
	scheduler   *cron.Cron
	jwtSecret   string
}

type ICronService interface {
	LogCronJobStart(job *models.CronJob) (uid string, err error)
	LogCronJobEnd(uid, jobStatus, jobMessage string) (err error)
	GetCronJobHistory() (cronJobHistory []models.CronJobHistory, err error)
}

func NewCronService(settings *config.Config, driver *neo4j.Driver) ICronService {

	// Schedule a job for Sync eli-bm employees
	c := createScheduler(driver, settings)
	if c == nil {
		log.Error().Msgf("Cron service initialization failed: %s", "Cron service is nil")
	}

	return &CronService{neo4jDriver: driver, jwtSecret: settings.JwtSecret, scheduler: c}
}

func createScheduler(driver *neo4j.Driver, settings *config.Config) *cron.Cron {
	c := cron.New()
	if c == nil {
		return nil
	}

	_, err := c.AddFunc("0 2 * * *", func() {
		log.Info().Msgf("Cron service started: %v", "SyncEliBeamlinesEmployees")
		errs := SyncEliBeamlinesEmployees(driver, settings)

		for _, err := range errs {
			log.Error().Msgf("Cron service SyncEliBeamlinesEmployees error: %v", err.Error())
		}

	})

	if err != nil {
		log.Error().Msgf("Cron service SyncEliBeamlinesEmployees initialization failed: %v", err.Error())
		return nil
	}

	_, err = c.AddFunc("@every 10s", func() {
		log.Info().Msgf("Cron service %v started", "SetAllSystemTypes")
		err := SetAllSystemTypes(driver, settings)

		if err != nil {
			log.Error().Msgf("Cron service SetAllSystemTypes error: %v", err.Error())
		} else {
			log.Info().Msgf("Cron service SetAllSystemTypes finished")
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
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := LogCronJobStartQuery(job.JobUID, job.JobType, job.JobName)
	uid, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return uid, err
}

func (svc *CronService) LogCronJobEnd(uid, jobStatus, jobMessage string) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := LogCronJobEndQuery(uid, jobStatus, jobMessage)
	_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return err
}

func (svc *CronService) GetCronJobHistory() (cronJobHistory []models.CronJobHistory, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetCronJobHistoryQuery()
	cronJobHistory, err = helpers.GetNeo4jArrayOfNodes[models.CronJobHistory](session, query)

	helpers.ProcessArrayResult(&cronJobHistory, err)

	return cronJobHistory, err
}

func SyncEliBeamlinesEmployees(neo4jDriver *neo4j.Driver, settings *config.Config) (errs []error) {
	session, _ := helpers.NewNeo4jSession(*neo4jDriver)

	// call the API to get the employees - POST method and result is JSON

	url := settings.ApiIntegrationBeamlinesOKBaseUrl
	apiKey := settings.ApiIntegrationBeamlinesOKBaseApiKey

	// get the employees
	employees := GetEmployeesFromApi(url, apiKey)
	succesCount := 0
	// save the employees to the database
	for _, employee := range employees.Rows {
		// merge the employee to the database
		query := SyncEliBeamlinesEmployeeQuery(employee, "B")
		err := helpers.WriteNeo4jAndReturnNothing(session, query)
		if err != nil {
			errs = append(errs, err)
		} else {
			succesCount++
		}
	}

	log.Info().Msgf("SyncEliBeamlinesEmployees finished: %v employees saved. Errors: %v", succesCount, len(errs))

	return errs
}

// use http post to get the employees from the API
func GetEmployeesFromApi(url, apiKey string) models.SyncEliBeamlinesEmployeesResponse {

	// JSON body
	body := []byte(`{
		"kod": "EXT_WS_LDAP",
		"uzivatelska": true,
		"export": false   
	}`)

	contentType := "application/json"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		log.Error().Msgf("Error while getting employees from API: %s", err.Error())
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Error().Msgf("Error while getting employees from API: %s", err.Error())
	} else {

		var data models.SyncEliBeamlinesEmployeesResponse
		err := json.NewDecoder(resp.Body).Decode(&data)

		if err != nil {
			log.Error().Msgf("Error while getting employees from API: %s", err.Error())
		}
		return data
	}

	defer resp.Body.Close()

	return models.SyncEliBeamlinesEmployeesResponse{}
}

func SetAllSystemTypes(neo4jDriver *neo4j.Driver, settings *config.Config) (err error) {

	session, _ := helpers.NewNeo4jSession(*neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothing(session, SetSystemTypesQuery())

	return err
}
