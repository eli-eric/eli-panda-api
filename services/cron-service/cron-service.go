package cronservice

import (
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
	c := createScheduler(driver)
	if c == nil {
		log.Error().Msgf("Cron service initialization failed: %s", "Cron service is nil")
	}
	defer c.Stop()

	return &CronService{neo4jDriver: driver, jwtSecret: settings.JwtSecret, scheduler: c}
}

func createScheduler(driver *neo4j.Driver) *cron.Cron {
	c := cron.New()
	if c == nil {
		return nil
	}

	// schedule a job every day at 02:00
	_, err := c.AddFunc("0 2 * * *", func() {
		log.Info().Msgf("Cron service started: %s", "SyncEliBeamlinesEmployees")
		err := SyncEliBeamlinesEmployees(driver)
		if err != nil {
			log.Error().Msgf("Cron service error: %s", err.Error())
		}
	})

	if err != nil {
		log.Error().Msgf("Cron service initialization failed: %s", err.Error())
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

	helpers.ProcessArrayResult[models.CronJobHistory](&cronJobHistory, err)

	return cronJobHistory, err
}

func SyncEliBeamlinesEmployees(neo4jDriver *neo4j.Driver) (err error) {
	//session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	// query := SyncEliBeamlinesEmployeesQuery()
	// _, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return err
}
