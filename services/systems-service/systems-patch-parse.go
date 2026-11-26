package systemsService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"panda/apigateway/services/systems-service/models"
)

var patchSystemAllowedKeys = map[string]struct{}{
	"lastUpdateTime": {}, "name": {}, "description": {}, "systemCode": {}, "systemAlias": {},
	"systemLevel": {}, "systemType": {}, "location": {}, "zone": {}, "importance": {},
	"responsible": {}, "responsibleTeam": {}, "owner": {}, "operators": {}, "maintainedBy": {},
	"minimalSpareParstCount": {},
}

var patchPhysicalItemAllowedKeys = map[string]struct{}{
	"lastUpdateTime": {}, "serialNumber": {}, "notes": {}, "itemUsage": {}, "conditionStatus": {},
}

func rawIsNull(r json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(r), []byte("null"))
}

func rejectUnknownKeys(raw map[string]json.RawMessage, allowed map[string]struct{}) error {
	unknown := []string{}
	for k := range raw {
		if _, ok := allowed[k]; !ok {
			unknown = append(unknown, k)
		}
	}
	if len(unknown) == 0 {
		return nil
	}
	sort.Strings(unknown)
	return fmt.Errorf("unknown fields: %v", unknown)
}

// parseOptional decodes a nullable PATCH key: nil = absent, Value nil = explicit null.
func parseOptional[T any](raw map[string]json.RawMessage, key string) (*models.Optional[T], error) {
	r, ok := raw[key]
	if !ok {
		return nil, nil
	}
	if rawIsNull(r) {
		return &models.Optional[T]{}, nil
	}
	var v T
	if err := json.Unmarshal(r, &v); err != nil {
		return nil, fmt.Errorf("invalid %s", key)
	}
	return &models.Optional[T]{Value: &v}, nil
}

func parseOptionalRef(raw map[string]json.RawMessage, key string) (*models.Optional[models.UIDRef], error) {
	ref, err := parseOptional[models.UIDRef](raw, key)
	if err != nil {
		return nil, err
	}
	if ref != nil && ref.Value != nil && ref.Value.UID == "" {
		return nil, fmt.Errorf("invalid %s: uid is required", key)
	}
	return ref, nil
}

// parseRefList decodes a full-replace list of {uid}; null is rejected (use [] to clear).
func parseRefList(raw map[string]json.RawMessage, key string) (*[]models.UIDRef, error) {
	r, ok := raw[key]
	if !ok {
		return nil, nil
	}
	if rawIsNull(r) {
		return nil, fmt.Errorf("invalid %s: use [] to clear", key)
	}
	var refs []models.UIDRef
	if err := json.Unmarshal(r, &refs); err != nil {
		return nil, fmt.Errorf("invalid %s", key)
	}
	seen := map[string]struct{}{}
	unique := make([]models.UIDRef, 0, len(refs))
	for _, ref := range refs {
		if ref.UID == "" {
			return nil, fmt.Errorf("invalid %s: uid is required", key)
		}
		if _, dup := seen[ref.UID]; dup {
			continue
		}
		seen[ref.UID] = struct{}{}
		unique = append(unique, ref)
	}
	return &unique, nil
}

func parsePatchToken(raw map[string]json.RawMessage) (models.PatchToken, error) {
	r, ok := raw["lastUpdateTime"]
	if !ok {
		return models.PatchToken{}, nil
	}
	if rawIsNull(r) {
		return models.PatchToken{Present: true}, nil
	}
	var t time.Time
	if err := json.Unmarshal(r, &t); err != nil {
		return models.PatchToken{}, fmt.Errorf("invalid lastUpdateTime")
	}
	return models.PatchToken{Present: true, Value: &t}, nil
}

func parsePatchSystemPayload(raw map[string]json.RawMessage) (*models.PatchSystemFields, error) {
	if err := rejectUnknownKeys(raw, patchSystemAllowedKeys); err != nil {
		return nil, err
	}
	f := &models.PatchSystemFields{}
	var err error

	if f.Token, err = parsePatchToken(raw); err != nil {
		return nil, err
	}
	if r, ok := raw["name"]; ok {
		var v string
		if rawIsNull(r) || json.Unmarshal(r, &v) != nil {
			return nil, fmt.Errorf("invalid name")
		}
		f.Name = &v
	}
	strFields := []struct {
		key string
		dst **models.Optional[string]
	}{
		{"description", &f.Description}, {"systemCode", &f.SystemCode},
		{"systemAlias", &f.SystemAlias}, {"systemLevel", &f.SystemLevel},
	}
	for _, sf := range strFields {
		if *sf.dst, err = parseOptional[string](raw, sf.key); err != nil {
			return nil, err
		}
	}
	refFields := []struct {
		key string
		dst **models.Optional[models.UIDRef]
	}{
		{"systemType", &f.SystemType}, {"location", &f.Location}, {"zone", &f.Zone},
		{"importance", &f.Importance}, {"responsible", &f.Responsible},
		{"responsibleTeam", &f.ResponsibleTeam}, {"owner", &f.Owner},
	}
	for _, rf := range refFields {
		if *rf.dst, err = parseOptionalRef(raw, rf.key); err != nil {
			return nil, err
		}
	}
	if f.Operators, err = parseRefList(raw, "operators"); err != nil {
		return nil, err
	}
	if f.MaintainedBy, err = parseRefList(raw, "maintainedBy"); err != nil {
		return nil, err
	}
	if f.MinimalSpareParstCount, err = parseOptional[float64](raw, "minimalSpareParstCount"); err != nil {
		return nil, err
	}
	return f, nil
}

func parsePatchPhysicalItemPayload(raw map[string]json.RawMessage) (*models.PatchPhysicalItemFields, error) {
	if err := rejectUnknownKeys(raw, patchPhysicalItemAllowedKeys); err != nil {
		return nil, err
	}
	f := &models.PatchPhysicalItemFields{}
	var err error
	if f.Token, err = parsePatchToken(raw); err != nil {
		return nil, err
	}
	if f.SerialNumber, err = parseOptional[string](raw, "serialNumber"); err != nil {
		return nil, err
	}
	if f.Notes, err = parseOptional[string](raw, "notes"); err != nil {
		return nil, err
	}
	if f.ItemUsage, err = parseOptionalRef(raw, "itemUsage"); err != nil {
		return nil, err
	}
	if f.ConditionStatus, err = parseOptionalRef(raw, "conditionStatus"); err != nil {
		return nil, err
	}
	return f, nil
}

var createSystemAllowedKeys = map[string]struct{}{
	"parentUid": {}, "name": {}, "description": {}, "systemCode": {}, "systemAlias": {},
	"systemLevel": {}, "systemType": {}, "location": {}, "zone": {}, "importance": {},
	"responsible": {}, "responsibleTeam": {}, "owner": {}, "operators": {}, "maintainedBy": {},
	"minimalSpareParstCount": {},
}

func parseCreateSystemPayload(raw map[string]json.RawMessage) (*models.CreateSystemFields, error) {
	if err := rejectUnknownKeys(raw, createSystemAllowedKeys); err != nil {
		return nil, err
	}
	result := &models.CreateSystemFields{}
	if r, ok := raw["parentUid"]; ok && !rawIsNull(r) {
		var parentUID string
		if err := json.Unmarshal(r, &parentUID); err != nil || parentUID == "" {
			return nil, fmt.Errorf("invalid parentUid")
		}
		result.ParentUID = &parentUID
	}
	delete(raw, "parentUid")
	fields, err := parsePatchSystemPayload(raw)
	if err != nil {
		return nil, err
	}
	result.Fields = *fields
	return result, nil
}
