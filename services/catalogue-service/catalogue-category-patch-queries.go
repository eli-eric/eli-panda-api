package catalogueService

import (
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

