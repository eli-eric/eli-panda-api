package helpers

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type projectionTestStruct struct {
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Location *struct {
		UID string `json:"uid"`
	} `json:"location"`
	Ignored string `json:"-"`
	// unexported should not be selectable
	internal string `json:"internal"`
}

func TestParseFieldsParam(t *testing.T) {
	assert.Nil(t, ParseFieldsParam(""))
	assert.Nil(t, ParseFieldsParam("   "))
	assert.Equal(t, []string{"uid", "name"}, ParseFieldsParam("uid,name"))
	assert.Equal(t, []string{"uid", "name"}, ParseFieldsParam(" uid , name "))
	assert.Equal(t, []string{"uid", "name"}, ParseFieldsParam("uid,uid,name,name"))
	assert.Equal(t, []string{"uid", "name"}, ParseFieldsParam(",uid,,name,"))
}

func TestProjectJSONFields_Success(t *testing.T) {
	s := projectionTestStruct{
		UID:  "U1",
		Name: "Room A",
	}

	projected, err := ProjectJSONFields(s, []string{"uid", "name"})
	assert.NoError(t, err)
	assert.Equal(t, "U1", projected["uid"])
	assert.Equal(t, "Room A", projected["name"])
	assert.Len(t, projected, 2)

	// ensure JSON marshaling is stable
	b, marshalErr := json.Marshal(projected)
	assert.NoError(t, marshalErr)
	assert.Contains(t, string(b), "\"uid\"")
	assert.Contains(t, string(b), "\"name\"")
}

func TestProjectJSONFields_NilPointerField_EncodesToNull(t *testing.T) {
	s := projectionTestStruct{
		UID:      "U1",
		Name:     "Room A",
		Location: nil,
	}

	projected, err := ProjectJSONFields(&s, []string{"location"})
	assert.NoError(t, err)

	v := reflect.ValueOf(projected["location"])
	assert.Equal(t, reflect.Pointer, v.Kind())
	assert.True(t, v.IsNil())

	b, marshalErr := json.Marshal(projected)
	assert.NoError(t, marshalErr)
	assert.Contains(t, string(b), "\"location\":null")
}

func TestProjectJSONFields_InvalidFields(t *testing.T) {
	s := projectionTestStruct{UID: "U1", Name: "Room A"}

	projected, err := ProjectJSONFields(s, []string{"uid", "doesNotExist", "location.uid", "-"})
	assert.Nil(t, projected)
	assert.Error(t, err)

	projErr, ok := err.(*FieldProjectionError)
	assert.True(t, ok)
	assert.Equal(t, []string{"-", "doesNotExist", "location.uid"}, projErr.InvalidFields)

	// allowed fields should include exported+tagged fields only (no '-' and no unexported)
	assert.Contains(t, projErr.AllowedFields, "uid")
	assert.Contains(t, projErr.AllowedFields, "name")
	assert.Contains(t, projErr.AllowedFields, "location")
	assert.NotContains(t, projErr.AllowedFields, "internal")
}

func TestProjectJSONFields_RejectsNonStruct(t *testing.T) {
	_, err := ProjectJSONFields("not-a-struct", []string{"uid"})
	assert.Error(t, err)
}
