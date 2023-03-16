package systemsService

import (
	"fmt"
	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

func GetSystemTypesCodebookQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = fmt.Sprintf(`MATCH (n:SystemTypeGroup)-[:CONTAINS_SYSTEM_TYPE]->(st) with n, st order by st.name 
	return {uid:st.uid, name: n.name+ " > "+ st.name + case when st.mask%v is null then "" else  " (" + st.mask%v  + ")" end } as result order by result.name`, facilityCode, facilityCode)
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetSystemImportancesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:SystemImportance) RETURN {uid: r.uid,name:r.name} as result ORDER BY result.code`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetSystemCriticalityCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:SystemCriticality) RETURN {uid: r.uid,name:r.name} as result ORDER BY result.code`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetItemUsagesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:ItemUsage) RETURN {uid: r.uid,name:r.name} as result ORDER BY result.code`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetItemConditionsCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:ItemCondition) RETURN {uid: r.uid,name:r.name} as result ORDER BY result.code`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetLocationsBySearchTextQuery(searchText string, limit int, facilityCode string) (result helpers.DatabaseQuery) {
	searchText = strings.ToLower(searchText)
	result.Query = `
	MATCH (n:Location)-[:BELONGS_TO_FACILITY]->(f) where f.code = $facilityCode and n.code is not null and not (n)-[:HAS_SUBLOCATION]->()
	with n 
	where (toLower(n.code) contains $searchText or toLower(n.name) contains $searchText) 
	optional match (parent)-[:HAS_SUBLOCATION*1..50]->(n) 
	with n, collect(parent.name) as parentNames
	return {uid: n.code, name: n.code + " - " +  n.name + " - " + apoc.text.join(reverse(parentNames), " > ")} as result
	order by result.name 
	limit $limit`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = searchText
	result.Parameters["limit"] = limit
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetZonesCodebookQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(f:Facility{code:$facilityCode}) WITH f
	MATCH(z:Zone)-[:HAS_SUBZONE]->(sz)-[:BELONGS_TO_FACILITY]->(f) return {uid:sz.uid, name: z.code+"-"+sz.code + " - " + sz.name + " ("+  z.name + ")"} as zone order by z.code, sz.code
		UNION
		MATCH(f:Facility{code:$facilityCode}) WITH f
		WITH f
		MATCH(z:Zone)-[:BELONGS_TO_FACILITY]->(f) where not ()-[:HAS_SUBZONE]->(z)  return {uid:z.uid, name:z.code + " - " +z.name } as zone order by z.code`
	result.ReturnAlias = "zone"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetSubSystemsQuery(parentUID string, facilityCode string) (result helpers.DatabaseQuery) {

	//we have to diff queries depend if it is or not a root parent
	if parentUID != "" {
		result.Query = `
		MATCH(r:System)-[:BELONGS_TO_FACILITY]->(f) WHERE f.code = $facilityCode WITH r			
		MATCH (parent)-[:HAS_SUBSYSTEM]->(r)
		where parent.uid = $parentUID
		return {uid: r.uid, name: r.name} as result;`
	} else {
		result.Query = `
		MATCH(r:System)-[:BELONGS_TO_FACILITY]->(f)			
		where not ()-[:HAS_SUBSYSTEM]->(r) and f.code = $facilityCode
		return {uid: r.uid, name: r.name} as result;`
	}

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["parentUID"] = parentUID
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func SystemImageByUidQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `match(r:System{uid: $uid})	
	return r.image as image`
	result.ReturnAlias = "image"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func SystemDetailQuery(uid string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:System{uid: $uid})-[:BELONGS_TO_FACILITY]->(f) WHERE f.code = $facilityCode
	WITH r,f
OPTIONAL MATCH(r)-[:HAS_LOCATION]->(l)
OPTIONAL MATCH(r)-[:HAS_ZONE]->(z)
OPTIONAL MATCH(r)-[:HAS_SYSTEM_TYPE]->(st)
OPTIONAL MATCH(r)-[:HAS_IMPORTANCE]->(imp)
OPTIONAL MATCH(r)-[:HAS_OWNER]->(own)
OPTIONAL MATCH(r)-[:HAS_CRITICALITY]->(cc)
OPTIONAL MATCH(r)-[:CONTAINS_ITEM]->(itm)
OPTIONAL MATCH(parent)-[:HAS_SUBSYSTEM*1..50]->(r)
WITH r,l, z, st,itm, imp, own,cc, case when parent is not null then collect({uid: parent.uid, name: parent.name}) else null end as parents
WITH r,l, z, st,itm, imp, own,cc, reverse(parents) as parents
RETURN {
    uid: r.uid, 
    name: r.name, 
    description: r.description,
    location: case when l is not null then {uid: l.code, name: l.name} else null end, 
    systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,
    systemCode: r.systemCode,
    systemALias: r.systemAlias,
    importance: case when imp is not null then {uid: imp.uid, name: imp.name} else null end,
    owner: case when own is not null then {uid: own.uid, name: own.lastName + " " + own.firstName} else null end,
    zone: case when z is not null then {uid: z.uid, name: z.name} else null end,
    parentPath: parents,
	itemUID: itm.uid    
    } as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func SystemFormDetailQuery(uid string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:System{uid: $uid})-[:BELONGS_TO_FACILITY]->(f) WHERE f.code = $facilityCode
	WITH r,f
OPTIONAL MATCH(r)-[:HAS_LOCATION]->(l)
OPTIONAL MATCH(r)-[:HAS_ZONE]->(z)
OPTIONAL MATCH(r)-[:HAS_SYSTEM_TYPE]->(st)
OPTIONAL MATCH(r)-[:HAS_IMPORTANCE]->(imp)
OPTIONAL MATCH(r)-[:HAS_OWNER]->(own)
OPTIONAL MATCH(r)-[:HAS_CRITICALITY]->(cc)
OPTIONAL MATCH(r)-[:CONTAINS_ITEM]->(itm)
OPTIONAL MATCH(parent)-[:HAS_SUBSYSTEM]->(r)
WITH r,l, z, st,itm, imp, own,cc, parent
RETURN {
    uid: r.uid, 
    name: r.name, 
    description: r.description,
    locationUID: case when l is not null then l.code else null end, 
    systemTypeUID: case when st is not null then st.uid else null end,
    systemCode: r.systemCode,
    systemAlias: r.systemAlias,
    importanceUID: case when imp is not null then imp.uid else null end,
    ownerUID: case when own is not null then own.uid else null end,
    zoneUID: case when z is not null then  z.uid else null end,
    parentUID: case when parent is not null then parent.uid else null end,
	itemUID: itm.uid    
    } as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func CreateNewSystemQuery(newSystem *models.SystemForm, facilityCode string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["uid"] = uuid.NewString()
	result.Parameters["name"] = newSystem.Name
	result.Parameters["description"] = newSystem.Description
	result.Parameters["systemCode"] = newSystem.SystemCode
	result.Parameters["systemAlias"] = newSystem.SystemAlias

	result.Query = `
	CREATE(s:System{uid: $uid})
	SET 
	s.name = $name, 
	s.description = $description,
	s.systemCode = $systemCode,
	s.systemAlias = $systemAlias
	WITH s
	MATCH(f:Facility{code: $facilityCode})
	CREATE(s)-[:BELONGS_TO_FACILITY]->(f)
	WITH s
	`

	if newSystem.ParentUID != nil && len(*newSystem.ParentUID) > 0 {
		result.Query += `WITH s MATCH(parent:System{uid:$parentUID}) MERGE(parent)-[:HAS_SUBSYSTEM]->(s) `
		result.Parameters["parentUID"] = newSystem.ParentUID
	}

	if newSystem.ZoneUID != nil && len(*newSystem.ZoneUID) > 0 {
		result.Query += `WITH s MATCH(z:Zone{uid:$zoneUID}) MERGE(s)-[:HAS_ZONE]->(z) `
		result.Parameters["zoneUID"] = newSystem.ZoneUID
	}

	if newSystem.LocationUID != nil && len(*newSystem.LocationUID) > 0 {
		result.Query += `WITH s MATCH(l:Location{code:$locationUID})-[:BELONGS_TO_FACILITY]->(f{code:$facilityCode}) MERGE(s)-[:HAS_LOCATION]->(l) `
		result.Parameters["locationUID"] = newSystem.LocationUID
	}

	if newSystem.SystemTypeUID != nil && len(*newSystem.SystemTypeUID) > 0 {
		result.Query += `WITH s MATCH(st:SystemType{uid:$systemTypeUID}) MERGE(s)-[:HAS_SYSTEM_TYPE]->(st) `
		result.Parameters["systemTypeUID"] = newSystem.SystemTypeUID
	}

	if newSystem.OwnerUID != nil && len(*newSystem.OwnerUID) > 0 {
		result.Query += `WITH s MATCH(owner:User{uid:$ownerUID}) MERGE(s)-[:HAS_OWNER]->(owner) `
		result.Parameters["ownerUID"] = newSystem.OwnerUID
	}

	if newSystem.ImportanceUID != nil && len(*newSystem.ImportanceUID) > 0 {
		result.Query += `WITH s MATCH(imp:SystemImportance{uid:$importanceUID}) MERGE(s)-[:HAS_IMPORTANCE]->(imp) `
		result.Parameters["importanceUID"] = newSystem.ImportanceUID
	}

	if newSystem.Image != nil && len(*newSystem.Image) > 0 {
		result.Query += `WITH s SET s.image = $image `
		result.Parameters["image"] = newSystem.Image
	}

	result.Query += `RETURN s.uid as result`

	result.ReturnAlias = "result"

	return result
}

func testr(newSystem any, oldSystem any, facilityCode string) {

	newObj := reflect.TypeOf(newSystem)

	for i := 0; i < newObj.NumField(); i++ {
		fld := newObj.Field(i)
		fmt.Println(fld.Name, " ", fld.Tag.Get("neo4j"), " ", fld.Type)
	}

}

func UpdateSystemQuery(newSystem *models.SystemForm, oldSystem *models.SystemForm, facilityCode string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["uid"] = oldSystem.UID

	result.Query = `MATCH(s:System{uid:$uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode}) `

	newObj := reflect.TypeOf(*newSystem)
	oldObj := reflect.TypeOf(*oldSystem)
	newValObj := reflect.ValueOf(*newSystem)
	oldValObj := reflect.ValueOf(*oldSystem)

	if newObj == oldObj {
		for i := 0; i < newObj.NumField(); i++ {

			newField := newObj.Field(i)
			oldField := oldObj.Field(i)
			neo4jTags := strings.Split(newField.Tag.Get("neo4j"), ",")
			fieldType := newField.Type.String()

			if len(neo4jTags) > 0 {

				neo4jPropType := neo4jTags[0]
				if neo4jPropType == "prop" {
					neo4jPropName := neo4jTags[1]

					if fieldType == "string" {
						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name).String()
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name).String()

						if newValue != oldValue {
							result.Parameters[neo4jPropName] = newValue
							result.Query += fmt.Sprintf(`WITH s SET s.%v=$%v `, neo4jPropName, neo4jPropName)
						}

					} else if fieldType == "*string" {
						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name)

						if newValue != oldValue {
							if newValue.IsNil() {
								result.Parameters[neo4jPropName] = nil
							} else {
								result.Parameters[neo4jPropName] = newValue.Elem().String()
							}

							result.Query += fmt.Sprintf(`WITH s SET s.%v=$%v `, neo4jPropName, neo4jPropName)
						}
					}
				} else if neo4jPropType == "rel" {
					neo4jLabel := neo4jTags[1]
					neo4jRelationType := neo4jTags[2]
					neo4jID := neo4jTags[3]
					neo4jAlias := neo4jTags[4]

					if fieldType == "*string" {
						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name)

						if !newValue.IsNil() && newValue.Elem().String() != "" && oldValue.IsNil() {
							result.Query += fmt.Sprintf(`WITH s MATCH(%v:%v{%v:$%v}) MERGE(s)-[:%v]->(%v) `, neo4jAlias, neo4jLabel, neo4jID, newField.Name, neo4jRelationType, neo4jAlias)
							result.Parameters[newField.Name] = newValue.Elem().String()
						} else if !newValue.IsNil() && newValue.Elem().String() != "" && !oldValue.IsNil() && newValue.Elem().String() != oldValue.Elem().String() {
							result.Query += fmt.Sprintf(`WITH s MATCH(s)-[r%v:%v]->(%v) delete r%v `, neo4jAlias, neo4jRelationType, neo4jAlias, neo4jAlias)
							result.Query += fmt.Sprintf(`WITH s MATCH(%v:%v{%v:$%v}) MERGE(s)-[:%v]->(%v) `, neo4jAlias, neo4jLabel, neo4jID, newField.Name, neo4jRelationType, neo4jAlias)
							result.Parameters[newField.Name] = newValue.Elem().String()
						} else if (newValue.IsNil() || newValue.Elem().String() == "") && !oldValue.IsNil() {
							result.Query += fmt.Sprintf(`WITH s MATCH(s)-[r%v:%v]->(%v) delete r%v `, neo4jAlias, neo4jRelationType, neo4jAlias, neo4jAlias)
						}
					}

				}
			}
		}
	}

	// if newSystem.Name != oldSystem.Name {
	// 	result.Parameters["name"] = newSystem.Name
	// 	result.Query += `WITH s SET s.name=$name `
	// }

	// if *newSystem.Description != *oldSystem.Description {
	// 	result.Parameters["description"] = newSystem.Description
	// 	result.Query += `WITH s SET s.description=$description `
	// }

	// if *newSystem.SystemCode != *oldSystem.SystemCode {
	// 	result.Parameters["systemCode"] = newSystem.SystemCode
	// 	result.Query += `WITH s SET s.systemCode=$systemCode `
	// }

	// if *newSystem.SystemAlias != *oldSystem.SystemAlias {
	// 	result.Parameters["systemAlias"] = newSystem.SystemAlias
	// 	result.Query += `WITH s SET s.systemAlias=$systemAlias `
	// }

	//to do later
	// if *newSystem.ParentUID != *oldSystem.ParentUID {
	// 	result.Query += `WITH s MATCH(parent:System{uid:$parentUID}) MERGE(parent)-[:HAS_SUBSYSTEM]->(s) `
	// 	result.Parameters["parentUID"] = newSystem.ParentUID
	// }

	// if newSystem.ZoneUID != nil && *newSystem.ZoneUID != "" && oldSystem.ZoneUID == nil {
	// 	result.Query += `WITH s MATCH(z:Zone{uid:$zoneUID}) MERGE(s)-[:HAS_ZONE]->(z) `
	// 	result.Parameters["zoneUID"] = newSystem.ZoneUID
	// } else if newSystem.ZoneUID != nil && *newSystem.ZoneUID != "" && oldSystem.ZoneUID != nil && *newSystem.ZoneUID != *oldSystem.ZoneUID {
	// 	result.Query += `WITH s MATCH(s)-[rz:HAS_ZONE]->(z) delete rz
	// 					 WITH s MATCH(z:Zone{uid:$zoneUID}) MERGE(s)-[:HAS_ZONE]->(z) `
	// 	result.Parameters["zoneUID"] = newSystem.ZoneUID
	// } else if (newSystem.ZoneUID == nil || *newSystem.ZoneUID == "") && oldSystem.ZoneUID != nil {
	// 	result.Query += `WITH s MATCH(s)-[rz:HAS_ZONE]->(z) delete rz `
	// }

	// if newSystem.LocationUID != nil && *newSystem.LocationUID != "" && oldSystem.LocationUID == nil {
	// 	result.Query += `WITH s MATCH(l:Location{uid:$LocationUID}) MERGE(s)-[:HAS_LOCATION]->(l) `
	// 	result.Parameters["LocationUID"] = newSystem.LocationUID
	// } else if newSystem.LocationUID != nil && *newSystem.LocationUID != "" && oldSystem.LocationUID != nil && *newSystem.LocationUID != *oldSystem.LocationUID {
	// 	result.Query += `WITH s MATCH(s)-[rl:HAS_LOCATION]->(l) delete rl
	// 					 WITH s MATCH(l:Location{uid:$LocationUID}) MERGE(s)-[:HAS_LOCATION]->(l) `
	// 	result.Parameters["LocationUID"] = newSystem.LocationUID
	// } else if (newSystem.LocationUID == nil || *newSystem.ZoneUID == "") && oldSystem.LocationUID != nil {
	// 	result.Query += `WITH s MATCH(s)-[rl:HAS_LOCATION]->(l) delete rl `
	// }

	// if newSystem.SystemTypeUID != nil && oldSystem.SystemTypeUID == nil {
	// 	result.Query += `WITH s MATCH(st:SystemType{uid:$SystemTypeUID}) MERGE(s)-[:HAS_SYSTEM_TYPE]->(st) `
	// 	result.Parameters["SystemTypeUID"] = newSystem.SystemTypeUID
	// } else if newSystem.SystemTypeUID != nil && oldSystem.SystemTypeUID != nil && *newSystem.SystemTypeUID != *oldSystem.SystemTypeUID {
	// 	result.Query += `WITH s MATCH(s)-[rst:HAS_SYSTEM_TYPE]->(st) delete rst
	//                      WITH s MATCH(st:SystemType{uid:$SystemTypeUID}) MERGE(s)-[:HAS_SYSTEM_TYPE]->(st) `
	// 	result.Parameters["SystemTypeUID"] = newSystem.SystemTypeUID
	// } else if newSystem.SystemTypeUID == nil && oldSystem.SystemTypeUID != nil {
	// 	result.Query += `WITH s MATCH(s)-[rst:HAS_SYSTEM_TYPE]->(st) delete rst `
	// }

	// if newSystem.OwnerUID != nil && oldSystem.OwnerUID == nil {
	// 	result.Query += `WITH s MATCH(own:User{uid:$OwnerUID}) MERGE(s)-[:HAS_OWNER]->(own) `
	// 	result.Parameters["OwnerUID"] = newSystem.OwnerUID
	// } else if newSystem.OwnerUID != nil && oldSystem.OwnerUID != nil && *newSystem.OwnerUID != *oldSystem.OwnerUID {
	// 	result.Query += `WITH s MATCH(s)-[rown:HAS_OWNER]->(own) delete rown
	// 					 WITH s MATCH(own:Owner{uid:$OwnerUID}) MERGE(s)-[:HAS_OWNER]->(own) `
	// 	result.Parameters["OwnerUID"] = newSystem.OwnerUID
	// } else if newSystem.OwnerUID == nil && oldSystem.OwnerUID != nil {
	// 	result.Query += `WITH s MATCH(s)-[rown:HAS_OWNER]->(own) delete rown `
	// }

	// if newSystem.ImportanceUID != nil && oldSystem.ImportanceUID == nil {
	// 	result.Query += `WITH s MATCH(imp:SystemImportance{uid:$ImportanceUID}) MERGE(s)-[:HAS_IMPORTANCE]->(imp) `
	// 	result.Parameters["ImportanceUID"] = newSystem.ImportanceUID
	// } else if newSystem.ImportanceUID != nil && oldSystem.ImportanceUID != nil && *newSystem.ImportanceUID != *oldSystem.ImportanceUID {
	// 	result.Query += `WITH s MATCH(s)-[rimp:HAS_IMPORTANCE]->(imp) delete rimp
	// 					 WITH s MATCH(imp:SystemImportance{uid:$ImportanceUID}) MERGE(s)-[:HAS_IMPORTANCE]->(imp) `
	// 	result.Parameters["ImportanceUID"] = newSystem.ImportanceUID
	// } else if newSystem.ImportanceUID == nil && oldSystem.ImportanceUID != nil {
	// 	result.Query += `WITH s MATCH(s)-[rimp:HAS_IMPORTANCE]->(imp) delete rimp `
	// }

	if newSystem.Image != nil {
		if *newSystem.Image == "deleted" {
			result.Query += `WITH s SET s.image = null `
			result.Parameters["image"] = newSystem.Image
		} else {
			result.Query += `WITH s SET s.image = $image `
			result.Parameters["image"] = newSystem.Image
		}
	}

	result.Query += `RETURN s.uid as result`

	result.ReturnAlias = "result"

	return result
}
