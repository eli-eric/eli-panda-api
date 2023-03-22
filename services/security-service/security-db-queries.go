package securityService

import (
	"panda/apigateway/helpers"
	"strings"
)

func UserWithRolesAndFailityQuery(username string) (result helpers.DatabaseQuery) {

	result.Query = `match(u:User{username: $userName})-[:HAS_ROLE]->(r:Role) 
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
	result.Parameters["userName"] = username

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
	searchText = strings.ToLower(searchText)
	result.Query = `
	MATCH(f:Facility{code:$facilityCode})
	WITH f
	MATCH(r:User)-[:BELONGS_TO_FACILITY]->(f)
	where toLower(r.lastName) contains $searchText or toLower(r.firstName) contains $searchText 
	RETURN {uid: r.uid,name: r.lastName + " " + r.firstName} as result 
	ORDER BY result.name limit $limit`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = searchText
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["limit"] = limit
	return result
}

func ChangeUserPasswordQuery(userUID string, newPassword string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(u:User{uid:$userUID}) SET u.passwordHash = $newPasswordHash RETURN u.uid AS result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["userUID"] = userUID
	return result
}
