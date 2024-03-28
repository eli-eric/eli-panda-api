package cronservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/cron-service/models"
)

func LogCronJobStartQuery(jobUid, jobType, jobName string) (result helpers.DatabaseQuery) {

	result.Query = `
	CREATE (c:CronJobHistory{
		uid: apoc.create.uuid(),
		jobUid: $jobUid,
		jobName: $jobName,
		jobType: $jobType,
		jobStatus: $jobStatus,
		jobMessage: 'Job started',
		startTime: timestamp(),
		endTime: null})
	RETURN c.uid AS result`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["jobName"] = jobName
	result.Parameters["jobType"] = jobType
	result.Parameters["jobUid"] = jobUid
	result.Parameters["jobStatus"] = models.CRON_JOB_STATUS_RUNNING

	return result
}

func LogCronJobEndQuery(uid, jobStatus, jobMessage string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH (c:CronJobHistory{uid: $uid})
	WITH c, timestamp() AS endTime
	SET c.jobStatus = $jobStatus, c.jobMessage = $jobMessage, c.endTime = endTime, c.totalTimeMs = endTime - c.startTime
	RETURN c.uid AS result`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["jobStatus"] = jobStatus
	result.Parameters["jobMessage"] = jobMessage

	return result
}

func GetCronJobHistoryQuery() (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH (c:CronJobHistory)
	RETURN 
	{
		uid: c.uid,
		jobUid: c.jobUid,
		jobName: c.jobName,
		jobType: c.jobType,
		jobStatus: c.jobStatus,
		jobMessage: c.jobMessage,
		startTime: c.startTime,
		endTime: c.endTime,
		totalTimeMs: c.totalTimeMs
	} AS result
	ORDER BY c.startTime DESC
	LIMIT 100`

	result.ReturnAlias = "result"

	return result
}
