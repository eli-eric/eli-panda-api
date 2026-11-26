package systemsService

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var (
	ErrStaleSystem       = errors.New("system was modified by someone else")
	ErrNoPhysicalItem    = errors.New("system has no physical item")
	ErrRelationNotFound  = errors.New("relationship not found")
	ErrSystemTypeMissing = errors.New("system has no system type")
)

// ValidationError is a client error (HTTP 400) raised by the PATCH family.
type ValidationError struct{ Msg string }

func (e *ValidationError) Error() string { return e.Msg }

// mapWriteError turns the in-transaction duplicate-code guard (apoc.util.validate) into
// ERR_DUPLICATE_SYSTEM_CODE; the pre-check only gives the common case a cheap early 409.
func mapWriteError(err error) error {
	if err != nil && strings.Contains(err.Error(), duplicateSystemCodeMarker) {
		return helpers.ERR_DUPLICATE_SYSTEM_CODE
	}
	return err
}

func validationErr(format string, args ...interface{}) error {
	return &ValidationError{Msg: fmt.Sprintf(format, args...)}
}

var systemLevels = map[string]struct{}{
	"SYSTEM_DOMAIN": {}, "TECHNOLOGY_UNIT": {}, "KEY_SYSTEMS": {}, "SUBSYSTEMS_AND_PARTS": {}, "TRASH": {},
}

// allowedChildLevels mirrors the FE create-subsystem rules (parent level -> allowed child levels).
var allowedChildLevels = map[string][]string{
	"SYSTEM_DOMAIN":        {"TECHNOLOGY_UNIT"},
	"TECHNOLOGY_UNIT":      {"TECHNOLOGY_UNIT", "KEY_SYSTEMS", "TRASH"},
	"KEY_SYSTEMS":          {"KEY_SYSTEMS", "SUBSYSTEMS_AND_PARTS", "TRASH"},
	"SUBSYSTEMS_AND_PARTS": {"SUBSYSTEMS_AND_PARTS", "TRASH"},
	"TRASH":                {},
}

// EditableRelationshipTypes may be removed via DELETE /v1/system/{uid}/relationships/{type}/{otherUid}.
var EditableRelationshipTypes = map[string]struct{}{
	"IS_SPARE_FOR": {}, "IS_COOLED_FROM": {}, "IS_POWERED_FROM": {}, "IS_CONTROLLED_BY": {},
	"IS_INTERLOCKED_BY": {}, "PROVIDES_DATA_TO": {}, "DIRECTS_BEAM_TO": {}, "PROVIDES_VACUUM_FOR": {},
}

type codebookRelSpec struct {
	field          string
	label          string
	relType        string
	employee       bool
	facilityScoped bool // label has BELONGS_TO_FACILITY; refs must belong to the caller's facility
}

var (
	specSystemType      = codebookRelSpec{"systemType", "SystemType", "HAS_SYSTEM_TYPE", false, false}
	specLocation        = codebookRelSpec{"location", "Location", "HAS_LOCATION", false, true}
	specZone            = codebookRelSpec{"zone", "Zone", "HAS_ZONE", false, true}
	specImportance      = codebookRelSpec{"importance", "SystemImportance", "HAS_IMPORTANCE", false, false}
	specAttribute       = codebookRelSpec{"attribute", "SystemAttribute", "HAS_SYSTEM_ATTRIBUTE", false, true}
	specResponsible     = codebookRelSpec{"responsible", "Employee", "HAS_RESPONSIBLE", true, false}
	specResponsibleTeam = codebookRelSpec{"responsibleTeam", "Team", "HAS_RESPONSIBLE_TEAM", false, true}
	specOwner           = codebookRelSpec{"owner", "Employee", "HAS_OWNER", true, false}
	specOperators       = codebookRelSpec{"operators", "Employee", "HAS_OPERATOR", true, false}
	specMaintainedBy    = codebookRelSpec{"maintainedBy", "Employee", "IS_MAINTAINED_BY", true, false}
	specItemUsage       = codebookRelSpec{"itemUsage", "ItemUsage", "HAS_ITEM_USAGE", false, false}
	specConditionStatus = codebookRelSpec{"conditionStatus", "ItemCondition", "HAS_CONDITION_STATUS", false, false}
)

// normalizeOptString trims and maps "" to nil so clearing and emptying are the same change.
func normalizeOptString(v *string) *string {
	if v == nil {
		return nil
	}
	t := strings.TrimSpace(*v)
	if t == "" {
		return nil
	}
	return &t
}

func tokensEqual(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Equal(*b)
}

func (svc *SystemsService) GetSystemDetail(uid string, facilityCode string) (models.SystemDetail, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()
	return svc.getSystemDetail(session, uid, facilityCode)
}

func (svc *SystemsService) getSystemDetail(session neo4j.Session, uid, facilityCode string) (models.SystemDetail, error) {
	detail, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.SystemDetail](session, GetSystemDetailQuery(uid, facilityCode))
	if errors.Is(err, helpers.ERR_NO_ROWS) {
		return detail, helpers.ERR_NOT_FOUND
	}
	if detail.ParentPath == nil {
		detail.ParentPath = []models.SystemPathItem{}
	}
	if detail.Operators == nil {
		detail.Operators = []models.EmployeeRef{}
	}
	if detail.MaintainedBy == nil {
		detail.MaintainedBy = []models.EmployeeRef{}
	}
	return detail, err
}

// isOutsideResponsibility reports whether userUID is not responsible for systemUID. It only
// feeds the audit flag — editing is never blocked on it (soft lock).
func (svc *SystemsService) isOutsideResponsibility(systemUID, userUID string) (bool, error) {
	res, err := svc.CanEditSystem(systemUID, userUID)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return false, nil
		}
		return false, err
	}
	return !res.Result, nil
}

// outsideResponsibilityByUID evaluates isOutsideResponsibility for each system uid.
func (svc *SystemsService) outsideResponsibilityByUID(userUID string, uids ...string) (map[string]bool, error) {
	result := make(map[string]bool, len(uids))
	for _, uid := range uids {
		if _, done := result[uid]; done || uid == "" {
			continue
		}
		outside, err := svc.isOutsideResponsibility(uid, userUID)
		if err != nil {
			return nil, err
		}
		result[uid] = outside
	}
	return result, nil
}

func (svc *SystemsService) resolveRefs(session neo4j.Session, spec codebookRelSpec, uids []string, facilityCode string) ([]models.CodebookRef, error) {
	if len(uids) == 0 {
		return []models.CodebookRef{}, nil
	}
	refs, err := helpers.GetNeo4jArrayOfNodes[models.CodebookRef](session, ResolveRefsQuery(spec, uids, facilityCode))
	if err != nil {
		return nil, err
	}
	byUID := make(map[string]models.CodebookRef, len(refs))
	for _, r := range refs {
		byUID[r.UID] = r
	}
	ordered := make([]models.CodebookRef, 0, len(uids))
	for _, uid := range uids {
		ref, ok := byUID[uid]
		if !ok {
			return nil, validationErr("unknown %s uid: %s", spec.field, uid)
		}
		ordered = append(ordered, ref)
	}
	return ordered, nil
}

// planSingleRef diffs one {uid}|null codebook field and adds the relation + change to the plan.
func (svc *SystemsService) planSingleRef(session neo4j.Session, facilityCode string, plan *patchPlan, spec codebookRelSpec, patch *models.Optional[models.UIDRef], current *models.CodebookRef) error {
	if patch == nil {
		return nil
	}
	var next *models.CodebookRef
	uids := []string{}
	if patch.Value != nil {
		refs, err := svc.resolveRefs(session, spec, []string{patch.Value.UID}, facilityCode)
		if err != nil {
			return err
		}
		next = &refs[0]
		uids = []string{next.UID}
	}
	before := len(plan.changes)
	plan.changes = helpers.AppendIfChanged(plan.changes, spec.field, helpers.ChangeTypeCodebook, current, next)
	if len(plan.changes) > before {
		plan.relations = append(plan.relations, relationPatch{relType: spec.relType, label: spec.label, uids: uids})
	}
	return nil
}

func employeeToCodebook(e *models.EmployeeRef) *models.CodebookRef {
	if e == nil {
		return nil
	}
	return &models.CodebookRef{UID: e.UID, Name: e.Name}
}

func employeesToRefs(list []models.EmployeeRef) []models.EmployeeRef {
	if list == nil {
		return []models.EmployeeRef{}
	}
	return list
}

func (svc *SystemsService) planRefList(session neo4j.Session, facilityCode string, plan *patchPlan, spec codebookRelSpec, patch *[]models.UIDRef, current []models.EmployeeRef) error {
	if patch == nil {
		return nil
	}
	uids := make([]string, 0, len(*patch))
	for _, r := range *patch {
		uids = append(uids, r.UID)
	}
	refs, err := svc.resolveRefs(session, spec, uids, facilityCode)
	if err != nil {
		return err
	}
	next := make([]models.EmployeeRef, 0, len(refs))
	for _, r := range refs {
		next = append(next, models.EmployeeRef{UID: r.UID, Name: r.Name})
	}
	sort.Slice(next, func(i, j int) bool { return next[i].Name < next[j].Name })
	before := len(plan.changes)
	plan.changes = helpers.AppendIfChanged(plan.changes, spec.field, helpers.ChangeTypeCodebookList, employeesToRefs(current), next)
	if len(plan.changes) > before {
		plan.relations = append(plan.relations, relationPatch{relType: spec.relType, label: spec.label, uids: uids})
	}
	return nil
}

// planOptString diffs a nullable string property.
func planOptString(plan *patchPlan, field string, patch *models.Optional[string], current *string) {
	if patch == nil {
		return
	}
	next := normalizeOptString(patch.Value)
	before := len(plan.changes)
	plan.changes = helpers.AppendIfChanged(plan.changes, field, helpers.ChangeTypeString, normalizeOptString(current), next)
	if len(plan.changes) > before {
		if next == nil {
			plan.scalars[field] = nil
		} else {
			plan.scalars[field] = *next
		}
	}
}

func (svc *SystemsService) buildSystemPlan(session neo4j.Session, current *models.SystemDetail, f *models.PatchSystemFields, facilityCode string) (*patchPlan, error) {
	plan := newPatchPlan()

	if f.Name != nil {
		name := strings.TrimSpace(*f.Name)
		if name == "" {
			return nil, validationErr("name must not be empty")
		}
		before := len(plan.changes)
		plan.changes = helpers.AppendIfChanged(plan.changes, "name", helpers.ChangeTypeString, current.Name, name)
		if len(plan.changes) > before {
			plan.scalars["name"] = name
		}
	}
	planOptString(plan, "description", f.Description, current.Description)
	planOptString(plan, "systemAlias", f.SystemAlias, current.SystemAlias)

	if f.SystemLevel != nil && f.SystemLevel.Value == nil {
		return nil, validationErr("systemLevel cannot be cleared")
	}
	if f.SystemLevel != nil {
		if _, ok := systemLevels[*f.SystemLevel.Value]; !ok {
			return nil, validationErr("invalid systemLevel: %s", *f.SystemLevel.Value)
		}
	}
	planOptString(plan, "systemLevel", f.SystemLevel, current.SystemLevel)

	if f.SystemCode != nil {
		if code := normalizeOptString(f.SystemCode.Value); code != nil {
			exists, err := helpers.GetNeo4jSingleRecordSingleValue[bool](session, SystemCodeExistsQuery(*code, current.UID, facilityCode))
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, helpers.ERR_DUPLICATE_SYSTEM_CODE
			}
		}
		planOptString(plan, "systemCode", f.SystemCode, current.SystemCode)
	}

	singles := []struct {
		spec    codebookRelSpec
		patch   *models.Optional[models.UIDRef]
		current *models.CodebookRef
	}{
		{specSystemType, f.SystemType, current.SystemType},
		{specLocation, f.Location, current.Location},
		{specZone, f.Zone, current.Zone},
		{specImportance, f.Importance, current.Importance},
		{specAttribute, f.Attribute, current.Attribute},
		{specResponsible, f.Responsible, employeeToCodebook(current.Responsible)},
		{specResponsibleTeam, f.ResponsibleTeam, current.ResponsibleTeam},
		{specOwner, f.Owner, employeeToCodebook(current.Owner)},
	}
	for _, s := range singles {
		if err := svc.planSingleRef(session, facilityCode, plan, s.spec, s.patch, s.current); err != nil {
			return nil, err
		}
	}
	if err := svc.planRefList(session, facilityCode, plan, specOperators, f.Operators, current.Operators); err != nil {
		return nil, err
	}
	if err := svc.planRefList(session, facilityCode, plan, specMaintainedBy, f.MaintainedBy, current.MaintainedBy); err != nil {
		return nil, err
	}

	if f.MinimalSpareParstCount != nil {
		var next *float64
		if v := f.MinimalSpareParstCount.Value; v != nil {
			if *v < 0 || *v != math.Trunc(*v) {
				return nil, validationErr("minimalSpareParstCount must be a non-negative integer")
			}
			next = v
		}
		before := len(plan.changes)
		plan.changes = helpers.AppendIfChanged(plan.changes, "minimalSpareParstCount", helpers.ChangeTypeNumber, current.MinimalSpareParstCount, next)
		if len(plan.changes) > before {
			if next == nil {
				plan.scalars["minimalSpareParstCount"] = nil
			} else {
				plan.scalars["minimalSpareParstCount"] = int64(*next)
			}
			plan.recalcCoverage = true
		}
	}

	return plan, nil
}

func (svc *SystemsService) buildItemPlan(session neo4j.Session, facilityCode string, item *models.SystemDetailPhysicalItem, f *models.PatchPhysicalItemFields) (*patchPlan, error) {
	plan := newPatchPlan()
	planOptString(plan, "serialNumber", f.SerialNumber, item.SerialNumber)
	planOptString(plan, "notes", f.Notes, item.Notes)
	if err := svc.planSingleRef(session, facilityCode, plan, specItemUsage, f.ItemUsage, item.ItemUsage); err != nil {
		return nil, err
	}
	if err := svc.planSingleRef(session, facilityCode, plan, specConditionStatus, f.ConditionStatus, item.ConditionStatus); err != nil {
		return nil, err
	}
	return plan, nil
}

// writePatch executes the plans guarded by the token read in current and returns fresh detail.
func (svc *SystemsService) writePatch(session neo4j.Session, current *models.SystemDetail, facilityCode, userUID string, systemPlan, itemPlan *patchPlan, changes []helpers.ChangeEntry) (models.SystemDetail, error) {
	outside, err := svc.isOutsideResponsibility(current.UID, userUID)
	if err != nil {
		return models.SystemDetail{}, err
	}
	query := PatchSystemQuery(current.UID, facilityCode, userUID, current.LastUpdateTime, systemPlan, itemPlan, helpers.MarshalChanges(changes), outside)
	if _, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query); err != nil {
		// no rows = the token changed between our read and the write
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return models.SystemDetail{}, ErrStaleSystem
		}
		return models.SystemDetail{}, mapWriteError(err)
	}
	return svc.getSystemDetail(session, current.UID, facilityCode)
}

// loadForPatch reads the current detail and checks the client token (when sent).
func (svc *SystemsService) loadForPatch(session neo4j.Session, uid, facilityCode string, token models.PatchToken) (models.SystemDetail, error) {
	current, err := svc.getSystemDetail(session, uid, facilityCode)
	if err != nil {
		return current, err
	}
	if token.Present && !tokensEqual(token.Value, current.LastUpdateTime) {
		return current, ErrStaleSystem
	}
	return current, nil
}

func (svc *SystemsService) PatchSystem(uid, facilityCode, userUID string, f *models.PatchSystemFields) (models.SystemDetail, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	current, err := svc.loadForPatch(session, uid, facilityCode, f.Token)
	if err != nil {
		return current, err
	}
	plan, err := svc.buildSystemPlan(session, &current, f, facilityCode)
	if err != nil {
		return current, err
	}
	if plan.isEmpty() {
		return current, nil
	}
	return svc.writePatch(session, &current, facilityCode, userUID, plan, nil, plan.changes)
}

func (svc *SystemsService) PatchSystemPhysicalItem(uid, facilityCode, userUID string, f *models.PatchPhysicalItemFields) (models.SystemDetail, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	current, err := svc.loadForPatch(session, uid, facilityCode, f.Token)
	if err != nil {
		return current, err
	}
	if current.PhysicalItem == nil {
		return current, ErrNoPhysicalItem
	}
	plan, err := svc.buildItemPlan(session, facilityCode, current.PhysicalItem, f)
	if err != nil {
		return current, err
	}
	if plan.isEmpty() {
		return current, nil
	}
	return svc.writePatch(session, &current, facilityCode, userUID, nil, plan, plan.changes)
}

func isAllowedChildLevel(parentLevel *string, childLevel string) bool {
	if parentLevel == nil {
		return true
	}
	allowed, known := allowedChildLevels[*parentLevel]
	if !known {
		return true
	}
	for _, l := range allowed {
		if l == childLevel {
			return true
		}
	}
	return false
}

// inheritFromParent fills absent responsible/location/zone keys with the parent's values.
// An explicit null in the request is kept (no inheritance).
func inheritFromParent(f *models.PatchSystemFields, parent *models.SystemDetail) {
	ref := func(uid string) *models.Optional[models.UIDRef] {
		return &models.Optional[models.UIDRef]{Value: &models.UIDRef{UID: uid}}
	}
	if f.Responsible == nil && parent.Responsible != nil {
		f.Responsible = ref(parent.Responsible.UID)
	}
	if f.Location == nil && parent.Location != nil {
		f.Location = ref(parent.Location.UID)
	}
	if f.Zone == nil && parent.Zone != nil {
		f.Zone = ref(parent.Zone.UID)
	}
}

// CreateSystem creates a root system or a child of f.ParentUID with the given fields.
func (svc *SystemsService) CreateSystem(facilityCode, userUID string, f *models.CreateSystemFields) (models.SystemDetail, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	fields := &f.Fields
	if fields.Name == nil || strings.TrimSpace(*fields.Name) == "" {
		return models.SystemDetail{}, validationErr("name must not be empty")
	}
	if fields.SystemLevel == nil || fields.SystemLevel.Value == nil {
		return models.SystemDetail{}, validationErr("systemLevel is required")
	}
	level := *fields.SystemLevel.Value

	outside := false
	if f.ParentUID != nil {
		parent, err := svc.getSystemDetail(session, *f.ParentUID, facilityCode)
		if err != nil {
			return parent, err
		}
		if !isAllowedChildLevel(parent.SystemLevel, level) {
			return models.SystemDetail{}, validationErr("systemLevel %s is not allowed under %s", level, *parent.SystemLevel)
		}
		inheritFromParent(fields, &parent)
		if outside, err = svc.isOutsideResponsibility(parent.UID, userUID); err != nil {
			return models.SystemDetail{}, err
		}
	}

	newUID := uuid.NewString()
	empty := models.SystemDetail{UID: newUID, Operators: []models.EmployeeRef{}, MaintainedBy: []models.EmployeeRef{}}
	plan, err := svc.buildSystemPlan(session, &empty, fields, facilityCode)
	if err != nil {
		return models.SystemDetail{}, err
	}
	query := CreateSystemQuery(newUID, f.ParentUID, facilityCode, userUID, plan, helpers.MarshalChanges(plan.changes), outside)
	if _, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query); err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return models.SystemDetail{}, helpers.ERR_NOT_FOUND
		}
		return models.SystemDetail{}, mapWriteError(err)
	}
	return svc.getSystemDetail(session, newUID, facilityCode)
}

// CreateSubsystem creates a child of parentUID; responsible/location/zone are inherited.
func (svc *SystemsService) CreateSubsystem(parentUID, facilityCode, userUID string, req *models.CreateSubsystemRequest) (models.SystemDetail, error) {
	name, level := req.Name, req.SystemLevel
	return svc.CreateSystem(facilityCode, userUID, &models.CreateSystemFields{
		ParentUID: &parentUID,
		Fields:    models.PatchSystemFields{Name: &name, SystemLevel: &models.Optional[string]{Value: &level}},
	})
}

// GenerateSystemCode generates a code from the system's current type/zone/location/parent and
// saves it through the regular PATCH path (token check, uniqueness, audit).
func (svc *SystemsService) GenerateSystemCode(uid, facilityCode, userUID string, token models.PatchToken) (models.SystemDetail, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	current, err := svc.loadForPatch(session, uid, facilityCode, token)
	session.Close()
	if err != nil {
		return current, err
	}
	if current.SystemType == nil {
		return current, ErrSystemTypeMissing
	}
	codebookUIDOrEmpty := func(c *models.CodebookRef) string {
		if c == nil {
			return ""
		}
		return c.UID
	}
	parentUID := ""
	if current.ParentUID != nil {
		parentUID = *current.ParentUID
	}
	code, err := svc.GetSystemCode(current.SystemType.UID, codebookUIDOrEmpty(current.Zone), codebookUIDOrEmpty(current.Location), parentUID, facilityCode)
	if err != nil {
		// GetSystemCode reports user-fixable gaps ("missing zone", "missing location", ...)
		if strings.HasPrefix(err.Error(), "missing") {
			return current, validationErr("cannot generate system code: %s", err.Error())
		}
		return current, err
	}
	return svc.PatchSystem(uid, facilityCode, userUID, &models.PatchSystemFields{
		Token:      models.PatchToken{Present: true, Value: current.LastUpdateTime},
		SystemCode: &models.Optional[string]{Value: &code},
	})
}

type systemNameRef struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

func systemNameQuery(uid string) helpers.DatabaseQuery {
	return helpers.DatabaseQuery{
		Query:       `MATCH (s:System{uid: $uid}) RETURN {uid: s.uid, name: s.name} AS system`,
		ReturnAlias: "system",
		Parameters:  map[string]interface{}{"uid": uid},
	}
}

type relationshipChangeValue struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	Direction string `json:"direction"`
}

// relationshipAudits builds both audit sides for an add (added=true) or removal of (from)-[relType]->(to).
func (svc *SystemsService) relationshipAudits(session neo4j.Session, fromUID, relType, toUID, userUID string, added bool) (relationshipSideAudit, relationshipSideAudit, error) {
	from, err := helpers.GetNeo4jSingleRecordAndMapToStruct[systemNameRef](session, systemNameQuery(fromUID))
	if err != nil {
		return relationshipSideAudit{}, relationshipSideAudit{}, err
	}
	to, err := helpers.GetNeo4jSingleRecordAndMapToStruct[systemNameRef](session, systemNameQuery(toUID))
	if err != nil {
		return relationshipSideAudit{}, relationshipSideAudit{}, err
	}
	side := func(self string, other systemNameRef, direction string) (relationshipSideAudit, error) {
		value := &relationshipChangeValue{UID: other.UID, Name: other.Name, Direction: direction}
		var oldV, newV interface{}
		if added {
			newV = value
		} else {
			oldV = value
		}
		outside, err := svc.isOutsideResponsibility(self, userUID)
		if err != nil {
			return relationshipSideAudit{}, err
		}
		changes := helpers.AppendIfChanged(nil, relType, helpers.ChangeTypeRelationship, oldV, newV)
		return relationshipSideAudit{systemUID: self, changesJSON: helpers.MarshalChanges(changes), outsideResponsibility: outside}, nil
	}
	fromSide, err := side(fromUID, to, "outbound")
	if err != nil {
		return relationshipSideAudit{}, relationshipSideAudit{}, err
	}
	toSide, err := side(toUID, from, "inbound")
	return fromSide, toSide, err
}

func (svc *SystemsService) DeleteRelationship(fromUID, relType, toUID, facilityCode, userUID string) error {
	if _, ok := EditableRelationshipTypes[relType]; !ok {
		return validationErr("relationship type %s cannot be removed", relType)
	}
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	exists, err := helpers.GetNeo4jSingleRecordSingleValue[bool](session, RelationshipExistsQuery(fromUID, relType, toUID))
	if err != nil {
		return err
	}
	if !exists {
		return ErrRelationNotFound
	}
	fromSide, toSide, err := svc.relationshipAudits(session, fromUID, relType, toUID, userUID, false)
	if err != nil {
		return err
	}
	_, err = helpers.WriteNeo4jAndReturnSingleValue[bool](session, DeleteRelationshipQuery(fromUID, relType, toUID, facilityCode, userUID, fromSide, toSide))
	if errors.Is(err, helpers.ERR_NO_ROWS) {
		return ErrRelationNotFound
	}
	return err
}

func (svc *SystemsService) GetSystemSpareParts(uid, facilityCode string) (models.SystemSparePartsResponse, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()
	result, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.SystemSparePartsResponse](session, GetSystemSparePartsQuery(uid, facilityCode))
	if errors.Is(err, helpers.ERR_NO_ROWS) {
		return result, helpers.ERR_NOT_FOUND
	}
	if result.SpareParts == nil {
		result.SpareParts = []models.SpareRelationRef{}
	}
	if result.SpareFor == nil {
		result.SpareFor = []models.SpareRelationRef{}
	}
	for _, list := range [][]models.SpareRelationRef{result.SpareParts, result.SpareFor} {
		for i := range list {
			if list[i].System.ParentPath == nil {
				list[i].System.ParentPath = []models.SystemSimpleResponse{}
			}
		}
	}
	return result, err
}

func (svc *SystemsService) GetPhysicalItemCatalogueProperties(itemUID, facilityCode string) (models.PhysicalItemCatalogueProperties, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()
	raw, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.RawPhysicalItemProperties](session, GetPhysicalItemRawPropertiesQuery(itemUID, facilityCode))
	if errors.Is(err, helpers.ERR_NO_ROWS) {
		return models.PhysicalItemCatalogueProperties{}, helpers.ERR_NOT_FOUND
	}
	if err != nil {
		return models.PhysicalItemCatalogueProperties{}, err
	}
	return resolveCatalogueProperties(raw), nil
}

func (svc *SystemsService) GetPhysicalItemOrders(itemUID, facilityCode string) (models.PhysicalItemOrders, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()
	result, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.PhysicalItemOrders](session, GetPhysicalItemOrdersQuery(itemUID, facilityCode))
	if errors.Is(err, helpers.ERR_NO_ROWS) {
		return result, helpers.ERR_NOT_FOUND
	}
	if result.ServiceItems == nil {
		result.ServiceItems = []models.PhysicalItemServiceRef{}
	}
	return result, err
}
