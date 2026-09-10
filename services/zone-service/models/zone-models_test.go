package models

import (
	"encoding/json"
	"testing"

	"panda/apigateway/helpers"

	"github.com/stretchr/testify/assert"
)

// zoneRecord mirrors the map the defaultParentSystem projection returns for a system created by
// the CS zone migration: those have no systemCode, so `code` comes back as null.
func zoneRecord(defaultParentSystem interface{}) map[string]interface{} {
	return map[string]interface{}{
		"uid":                 "zone-1",
		"name":                "L1",
		"code":                "01",
		"notes":               "",
		"parentZone":          nil,
		"defaultParentSystem": defaultParentSystem,
	}
}

// TestMapZone_DefaultParentSystemNullCode pins the answer to the review comment claiming that
// `code: null` fails to unmarshal into the non-pointer Codebook.Code and 500s the zone reads:
// encoding/json treats null as a no-op for non-pointer targets, so it maps to "" and is then
// dropped by omitempty.
func TestMapZone_DefaultParentSystemNullCode(t *testing.T) {
	zone, err := helpers.MapStruct[Zone](zoneRecord(map[string]interface{}{
		"uid":  "sys-1",
		"name": "CS Zone Default System",
		"code": nil,
	}))

	assert.NoError(t, err)
	assert.NotNil(t, zone.DefaultParentSystem)
	assert.Equal(t, "sys-1", zone.DefaultParentSystem.UID)
	assert.Equal(t, "CS Zone Default System", zone.DefaultParentSystem.Name)
	assert.Equal(t, "", zone.DefaultParentSystem.Code)

	// and the response the client sees carries the system without a code, not code: null
	var payload map[string]interface{}
	body, marshalErr := json.Marshal(zone)
	assert.NoError(t, marshalErr)
	assert.NoError(t, json.Unmarshal(body, &payload))

	dps, ok := payload["defaultParentSystem"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotContains(t, dps, "code")
}

func TestMapZone_DefaultParentSystemPresentCode(t *testing.T) {
	zone, err := helpers.MapStruct[Zone](zoneRecord(map[string]interface{}{
		"uid":  "sys-1",
		"name": "01 - L1 laser system",
		"code": "PLC01-001",
	}))

	assert.NoError(t, err)
	assert.NotNil(t, zone.DefaultParentSystem)
	assert.Equal(t, "PLC01-001", zone.DefaultParentSystem.Code)
}

// A zone with no default parent system projects null, which must map to a nil pointer and be
// omitted from the response - the shape the handler docs and both specs describe.
func TestMapZone_DefaultParentSystemNull(t *testing.T) {
	zone, err := helpers.MapStruct[Zone](zoneRecord(nil))

	assert.NoError(t, err)
	assert.Nil(t, zone.DefaultParentSystem)

	var payload map[string]interface{}
	body, marshalErr := json.Marshal(zone)
	assert.NoError(t, marshalErr)
	assert.NoError(t, json.Unmarshal(body, &payload))
	assert.NotContains(t, payload, "defaultParentSystem")
}
