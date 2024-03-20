package securityService

import (
	"panda/apigateway/helpers"
	"strings"
)

func UserWithRolesAndFailityQuery(username string) (result helpers.DatabaseQuery) {

	result.Query = `match(u:User)-[:HAS_ROLE]->(r:Role) WHERE toLower(u.username) = $userName
	optional match(u)-[:BELONGS_TO_FACILITY]->(f)
	return {
		uid: u.uid,
		passwordHash: u.passwordHash, 
		lastName: u.lastName ,
		firstName: u.firstName,
		email: u.email, 
		facility: f.name,
		facilityCode: f.code,
		isEnabled: u.isEnabled,
		roles: collect(r.code)} as userInfo`

	result.ReturnAlias = "userInfo"

	result.Parameters = make(map[string]interface{})
	result.Parameters["userName"] = strings.ToLower(username)

	return result
}

func GetUsersCodebookQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(f:Facility{code:$facilityCode})
	WITH f MATCH(r:User)-[:BELONGS_TO_FACILITY]->(f) RETURN {uid: r.uid,name: r.lastName + " " + r.firstName} as result ORDER BY result.name`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetUsersAutocompleteCodebookQuery(searchText string, limit int, facilityCode string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH(f:Facility{code:$facilityCode})
	WITH f
	MATCH(r:User)-[:BELONGS_TO_FACILITY]->(f)
	where apoc.text.clean(r.lastName) contains apoc.text.clean($searchText) or apoc.text.clean(r.firstName) contains apoc.text.clean($searchText) 
	RETURN {uid: r.uid,name: r.lastName + " " + r.firstName} as result 
	ORDER BY result.name limit $limit`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = searchText
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["limit"] = limit
	return result
}

func ChangeUserPasswordQuery(userUID string, newPasswordHash string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(u:User{uid:$userUID}) SET u.passwordHash = $newPasswordHash RETURN u.uid AS result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["userUID"] = userUID
	return result
}

func GetEmployeesAutocompleteCodebookQuery(searchText string, limit int, facilityCode string, getAllEmployees bool, flags ...string) (result helpers.DatabaseQuery) {

	if getAllEmployees {
		result.Query = `
		MATCH(r:Employee)-[:AFFILIATED_WITH_FACILITY]->(f)
		where (apoc.text.clean(r.lastName) contains apoc.text.clean($searchText) or apoc.text.clean(r.firstName) contains apoc.text.clean($searchText)) `

	} else {
		result.Query = `
		MATCH(f:Facility{code:$facilityCode})
		WITH f
		MATCH(r:Employee)-[:AFFILIATED_WITH_FACILITY]->(f)
		where (apoc.text.clean(r.lastName) contains apoc.text.clean($searchText) or apoc.text.clean(r.firstName) contains apoc.text.clean($searchText)) `
	}

	for _, flag := range flags {
		result.Query += ` AND r.` + flag + ` = true `
	}

	if getAllEmployees {
		result.Query += `
		RETURN {uid: r.uid,name: r.lastName + " " + r.firstName + " (" + f.name + ")" } as result 
		ORDER BY f.name, result.name limit $limit`
	} else {
		result.Query += `
		RETURN {uid: r.uid,name: r.lastName + " " + r.firstName} as result 
		ORDER BY result.name limit $limit`
	}

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = searchText
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["limit"] = limit
	return result
}

func GetTeamsAutocompleteCodebookQuery(searchText string, limit int, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(f:Facility{code:$facilityCode})
	WITH f
	MATCH(r:Team)-[:BELONGS_TO_FACILITY]->(f)
	where apoc.text.clean(r.name) contains apoc.text.clean($searchText)
	RETURN {uid: r.uid,name: r.name} as result
	ORDER BY result.name limit $limit`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})

	result.Parameters["searchText"] = searchText
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["limit"] = limit
	return result
}

func GetContactPersonRolesAutocompleteCodebookQuery(searchText string, limit int, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(f:Facility{code:$facilityCode})
	WITH f
	MATCH(r:ContactPersonRole)-[:BELONGS_TO_FACILITY]->(f)
	where apoc.text.clean(r.name) contains apoc.text.clean($searchText)
	RETURN {uid: r.uid,name: r.name} as result
	ORDER BY result.name limit $limit`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})

	result.Parameters["searchText"] = searchText
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["limit"] = limit
	return result
}
