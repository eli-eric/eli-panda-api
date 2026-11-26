//go:build neo4jintegration

// Integration tests for the system PATCH family against a disposable Neo4j 4.4 + APOC.
// They seed and wipe their own fixture, so NEVER point them at a shared database:
//
//	docker run -d --rm --name panda-verify-neo4j -p 7699:7687 -e NEO4J_AUTH=neo4j/verifypass123 \
//	  -e 'NEO4JLABS_PLUGINS=["apoc"]' -e NEO4J_dbms_security_procedures_unrestricted='apoc.*' neo4j:4.4.15
//	PANDA_INTEGRATION_NEO4J_URI=bolt://localhost:7699 PANDA_INTEGRATION_NEO4J_PASSWORD=verifypass123 \
//	  go test -tags neo4jintegration ./services/systems-service -run Integration -v
package systemsService

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const integrationFixture = `
CREATE (f:Facility{code: 'B', name: 'Beamlines'})
CREATE (editor:User{uid: 'u-editor', firstName: 'Eva', lastName: 'Editor', username: 'eva'})
CREATE (owner:User{uid: 'u-owner', firstName: 'Olga', lastName: 'Owner', username: 'olga'})
CREATE (editR:Role{code: 'systems-edit'})
CREATE (editor)-[:HAS_ROLE]->(editR)
CREATE (owner)-[:HAS_ROLE]->(editR)
CREATE (eOwner:Employee{uid: 'e-owner', fullName: 'Olga Owner', firstName: 'Olga', lastName: 'Owner'})-[:HAS_USER]->(owner)
CREATE (eAnn:Employee{uid: 'e-ann', firstName: 'Ann', lastName: 'Bell'})
CREATE (eBob:Employee{uid: 'e-bob', fullName: 'Bob Carr'})
CREATE (team:Team{uid: 't-1', name: 'Vacuum', code: 'VAC'})
CREATE (zone:Zone{uid: 'z-1', name: 'Zone 1', code: 'Z1'})
CREATE (loc:Location{uid: 'l-1', name: 'Hall', code: 'H1'})-[:BELONGS_TO_FACILITY]->(f)
CREATE (st:SystemType{uid: 'st-1', name: 'Pump', code: 'PMP'})
CREATE (usage:ItemUsage{uid: 'iu-1', name: 'In use'})
CREATE (usage2:ItemUsage{uid: 'iu-2', name: 'Spare'})
CREATE (cond:ItemCondition{uid: 'ic-1', name: 'OK'})
CREATE (root:System{uid: 's-root', name: 'Root', systemLevel: 'SYSTEM_DOMAIN', deleted: false})-[:BELONGS_TO_FACILITY]->(f)
CREATE (tu:System{uid: 's-tu', name: 'TU', systemLevel: 'TECHNOLOGY_UNIT', deleted: false, lastUpdateTime: datetime('2026-01-01T10:00:00.123456789Z')})-[:BELONGS_TO_FACILITY]->(f)
CREATE (ks:System{uid: 's-ks', name: 'Key', systemLevel: 'KEY_SYSTEMS', deleted: false, minimalSpareParstCount: 2})-[:BELONGS_TO_FACILITY]->(f)
CREATE (spare:System{uid: 's-spare', name: 'Spare pump', systemLevel: 'SUBSYSTEMS_AND_PARTS', deleted: false})-[:BELONGS_TO_FACILITY]->(f)
CREATE (other:System{uid: 's-other', name: 'Other', deleted: false, systemCode: 'TAKEN-1'})-[:BELONGS_TO_FACILITY]->(f)
CREATE (root)-[:HAS_SUBSYSTEM]->(tu)
CREATE (tu)-[:HAS_SUBSYSTEM]->(ks)
CREATE (tu)-[:HAS_SUBSYSTEM]->(spare)
CREATE (root)-[:HAS_RESPONSIBLE]->(eOwner)
CREATE (tu)-[:HAS_LOCATION]->(loc)
CREATE (tu)-[:HAS_ZONE]->(zone)
CREATE (tu)-[:HAS_RESPONSIBLE]->(eOwner)
CREATE (spare)-[:IS_SPARE_FOR{coverage: 1.0}]->(ks)
CREATE (spare)-[:IS_POWERED_FROM]->(tu)
CREATE (cat:CatalogueCategory{uid: 'cc-1', name: 'Pumps'})
CREATE (grp:CatalogueCategoryPropertyGroup{uid: 'g-1', name: 'Physical'})
CREATE (pType:CatalogueCategoryPropertyType{uid: 'pt-1', name: 'Text'})
CREATE (rType:CatalogueCategoryPropertyType{uid: 'pt-2', name: 'Range'})
CREATE (pWeight:CatalogueCategoryProperty{uid: 'p-weight', name: 'Weight'})-[:IS_PROPERTY_TYPE]->(pType)
CREATE (pRange:CatalogueCategoryProperty{uid: 'p-range', name: 'Range'})-[:IS_PROPERTY_TYPE]->(rType)
CREATE (pColor:CatalogueCategoryProperty{uid: 'p-color', name: 'Color'})-[:IS_PROPERTY_TYPE]->(pType)
CREATE (grp)-[:CONTAINS_PROPERTY]->(pWeight)
CREATE (ci:CatalogueItem{uid: 'ci-1', name: 'Pump X', catalogueNumber: 'PX-1'})-[:BELONGS_TO_CATEGORY]->(cat)
CREATE (ci)-[:HAS_CATALOGUE_PROPERTY{value: '10'}]->(pWeight)
CREATE (ci)-[:HAS_CATALOGUE_PROPERTY{value: '{"min":1,"max":5}'}]->(pRange)
CREATE (itm:Item{uid: 'i-1', name: 'Pump X #1', eun: 'EUN1', serialNumber: 'SN1'})-[:IS_BASED_ON]->(ci)
CREATE (itm)-[:HAS_ITEM_USAGE]->(usage)
CREATE (ks)-[:CONTAINS_ITEM]->(itm)
CREATE (o:Order{uid: 'o-1', name: 'PO-1', orderDate: datetime('2025-05-01T00:00:00Z')})-[:HAS_ORDER_LINE{isDelivered: true}]->(itm)
CREATE (svcOld:ServiceItem{uid: 'si-old', name: 'Service old', isDelivered: true})
CREATE (svcNew:ServiceItem{uid: 'si-new', name: 'Service new', isDelivered: false})
CREATE (itm)-[:IS_SERVICED_BY{created: datetime('2025-06-01T00:00:00Z')}]->(svcOld)
CREATE (itm)-[:IS_SERVICED_BY{created: datetime('2025-07-01T00:00:00Z')}]->(svcNew)
CREATE (svcOld)-[:HAS_CATALOGUE_PROPERTY{value: '11'}]->(pWeight)
CREATE (svcNew)-[:HAS_CATALOGUE_PROPERTY{value: '12'}]->(pWeight)
CREATE (svcNew)-[:HAS_CATALOGUE_PROPERTY{value: 'red'}]->(pColor)
CREATE (o)-[:HAS_SERVICE_LINE]->(svcNew)
`

func integrationService(t *testing.T) (*SystemsService, neo4j.Driver) {
	t.Helper()
	uri := os.Getenv("PANDA_INTEGRATION_NEO4J_URI")
	if uri == "" {
		t.Skip("PANDA_INTEGRATION_NEO4J_URI not set")
	}
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth("neo4j", os.Getenv("PANDA_INTEGRATION_NEO4J_PASSWORD"), ""))
	require.NoError(t, err)
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err = session.Run(`MATCH (n) DETACH DELETE n`, nil)
	require.NoError(t, err)
	_, err = session.Run(integrationFixture, nil)
	require.NoError(t, err)
	t.Cleanup(func() { driver.Close() })
	return &SystemsService{neo4jDriver: &driver}, driver
}

func lastAudit(t *testing.T, driver neo4j.Driver, systemUID string) map[string]interface{} {
	t.Helper()
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	rec, err := session.Run(`MATCH (:System{uid: $uid})-[r:WAS_UPDATED_BY]->(u:User)
		RETURN {action: r.action, changes: r.changes, outside: r.outsideResponsibility, user: u.uid} AS a
		ORDER BY r.at DESC LIMIT 1`, map[string]interface{}{"uid": systemUID})
	require.NoError(t, err)
	single, err := rec.Single()
	require.NoError(t, err)
	a, _ := single.Get("a")
	return a.(map[string]interface{})
}

func decodeChanges(t *testing.T, audit map[string]interface{}) []map[string]interface{} {
	t.Helper()
	var changes []map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(audit["changes"].(string)), &changes))
	return changes
}

func strPtr(s string) *string { return &s }

func TestIntegration_GetSystemDetail(t *testing.T) {
	svc, _ := integrationService(t)

	d, err := svc.GetSystemDetail("s-ks", "B")
	require.NoError(t, err)
	assert.Equal(t, "Key", d.Name)
	assert.Equal(t, "s-tu", *d.ParentUID)
	require.Len(t, d.ParentPath, 2)
	assert.Equal(t, "s-root", d.ParentPath[0].UID)
	assert.Equal(t, "TECHNOLOGY_UNIT", *d.ParentPath[1].SystemLevel)
	assert.Nil(t, d.LastUpdateTime)
	require.NotNil(t, d.PhysicalItem)
	assert.Equal(t, "SN1", *d.PhysicalItem.SerialNumber)
	assert.Equal(t, "iu-1", d.PhysicalItem.ItemUsage.UID)
	assert.Equal(t, "PX-1", *d.PhysicalItem.CatalogueItem.CatalogueNumber)
	assert.Equal(t, int64(1), d.Statistics.SparePartsCount)
	assert.Equal(t, []models.EmployeeRef{}, d.Operators)

	tu, err := svc.GetSystemDetail("s-tu", "B")
	require.NoError(t, err)
	assert.Equal(t, int64(2), tu.Statistics.DirectSubsystemsCount)
	assert.Equal(t, "Olga Owner", tu.Responsible.Name)
	assert.Equal(t, "H1", *tu.Location.Code)
	require.NotNil(t, tu.LastUpdateTime)

	_, err = svc.GetSystemDetail("missing", "B")
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND)
}

func TestIntegration_PatchSystemScalarsCodebooksAndAudit(t *testing.T) {
	svc, driver := integrationService(t)

	operators := []models.UIDRef{{UID: "e-bob"}, {UID: "e-ann"}}
	d, err := svc.PatchSystem("s-ks", "B", "u-editor", &models.PatchSystemFields{
		Token:       models.PatchToken{Present: true}, // legacy null token
		Name:        strPtr("Key renamed"),
		Description: &models.Optional[string]{Value: strPtr("desc")},
		Zone:        &models.Optional[models.UIDRef]{Value: &models.UIDRef{UID: "z-1"}},
		Responsible: &models.Optional[models.UIDRef]{Value: &models.UIDRef{UID: "e-ann"}},
		Operators:   &operators,
	})
	require.NoError(t, err)
	assert.Equal(t, "Key renamed", d.Name)
	assert.Equal(t, "z-1", d.Zone.UID)
	assert.Equal(t, "Bell Ann", d.Responsible.Name, "no fullName -> \"last first\"")
	assert.Equal(t, []models.EmployeeRef{{UID: "e-ann", Name: "Bell Ann"}, {UID: "e-bob", Name: "Bob Carr"}}, d.Operators)
	require.NotNil(t, d.LastUpdateTime)

	audit := lastAudit(t, driver, "s-ks")
	assert.Equal(t, "UPDATE", audit["action"])
	assert.Equal(t, true, audit["outside"], "editor is not responsible for s-ks (owner is via root)")
	changes := decodeChanges(t, audit)
	fields := map[string]map[string]interface{}{}
	for _, c := range changes {
		fields[c["field"].(string)] = c
	}
	assert.Equal(t, "string", fields["name"]["type"])
	assert.Equal(t, "Key", fields["name"]["oldValue"])
	assert.Equal(t, "codebook", fields["zone"]["type"])
	assert.Equal(t, "Zone 1", fields["zone"]["newValue"].(map[string]interface{})["name"])
	assert.Equal(t, "codebookList", fields["operators"]["type"])
	assert.Len(t, fields["operators"]["newValue"], 2)

	// clear via null; stale token rejected
	_, err = svc.PatchSystem("s-ks", "B", "u-editor", &models.PatchSystemFields{
		Token: models.PatchToken{Present: true},
		Zone:  &models.Optional[models.UIDRef]{},
	})
	assert.ErrorIs(t, err, ErrStaleSystem)

	d2, err := svc.PatchSystem("s-ks", "B", "u-owner", &models.PatchSystemFields{
		Token: models.PatchToken{Present: true, Value: d.LastUpdateTime},
		Zone:  &models.Optional[models.UIDRef]{},
	})
	require.NoError(t, err)
	assert.Nil(t, d2.Zone)
	assert.Nil(t, lastAudit(t, driver, "s-ks")["outside"], "owner is responsible via root")

	// no-op: no write, token unchanged
	d3, err := svc.PatchSystem("s-ks", "B", "u-owner", &models.PatchSystemFields{Name: strPtr("Key renamed")})
	require.NoError(t, err)
	assert.True(t, d3.LastUpdateTime.Equal(*d2.LastUpdateTime))
}

func TestIntegration_PatchSystemValidation(t *testing.T) {
	svc, _ := integrationService(t)

	_, err := svc.PatchSystem("s-ks", "B", "u-editor", &models.PatchSystemFields{SystemCode: &models.Optional[string]{Value: strPtr("TAKEN-1")}})
	assert.ErrorIs(t, err, helpers.ERR_DUPLICATE_SYSTEM_CODE)

	_, err = svc.PatchSystem("s-ks", "B", "u-editor", &models.PatchSystemFields{Zone: &models.Optional[models.UIDRef]{Value: &models.UIDRef{UID: "nope"}}})
	var ve *ValidationError
	assert.ErrorAs(t, err, &ve)

	_, err = svc.PatchSystem("s-ks", "B", "u-editor", &models.PatchSystemFields{SystemLevel: &models.Optional[string]{Value: strPtr("BOGUS")}})
	assert.ErrorAs(t, err, &ve)
}

func TestIntegration_PatchTokenNanosecondRoundTrip(t *testing.T) {
	svc, _ := integrationService(t)

	tu, err := svc.GetSystemDetail("s-tu", "B")
	require.NoError(t, err)
	// simulate the FE: token travels as JSON and back
	b, _ := json.Marshal(tu.LastUpdateTime)
	var token time.Time
	require.NoError(t, json.Unmarshal(b, &token))
	_, err = svc.PatchSystem("s-tu", "B", "u-owner", &models.PatchSystemFields{
		Token: models.PatchToken{Present: true, Value: &token},
		Name:  strPtr("TU 2"),
	})
	assert.NoError(t, err)
}

func TestIntegration_MinimalSparesRecalculatesCoverage(t *testing.T) {
	svc, _ := integrationService(t)

	four := 4.0
	d, err := svc.PatchSystem("s-ks", "B", "u-owner", &models.PatchSystemFields{MinimalSpareParstCount: &models.Optional[float64]{Value: &four}})
	require.NoError(t, err)
	assert.Equal(t, 4.0, *d.MinimalSpareParstCount)
	assert.Equal(t, 1.0, *d.SparePartsCoverageSum)
	assert.Equal(t, 0.25, *d.SpCoverage)
}

func TestIntegration_PatchPhysicalItem(t *testing.T) {
	svc, driver := integrationService(t)

	d, err := svc.PatchSystemPhysicalItem("s-ks", "B", "u-owner", &models.PatchPhysicalItemFields{
		SerialNumber:    &models.Optional[string]{Value: strPtr("SN2")},
		Notes:           &models.Optional[string]{Value: strPtr("note")},
		ItemUsage:       &models.Optional[models.UIDRef]{Value: &models.UIDRef{UID: "iu-2"}},
		ConditionStatus: &models.Optional[models.UIDRef]{Value: &models.UIDRef{UID: "ic-1"}},
	})
	require.NoError(t, err)
	assert.Equal(t, "SN2", *d.PhysicalItem.SerialNumber)
	assert.Equal(t, "iu-2", d.PhysicalItem.ItemUsage.UID)
	assert.Equal(t, "ic-1", d.PhysicalItem.ConditionStatus.UID)
	require.NotNil(t, d.LastUpdateTime, "system token bumped by item edit")
	assert.Len(t, decodeChanges(t, lastAudit(t, driver, "s-ks")), 4)

	_, err = svc.PatchSystemPhysicalItem("s-tu", "B", "u-owner", &models.PatchPhysicalItemFields{Notes: &models.Optional[string]{}})
	assert.ErrorIs(t, err, ErrNoPhysicalItem)
}

func TestIntegration_CreateSubsystem(t *testing.T) {
	svc, driver := integrationService(t)

	d, err := svc.CreateSubsystem("s-tu", "B", "u-editor", &models.CreateSubsystemRequest{Name: " New ", SystemLevel: "KEY_SYSTEMS"})
	require.NoError(t, err)
	assert.Equal(t, "New", d.Name)
	assert.Equal(t, "s-tu", *d.ParentUID)
	assert.Equal(t, "z-1", d.Zone.UID)
	assert.Equal(t, "l-1", d.Location.UID)
	assert.Equal(t, "e-owner", d.Responsible.UID)
	audit := lastAudit(t, driver, d.UID)
	assert.Equal(t, "INSERT", audit["action"])
	assert.Equal(t, true, audit["outside"])

	_, err = svc.CreateSubsystem("s-root", "B", "u-editor", &models.CreateSubsystemRequest{Name: "x", SystemLevel: "KEY_SYSTEMS"})
	var ve *ValidationError
	assert.ErrorAs(t, err, &ve)
}

func TestIntegration_DeleteRelationshipRecalculatesSpares(t *testing.T) {
	svc, driver := integrationService(t)

	err := svc.DeleteRelationship("s-spare", "IS_SPARE_FOR", "s-ks", "B", "u-owner")
	require.NoError(t, err)

	d, err := svc.GetSystemDetail("s-ks", "B")
	require.NoError(t, err)
	assert.Equal(t, int64(0), d.Statistics.SparePartsCount)
	assert.Equal(t, 0.0, *d.SpCoverage)

	fromChanges := decodeChanges(t, lastAudit(t, driver, "s-spare"))
	require.Len(t, fromChanges, 1)
	assert.Equal(t, "relationship", fromChanges[0]["type"])
	assert.Equal(t, "outbound", fromChanges[0]["oldValue"].(map[string]interface{})["direction"])
	toChanges := decodeChanges(t, lastAudit(t, driver, "s-ks"))
	assert.Equal(t, "inbound", toChanges[0]["oldValue"].(map[string]interface{})["direction"])

	assert.ErrorIs(t, svc.DeleteRelationship("s-spare", "IS_SPARE_FOR", "s-ks", "B", "u-owner"), ErrRelationNotFound)
	var ve *ValidationError
	assert.ErrorAs(t, svc.DeleteRelationship("s-tu", "HAS_SUBSYSTEM", "s-ks", "B", "u-owner"), &ve)
	assert.NoError(t, svc.DeleteRelationship("s-spare", "IS_POWERED_FROM", "s-tu", "B", "u-owner"))
}

func TestIntegration_CreateRelationshipAudit(t *testing.T) {
	svc, driver := integrationService(t)

	_, err := svc.CreateNewSystemRelationship(&models.SystemRelationshipRequest{RelationTypeCode: "IS_SPARE_FOR", SystemFromUID: "s-other", SystemToUID: "s-ks"}, "B", "u-editor")
	require.NoError(t, err)
	audit := lastAudit(t, driver, "s-other")
	assert.Nil(t, audit["outside"], "s-other has no responsibles anywhere -> orphan, editable by anyone")
	assert.Equal(t, true, lastAudit(t, driver, "s-ks")["outside"], "s-ks is owned via root")
	changes := decodeChanges(t, audit)
	assert.Equal(t, "IS_SPARE_FOR", changes[0]["field"])
	assert.Equal(t, "Key", changes[0]["newValue"].(map[string]interface{})["name"])

	_, err = svc.CreateBatchRelationships(&models.BatchRelationshipRequest{SourceUids: []string{"s-other"}, TargetUids: []string{"s-tu"}, RelationshipType: "IS_COOLED_FROM"}, "B", "u-owner")
	require.NoError(t, err)
	changes = decodeChanges(t, lastAudit(t, driver, "s-tu"))
	assert.Equal(t, "inbound", changes[0]["newValue"].(map[string]interface{})["direction"])
}

func TestIntegration_SpareParts(t *testing.T) {
	svc, _ := integrationService(t)

	res, err := svc.GetSystemSpareParts("s-ks", "B")
	require.NoError(t, err)
	require.Len(t, res.SpareParts, 1)
	assert.Equal(t, "s-spare", res.SpareParts[0].System.UID)
	assert.Equal(t, 1.0, *res.SpareParts[0].Coverage)
	require.Len(t, res.SpareParts[0].System.ParentPath, 2)
	assert.Empty(t, res.SpareFor)

	res, err = svc.GetSystemSpareParts("s-spare", "B")
	require.NoError(t, err)
	require.Len(t, res.SpareFor, 1)
	assert.Equal(t, "EUN1", *res.SpareFor[0].System.PhysicalItem.EUN)
}

func TestIntegration_CatalogueProperties(t *testing.T) {
	svc, _ := integrationService(t)

	res, err := svc.GetPhysicalItemCatalogueProperties("i-1")
	require.NoError(t, err)
	assert.True(t, res.HasOverriddenProperties)
	byName := map[string]models.ResolvedCatalogueProperty{}
	for _, g := range res.Groups {
		for _, p := range g.Properties {
			byName[g.Key+"/"+p.Name] = p
		}
	}
	weight := byName["g-1/Weight"]
	assert.Equal(t, "10", *weight.Value)
	assert.Equal(t, "12", *weight.ServiceValue, "newest service item wins")
	assert.True(t, weight.IsOverridden)
	assert.Equal(t, "1 - 5", *byName["no-group/Range"].Value)
	color := byName["no-group/Color"]
	assert.Nil(t, color.Value)
	assert.Equal(t, "red", *color.ServiceValue)
	assert.Equal(t, "no-group", res.Groups[0].Key)
}

func TestIntegration_Orders(t *testing.T) {
	svc, _ := integrationService(t)

	res, err := svc.GetPhysicalItemOrders("i-1")
	require.NoError(t, err)
	require.NotNil(t, res.Order)
	assert.Equal(t, "PO-1", res.Order.Name)
	assert.True(t, *res.Order.IsDelivered)
	require.Len(t, res.ServiceItems, 2)
	assert.Equal(t, "si-new", res.ServiceItems[0].UID)
	assert.Equal(t, "PO-1", res.ServiceItems[0].Order.Name)
}

func TestIntegration_DeleteSystemRecursiveAudit(t *testing.T) {
	svc, driver := integrationService(t)

	require.NoError(t, svc.DeleteSystemRecursive("s-tu", "u-editor"))
	for _, uid := range []string{"s-tu", "s-ks", "s-spare"} {
		audit := lastAudit(t, driver, uid)
		assert.Equal(t, "DELETE", audit["action"], uid)
		assert.Equal(t, true, audit["outside"], uid)
	}
	_, err := svc.GetSystemDetail("s-ks", "B")
	assert.ErrorIs(t, err, helpers.ERR_NOT_FOUND)
}

func TestIntegration_GraphNodeItemUsage(t *testing.T) {
	_, driver := integrationService(t)

	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	nodes, err := helpers.GetNeo4jArrayOfNodes[models.SystemGraphNode](session, GetSystemGraphNodesByUidsQuery([]string{"s-ks", "s-tu"}, "B"))
	require.NoError(t, err)
	byUID := map[string]models.SystemGraphNode{}
	for _, n := range nodes {
		byUID[n.Uid] = n
	}
	assert.Equal(t, "iu-1", byUID["s-ks"].Properties["itemUsageUid"])
	_, has := byUID["s-tu"].Properties["itemUsageUid"]
	assert.False(t, has)
}

func TestIntegration_GenerateSystemCodeRequiresType(t *testing.T) {
	svc, _ := integrationService(t)

	_, err := svc.GenerateSystemCode("s-ks", "B", "u-owner", models.PatchToken{})
	assert.ErrorIs(t, err, ErrSystemTypeMissing)

	// type without a serial mask -> user-fixable 400, not a 500
	_, err = svc.PatchSystem("s-ks", "B", "u-owner", &models.PatchSystemFields{SystemType: &models.Optional[models.UIDRef]{Value: &models.UIDRef{UID: "st-1"}}})
	require.NoError(t, err)
	_, err = svc.GenerateSystemCode("s-ks", "B", "u-owner", models.PatchToken{})
	var ve *ValidationError
	assert.ErrorAs(t, err, &ve)
}

func TestIntegration_HistoryExposesOutsideResponsibility(t *testing.T) {
	svc, _ := integrationService(t)

	_, err := svc.PatchSystem("s-ks", "B", "u-editor", &models.PatchSystemFields{Name: strPtr("by editor")})
	require.NoError(t, err)
	history, err := svc.GetSystemHistory("s-ks")
	require.NoError(t, err)
	require.NotEmpty(t, history)
	found := false
	for _, h := range history {
		if h.Action == "UPDATE" && h.OutsideResponsibility {
			found = true
		}
	}
	assert.True(t, found)
}
