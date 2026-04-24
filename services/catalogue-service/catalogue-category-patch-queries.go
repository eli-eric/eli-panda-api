package catalogueService

import (
	"fmt"
	"strings"

	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
)

// =====  Shared helpers for category granular PATCH queries  =====

// categoryAuditSuffix is the last two Cypher phases common to every granular category
// mutation: the WAS_UPDATED_BY audit edge carrying the JSON-stringified changes and the
// RETURN that produces the category UID the service uses to detect zero-row failures.
const categoryAuditSuffix = `
	CREATE(category)-[:WAS_UPDATED_BY{ at: datetime(), action: $action, changes: $changes }]->(u)
	RETURN category.uid as uid
	`

// initCategoryPatchQuery seeds the common parameters (uid, userUID, action) and returns
// the phase-1 skeleton that every category mutation shares: MATCH(user) → MATCH(category).
// Additional phase-1 MATCHes (supplier/category/property-specific) are appended by each
// builder. Order of operations: all MATCHes first, then all writes, then the audit edge.
func initCategoryPatchQuery(categoryUID, userUID, action string) (params map[string]interface{}, skeleton string) {
	params = make(map[string]interface{})
	params["uid"] = categoryUID
	params["userUID"] = userUID
	params["action"] = action

	skeleton = `
	MATCH(u:User{uid: $userUID})
	WITH u
	MATCH(category:CatalogueCategory{uid: $uid})
	`
	return
}

// =====  PATCH /v1/catalogue/category/:uid  =====

// PatchCatalogueCategoryQuery builds the Cypher for a scalar-field PATCH on a category.
// Only name, code, and systemType are in scope — image is handled out-of-band via Minio.
//
// Structure: Phase 1 matches user, category, and (when set) new systemType. Phase 2 applies
// SETs and the optional HAS_SYSTEM_TYPE relationship swap. Phase 3 writes the audit edge.
func PatchCatalogueCategoryQuery(uid string, fields *models.PatchCatalogueCategoryFields, original *models.CatalogueCategory, userUID string) (result helpers.DatabaseQuery) {
	params, skeleton := initCategoryPatchQuery(uid, userUID, "UPDATE")
	result.Parameters = params
	result.Query = skeleton

	var changes []helpers.ChangeEntry
	carry := []string{"u", "category"}

	// -- PHASE 1 — additional MATCHes --

	settingSystemType := fields.SystemType != nil && fields.SystemType.Value != nil && fields.SystemType.Value.UID != ""
	if settingSystemType {
		result.Parameters["systemTypeUid"] = fields.SystemType.Value.UID
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` MATCH(newSystemType:SystemType{uid: $systemTypeUid}) `
		carry = append(carry, "newSystemType")
	}

	// Consolidating WITH so phase-2 sees the full variable set.
	result.Query += ` WITH ` + strings.Join(carry, ", ") + ` `

	// -- PHASE 2 — writes --

	if fields.Name != nil {
		result.Parameters["name"] = *fields.Name
		result.Query += ` SET category.name = $name `
		changes = helpers.AppendIfChanged(changes, "name", helpers.ChangeTypeString, original.Name, *fields.Name)
	}

	if fields.Code != nil {
		result.Parameters["code"] = *fields.Code
		result.Query += ` SET category.code = $code `
		changes = helpers.AppendIfChanged(changes, "code", helpers.ChangeTypeString, original.Code, *fields.Code)
	}

	if fields.SystemType != nil {
		// Always drop any existing HAS_SYSTEM_TYPE; MERGE the new one only when setting.
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` OPTIONAL MATCH (category)-[oldST:HAS_SYSTEM_TYPE]->() DELETE oldST `
		if settingSystemType {
			result.Query += ` WITH ` + strings.Join(carry, ", ") + ` MERGE(category)-[:HAS_SYSTEM_TYPE]->(newSystemType) `
		}
		var oldST, newST interface{}
		if original.SystemType != nil {
			oldST = original.SystemType
		}
		if fields.SystemType.Value != nil {
			newST = fields.SystemType.Value
		}
		changes = helpers.AppendIfChanged(changes, "systemType", helpers.ChangeTypeCodebook, oldST, newST)
	}

	// -- PHASE 3 — audit + return --

	result.Parameters["changes"] = helpers.MarshalChanges(changes)
	result.Query += categoryAuditSuffix
	result.ReturnAlias = "uid"
	return result
}

// =====  POST /v1/catalogue/category/:uid/group  =====

// CreateCatalogueCategoryGroupQuery creates a new group under a category. The caller
// (service layer) is responsible for computing newGroupUID and resolvedOrder (payload
// order, or max(siblings)+10 when absent). Single Cypher transaction: MATCH category,
// CREATE new group + HAS_GROUP edge, emit audit.
func CreateCatalogueCategoryGroupQuery(categoryUID, newGroupUID, userUID, name string, resolvedOrder int) (result helpers.DatabaseQuery) {
	params, skeleton := initCategoryPatchQuery(categoryUID, userUID, "INSERT")
	result.Parameters = params
	result.Parameters["newGroupUid"] = newGroupUID
	result.Parameters["name"] = name
	result.Parameters["order"] = resolvedOrder
	result.Query = skeleton

	result.Query += `
	WITH u, category
	CREATE(newGroup:CatalogueCategoryPropertyGroup{uid: $newGroupUid, name: $name, order: $order})
	CREATE(category)-[:HAS_GROUP]->(newGroup)
	WITH u, category
	`

	changes := []helpers.ChangeEntry{
		{Field: "group", Type: string(helpers.ChangeTypeString), OldValue: nil, NewValue: name},
	}
	result.Parameters["changes"] = helpers.MarshalChanges(changes)
	result.Query += categoryAuditSuffix
	result.ReturnAlias = "uid"
	return result
}

// =====  PATCH /v1/catalogue/category/:uid/group/:gid  =====

// PatchCatalogueCategoryGroupQuery updates an existing group's name and/or order.
// Phase-1 MATCH enforces category↔group consistency — category UID mismatch yields 0 rows
// and the service maps to 404.
func PatchCatalogueCategoryGroupQuery(categoryUID, groupUID string, fields *models.PatchCatalogueCategoryGroupFields, original *models.CatalogueCategoryPropertyGroup, userUID string) (result helpers.DatabaseQuery) {
	params, skeleton := initCategoryPatchQuery(categoryUID, userUID, "UPDATE")
	result.Parameters = params
	result.Parameters["groupUid"] = groupUID
	result.Query = skeleton
	result.Query += ` WITH u, category MATCH(category)-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup{uid: $groupUid}) WITH u, category, g `

	var changes []helpers.ChangeEntry

	if fields.Name != nil {
		result.Parameters["name"] = *fields.Name
		result.Query += ` SET g.name = $name `
		changes = helpers.AppendIfChanged(changes, fmt.Sprintf("group.%s.name", groupUID), helpers.ChangeTypeString, original.Name, *fields.Name)
	}

	if fields.Order != nil {
		result.Parameters["order"] = *fields.Order
		result.Query += ` SET g.order = $order `
		var oldOrder interface{}
		if original.Order != nil {
			oldOrder = *original.Order
		}
		changes = helpers.AppendIfChanged(changes, fmt.Sprintf("group.%s.order", groupUID), helpers.ChangeTypeNumber, oldOrder, *fields.Order)
	}

	result.Parameters["changes"] = helpers.MarshalChanges(changes)
	result.Query += categoryAuditSuffix
	result.ReturnAlias = "uid"
	return result
}

// =====  DELETE /v1/catalogue/category/:uid/group/:gid  =====

// DeleteCatalogueCategoryGroupQuery removes a group and all its CatalogueCategoryProperty
// nodes IF no CatalogueItem references any of those properties. The phase-1 count acts as
// the gate: refs>0 makes the WHERE filter out the only matching row and the query returns
// zero rows → service maps to ERR_DELETE_RELATED_ITEMS → 409.
func DeleteCatalogueCategoryGroupQuery(categoryUID, groupUID, userUID, originalName string) (result helpers.DatabaseQuery) {
	params, skeleton := initCategoryPatchQuery(categoryUID, userUID, "DELETE")
	result.Parameters = params
	result.Parameters["groupUid"] = groupUID
	result.Query = skeleton
	result.Query += `
	WITH u, category
	MATCH(category)-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup{uid: $groupUid})
	OPTIONAL MATCH(g)-[:CONTAINS_PROPERTY]->(p:CatalogueCategoryProperty)
	OPTIONAL MATCH(:CatalogueItem)-[itemRef:HAS_CATALOGUE_PROPERTY]->(p)
	WITH u, category, g, collect(DISTINCT p) as props, count(itemRef) as refs
	WHERE refs = 0
	FOREACH (prop IN props | DETACH DELETE prop)
	DETACH DELETE g
	WITH u, category
	`

	changes := []helpers.ChangeEntry{
		{Field: "group", Type: string(helpers.ChangeTypeString), OldValue: originalName, NewValue: nil},
	}
	result.Parameters["changes"] = helpers.MarshalChanges(changes)
	result.Query += categoryAuditSuffix
	result.ReturnAlias = "uid"
	return result
}

// =====  Lazy order seed helper  =====

// SeedCategoryGroupOrdersQuery renumbers all NULL-order groups under a category using
// their current id(n) ordering, with step=10. Invoked by the service layer the first
// time a group PATCH sets explicit order in a category with legacy (unseeded) groups.
// Safe to call repeatedly — already-ordered groups keep their value; seeds start
// at max(existing.order)+10 to avoid collisions with any partial seeding.
func SeedCategoryGroupOrdersQuery(categoryUID string) (result helpers.DatabaseQuery) {
	result.Parameters = map[string]interface{}{"uid": categoryUID}
	result.Query = `
	MATCH(c:CatalogueCategory{uid: $uid})
	OPTIONAL MATCH(c)-[:HAS_GROUP]->(ordered:CatalogueCategoryPropertyGroup) WHERE ordered.order IS NOT NULL
	WITH c, coalesce(max(ordered.order), 0) as maxExisting
	MATCH(c)-[:HAS_GROUP]->(unseeded:CatalogueCategoryPropertyGroup)
	WHERE unseeded.order IS NULL
	WITH maxExisting, unseeded ORDER BY id(unseeded)
	WITH maxExisting, collect(unseeded) as gs
	UNWIND range(0, size(gs)-1) as i
	WITH gs[i] as gn, maxExisting, i
	SET gn.order = maxExisting + (i + 1) * 10
	RETURN count(gn) as seeded
	`
	result.ReturnAlias = "seeded"
	return result
}

// NextGroupOrderQuery returns max(order of existing groups under this category) + 10,
// or 10 if no sibling has an order yet. Used by the create-group service to auto-assign
// a sensible default when the POST payload omits order.
func NextGroupOrderQuery(categoryUID string) (result helpers.DatabaseQuery) {
	result.Parameters = map[string]interface{}{"uid": categoryUID}
	result.Query = `
	MATCH(c:CatalogueCategory{uid: $uid})
	OPTIONAL MATCH(c)-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup)
	RETURN coalesce(max(g.order), 0) + 10 as nextOrder
	`
	result.ReturnAlias = "nextOrder"
	return result
}

// GetCatalogueCategoryGroupByUidsQuery fetches a single group under a category. The MATCH
// enforces category↔group consistency — if the group doesn't belong to this category (or
// either UID is unknown) the query returns zero rows and the caller gets ERR_NO_ROWS.
// Unlike the legacy CatalogueCategoryWithDetailsQuery, this does not require the group
// to contain properties — a freshly-created empty group still reads back.
func GetCatalogueCategoryGroupByUidsQuery(categoryUID, groupUID string) (result helpers.DatabaseQuery) {
	result.Parameters = map[string]interface{}{"uid": categoryUID, "groupUid": groupUID}
	result.Query = `
	MATCH(c:CatalogueCategory{uid: $uid})-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup{uid: $groupUid})
	RETURN { uid: g.uid, name: g.name, order: g.order } as group
	`
	result.ReturnAlias = "group"
	return result
}

// =====  Property CRUD — POST/PATCH/DELETE /v1/catalogue/category/:uid/property[/...]  =====

// GetCatalogueCategoryPropertyByUidsQuery fetches a single property scoped to a category
// via HAS_GROUP → CONTAINS_PROPERTY. Empty result = property doesn't exist OR doesn't
// belong to this category (caller maps to 404). Includes the parent groupUid so the
// service can detect move operations. Unit and type are returned as inline codebook-shape
// maps so the Go struct unmarshals into CatalogueCategoryProperty unchanged.
func GetCatalogueCategoryPropertyByUidsQuery(categoryUID, propertyUID string) (result helpers.DatabaseQuery) {
	result.Parameters = map[string]interface{}{"uid": categoryUID, "propertyUid": propertyUID}
	result.Query = `
	MATCH(c:CatalogueCategory{uid: $uid})-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup)-[:CONTAINS_PROPERTY]->(p:CatalogueCategoryProperty{uid: $propertyUid})
	OPTIONAL MATCH(p)-[:IS_PROPERTY_TYPE]->(t:CatalogueCategoryPropertyType)
	OPTIONAL MATCH(p)-[:HAS_UNIT]->(u:Unit)
	RETURN {
		uid: p.uid,
		name: p.name,
		defaultValue: p.defaultValue,
		order: p.order,
		listOfValues: case when p.listOfValues is not null and p.listOfValues <> '' then apoc.text.split(p.listOfValues, ';') else null end,
		type: case when t is not null then {uid: t.uid, name: t.name, code: t.code} else null end,
		unit: case when u is not null then {uid: u.uid, name: u.name} else null end,
		groupUid: g.uid
	} as property
	`
	result.ReturnAlias = "property"
	return result
}

// CategoryPropertyWithGroup is an internal shape returned by the query above — carries
// parent group UID alongside the property so the service can detect move operations.
type CategoryPropertyWithGroup struct {
	UID          string                                 `json:"uid"`
	Name         string                                 `json:"name"`
	DefaultValue string                                 `json:"defaultValue"`
	Order        *int                                   `json:"order"`
	ListOfValues []string                               `json:"listOfValues"`
	Type         *models.CatalogueCategoryPropertyType  `json:"type"`
	Unit         *codebookRef                           `json:"unit"`
	GroupUID     string                                 `json:"groupUid"`
}

// codebookRef mirrors models.Codebook for unmarshaling; keeping it local avoids pulling
// the codebook package path into every test that reads a property.
type codebookRef struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

// ListCatalogueCategoryPropertiesInGroupQuery returns every property under a specific
// group, used for NextPropertyOrderQuery's sibling lookup and for lazy-seed triggers.
func ListCatalogueCategoryPropertiesInGroupQuery(categoryUID, groupUID string) (result helpers.DatabaseQuery) {
	result.Parameters = map[string]interface{}{"uid": categoryUID, "groupUid": groupUID}
	result.Query = `
	MATCH(c:CatalogueCategory{uid: $uid})-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup{uid: $groupUid})-[:CONTAINS_PROPERTY]->(p:CatalogueCategoryProperty)
	WITH p ORDER BY coalesce(p.order, 2147483647), id(p)
	RETURN { uid: p.uid, name: p.name, order: p.order } as property
	`
	result.ReturnAlias = "property"
	return result
}

// NextPropertyOrderQuery returns max(order of existing properties in this group) + 10,
// or 10 when the group is empty. Used by create-property service when payload omits order.
func NextPropertyOrderQuery(categoryUID, groupUID string) (result helpers.DatabaseQuery) {
	result.Parameters = map[string]interface{}{"uid": categoryUID, "groupUid": groupUID}
	result.Query = `
	MATCH(c:CatalogueCategory{uid: $uid})-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup{uid: $groupUid})
	OPTIONAL MATCH(g)-[:CONTAINS_PROPERTY]->(p:CatalogueCategoryProperty)
	RETURN coalesce(max(p.order), 0) + 10 as nextOrder
	`
	result.ReturnAlias = "nextOrder"
	return result
}

// SeedCategoryPropertyOrdersQuery — same lazy-seed pattern as SeedCategoryGroupOrdersQuery
// but scoped to a single group's properties. Renumbers NULL-order properties starting at
// max(existing.order)+10 in id(n) order, step=10.
func SeedCategoryPropertyOrdersQuery(categoryUID, groupUID string) (result helpers.DatabaseQuery) {
	result.Parameters = map[string]interface{}{"uid": categoryUID, "groupUid": groupUID}
	result.Query = `
	MATCH(c:CatalogueCategory{uid: $uid})-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup{uid: $groupUid})
	OPTIONAL MATCH(g)-[:CONTAINS_PROPERTY]->(ordered:CatalogueCategoryProperty) WHERE ordered.order IS NOT NULL
	WITH g, coalesce(max(ordered.order), 0) as maxExisting
	MATCH(g)-[:CONTAINS_PROPERTY]->(unseeded:CatalogueCategoryProperty)
	WHERE unseeded.order IS NULL
	WITH maxExisting, unseeded ORDER BY id(unseeded)
	WITH maxExisting, collect(unseeded) as ps
	UNWIND range(0, size(ps)-1) as i
	WITH ps[i] as pn, maxExisting, i
	SET pn.order = maxExisting + (i + 1) * 10
	RETURN count(pn) as seeded
	`
	result.ReturnAlias = "seeded"
	return result
}

// CreateCatalogueCategoryPropertyQuery creates a new property under a group with its type
// and optional unit. Phase-1 MATCHes enforce category↔group consistency and validate
// type/unit UIDs. The service pre-computes resolvedOrder (payload value or max+10).
func CreateCatalogueCategoryPropertyQuery(categoryUID, groupUID, newPropertyUID, userUID string, fields *models.CreateCatalogueCategoryPropertyFields, resolvedOrder int) (result helpers.DatabaseQuery) {
	params, skeleton := initCategoryPatchQuery(categoryUID, userUID, "INSERT")
	result.Parameters = params
	result.Parameters["groupUid"] = groupUID
	result.Parameters["newPropUid"] = newPropertyUID
	result.Parameters["name"] = fields.Name
	result.Parameters["order"] = resolvedOrder
	result.Parameters["typeUid"] = fields.Type.UID

	defaultValue := ""
	if fields.DefaultValue != nil {
		defaultValue = *fields.DefaultValue
	}
	result.Parameters["defaultValue"] = defaultValue
	result.Parameters["listOfValues"] = strings.Join(fields.ListOfValues, ";")

	settingUnit := fields.Unit != nil && fields.Unit.UID != ""
	if settingUnit {
		result.Parameters["unitUid"] = fields.Unit.UID
	}

	result.Query = skeleton
	result.Query += `
	WITH u, category
	MATCH(category)-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup{uid: $groupUid})
	WITH u, category, g
	MATCH(propType:CatalogueCategoryPropertyType{uid: $typeUid})
	WITH u, category, g, propType
	`
	if settingUnit {
		result.Query += `
		MATCH(unit:Unit{uid: $unitUid})
		WITH u, category, g, propType, unit
		`
	}
	result.Query += `
	CREATE(newProp:CatalogueCategoryProperty{uid: $newPropUid, name: $name, defaultValue: $defaultValue, listOfValues: $listOfValues, order: $order})
	CREATE(g)-[:CONTAINS_PROPERTY]->(newProp)
	CREATE(newProp)-[:IS_PROPERTY_TYPE]->(propType)
	`
	if settingUnit {
		result.Query += `
		CREATE(newProp)-[:HAS_UNIT]->(unit)
		`
	}
	result.Query += ` WITH u, category `

	changes := []helpers.ChangeEntry{
		{Field: "property", Type: string(helpers.ChangeTypeString), OldValue: nil, NewValue: fields.Name},
	}
	result.Parameters["changes"] = helpers.MarshalChanges(changes)
	result.Query += categoryAuditSuffix
	result.ReturnAlias = "uid"
	return result
}

// =====  PATCH /v1/catalogue/category/:uid/property/:pid  =====

// PatchCatalogueCategoryPropertyQuery updates a property's fields and, optionally, moves
// it to a different group. Supports partial updates via the PatchCatalogueCategoryProperty
// Fields struct — each non-nil pointer produces a SET, MERGE, or swap segment.
//
// Callers pass `original` (with current values + current GroupUID) so change tracking can
// produce accurate oldValue entries and the move-detection logic can decide whether to
// swap CONTAINS_PROPERTY.
func PatchCatalogueCategoryPropertyQuery(categoryUID, propertyUID string, fields *models.PatchCatalogueCategoryPropertyFields, original *CategoryPropertyWithGroup, userUID string) (result helpers.DatabaseQuery) {
	params, skeleton := initCategoryPatchQuery(categoryUID, userUID, "UPDATE")
	result.Parameters = params
	result.Parameters["propertyUid"] = propertyUID
	result.Query = skeleton

	carry := []string{"u", "category"}

	// Phase-1 property match (via HAS_GROUP → CONTAINS_PROPERTY so cross-category UIDs fail).
	result.Query += ` WITH ` + strings.Join(carry, ", ") + ` MATCH(category)-[:HAS_GROUP]->(oldGroup:CatalogueCategoryPropertyGroup)-[:CONTAINS_PROPERTY]->(p:CatalogueCategoryProperty{uid: $propertyUid}) `
	carry = append(carry, "oldGroup", "p")

	// Optional MATCHes for swap targets — fail phase 1 if any is missing.
	movingGroup := fields.GroupUID != nil && *fields.GroupUID != "" && *fields.GroupUID != original.GroupUID
	if movingGroup {
		result.Parameters["newGroupUid"] = *fields.GroupUID
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` MATCH(category)-[:HAS_GROUP]->(newGroup:CatalogueCategoryPropertyGroup{uid: $newGroupUid}) `
		carry = append(carry, "newGroup")
	}

	settingType := fields.Type != nil && fields.Type.UID != ""
	if settingType {
		result.Parameters["typeUid"] = fields.Type.UID
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` MATCH(newPropType:CatalogueCategoryPropertyType{uid: $typeUid}) `
		carry = append(carry, "newPropType")
	}

	settingUnit := fields.Unit != nil && fields.Unit.Value != nil && fields.Unit.Value.UID != ""
	if settingUnit {
		result.Parameters["unitUid"] = fields.Unit.Value.UID
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` MATCH(newUnit:Unit{uid: $unitUid}) `
		carry = append(carry, "newUnit")
	}

	result.Query += ` WITH ` + strings.Join(carry, ", ") + ` `

	// -- PHASE 2 — writes --
	var changes []helpers.ChangeEntry

	if fields.Name != nil {
		result.Parameters["name"] = *fields.Name
		result.Query += ` SET p.name = $name `
		changes = helpers.AppendIfChanged(changes, propertyField(propertyUID, "name"), helpers.ChangeTypeString, original.Name, *fields.Name)
	}

	if fields.DefaultValue != nil {
		if fields.DefaultValue.Value != nil {
			result.Parameters["defaultValue"] = *fields.DefaultValue.Value
		} else {
			result.Parameters["defaultValue"] = ""
		}
		result.Query += ` SET p.defaultValue = $defaultValue `
		newVal := ""
		if fields.DefaultValue.Value != nil {
			newVal = *fields.DefaultValue.Value
		}
		changes = helpers.AppendIfChanged(changes, propertyField(propertyUID, "defaultValue"), helpers.ChangeTypeString, original.DefaultValue, newVal)
	}

	if fields.ListOfValues != nil {
		joined := strings.Join(*fields.ListOfValues, ";")
		result.Parameters["listOfValues"] = joined
		result.Query += ` SET p.listOfValues = $listOfValues `
		oldJoined := strings.Join(original.ListOfValues, ";")
		changes = helpers.AppendIfChanged(changes, propertyField(propertyUID, "listOfValues"), helpers.ChangeTypeString, oldJoined, joined)
	}

	if fields.Order != nil {
		result.Parameters["order"] = *fields.Order
		result.Query += ` SET p.order = $order `
		var oldOrder interface{}
		if original.Order != nil {
			oldOrder = *original.Order
		}
		changes = helpers.AppendIfChanged(changes, propertyField(propertyUID, "order"), helpers.ChangeTypeNumber, oldOrder, *fields.Order)
	}

	if settingType {
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` OPTIONAL MATCH (p)-[oldTypeRel:IS_PROPERTY_TYPE]->() DELETE oldTypeRel `
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` MERGE(p)-[:IS_PROPERTY_TYPE]->(newPropType) `
		var oldTy, newTy interface{}
		if original.Type != nil {
			oldTy = original.Type
		}
		newTy = fields.Type
		changes = helpers.AppendIfChanged(changes, propertyField(propertyUID, "type"), helpers.ChangeTypeCodebook, oldTy, newTy)
	}

	if fields.Unit != nil {
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` OPTIONAL MATCH (p)-[oldUnitRel:HAS_UNIT]->() DELETE oldUnitRel `
		if settingUnit {
			result.Query += ` WITH ` + strings.Join(carry, ", ") + ` MERGE(p)-[:HAS_UNIT]->(newUnit) `
		}
		var oldUnit, newUnit interface{}
		if original.Unit != nil {
			oldUnit = original.Unit
		}
		if fields.Unit.Value != nil {
			newUnit = fields.Unit.Value
		}
		changes = helpers.AppendIfChanged(changes, propertyField(propertyUID, "unit"), helpers.ChangeTypeCodebook, oldUnit, newUnit)
	}

	if movingGroup {
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` OPTIONAL MATCH (oldGroup)-[oldCP:CONTAINS_PROPERTY]->(p) DELETE oldCP `
		result.Query += ` WITH ` + strings.Join(carry, ", ") + ` MERGE(newGroup)-[:CONTAINS_PROPERTY]->(p) `
		changes = helpers.AppendIfChanged(changes, propertyField(propertyUID, "groupUid"), helpers.ChangeTypeString, original.GroupUID, *fields.GroupUID)
	}

	// -- PHASE 3 --
	result.Parameters["changes"] = helpers.MarshalChanges(changes)
	result.Query += categoryAuditSuffix
	result.ReturnAlias = "uid"
	return result
}

// =====  DELETE /v1/catalogue/category/:uid/property/:pid  =====

// DeleteCatalogueCategoryPropertyQuery removes a property IFF no CatalogueItem references
// it. The WHERE refs=0 gate propagates zero rows through the pipeline when refs>0, which
// the service maps to ERR_DELETE_RELATED_ITEMS → 409.
func DeleteCatalogueCategoryPropertyQuery(categoryUID, propertyUID, userUID, originalName string) (result helpers.DatabaseQuery) {
	params, skeleton := initCategoryPatchQuery(categoryUID, userUID, "DELETE")
	result.Parameters = params
	result.Parameters["propertyUid"] = propertyUID
	result.Query = skeleton
	result.Query += `
	WITH u, category
	MATCH(category)-[:HAS_GROUP]->(:CatalogueCategoryPropertyGroup)-[:CONTAINS_PROPERTY]->(p:CatalogueCategoryProperty{uid: $propertyUid})
	OPTIONAL MATCH(:CatalogueItem)-[itemRef:HAS_CATALOGUE_PROPERTY]->(p)
	WITH u, category, p, count(itemRef) as refs
	WHERE refs = 0
	DETACH DELETE p
	WITH u, category
	`

	changes := []helpers.ChangeEntry{
		{Field: "property", Type: string(helpers.ChangeTypeString), OldValue: originalName, NewValue: nil},
	}
	result.Parameters["changes"] = helpers.MarshalChanges(changes)
	result.Query += categoryAuditSuffix
	result.ReturnAlias = "uid"
	return result
}

// propertyField builds the path-like field name "property.<uid>.<scalar>" used in the
// audit changes array to identify which property + attribute was touched.
func propertyField(propertyUID, scalar string) string {
	return fmt.Sprintf("property.%s.%s", propertyUID, scalar)
}

// ListCatalogueCategoryGroupsQuery fetches every group (with its order) under a category.
// Used by the PATCH group service to decide whether lazy-seed is required — it must see
// ALL groups, including empty ones. Returns an array ordered by current sort to keep the
// lazy-seed renumbering deterministic.
func ListCatalogueCategoryGroupsQuery(categoryUID string) (result helpers.DatabaseQuery) {
	result.Parameters = map[string]interface{}{"uid": categoryUID}
	result.Query = `
	MATCH(c:CatalogueCategory{uid: $uid})-[:HAS_GROUP]->(g:CatalogueCategoryPropertyGroup)
	WITH g ORDER BY coalesce(g.order, 2147483647), id(g)
	RETURN { uid: g.uid, name: g.name, order: g.order } as group
	`
	result.ReturnAlias = "group"
	return result
}

