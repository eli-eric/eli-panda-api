package zoneservice

import (
	"fmt"
	"panda/apigateway/helpers"
)

func GetAllZonesQuery(facilityCode, search string, skip, limit int, sorting *[]helpers.Sorting) helpers.DatabaseQuery {
	query := `MATCH (z:Zone)-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				AND ($search = '' OR toLower(z.name) CONTAINS toLower($search) OR toLower(z.code) CONTAINS toLower($search))
				OPTIONAL MATCH (parent:Zone)-[:HAS_SUBZONE]->(z)
				WITH z, parent`

	query += getZoneSortingClause(sorting)
	query += fmt.Sprintf(" SKIP %d LIMIT %d ", skip, limit)

	query += ` RETURN {uid: z.uid, name: z.name, code: z.code,
						parentZone: CASE WHEN parent IS NOT NULL THEN {uid: parent.uid, name: parent.name, code: parent.code} ELSE null END} as zone`

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
				AND ($search = '' OR toLower(z.name) CONTAINS toLower($search) OR toLower(z.code) CONTAINS toLower($search))
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
		return " ORDER BY coalesce(parent.code, z.code), z.code "
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
	return "z." + fieldID
}

func GetZoneByUIDQuery(uid, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				OPTIONAL MATCH (parent:Zone)-[:HAS_SUBZONE]->(z)
				RETURN {uid: z.uid, name: z.name, code: z.code,
						parentZone: CASE WHEN parent IS NOT NULL THEN {uid: parent.uid, name: parent.name, code: parent.code} ELSE null END} as zone`,
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

func CheckParentIsRootQuery(parentUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (p:Zone{uid:$parentUID})
				WHERE (p.deleted IS NULL OR p.deleted <> true)
				OPTIONAL MATCH (gp:Zone)-[:HAS_SUBZONE]->(p)
				RETURN {uid: p.uid, hasParent: gp IS NOT NULL} as result`,
		ReturnAlias: "result",
		Parameters: map[string]interface{}{
			"parentUID": parentUID,
		},
	}
}

func CreateRootZoneQuery(uid, name, code, facilityCode, userUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (f:Facility{code:$facilityCode})
				CREATE (z:Zone{uid:$uid, name:$name, code:$code, deleted:false})-[:BELONGS_TO_FACILITY]->(f)
				WITH z
				MATCH (u:User{uid:$userUID})
				CREATE (z)-[:WAS_UPDATED_BY{at:datetime(), action:"INSERT"}]->(u)
				RETURN {uid:z.uid, name:z.name, code:z.code} as zone`,
		ReturnAlias: "zone",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"name":         name,
			"code":         code,
			"facilityCode": facilityCode,
			"userUID":      userUID,
		},
	}
}

func CreateSubZoneQuery(uid, name, code, facilityCode, parentUID, userUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (f:Facility{code:$facilityCode})
				MATCH (parent:Zone{uid:$parentUID})
				CREATE (z:Zone{uid:$uid, name:$name, code:$code, deleted:false})-[:BELONGS_TO_FACILITY]->(f)
				CREATE (parent)-[:HAS_SUBZONE]->(z)
				WITH z
				MATCH (u:User{uid:$userUID})
				CREATE (z)-[:WAS_UPDATED_BY{at:datetime(), action:"INSERT"}]->(u)
				RETURN {uid:z.uid, name:z.name, code:z.code, parentUid:$parentUID} as zone`,
		ReturnAlias: "zone",
		Parameters: map[string]interface{}{
			"uid":          uid,
			"name":         name,
			"code":         code,
			"facilityCode": facilityCode,
			"parentUID":    parentUID,
			"userUID":      userUID,
		},
	}
}

func UpdateZoneQuery(uid, name, code, userUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				SET z.name = $name, z.code = $code
				WITH z
				MATCH (u:User{uid:$userUID})
				CREATE (z)-[:WAS_UPDATED_BY{at:datetime(), action:"UPDATE"}]->(u)
				RETURN z.uid as uid`,
		ReturnAlias: "uid",
		Parameters: map[string]interface{}{
			"uid":     uid,
			"name":    name,
			"code":    code,
			"userUID": userUID,
		},
	}
}

func RemoveParentRelQuery(uid string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (parent:Zone)-[rel:HAS_SUBZONE]->(z:Zone{uid:$uid})
				DELETE rel`,
		Parameters: map[string]interface{}{
			"uid": uid,
		},
	}
}

func SetParentRelQuery(uid, parentUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})
				MATCH (parent:Zone{uid:$parentUID})
				MERGE (parent)-[:HAS_SUBZONE]->(z)`,
		Parameters: map[string]interface{}{
			"uid":       uid,
			"parentUID": parentUID,
		},
	}
}

func CheckZoneHasSubzonesQuery(uid string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{uid:$uid})-[:HAS_SUBZONE]->(sub:Zone)
				WHERE (sub.deleted IS NULL OR sub.deleted <> true)
				RETURN count(sub) as cnt`,
		ReturnAlias: "cnt",
		Parameters: map[string]interface{}{
			"uid": uid,
		},
	}
}

func CheckZoneHasSystemRefsQuery(uid string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (s:System)-[:HAS_ZONE]->(z:Zone{uid:$uid})
				WHERE (s.deleted IS NULL OR s.deleted <> true)
				RETURN count(s) as cnt`,
		ReturnAlias: "cnt",
		Parameters: map[string]interface{}{
			"uid": uid,
		},
	}
}

func GetZoneByCodeAndFacilityQuery(code, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (z:Zone{code:$code})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode})
				WHERE (z.deleted IS NULL OR z.deleted <> true)
				OPTIONAL MATCH (parent:Zone)-[:HAS_SUBZONE]->(z)
				RETURN {uid: z.uid, name: z.name, code: z.code,
						parentZone: CASE WHEN parent IS NOT NULL THEN {uid: parent.uid, name: parent.name, code: parent.code} ELSE null END} as zone`,
		ReturnAlias: "zone",
		Parameters: map[string]interface{}{
			"code":         code,
			"facilityCode": facilityCode,
		},
	}
}
