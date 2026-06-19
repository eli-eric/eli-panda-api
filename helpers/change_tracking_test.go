package helpers

import (
	"encoding/json"
	"testing"
	"time"

	codebookModels "panda/apigateway/services/codebook-service/models"

	"github.com/stretchr/testify/assert"
)

func TestAppendIfChanged_SameValue_NoAppend(t *testing.T) {
	entries := []ChangeEntry{}
	entries = AppendIfChanged(entries, "name", ChangeTypeString, "foo", "foo")
	assert.Empty(t, entries)
}

func TestAppendIfChanged_DifferentString(t *testing.T) {
	entries := AppendIfChanged(nil, "name", ChangeTypeString, "old", "new")
	assert.Len(t, entries, 1)
	assert.Equal(t, "name", entries[0].Field)
	assert.Equal(t, "string", entries[0].Type)
	assert.Equal(t, "old", entries[0].OldValue)
	assert.Equal(t, "new", entries[0].NewValue)
}

func TestAppendIfChanged_NilOldNonNilNew(t *testing.T) {
	var oldVal *string
	newVal := "new"
	entries := AppendIfChanged(nil, "description", ChangeTypeString, oldVal, &newVal)
	assert.Len(t, entries, 1)
	assert.Nil(t, entries[0].OldValue)
	assert.NotNil(t, entries[0].NewValue)
}

func TestAppendIfChanged_NonNilOldNilNew(t *testing.T) {
	oldVal := "old"
	var newVal *string
	entries := AppendIfChanged(nil, "description", ChangeTypeString, &oldVal, newVal)
	assert.Len(t, entries, 1)
	assert.NotNil(t, entries[0].OldValue)
	assert.Nil(t, entries[0].NewValue)
}

func TestAppendIfChanged_BothNil_NoAppend(t *testing.T) {
	var oldVal, newVal *string
	entries := AppendIfChanged(nil, "description", ChangeTypeString, oldVal, newVal)
	assert.Empty(t, entries)
}

func TestAppendIfChanged_Codebook_SameUID_DifferentName_NoAppend(t *testing.T) {
	oldVal := &codebookModels.Codebook{UID: "abc-123", Name: "Old Name"}
	newVal := &codebookModels.Codebook{UID: "abc-123", Name: "New Name"}
	entries := AppendIfChanged(nil, "supplier", ChangeTypeCodebook, oldVal, newVal)
	assert.Empty(t, entries, "same UID should not register as a change even when Name differs")
}

func TestAppendIfChanged_Codebook_DifferentUID(t *testing.T) {
	oldVal := &codebookModels.Codebook{UID: "abc-123", Name: "Supplier A"}
	newVal := &codebookModels.Codebook{UID: "xyz-789", Name: "Supplier B"}
	entries := AppendIfChanged(nil, "supplier", ChangeTypeCodebook, oldVal, newVal)
	assert.Len(t, entries, 1)
	assert.Equal(t, "codebook", entries[0].Type)
}

func TestAppendIfChanged_Codebook_NilToValue(t *testing.T) {
	var oldVal *codebookModels.Codebook
	newVal := &codebookModels.Codebook{UID: "abc", Name: "X"}
	entries := AppendIfChanged(nil, "supplier", ChangeTypeCodebook, oldVal, newVal)
	assert.Len(t, entries, 1)
	assert.Nil(t, entries[0].OldValue)
	assert.NotNil(t, entries[0].NewValue)
}

func TestAppendIfChanged_AllChangeTypes(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Hour)

	cases := []struct {
		name     string
		typ      ChangeType
		old, new interface{}
	}{
		{"string", ChangeTypeString, "a", "b"},
		{"number", ChangeTypeNumber, float64(1), float64(2)},
		{"boolean", ChangeTypeBoolean, true, false},
		{"date", ChangeTypeDate, now, later},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			entries := AppendIfChanged(nil, "f", c.typ, c.old, c.new)
			assert.Len(t, entries, 1)
			assert.Equal(t, string(c.typ), entries[0].Type)
		})
	}
}

func TestMarshalChanges_Empty(t *testing.T) {
	assert.Equal(t, "[]", MarshalChanges(nil))
	assert.Equal(t, "[]", MarshalChanges([]ChangeEntry{}))
}

func TestMarshalChanges_ProducesExpectedShape(t *testing.T) {
	entries := []ChangeEntry{
		{Field: "name", Type: "string", OldValue: "old", NewValue: "new"},
		{Field: "supplier", Type: "codebook", OldValue: nil, NewValue: map[string]string{"uid": "u", "name": "n"}},
	}
	out := MarshalChanges(entries)

	var parsed []map[string]interface{}
	err := json.Unmarshal([]byte(out), &parsed)
	assert.NoError(t, err)
	assert.Len(t, parsed, 2)
	assert.Equal(t, "name", parsed[0]["field"])
	assert.Equal(t, "string", parsed[0]["type"])
	assert.Equal(t, "old", parsed[0]["oldValue"])
	assert.Equal(t, "new", parsed[0]["newValue"])
	assert.Nil(t, parsed[1]["oldValue"])
}
