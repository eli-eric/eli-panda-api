package cronservice

import (
	"strings"

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

	// HAS_USER is created only when neither side is linked yet (one Employee = one User).
	// The best candidate is picked first (username = login, then username / email = any login)
	// and linked only when free, so a taken login account never falls back to another account.
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
	empl.superiorNumber = $superiorNumber,
	empl.superiorName = $superiorName,
	empl.jobPosition = $jobPosition,
	empl.workplaceName = $workplaceName,
	empl.jobPositionCode = $jobPositionCode,
	empl.orgUnitCode = $orgUnitCode,
	empl.ldapUser = $ldapUser
	WITH empl 
	MATCH(f:Facility{code: $facilityCode})
	MERGE(empl)-[:AFFILIATED_WITH_FACILITY]->(f)
	WITH empl 
	WHERE NOT (empl)-[:HAS_USER]->(:User)
	MATCH(usr:User)
	WHERE toLower(usr.username) IN $logins OR toLower(usr.email) IN $logins
	WITH empl, usr
	ORDER BY CASE
		WHEN toLower(usr.username) = $logins[0] THEN 0
		WHEN toLower(usr.username) IN $logins THEN 1
		ELSE 2
	END, usr.uid
	WITH empl, head(collect(usr)) AS usr
	WHERE NOT (usr)<-[:HAS_USER]-(:Employee)
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
	result.Parameters["orgUnitCode"] = employee.OJ_KOD
	result.Parameters["ldapUser"] = strings.TrimSpace(employee.LDAP_USER)
	result.Parameters["logins"] = EmployeeLogins(employee)
	result.Parameters["facilityCode"] = facilityCode

	fullName := employee.PRIJMENI + " " + employee.JMENO
	result.Parameters["fullName"] = fullName

	return result
}

func SetSystemTypesQuery() (result helpers.DatabaseQuery) {

	result.Query = `match (s:System)-[:CONTAINS_ITEM]->(itm:Item)-[:IS_BASED_ON]->(ci:CatalogueItem)-[:BELONGS_TO_CATEGORY]->(cc:CatalogueCategory)-[:HAS_SYSTEM_TYPE]->(st)
	with st, s where NOT (s)-[:HAS_SYSTEM_TYPE]->()
	MERGE(s)-[:HAS_SYSTEM_TYPE]->(st)
	return distinct count(s) as result;`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})

	return result
}

func CountEmployeesByNumberQuery(employeeNumber string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH (empl:Employee{employeeNumber: $employeeNumber}) RETURN count(empl) AS result`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["employeeNumber"] = employeeNumber

	return result
}

// old and misspelled domains users still have in HR, Azure AD logins use @eli-laser.eu
var legacyLoginDomains = map[string]bool{"eli-beams.eu": true, "eli-lasers.eu": true}

const loginDomain = "eli-laser.eu"

// NormalizeLogin lower-cases a login / e-mail and maps legacy domains to the current one
func NormalizeLogin(value string) string {
	login := strings.ToLower(strings.TrimSpace(value))

	name, domain, found := strings.Cut(login, "@")
	if found && legacyLoginDomains[domain] {
		return name + "@" + loginDomain
	}

	return login
}

// EmployeeLogins returns the identifiers to pair an Employee with a User, most preferred first:
// the Azure AD login (LDAP_USER, fallback EMAIL), then EMAIL. Both normalized, without blanks and duplicates.
func EmployeeLogins(employee models.SyncEliBeamlinesEmployee) []string {
	login := employee.LDAP_USER
	if strings.TrimSpace(login) == "" {
		login = employee.EMAIL
	}

	logins := []string{}
	for _, candidate := range []string{NormalizeLogin(login), NormalizeLogin(employee.EMAIL)} {
		if candidate != "" && (len(logins) == 0 || logins[0] != candidate) {
			logins = append(logins, candidate)
		}
	}

	return logins
}
