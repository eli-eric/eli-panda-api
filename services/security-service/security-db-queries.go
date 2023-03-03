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
	result.Query = `MATCH(r:User) RETURN {uid: r.uid,name: r.lastName + " " + r.firstName} as result ORDER BY result.name`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetUsersAutocompleteCodebookQuery(searchText string, limit int, facilityCode string) (result helpers.DatabaseQuery) {
	searchText = strings.ToLower(searchText)
	result.Query = `
	MATCH(r:User) 
	where toLower(r.lastName) contains $searchText or toLower(r.firstName) contains $searchText 
	RETURN {uid: r.uid,name: r.lastName + " " + r.firstName} as result 
	ORDER BY result.name limit $limit`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = searchText
	result.Parameters["limit"] = limit
	return result
}
