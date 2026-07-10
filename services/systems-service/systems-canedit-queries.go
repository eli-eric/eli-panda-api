package systemsService

import "panda/apigateway/helpers"

// CanEditSystemQuery decides whether userUID may edit the system, and returns the responsible
// users (contacts). Responsibility bubbles UP the HAS_SUBSYSTEM chain (union of the system and
// all non-deleted ancestors). A responsible is a login-capable User reached either via
// HAS_RESPONSIBLE->Employee->HAS_USER or via a HAS_RESPONSIBLE_TEAM team the user BELONGS_TO_TEAM.
//
// Result precedence: admin -> true; no systems-edit role -> false; no responsibles anywhere
// (orphan) -> true; otherwise true iff the caller is among the responsibles. Responsibles are
// always returned (even when result is false) so the FE can show who to contact.
//
// The caller is matched with OPTIONAL MATCH so a missing User node (stale JWT / deleted user)
// yields hasEdit=false -> result=false (fail-closed), while only a missing/deleted SYSTEM (the
// required MATCH below) produces ERR_NO_ROWS.
func CanEditSystemQuery(systemUID, userUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `OPTIONAL MATCH (caller:User{uid:$userUID})
				OPTIONAL MATCH (caller)-[:HAS_ROLE]->(adminR:Role{code:'admin'})
				OPTIONAL MATCH (caller)-[:HAS_ROLE]->(editR:Role{code:'systems-edit'})
				WITH caller, count(adminR) > 0 AS isAdmin, count(editR) > 0 AS hasEdit
				MATCH (sys:System{uid:$systemUID, deleted:false})
				OPTIONAL MATCH ancestorPath = (ancestor:System)-[:HAS_SUBSYSTEM*1..50]->(sys)
				WHERE all(n IN nodes(ancestorPath) WHERE n.deleted = false)
				WITH isAdmin, hasEdit, [sys] + collect(DISTINCT ancestor) AS chain
				UNWIND chain AS s
				OPTIONAL MATCH (s)-[:HAS_RESPONSIBLE]->(:Employee)-[:HAS_USER]->(ru:User)
				OPTIONAL MATCH (s)-[:HAS_RESPONSIBLE_TEAM]->(:Team)<-[:BELONGS_TO_TEAM]-(tu:User)
				WITH isAdmin, hasEdit,
					[u IN apoc.coll.toSet(collect(DISTINCT ru) + collect(DISTINCT tu)) WHERE u IS NOT NULL] AS respUsers
				RETURN {
					result: CASE
						WHEN isAdmin THEN true
						WHEN NOT hasEdit THEN false
						WHEN size(respUsers) = 0 THEN true
						ELSE any(u IN respUsers WHERE u.uid = $userUID)
					END,
					responsibles: [u IN respUsers | {
						uid: u.uid,
						firstName: coalesce(u.firstName, ''),
						lastName: coalesce(u.lastName, ''),
						username: coalesce(u.username, ''),
						email: coalesce(u.email, '')
					}]
				} AS canEdit`,
		ReturnAlias: "canEdit",
		Parameters: map[string]interface{}{
			"systemUID": systemUID,
			"userUID":   userUID,
		},
	}
}

// GetSystemUIDByItemUIDQuery resolves the System that contains a given physical item, so the
// edit-guard can be applied to physical-item endpoints that identify the item, not the system.
func GetSystemUIDByItemUIDQuery(itemUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (sys:System{deleted:false})-[:CONTAINS_ITEM]->(item:Item{uid:$itemUID})
				RETURN sys.uid AS systemUID`,
		ReturnAlias: "systemUID",
		Parameters: map[string]interface{}{
			"itemUID": itemUID,
		},
	}
}
