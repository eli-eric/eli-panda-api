package systemsService

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"panda/apigateway/helpers"
)

// relationPatch replaces all outgoing relType edges of the patch target with edges to uids
// (empty uids = clear). Label/relType come from the fixed specs below, never from input.
type relationPatch struct {
	relType string
	label   string
	uids    []string
}

// patchPlan is the diff-derived write set for one node (System or its Item).
type patchPlan struct {
	scalars        map[string]interface{}
	relations      []relationPatch
	recalcCoverage bool
	changes        []helpers.ChangeEntry
}

func newPatchPlan() *patchPlan {
	return &patchPlan{scalars: map[string]interface{}{}}
}

func (p *patchPlan) isEmpty() bool {
	return len(p.changes) == 0
}

// duplicateSystemCodeMarker is the apoc.util.validate message raised inside the write tx.
const duplicateSystemCodeMarker = "DUPLICATE_SYSTEM_CODE"

// duplicateCodeGuardCypher re-checks code uniqueness inside the write transaction while holding
// the facility lock, so two concurrent writes of the same code cannot both pass the pre-check.
// Empty when the plan does not set a code.
func duplicateCodeGuardCypher(params map[string]interface{}, plan *patchPlan, carry string) string {
	if plan == nil {
		return ""
	}
	code, ok := plan.scalars["systemCode"].(string)
	if !ok {
		return ""
	}
	params["guardSystemCode"] = code
	return fmt.Sprintf(`WITH %[1]s
		SET f._systemCodeLock = true
		REMOVE f._systemCodeLock
		WITH %[1]s
		CALL apoc.util.validate(
			size([(o:System{systemCode: $guardSystemCode, deleted: false})-[:BELONGS_TO_FACILITY]->(f) WHERE o.uid <> s.uid | o]) > 0,
			'%[2]s', [])
		`, carry, duplicateSystemCodeMarker)
}

// tokenCondition renders the optimistic-concurrency WHERE clause for alias s against the
// token read before the write. A legacy system without lastUpdateTime matches IS NULL.
func tokenCondition(params map[string]interface{}, token *time.Time) string {
	if token == nil {
		return `s.lastUpdateTime IS NULL`
	}
	params["tokenEpochSeconds"] = token.Unix()
	params["tokenNanosecond"] = token.Nanosecond()
	return `s.lastUpdateTime.epochSeconds = $tokenEpochSeconds AND s.lastUpdateTime.nanosecond = $tokenNanosecond`
}

// auditCypher writes the WAS_UPDATED_BY edge for target alias (requires u in scope).
func auditCypher(target, action string) string {
	return fmt.Sprintf(`CREATE (%s)-[:WAS_UPDATED_BY{at: datetime(), action: "%s", changes: $changes,
		outsideResponsibility: CASE WHEN $outsideResponsibility THEN true ELSE null END}]->(u) `, target, action)
}

// applyPlanCypher renders SET + relation replacement for the plan on alias target, keeping
// carry in scope. Parameters are namespaced by prefix so two plans can share one query.
func applyPlanCypher(params map[string]interface{}, target, prefix string, plan *patchPlan, carry string) string {
	var b strings.Builder

	props := make([]string, 0, len(plan.scalars))
	for prop := range plan.scalars {
		props = append(props, prop)
	}
	sort.Strings(props)
	if len(props) > 0 {
		sets := make([]string, 0, len(props))
		for _, prop := range props {
			key := prefix + "_" + prop
			params[key] = plan.scalars[prop]
			sets = append(sets, fmt.Sprintf("%s.%s = $%s", target, prop, key))
		}
		b.WriteString("SET " + strings.Join(sets, ", ") + " ")
	}

	for i, rel := range plan.relations {
		key := fmt.Sprintf("%s_rel%d", prefix, i)
		params[key] = rel.uids
		fmt.Fprintf(&b, `WITH %[1]s
		CALL { WITH %[2]s OPTIONAL MATCH (%[2]s)-[old:%[3]s]->() DELETE old }
		CALL { WITH %[2]s MATCH (n:%[4]s) WHERE n.uid IN $%[5]s MERGE (%[2]s)-[:%[3]s]->(n) }
		`, carry, target, rel.relType, rel.label, key)
	}

	if plan.recalcCoverage {
		fmt.Fprintf(&b, `WITH %[1]s
		CALL {
			WITH %[2]s
			OPTIONAL MATCH (:System)-[sp:IS_SPARE_FOR]->(%[2]s)
			WITH %[2]s, coalesce(sum(sp.coverage), 0) AS coverage
			SET %[2]s.sparePartsCoverageSum = CASE WHEN coalesce(%[2]s.minimalSpareParstCount, 0) > 0 THEN coverage ELSE %[2]s.sparePartsCoverageSum END,
				%[2]s.sp_coverage = CASE WHEN coalesce(%[2]s.minimalSpareParstCount, 0) > 0
					THEN toFloat(coverage) / %[2]s.minimalSpareParstCount ELSE null END
		}
		`, carry, target)
	}

	return b.String()
}

// PatchSystemQuery applies a system plan (and optionally an item plan) in one transaction,
// guarded by the concurrency token, bumps lastUpdateTime and writes one audit edge on the System.
// Returns no rows when the token no longer matches (concurrent write).
func PatchSystemQuery(uid, facilityCode, userUID string, token *time.Time, systemPlan, itemPlan *patchPlan, changesJSON string, outsideResponsibility bool) helpers.DatabaseQuery {
	params := map[string]interface{}{
		"uid":                   uid,
		"facilityCode":          facilityCode,
		"userUID":               userUID,
		"changes":               changesJSON,
		"outsideResponsibility": outsideResponsibility,
	}

	// The write lock is taken BEFORE the token is compared: Neo4j reads without locks, so two
	// writers could both pass an unlocked check and the second would silently overwrite. The
	// lock is a write (SET/REMOVE) rather than apoc.lock.nodes because the planner may push a
	// filter in front of a procedure call, but never across an updating clause.
	query := `MATCH (u:User{uid: $userUID})
	MATCH (s:System{uid: $uid, deleted: false})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode})
	SET s._patchLock = true
	REMOVE s._patchLock
	WITH s, u, f
	WHERE ` + tokenCondition(params, token) + `
	` + duplicateCodeGuardCypher(params, systemPlan, "s, u, f") + `WITH s, u
	`

	if itemPlan != nil {
		query += `MATCH (s)-[:CONTAINS_ITEM]->(itm:Item) `
		query += applyPlanCypher(params, "itm", "item", itemPlan, "s, u, itm")
		query += `WITH s, u, itm SET itm.lastUpdateTime = datetime(), itm.lastUpdatedBy = u.lastName + " " + u.firstName `
		query += `WITH s, u `
	}
	if systemPlan != nil {
		query += applyPlanCypher(params, "s", "sys", systemPlan, "s, u")
	}

	query += `WITH s, u
	SET s.lastUpdateTime = datetime(), s.lastUpdatedBy = u.lastName + " " + u.firstName
	` + auditCypher("s", "UPDATE") + `
	RETURN s.uid AS uid`

	return helpers.DatabaseQuery{Query: query, ReturnAlias: "uid", Parameters: params}
}

// ResolveRefsQuery looks up non-deleted nodes of a spec's fixed label by uid (for validation and
// audit snapshots); facility-scoped labels must belong to the caller's facility.
func ResolveRefsQuery(spec codebookRelSpec, uids []string, facilityCode string) helpers.DatabaseQuery {
	name := "n.name"
	if spec.employee {
		name = employeeNameCypher("n")
	}
	where := `n.uid IN $uids AND coalesce(n.deleted, false) = false`
	if spec.facilityScoped {
		where += ` AND (n)-[:BELONGS_TO_FACILITY]->(:Facility{code: $facilityCode})`
	}
	return helpers.DatabaseQuery{
		Query: fmt.Sprintf(`MATCH (n:%s) WHERE %s
			RETURN {uid: n.uid, name: coalesce(%s, ""), code: n.code} AS ref`, spec.label, where, name),
		ReturnAlias: "ref",
		Parameters:  map[string]interface{}{"uids": uids, "facilityCode": facilityCode},
	}
}

// SystemCodeExistsQuery reports whether another non-deleted system in the facility uses code.
func SystemCodeExistsQuery(code, excludeUID, facilityCode string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `OPTIONAL MATCH (s:System{systemCode: $code, deleted: false})-[:BELONGS_TO_FACILITY]->(:Facility{code: $facilityCode})
			WHERE s.uid <> $excludeUID
			RETURN count(s) > 0 AS result`,
		ReturnAlias: "result",
		Parameters:  map[string]interface{}{"code": code, "excludeUID": excludeUID, "facilityCode": facilityCode},
	}
}

// CreateSystemQuery creates a system (optionally under parentUID) and applies the create plan
// in one transaction, writing the INSERT audit edge on the new system.
func CreateSystemQuery(newUID string, parentUID *string, facilityCode, userUID string, plan *patchPlan, changesJSON string, outsideResponsibility bool) helpers.DatabaseQuery {
	params := map[string]interface{}{
		"newUID":                newUID,
		"facilityCode":          facilityCode,
		"userUID":               userUID,
		"changes":               changesJSON,
		"outsideResponsibility": outsideResponsibility,
	}
	query := `MATCH (u:User{uid: $userUID})
		MATCH (f:Facility{code: $facilityCode})
		CREATE (s:System{uid: $newUID, deleted: false,
			lastUpdateTime: datetime(), lastUpdatedBy: u.lastName + " " + u.firstName})-[:BELONGS_TO_FACILITY]->(f)
		` + duplicateCodeGuardCypher(params, plan, "s, u, f") + `WITH s, u, f
		`
	if parentUID != nil {
		params["parentUID"] = *parentUID
		query += `MATCH (parent:System{uid: $parentUID, deleted: false})-[:BELONGS_TO_FACILITY]->(f)
		CREATE (parent)-[:HAS_SUBSYSTEM]->(s)
		WITH s, u
		`
	}
	query += applyPlanCypher(params, "s", "sys", plan, "s, u") + `WITH s, u
		` + auditCypher("s", "INSERT") + `
		RETURN s.uid AS uid`
	return helpers.DatabaseQuery{Query: query, ReturnAlias: "uid", Parameters: params}
}

// relationshipSideAudit is one audit edge of a relationship add/remove.
type relationshipSideAudit struct {
	systemUID             string
	changesJSON           string
	outsideResponsibility bool
}

// relationshipAuditCypher writes WAS_UPDATED_BY edges on both systems (aliases a, b; u in scope).
func relationshipAuditCypher(params map[string]interface{}, from, to relationshipSideAudit) string {
	params["fromChanges"] = from.changesJSON
	params["toChanges"] = to.changesJSON
	params["fromOutside"] = from.outsideResponsibility
	params["toOutside"] = to.outsideResponsibility
	return `CREATE (a)-[:WAS_UPDATED_BY{at: datetime(), action: "UPDATE", changes: $fromChanges,
			outsideResponsibility: CASE WHEN $fromOutside THEN true ELSE null END}]->(u)
		CREATE (b)-[:WAS_UPDATED_BY{at: datetime(), action: "UPDATE", changes: $toChanges,
			outsideResponsibility: CASE WHEN $toOutside THEN true ELSE null END}]->(u) `
}

// sparesRecalcCypher recalculates coverage of IS_SPARE_FOR edges of the spare (a) and the
// coverage sums of every system the spare still covers plus the formerly covered system (b).
const sparesRecalcCypher = `
		WITH a, b
		CALL {
			WITH a
			MATCH (a)-[spf:IS_SPARE_FOR]->()
			WITH collect(spf) AS edges
			FOREACH (e IN edges | SET e.coverage = 1.0 / size(edges))
		}
		WITH a, b
		CALL {
			WITH a, b
			OPTIONAL MATCH (a)-[:IS_SPARE_FOR]->(covered:System)
			WITH b, collect(covered) AS covered
			UNWIND covered + [b] AS target
			WITH DISTINCT target
			OPTIONAL MATCH (:System)-[sp:IS_SPARE_FOR]->(target)
			WITH target, coalesce(sum(sp.coverage), 0) AS coverage
			SET target.sparePartsCoverageSum = CASE WHEN coalesce(target.minimalSpareParstCount, 0) > 0 THEN coverage ELSE target.sparePartsCoverageSum END,
				target.sp_coverage = CASE WHEN coalesce(target.minimalSpareParstCount, 0) > 0
					THEN toFloat(coverage) / target.minimalSpareParstCount ELSE null END
		}
	`

// DeleteRelationshipQuery removes (a)-[relType]->(b), audits both sides and, for IS_SPARE_FOR,
// recalculates affected coverages in the same transaction. Returns no rows when the edge is absent.
func DeleteRelationshipQuery(fromUID, relType, toUID, facilityCode, userUID string, from, to relationshipSideAudit) helpers.DatabaseQuery {
	params := map[string]interface{}{
		"fromUID":      fromUID,
		"toUID":        toUID,
		"facilityCode": facilityCode,
		"userUID":      userUID,
	}
	query := fmt.Sprintf(`MATCH (u:User{uid: $userUID})
		MATCH (a:System{uid: $fromUID})-[:BELONGS_TO_FACILITY]->(:Facility{code: $facilityCode})
		MATCH (a)-[r:%s]->(b:System{uid: $toUID})
		// legacy data may hold duplicate edges: remove them all, audit once
		WITH a, b, u, collect(r) AS edges
		FOREACH (e IN edges | DELETE e)
		WITH a, b, u
		`, relType) + relationshipAuditCypher(params, from, to)
	if relType == "IS_SPARE_FOR" {
		query += `WITH a, b ` + sparesRecalcCypher
	}
	query += `RETURN true AS result`
	return helpers.DatabaseQuery{Query: query, ReturnAlias: "result", Parameters: params}
}

// RelationshipExistsQuery reports whether (from)-[relType]->(to) exists.
func RelationshipExistsQuery(fromUID, relType, toUID string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query:       fmt.Sprintf(`OPTIONAL MATCH (:System{uid: $fromUID})-[r:%s]->(:System{uid: $toUID}) RETURN count(r) > 0 AS result`, relType),
		ReturnAlias: "result",
		Parameters:  map[string]interface{}{"fromUID": fromUID, "toUID": toUID},
	}
}
