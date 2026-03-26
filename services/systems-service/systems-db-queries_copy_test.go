package systemsService

import (
	"regexp"
	"testing"

	"panda/apigateway/services/systems-service/models"

	"github.com/stretchr/testify/assert"
)

func TestCopySystemsQuery_Parameters(t *testing.T) {
	request := &models.SystemCopyRequest{
		SourceSystemUID:              "src-uid",
		DestinationSystemUID:         "dst-uid",
		CopyOnlySourceSystemChildren: true,
		CopyRecursive:                true,
	}
	query := CopySystemsQuery(request, "ELI", "user-123")

	assert.Equal(t, "ELI", query.Parameters["facilityCode"])
	assert.Equal(t, "user-123", query.Parameters["userUID"])
	assert.Equal(t, "src-uid", query.Parameters["sourceUid"])
	assert.Equal(t, "dst-uid", query.Parameters["destinationUid"])
	assert.Equal(t, true, query.Parameters["copyOnlyChildren"])
	assert.Equal(t, true, query.Parameters["copyRecursive"])
	assert.Equal(t, "result", query.ReturnAlias)
}

func TestCopySystemsQuery_QueryStructure(t *testing.T) {
	request := &models.SystemCopyRequest{
		SourceSystemUID:      "src-uid",
		DestinationSystemUID: "dst-uid",
	}
	query := CopySystemsQuery(request, "ELI", "user-123")

	assert.Contains(t, query.Query, "MATCH (source:System{uid:$sourceUid")
	assert.Contains(t, query.Query, "MATCH (destination:System{uid:$destinationUid")
	assert.Contains(t, query.Query, "CASE WHEN $copyRecursive")
	assert.Contains(t, query.Query, "CREATE (c:System{uid:newUid")
	assert.Contains(t, query.Query, "MERGE (copy)-[:HAS_SYSTEM_TYPE]->(st)")
	assert.Contains(t, query.Query, "MERGE (destination)-[:HAS_SUBSYSTEM]->(rootCopy)")
	assert.Contains(t, query.Query, "CASE WHEN $copyOnlyChildren")

	// Destination linking must come BEFORE type/hierarchy subqueries
	destLinkIdx := 0
	for i := range query.Query {
		if i+len("MERGE (destination)") <= len(query.Query) && query.Query[i:i+len("MERGE (destination)")] == "MERGE (destination)" {
			destLinkIdx = i
			break
		}
	}
	typeLinkIdx := 0
	for i := range query.Query {
		if i+len("MERGE (copy)-[:HAS_SYSTEM_TYPE]") <= len(query.Query) && query.Query[i:i+len("MERGE (copy)-[:HAS_SYSTEM_TYPE]")] == "MERGE (copy)-[:HAS_SYSTEM_TYPE]" {
			typeLinkIdx = i
			break
		}
	}
	assert.True(t, destLinkIdx < typeLinkIdx, "destination linking must happen before type copying to avoid blocking on leaf nodes")
}

func TestCopySystemsQuery_IsolatedSubqueries(t *testing.T) {
	request := &models.SystemCopyRequest{
		SourceSystemUID:      "src-uid",
		DestinationSystemUID: "dst-uid",
	}
	query := CopySystemsQuery(request, "ELI", "user-123")

	// Type copying and hierarchy must be in isolated CALL subqueries with count(*) return
	// so they don't block the main flow when no matches are found
	assert.Contains(t, query.Query, "RETURN count(*) AS typesLinked")
	assert.Contains(t, query.Query, "RETURN count(*) AS linksCreated")
}

func TestCopySystemsQuery_NoBrokenWithWhere(t *testing.T) {
	request := &models.SystemCopyRequest{
		SourceSystemUID:      "src-uid",
		DestinationSystemUID: "dst-uid",
	}
	query := CopySystemsQuery(request, "ELI", "user-123")

	// Regression: importing WITH inside CALL must not be followed by WHERE
	re := regexp.MustCompile(`(?i)WITH\s+\w+[\w\s,]*\n\s*WHERE\s+\$copyRecursive`)
	assert.False(t, re.MatchString(query.Query), "query must not use WHERE after importing WITH inside CALL subquery")
}

func TestCopySystemsQuery_RootOldUidSplitWith(t *testing.T) {
	request := &models.SystemCopyRequest{
		SourceSystemUID:      "src-uid",
		DestinationSystemUID: "dst-uid",
	}
	query := CopySystemsQuery(request, "ELI", "user-123")

	// Regression: rootOldUid must be defined in separate WITH before being referenced
	assert.Contains(t, query.Query, "root.uid AS rootOldUid")
	// Must NOT define and use rootOldUid in same WITH clause
	re := regexp.MustCompile(`WITH.*root\.uid AS rootOldUid.*rootOldUid\]`)
	assert.False(t, re.MatchString(query.Query), "rootOldUid must not be defined and referenced in the same WITH clause")
}

func TestValidateSystemCopyRequest_Valid(t *testing.T) {
	request := &models.SystemCopyRequest{
		SourceSystemUID:      "src-uid",
		DestinationSystemUID: "dst-uid",
	}
	err := ValidateSystemCopyRequest(request)
	assert.Nil(t, err)
}

func TestValidateSystemCopyRequest_MissingSource(t *testing.T) {
	request := &models.SystemCopyRequest{
		DestinationSystemUID: "dst-uid",
	}
	err := ValidateSystemCopyRequest(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing source system")
}

func TestValidateSystemCopyRequest_MissingDestination(t *testing.T) {
	request := &models.SystemCopyRequest{
		SourceSystemUID: "src-uid",
	}
	err := ValidateSystemCopyRequest(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing destination parent system")
}

func TestValidateSystemCopyRequest_MissingBoth(t *testing.T) {
	request := &models.SystemCopyRequest{}
	err := ValidateSystemCopyRequest(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing source system")
	assert.Contains(t, err.Error(), "missing destination parent system")
}

func TestValidateSystemCopyRequest_Nil(t *testing.T) {
	err := ValidateSystemCopyRequest(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing request")
}
