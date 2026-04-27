package helpers

import (
	"encoding/json"
	"reflect"
)

type ChangeType string

const (
	ChangeTypeString   ChangeType = "string"
	ChangeTypeNumber   ChangeType = "number"
	ChangeTypeBoolean  ChangeType = "boolean"
	ChangeTypeDate     ChangeType = "date"
	ChangeTypeCodebook ChangeType = "codebook"
)

// ChangeEntity identifies the entity whose attribute changed. Optional — top-level
// scalar changes (e.g. category.name) typically include it pointing at the parent;
// nested changes (e.g. property under category) point at the nested entity so the FE
// can render "<entity.name>: <field>" without parsing the field string.
//
// Name is captured at change time (snapshot semantics) — if the entity is later renamed,
// historical audit entries still display the name as it was at the moment of change.
type ChangeEntity struct {
	Type string `json:"type"` // e.g. "category", "group", "property", "physicalProperty", "item"
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type ChangeEntry struct {
	Field    string        `json:"field"`
	Type     string        `json:"type"`
	OldValue interface{}   `json:"oldValue"`
	NewValue interface{}   `json:"newValue"`
	Entity   *ChangeEntity `json:"entity,omitempty"`
}

// AppendIfChanged appends a ChangeEntry to entries when oldVal and newVal differ.
// Codebook values are compared by UID field; everything else via reflect.DeepEqual.
// Typed-nil pointers are normalized to untyped nil in the serialized output.
func AppendIfChanged(entries []ChangeEntry, field string, t ChangeType, oldVal, newVal interface{}) []ChangeEntry {
	return AppendIfChangedFor(entries, nil, field, t, oldVal, newVal)
}

// AppendIfChangedFor is the entity-aware variant. Pass nil entity for top-level scalar
// changes; pass a ChangeEntity describing the nested target (group, property, physical
// property) so the audit row carries the entity identity for FE display.
func AppendIfChangedFor(entries []ChangeEntry, entity *ChangeEntity, field string, t ChangeType, oldVal, newVal interface{}) []ChangeEntry {
	if valuesEqual(t, oldVal, newVal) {
		return entries
	}
	return append(entries, ChangeEntry{
		Field:    field,
		Type:     string(t),
		OldValue: normalizeNil(oldVal),
		NewValue: normalizeNil(newVal),
		Entity:   entity,
	})
}

// MarshalChanges serializes a ChangeEntry slice to the JSON string format
// stored on WAS_UPDATED_BY.changes. An empty/nil slice yields "[]".
func MarshalChanges(entries []ChangeEntry) string {
	if entries == nil {
		return "[]"
	}
	b, err := json.Marshal(entries)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func valuesEqual(t ChangeType, a, b interface{}) bool {
	aNil, bNil := isNil(a), isNil(b)
	if aNil && bNil {
		return true
	}
	if aNil != bNil {
		return false
	}
	if t == ChangeTypeCodebook {
		return codebookUID(a) == codebookUID(b)
	}
	return reflect.DeepEqual(a, b)
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	}
	return false
}

func normalizeNil(v interface{}) interface{} {
	if isNil(v) {
		return nil
	}
	return v
}

func codebookUID(v interface{}) string {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if !rv.IsValid() || rv.Kind() != reflect.Struct {
		return ""
	}
	f := rv.FieldByName("UID")
	if !f.IsValid() || f.Kind() != reflect.String {
		return ""
	}
	return f.String()
}
