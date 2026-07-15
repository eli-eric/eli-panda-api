package teamsservice

import (
	"errors"
	"panda/apigateway/helpers"
	"panda/apigateway/services/teams-service/models"
	"strings"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var (
	ErrNotFound      = errors.New("team not found")
	ErrDuplicateCode = errors.New("team with this code already exists in this facility")
	ErrValidation    = errors.New("invalid team payload")
)

// ReferencedError is returned when a team cannot be deleted because Systems/RoomCards
// still reference it. It carries per-label counts for the 409 response.
type ReferencedError struct {
	SystemCount   int
	RoomCardCount int
}

func (e *ReferencedError) Error() string {
	return "team is referenced by other entities and cannot be deleted"
}

// InvalidMembersError is returned when a membership request contains uids that are not
// Users in the team's facility. It carries the offending uids for the 400 message.
type InvalidMembersError struct {
	Uids []string
}

func (e *InvalidMembersError) Error() string {
	return "invalid member uids (not users in this facility): " + strings.Join(e.Uids, ", ")
}

type TeamsService struct {
	neo4jDriver *neo4j.Driver
}

type ITeamsService interface {
	GetAllTeams(facilityCode string) ([]models.TeamListItem, error)
	GetAssignableUsers(facilityCode, search string) ([]models.TeamMember, error)
	GetTeamByUID(uid, facilityCode string) (models.TeamDetail, error)
	CreateTeam(facilityCode, userUID string, req *models.TeamCreateRequest) (models.Team, error)
	UpdateTeam(uid, facilityCode, userUID string, req *models.TeamUpdateRequest) (models.Team, error)
	PatchTeam(uid, facilityCode, userUID string, fields *models.PatchTeamFields) (models.Team, error)
	DeleteTeam(uid, facilityCode, userUID string) error
	AddMembers(teamUID, facilityCode, userUID string, userUids []string) (models.TeamDetail, error)
	RemoveMember(teamUID, facilityCode, userUID, memberUID string) (models.TeamDetail, error)
	ReplaceMembers(teamUID, facilityCode, userUID string, userUids []string) (models.TeamDetail, error)
}

func NewTeamsService(driver *neo4j.Driver) ITeamsService {
	return &TeamsService{neo4jDriver: driver}
}

func (svc *TeamsService) GetAllTeams(facilityCode string) (result []models.TeamListItem, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	result, err = helpers.GetNeo4jArrayOfNodes[models.TeamListItem](session, GetAllTeamsQuery(facilityCode))
	helpers.ProcessArrayResult(&result, err)
	return result, err
}

// GetAssignableUsers returns enabled facility users for the team member picker,
// optionally filtered by a case-insensitive search across name/username/email.
func (svc *TeamsService) GetAssignableUsers(facilityCode, search string) (result []models.TeamMember, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	result, err = helpers.GetNeo4jArrayOfNodes[models.TeamMember](session, GetAssignableUsersQuery(facilityCode, search))
	helpers.ProcessArrayResult(&result, err)
	return result, err
}

func (svc *TeamsService) GetTeamByUID(uid, facilityCode string) (result models.TeamDetail, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.TeamDetail](session, GetTeamByUIDQuery(uid, facilityCode))
	if isNoRecords(err) {
		return result, ErrNotFound
	}
	return result, err
}

func (svc *TeamsService) CreateTeam(facilityCode, userUID string, req *models.TeamCreateRequest) (result models.Team, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" {
		return result, ErrValidation
	}

	if err = svc.ensureCodeUnique(session, code, facilityCode, ""); err != nil {
		return result, err
	}

	uid := uuid.New().String()
	query := CreateTeamQuery(uid, name, code, strings.TrimSpace(req.Description), facilityCode, userUID)
	return helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.Team](session, query)
}

func (svc *TeamsService) UpdateTeam(uid, facilityCode, userUID string, req *models.TeamUpdateRequest) (result models.Team, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	if _, err = svc.getTeam(session, uid, facilityCode); err != nil {
		return result, err
	}

	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" {
		return result, ErrValidation
	}

	if err = svc.ensureCodeUnique(session, code, facilityCode, uid); err != nil {
		return result, err
	}

	if err = helpers.WriteNeo4jAndReturnNothing(session, UpdateTeamQuery(uid, name, code, req.Description, facilityCode, userUID)); err != nil {
		return result, err
	}

	return svc.getTeam(session, uid, facilityCode)
}

func (svc *TeamsService) PatchTeam(uid, facilityCode, userUID string, fields *models.PatchTeamFields) (result models.Team, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	current, err := svc.getTeam(session, uid, facilityCode)
	if err != nil {
		return result, err
	}

	changes := []helpers.ChangeEntry{}

	if fields.Name != nil {
		newName := strings.TrimSpace(*fields.Name)
		if newName == "" {
			return result, ErrValidation
		}
		changes = helpers.AppendIfChanged(changes, "name", helpers.ChangeTypeString, current.Name, newName)
	}

	if fields.Code != nil {
		var newCode *string
		if fields.Code.Value != nil {
			trimmed := strings.TrimSpace(*fields.Code.Value)
			newCode = &trimmed
		}
		if newCode != nil && *newCode != "" {
			if err = svc.ensureCodeUnique(session, *newCode, facilityCode, uid); err != nil {
				return result, err
			}
		}
		changes = helpers.AppendIfChanged(changes, "code", helpers.ChangeTypeString, current.Code, optStringValue(newCode))
	}

	if fields.Description != nil {
		var newDesc *string
		if fields.Description.Value != nil {
			newDesc = fields.Description.Value
		}
		changes = helpers.AppendIfChanged(changes, "description", helpers.ChangeTypeString, current.Description, optStringValue(newDesc))
	}

	// nothing supplied or nothing actually changed -> no write, no audit noise
	if len(changes) == 0 {
		return current, nil
	}

	query := PatchTeamQuery(uid, facilityCode, userUID, fields, helpers.MarshalChanges(changes))
	if err = helpers.WriteNeo4jAndReturnNothing(session, query); err != nil {
		return result, err
	}

	return svc.getTeam(session, uid, facilityCode)
}

func (svc *TeamsService) DeleteTeam(uid, facilityCode, userUID string) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	if _, err := svc.getTeam(session, uid, facilityCode); err != nil {
		return err
	}

	refs, err := helpers.GetNeo4jSingleRecordAndMapToStruct[teamReferences](session, CheckTeamReferencesQuery(uid, facilityCode))
	if err != nil {
		return err
	}
	if refs.SystemCount > 0 || refs.RoomCardCount > 0 {
		return &ReferencedError{SystemCount: refs.SystemCount, RoomCardCount: refs.RoomCardCount}
	}

	return helpers.WriteNeo4jAndReturnNothing(session, DeleteTeamQuery(uid, facilityCode))
}

// getTeam fetches the plain team (name/code/description), translating no-rows to ErrNotFound.
func (svc *TeamsService) getTeam(session neo4j.Session, uid, facilityCode string) (models.Team, error) {
	detail, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.TeamDetail](session, GetTeamByUIDQuery(uid, facilityCode))
	if err != nil {
		if isNoRecords(err) {
			return models.Team{}, ErrNotFound
		}
		return models.Team{}, err
	}
	return models.Team{UID: detail.UID, Name: detail.Name, Code: detail.Code, Description: detail.Description}, nil
}

// ensureCodeUnique rejects a duplicate non-empty code within the facility (excluding excludeUID).
func (svc *TeamsService) ensureCodeUnique(session neo4j.Session, code, facilityCode, excludeUID string) error {
	if code == "" {
		return nil
	}
	cnt, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CheckTeamCodeExistsQuery(code, facilityCode, excludeUID))
	if err != nil {
		return err
	}
	if cnt > 0 {
		return ErrDuplicateCode
	}
	return nil
}

type teamReferences struct {
	SystemCount   int `json:"systemCount"`
	RoomCardCount int `json:"roomCardCount"`
}

func isNoRecords(err error) bool {
	return errors.Is(err, helpers.ERR_NO_ROWS)
}

func optStringValue(v *string) interface{} {
	if v == nil {
		return nil
	}
	return *v
}
