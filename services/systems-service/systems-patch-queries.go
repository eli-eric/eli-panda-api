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
			SET %[2]s.sparePartsCoverageSum = coverage,
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

	query := `MATCH (u:User{uid: $userUID})
	MATCH (s:System{uid: $uid, deleted: false})-[:BELONGS_TO_FACILITY]->(:Facility{code: $facilityCode})
	WHERE ` + tokenCondition(params, token) + ` `

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

// ResolveRefsQuery looks up nodes of a fixed label by uid for validation and audit snapshots.
func ResolveRefsQuery(label string, uids []string, employee bool) helpers.DatabaseQuery {
	name := "n.name"
	if employee {
		name = employeeNameCypher("n")
	}
	return helpers.DatabaseQuery{
		Query: fmt.Sprintf(`MATCH (n:%s) WHERE n.uid IN $uids
			RETURN {uid: n.uid, name: coalesce(%s, ""), code: n.code} AS ref`, label, name),
		ReturnAlias: "ref",
		Parameters:  map[string]interface{}{"uids": uids},
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

// CreateSubsystemQuery creates a child of parentUID inheriting responsible, location and zone,
// and writes the INSERT audit edge on the new system.
func CreateSubsystemQuery(parentUID, newUID, name, systemLevel, facilityCode, userUID, changesJSON string, outsideResponsibility bool) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query: `MATCH (u:User{uid: $userUID})
		MATCH (parent:System{uid: $parentUID, deleted: false})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode})
		CREATE (s:System{uid: $newUID, name: $name, systemLevel: $systemLevel, deleted: false,
			lastUpdateTime: datetime(), lastUpdatedBy: u.lastName + " " + u.firstName})-[:BELONGS_TO_FACILITY]->(f)
		CREATE (parent)-[:HAS_SUBSYSTEM]->(s)
		WITH s, u, parent
		CALL { WITH s, parent MATCH (parent)-[:HAS_RESPONSIBLE]->(x) MERGE (s)-[:HAS_RESPONSIBLE]->(x) }
		CALL { WITH s, parent MATCH (parent)-[:HAS_LOCATION]->(x) MERGE (s)-[:HAS_LOCATION]->(x) }
		CALL { WITH s, parent MATCH (parent)-[:HAS_ZONE]->(x) MERGE (s)-[:HAS_ZONE]->(x) }
		` + auditCypher("s", "INSERT") + `
		RETURN s.uid AS uid`,
		ReturnAlias: "uid",
		Parameters: map[string]interface{}{
			"parentUID":             parentUID,
			"newUID":                newUID,
			"name":                  name,
			"systemLevel":           systemLevel,
			"facilityCode":          facilityCode,
			"userUID":               userUID,
			"changes":               changesJSON,
			"outsideResponsibility": outsideResponsibility,
		},
	}
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
		DELETE r
		WITH a, b, u
		`, relType) + relationshipAuditCypher(params, from, to)
	if relType == "IS_SPARE_FOR" {
		query += `WITH a, b ` + sparesRecalcCypher
	}
	query += `RETURN true AS result`
	return helpers.DatabaseQuery{Query: query, ReturnAlias: "result", Parameters: params}
}

// RelationshipAuditQuery writes the add-relationship audit on both systems (used after create).
func RelationshipAuditQuery(fromUID, toUID, userUID string, from, to relationshipSideAudit) helpers.DatabaseQuery {
	params := map[string]interface{}{"fromUID": fromUID, "toUID": toUID, "userUID": userUID}
	query := `MATCH (u:User{uid: $userUID})
		MATCH (a:System{uid: $fromUID})
		MATCH (b:System{uid: $toUID})
		` + relationshipAuditCypher(params, from, to) + `RETURN true AS result`
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
