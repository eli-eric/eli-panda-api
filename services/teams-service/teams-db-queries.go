package teamsservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/teams-service/models"
	"strings"
)

// GetAllTeamsQuery returns every team in the facility with its member count, sorted by name.
func GetAllTeamsQuery(facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (t:Team)-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				OPTIONAL MATCH (t)<-[:BELONGS_TO_TEAM]-(u:User)-[:BELONGS_TO_FACILITY]->(f)
				WITH t, count(u) AS memberCount
				RETURN {uid: t.uid, name: t.name, code: coalesce(t.code, ''),
						description: coalesce(t.description, ''), memberCount: memberCount} AS team
				ORDER BY team.name`,
		ReturnAlias: "team",
		Parameters: map[string]interface{}{
			"facilityCode": facilityCode,
		},
	}
}

// GetTeamByUIDQuery returns a facility-scoped team with its full inline member list
// (including disabled users, sorted by lastName then firstName).
func GetTeamByUIDQuery(uid, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (t:Team{uid:$uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				OPTIONAL MATCH (t)<-[:BELONGS_TO_TEAM]-(u:User)-[:BELONGS_TO_FACILITY]->(f)
				WITH t, u ORDER BY u.lastName, u.firstName
				WITH t, collect(CASE WHEN u IS NULL THEN null ELSE {
						uid: u.uid, firstName: coalesce(u.firstName, ''), lastName: coalesce(u.lastName, ''),
						username: coalesce(u.username, ''), email: coalesce(u.email, ''),
						isEnabled: coalesce(u.isEnabled, false)} END) AS members
				RETURN {uid: t.uid, name: t.name, code: coalesce(t.code, ''),
						description: coalesce(t.description, ''),
						members: [m IN members WHERE m IS NOT NULL]} AS team`,
		ReturnAlias: "team",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"facilityCode": facilityCode,
		},
	}
}

// CheckTeamCodeExistsQuery counts teams (excluding excludeUID) with the given code in the
// facility. Only invoked when code is non-empty.
func CheckTeamCodeExistsQuery(code, facilityCode, excludeUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (t:Team{code:$code})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				WHERE t.uid <> $excludeUID
				RETURN count(t) AS cnt`,
		ReturnAlias: "cnt",
		Parameters: map[string]interface{}{
			"code":         code,
			"facilityCode": facilityCode,
			"excludeUID":   excludeUID,
		},
	}
}

// CreateTeamQuery creates a facility-scoped team and its INSERT audit edge.
func CreateTeamQuery(uid, name, code, description, facilityCode, userUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (f:Facility{code:$facilityCode})
				MATCH (u:User{uid:$userUID})
				CREATE (t:Team{uid:$uid, name:$name, code:$code, description:$description})-[:BELONGS_TO_FACILITY]->(f)
				CREATE (t)-[:WAS_UPDATED_BY{at:datetime(), action:"INSERT"}]->(u)
				RETURN {uid: t.uid, name: t.name, code: coalesce(t.code, ''),
						description: coalesce(t.description, '')} AS team`,
		ReturnAlias: "team",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"name":         name,
			"code":         code,
			"description":  description,
			"facilityCode": facilityCode,
			"userUID":      userUID,
		},
	}
}

// UpdateTeamQuery is the PUT full replace. description may be nil (clears the property).
func UpdateTeamQuery(uid, name, code string, description *string, facilityCode, userUID string) helpers.DatabaseQuery {
	q := helpers.DatabaseQuery{
		Query: `MATCH (t:Team{uid:$uid})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				SET t.name = $name, t.code = $code, t.description = $description
				WITH t
				MATCH (u:User{uid:$userUID})
				CREATE (t)-[:WAS_UPDATED_BY{at:datetime(), action:"UPDATE"}]->(u)
				RETURN t.uid AS uid`,
		ReturnAlias: "uid",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"name":         name,
			"code":         code,
			"description":  nil,
			"facilityCode": facilityCode,
			"userUID":      userUID,
		},
	}
	if description != nil {
		q.Parameters["description"] = *description
	}
	return q
}

// PatchTeamQuery builds a partial-update query from the parsed PATCH fields. Only supplied
// keys are SET (absent = untouched; explicit null clears). Writes a UPDATE audit edge with
// the change delta. Callers must ensure at least one field is present.
func PatchTeamQuery(uid, facilityCode, userUID string, fields *models.PatchTeamFields, changesJSON string) helpers.DatabaseQuery {
	params := map[string]interface{}{
		"uid":          uid,
		"facilityCode": facilityCode,
		"userUID":      userUID,
		"changes":      changesJSON,
	}

	setClauses := make([]string, 0, 3)
	if fields.Name != nil {
		params["name"] = strings.TrimSpace(*fields.Name)
		setClauses = append(setClauses, "t.name = $name")
	}
	if fields.Code != nil {
		if fields.Code.Value != nil {
			params["code"] = strings.TrimSpace(*fields.Code.Value)
		} else {
			params["code"] = nil
		}
		setClauses = append(setClauses, "t.code = $code")
	}
	if fields.Description != nil {
		if fields.Description.Value != nil {
			params["description"] = *fields.Description.Value
		} else {
			params["description"] = nil
		}
		setClauses = append(setClauses, "t.description = $description")
	}

	query := `MATCH (t:Team{uid:$uid})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode}) `
	if len(setClauses) > 0 {
		query += "SET " + strings.Join(setClauses, ", ") + " "
	}
	query += `WITH t
			MATCH (u:User{uid:$userUID})
			CREATE (t)-[:WAS_UPDATED_BY{at:datetime(), action:"UPDATE", changes:$changes}]->(u)
			RETURN t.uid AS uid`

	return helpers.DatabaseQuery{
		Query:       query,
		ReturnAlias: "uid",
		Parameters:  params,
	}
}

// CheckTeamReferencesQuery returns the counts of Systems (HAS_RESPONSIBLE_TEAM) and RoomCards
// (HAS_TEAM) still referencing the team. Sequential OPTIONAL MATCH avoids cartesian inflation.
func CheckTeamReferencesQuery(uid, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (t:Team{uid:$uid})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				OPTIONAL MATCH (t)<-[:HAS_RESPONSIBLE_TEAM]-(s:System)
				WITH t, count(DISTINCT s) AS systemCount
				OPTIONAL MATCH (t)<-[:HAS_TEAM]-(rc:RoomCard)
				RETURN {systemCount: systemCount, roomCardCount: count(DISTINCT rc)} AS refs`,
		ReturnAlias: "refs",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"facilityCode": facilityCode,
		},
	}
}

// DeleteTeamQuery hard-deletes a facility-scoped team; DETACH removes membership + audit edges.
func DeleteTeamQuery(uid, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (t:Team{uid:$uid})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				DETACH DELETE t`,
		Parameters: map[string]interface{}{
			"uid":          uid,
			"facilityCode": facilityCode,
		},
	}
}

// ValidateTeamMemberUidsQuery returns the subset of userUids that are Users in the facility.
func ValidateTeamMemberUidsQuery(facilityCode string, userUids []string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (u:User)-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				WHERE u.uid IN $userUids
				RETURN collect(u.uid) AS validUids`,
		ReturnAlias: "validUids",
		Parameters: map[string]interface{}{
			"facilityCode": facilityCode,
			"userUids":     userUids,
		},
	}
}

// GetTeamMemberUidsQuery returns the current member uids of a team (for replace-all diffing).
func GetTeamMemberUidsQuery(teamUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (t:Team{uid:$teamUid})<-[:BELONGS_TO_TEAM]-(u:User)
				RETURN collect(u.uid) AS memberUids`,
		ReturnAlias: "memberUids",
		Parameters: map[string]interface{}{
			"teamUid": teamUID,
		},
	}
}

// MergeTeamMembersQuery idempotently links same-facility users to the team.
func MergeTeamMembersQuery(teamUID, facilityCode string, userUids []string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (t:Team{uid:$teamUid})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				MATCH (u:User)-[:BELONGS_TO_FACILITY]->(f)
				WHERE u.uid IN $userUids
				MERGE (u)-[:BELONGS_TO_TEAM]->(t)`,
		Parameters: map[string]interface{}{
			"teamUid":      teamUID,
			"facilityCode": facilityCode,
			"userUids":     userUids,
		},
	}
}

// RemoveTeamMembersQuery deletes membership edges for the given user uids.
func RemoveTeamMembersQuery(teamUID string, userUids []string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (u:User)-[r:BELONGS_TO_TEAM]->(t:Team{uid:$teamUid})
				WHERE u.uid IN $userUids
				DELETE r`,
		Parameters: map[string]interface{}{
			"teamUid":  teamUID,
			"userUids": userUids,
		},
	}
}

// TeamMembershipAuditQuery writes a single UPDATE audit edge on the team with the membership delta.
func TeamMembershipAuditQuery(teamUID, userUID, changesJSON string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (t:Team{uid:$teamUid})
				MATCH (u:User{uid:$userUID})
				CREATE (t)-[:WAS_UPDATED_BY{at:datetime(), action:"UPDATE", changes:$changes}]->(u)`,
		Parameters: map[string]interface{}{
			"teamUid": teamUID,
			"userUID": userUID,
			"changes": changesJSON,
		},
	}
}
