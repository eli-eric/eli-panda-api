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

