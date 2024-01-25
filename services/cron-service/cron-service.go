package cronservice

import (
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	"panda/apigateway/services/cron-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type CronService struct {
	neo4jDriver *neo4j.Driver
	jwtSecret   string
}

type ICronService interface {
	LogCronJobStart(job *models.CronJob) (uid string, err error)
	LogCronJobEnd(uid, jobStatus, jobMessage string) (err error)
	GetCronJobHistory() (cronJobHistory []models.CronJobHistory, err error)
}

func NewCronService(settings *config.Config, driver *neo4j.Driver) ICronService {
	return &CronService{neo4jDriver: driver, jwtSecret: settings.JwtSecret}
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
