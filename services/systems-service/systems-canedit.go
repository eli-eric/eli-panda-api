package systemsService

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"
)

// CanEditSystem is the single source of truth for the system edit-guard. It bubbles up the
// HAS_SUBSYSTEM chain and returns whether userUID may edit systemUID plus the responsible users.
// A non-existent/deleted system yields helpers.ERR_NO_ROWS (callers decide 404 vs. pass-through).
func (svc *SystemsService) CanEditSystem(systemUID, userUID string) (models.CanEditSystemResult, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	return helpers.GetNeo4jSingleRecordAndMapToStruct[models.CanEditSystemResult](session, CanEditSystemQuery(systemUID, userUID))
}
