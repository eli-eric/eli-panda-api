package securityService

import "panda/apigateway/helpers"

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

	result.Parameters["userName"] = username

	return result
}
