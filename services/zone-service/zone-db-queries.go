package zoneservice

import (
	"fmt"
	"panda/apigateway/helpers"
)

// defaultParentSystemSubQuery resolves the zone's default parent system into the `dps` variable.
// It is a CALL subquery (not a plain OPTIONAL MATCH) so a duplicated relationship can never
// multiply the result rows of the outer query. The system is facility-scoped via $facilityCode,
// the same way every write/validation path is, so a cross-facility relationship never leaks
// into the response. Ordering is only a tiebreak for duplicated relationships - systemCode is
// null on the systems the CS zone migration creates, hence the uid fallback for determinism.
const defaultParentSystemSubQuery = `
				CALL {
					WITH z
					OPTIONAL MATCH (z)-[:HAS_DEFAULT_PARENT_SYSTEM]->(dps:System)-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
					WHERE (dps.deleted IS NULL OR dps.deleted <> true)
					RETURN dps ORDER BY dps.systemCode, dps.uid LIMIT 1
				}`

// System code is stored in the systemCode property, not code.
const defaultParentSystemProjection = `CASE WHEN dps IS NOT NULL THEN {uid: dps.uid, name: dps.name, code: dps.systemCode} ELSE null END`

func GetAllZonesQuery(facilityCode, search string, skip, limit int, sorting *[]helpers.Sorting) helpers.DatabaseQuery {
	query := `MATCH (z:Zone)-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				OPTIONAL MATCH (parent:Zone)-[:HAS_SUBZONE]->(z)
				WITH z, parent
				WHERE ($search = '' OR toLower(z.name) CONTAINS toLower($search) OR toLower(z.code) CONTAINS toLower($search)
					OR toLower(parent.name) CONTAINS toLower($search) OR toLower(parent.code) CONTAINS toLower($search))
				WITH z, parent`

	query += getZoneSortingClause(sorting)
	query += fmt.Sprintf(" SKIP %d LIMIT %d ", skip, limit)

	query += defaultParentSystemSubQuery

	query += ` RETURN {uid: z.uid, name: z.name, code: z.code, notes: coalesce(z.notes, ''),
						parentZone: CASE WHEN parent IS NOT NULL THEN {uid: parent.uid, name: parent.name, code: parent.code} ELSE null END,
						defaultParentSystem: ` + defaultParentSystemProjection + `} as zone`

	// the paging ORDER BY above is separated from this RETURN by the CALL subquery, and Cypher
	// only guarantees row order for a projection that carries its own ORDER BY - so repeat it
	query += getZoneSortingClause(sorting)

	return helpers.DatabaseQuery{
		Query:       query,
		ReturnAlias: "zone",
		Parameters: map[string]interface{}{
			"facilityCode": facilityCode,
			"search":       search,
		},
	}
}

func GetAllZonesCountQuery(facilityCode, search string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone)-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				OPTIONAL MATCH (parent:Zone)-[:HAS_SUBZONE]->(z)
				WITH z, parent
				WHERE ($search = '' OR toLower(z.name) CONTAINS toLower($search) OR toLower(z.code) CONTAINS toLower($search)
					OR toLower(parent.name) CONTAINS toLower($search) OR toLower(parent.code) CONTAINS toLower($search))
				RETURN count(z) as totalCount`,
		ReturnAlias: "totalCount",
		Parameters: map[string]interface{}{
			"facilityCode": facilityCode,
			"search":       search,
		},
	}
}

func getZoneSortingClause(sorting *[]helpers.Sorting) string {
	if sorting == nil || len(*sorting) == 0 {
		return " ORDER BY parent.code DESC, z.code ASC "
	}

	orderBy := " ORDER BY "
	for i, sort := range *sorting {
		sortField := mapZoneSortField(sort.ID)
		direction := helpers.GetSortingDirectionString(sort.DESC)
		if i > 0 {
			orderBy += ", "
		}
		orderBy += fmt.Sprintf("%s %s", sortField, direction)
	}
	return orderBy
}

func mapZoneSortField(fieldID string) string {
	fieldMap := map[string]string{
		"name":       "z.name",
		"code":       "z.code",
		"parentZone": "parent.name",
	}
	if mapped, ok := fieldMap[fieldID]; ok {
		return mapped
	}
	return "z.name"
}

func GetZoneByUIDQuery(uid, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				OPTIONAL MATCH (parent:Zone)-[:HAS_SUBZONE]->(z)` +
			defaultParentSystemSubQuery + `
				RETURN {uid: z.uid, name: z.name, code: z.code, notes: coalesce(z.notes, ''),
						parentZone: CASE WHEN parent IS NOT NULL THEN {uid: parent.uid, name: parent.name, code: parent.code} ELSE null END,
						defaultParentSystem: ` + defaultParentSystemProjection + `} as zone`,
		ReturnAlias: "zone",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"facilityCode": facilityCode,
		},
	}
}

func CheckZoneCodeExistsQuery(code, facilityCode, excludeUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{code:$code})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true) AND z.uid <> $excludeUID
				RETURN count(z) as cnt`,
		ReturnAlias: "cnt",
		Parameters: map[string]interface{}{
			"code":         code,
			"facilityCode": facilityCode,
			"excludeUID":   excludeUID,
		},
	}
}

func CheckParentIsRootQuery(parentUID, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (p:Zone{uid:$parentUID})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				WHERE (p.deleted IS NULL OR p.deleted <> true)
				OPTIONAL MATCH (gp:Zone)-[:HAS_SUBZONE]->(p)
				RETURN {uid: p.uid, hasParent: gp IS NOT NULL} as result`,
		ReturnAlias: "result",
		Parameters: map[string]interface{}{
			"parentUID":    parentUID,
			"facilityCode": facilityCode,
		},
	}
}

func CreateRootZoneQuery(uid, name, code, notes, facilityCode, userUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (f:Facility{code:$facilityCode})
				MATCH (u:User{uid:$userUID})
				CREATE (z:Zone{uid:$uid, name:$name, code:$code, notes:$notes, deleted:false})-[:BELONGS_TO_FACILITY]->(f)
				CREATE (z)-[:WAS_UPDATED_BY{at:datetime(), action:"INSERT"}]->(u)
				RETURN {uid:z.uid, name:z.name, code:z.code, notes:z.notes} as zone`,
		ReturnAlias: "zone",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"name":         name,
			"code":         code,
			"notes":        notes,
			"facilityCode": facilityCode,
			"userUID":      userUID,
		},
	}
}

func CreateSubZoneQuery(uid, name, code, notes, facilityCode, parentUID, userUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (f:Facility{code:$facilityCode})
				MATCH (parent:Zone{uid:$parentUID})-[:BELONGS_TO_FACILITY]->(f)
				WHERE (parent.deleted IS NULL OR parent.deleted <> true)
				MATCH (u:User{uid:$userUID})
				CREATE (z:Zone{uid:$uid, name:$name, code:$code, notes:$notes, deleted:false})-[:BELONGS_TO_FACILITY]->(f)
				CREATE (parent)-[:HAS_SUBZONE]->(z)
				CREATE (z)-[:WAS_UPDATED_BY{at:datetime(), action:"INSERT"}]->(u)
				RETURN {uid:z.uid, name:z.name, code:z.code, notes:z.notes,
						parentZone: {uid:parent.uid, name:parent.name, code:parent.code}} as zone`,
		ReturnAlias: "zone",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"name":         name,
			"code":         code,
			"notes":        notes,
			"facilityCode": facilityCode,
			"parentUID":    parentUID,
			"userUID":      userUID,
		},
	}
}

func UpdateZoneQuery(uid, name, code string, notes *string, facilityCode, userUID string) helpers.DatabaseQuery {
	q := helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				SET z.name = $name, z.code = $code, z.notes = coalesce($notes, z.notes)
				WITH z
				MATCH (u:User{uid:$userUID})
				CREATE (z)-[:WAS_UPDATED_BY{at:datetime(), action:"UPDATE"}]->(u)
				RETURN z.uid as uid`,
		ReturnAlias: "uid",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"name":         name,
			"code":         code,
			"notes":        nil,
			"facilityCode": facilityCode,
			"userUID":      userUID,
		},
	}

	if notes != nil {
		q.Parameters["notes"] = *notes
	}

	return q
}

func RemoveParentRelQuery(uid, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (parent:Zone)-[rel:HAS_SUBZONE]->(z:Zone{uid:$uid})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				DELETE rel`,
		Parameters: map[string]interface{}{
			"uid":          uid,
			"facilityCode": facilityCode,
		},
	}
}

func SetParentRelQuery(uid, parentUID, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				MATCH (parent:Zone{uid:$parentUID})-[:BELONGS_TO_FACILITY]->(f)
				WHERE (parent.deleted IS NULL OR parent.deleted <> true)
				MERGE (parent)-[:HAS_SUBZONE]->(z)`,
		Parameters: map[string]interface{}{
			"uid":          uid,
			"parentUID":    parentUID,
			"facilityCode": facilityCode,
		},
	}
}

func RemoveDefaultParentSystemRelQuery(uid, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				MATCH (z)-[rel:HAS_DEFAULT_PARENT_SYSTEM]->(:System)
				DELETE rel`,
		Parameters: map[string]interface{}{
			"uid":          uid,
			"facilityCode": facilityCode,
		},
	}
}

// SetDefaultParentSystemRelQuery replaces the zone's default parent system in a single query:
// the old relationship is only deleted once the new system has matched, so a soft-deleted or
// foreign system leaves the existing value untouched instead of silently detaching it.
func SetDefaultParentSystemRelQuery(uid, systemUID, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				MATCH (s:System{uid:$systemUID})-[:BELONGS_TO_FACILITY]->(f)
				WHERE (s.deleted IS NULL OR s.deleted <> true)
				WITH z, s
				OPTIONAL MATCH (z)-[rel:HAS_DEFAULT_PARENT_SYSTEM]->(other:System)
				WHERE other <> s
				DELETE rel
				WITH DISTINCT z, s
				MERGE (z)-[:HAS_DEFAULT_PARENT_SYSTEM]->(s)`,
		Parameters: map[string]interface{}{
			"uid":          uid,
			"systemUID":    systemUID,
			"facilityCode": facilityCode,
		},
	}
}

func CheckSystemExistsInFacilityQuery(systemUID, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (s:System{uid:$systemUID})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				WHERE (s.deleted IS NULL OR s.deleted <> true)
				RETURN count(s) as cnt`,
		ReturnAlias: "cnt",
		Parameters: map[string]interface{}{
			"systemUID":    systemUID,
			"facilityCode": facilityCode,
		},
	}
}

func CheckZoneHasSubzonesQuery(uid, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				MATCH (z)-[:HAS_SUBZONE]->(sub:Zone)
				WHERE (sub.deleted IS NULL OR sub.deleted <> true)
				RETURN count(sub) as cnt`,
		ReturnAlias: "cnt",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"facilityCode": facilityCode,
		},
	}
}

func CheckZoneHasSystemRefsQuery(uid, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})-[:BELONGS_TO_FACILITY]->(:Facility{code:$facilityCode})
				MATCH (s:System)-[:HAS_ZONE]->(z)
				WHERE (s.deleted IS NULL OR s.deleted <> true)
				RETURN count(s) as cnt`,
		ReturnAlias: "cnt",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"facilityCode": facilityCode,
		},
	}
}

func GetZoneByCodeAndFacilityQuery(code, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{code:$code})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				OPTIONAL MATCH (parent:Zone)-[:HAS_SUBZONE]->(z)` +
			defaultParentSystemSubQuery + `
				RETURN {uid: z.uid, name: z.name, code: z.code, notes: coalesce(z.notes, ''),
						parentZone: CASE WHEN parent IS NOT NULL THEN {uid: parent.uid, name: parent.name, code: parent.code} ELSE null END,
						defaultParentSystem: ` + defaultParentSystemProjection + `} as zone`,
		ReturnAlias: "zone",
		Parameters: map[string]interface{}{
			"code":         code,
			"facilityCode": facilityCode,
		},
	}
}
