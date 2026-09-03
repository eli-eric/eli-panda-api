package systemsService

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"panda/apigateway/helpers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// systemCodesClientError decides whether a system code endpoint answers 400 with a readable
// message or falls through to a 500, so both directions are pinned down here.
func TestSystemCodesClientError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantMessage string // empty = must not be treated as a client error
	}{
		{
			name:        "missing default parent system",
			err:         ErrMissingDefaultParentSystem,
			wantMessage: "Bad request: missing default parent system for selected zone (set it on the zone via PUT /v1/zones/{uid})",
		},
		{
			name:        "wrapped invalid input keeps only the readable part",
			err:         invalidSystemCodesInput("batch must be greater than zero"),
			wantMessage: "Bad request: batch must be greater than zero",
		},
		{
			name:        "zone not found",
			err:         invalidSystemCodesInput("zone not found"),
			wantMessage: "Bad request: zone not found",
		},
		{
			name:        "wrapped sentinel is still recognised",
			err:         fmt.Errorf("save failed: %w", ErrMissingDefaultParentSystem),
			wantMessage: "Bad request: missing default parent system for selected zone (set it on the zone via PUT /v1/zones/{uid})",
		},
		{
			name:        "database error is not a client error",
			err:         errors.New("connection refused"),
			wantMessage: "",
		},
		{
			name:        "no rows is not a client error on its own",
			err:         helpers.ERR_NO_ROWS,
			wantMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := systemCodesClientError(tt.err)

			if tt.wantMessage == "" {
				assert.Nil(t, result)
				return
			}

			httpErr, ok := result.(*echo.HTTPError)
			assert.True(t, ok, "expected an echo.HTTPError")
			assert.Equal(t, http.StatusBadRequest, httpErr.Code)
			assert.Equal(t, tt.wantMessage, httpErr.Message)
		})
	}
}

// The zone match must stay mandatory and the parent match optional: an unknown zone then returns no
// rows and is reported as "zone not found" instead of "missing default parent system".
func TestCheckZoneHasDefaultParentSystemQuery(t *testing.T) {
	query := CheckZoneHasDefaultParentSystemQuery("zone-1", "B")

	assert.Equal(t, "cnt", query.ReturnAlias)
	assert.Contains(t, query.Query, "MATCH(z:Zone{uid: $zoneUID})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode})")
	assert.Contains(t, query.Query, "OPTIONAL MATCH(z)-[:HAS_DEFAULT_PARENT_SYSTEM]->(parent:System)-[:BELONGS_TO_FACILITY]->(f)")
	assert.Contains(t, query.Query, "coalesce(parent.deleted, false) = false")
	// z as grouping key: without it a bare count() returns a 0 row even for an unknown zone
	assert.Contains(t, query.Query, "WITH z, count(parent) as cnt")
	assert.Equal(t, "zone-1", query.Parameters["zoneUID"])
	assert.Equal(t, "B", query.Parameters["facilityCode"])
}
