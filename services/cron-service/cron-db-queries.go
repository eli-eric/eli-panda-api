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

func SyncEliBeamlinesEmployeeQuery(employee models.SyncEliBeamlinesEmployee, facilityCode string) (result helpers.DatabaseQuery) {

	result.Query = `MERGE(empl:Employee{employeeNumber: $employeeNumber})
	WITH empl, case when empl.uid IS NULL THEN apoc.create.uuid() ELSE empl.uid END as uid
	SET 
	empl.uid = uid, 
	empl.email = toLower($email),
	empl.firstName = $firstName,
	empl.lastName = $lastName,
	empl.fullName = $fullName,
	empl.phone1 = $phone1,
	empl.phone2 = $phone2,
	empl.phoneNumber = apoc.text.join($phoneNumbers, ","),	
	empl.superiorNumber = $superiorNumber,
	empl.superiorName = $superiorName,
	empl.jobPosition = $jobPosition,
	empl.workplaceName = $workplaceName,
	empl.jobPositionCode = $jobPositionCode
	WITH empl 
	MATCH(f:Facility{code: $facilityCode})
	MERGE(empl)-[:AFFILIATED_WITH_FACILITY]->(f)
	WITH empl 
	MATCH(usr:User{username: empl.email})
	MERGE(empl)-[:HAS_USER]->(usr) 
	RETURN true AS result`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["employeeNumber"] = employee.OSOBNI_CISLO
	result.Parameters["email"] = employee.EMAIL
	result.Parameters["firstName"] = employee.JMENO
	result.Parameters["lastName"] = employee.PRIJMENI
	result.Parameters["phone1"] = employee.PHONE4
	result.Parameters["phone2"] = employee.PHONE5
	result.Parameters["superiorNumber"] = employee.NADRIZENY_OS_CISLO
	result.Parameters["superiorName"] = employee.NADRIZENY_JMENO
	result.Parameters["jobPosition"] = employee.PRACOVNI_MISTO_NAZEV
	result.Parameters["workplaceName"] = employee.PRACOVISTE
	result.Parameters["jobPositionCode"] = employee.PRACOVNI_MISTO_KOD
	result.Parameters["facilityCode"] = facilityCode

	phoneNumbers := []string{}
	if employee.PHONE4 != "" {
		phoneNumbers = append(phoneNumbers, employee.PHONE4)
	}
	if employee.PHONE5 != "" {
		phoneNumbers = append(phoneNumbers, employee.PHONE5)
	}
	result.Parameters["phoneNumbers"] = phoneNumbers

	fullName := employee.PRIJMENI + " " + employee.JMENO
	result.Parameters["fullName"] = fullName

	return result
}
