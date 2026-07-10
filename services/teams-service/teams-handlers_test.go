package teamsservice

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"panda/apigateway/services/teams-service/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// teamsServiceMock is a hand-written ITeamsService stub; each method delegates to an
// optional func field so individual tests control behavior and capture arguments.
type teamsServiceMock struct {
	getAllFn   func(facilityCode string) ([]models.TeamListItem, error)
	getByUIDFn func(uid, facilityCode string) (models.TeamDetail, error)
	createFn   func(facilityCode, userUID string, req *models.TeamCreateRequest) (models.Team, error)
	updateFn   func(uid, facilityCode, userUID string, req *models.TeamUpdateRequest) (models.Team, error)
	patchFn    func(uid, facilityCode, userUID string, fields *models.PatchTeamFields) (models.Team, error)
	deleteFn   func(uid, facilityCode, userUID string) error
	addFn      func(teamUID, facilityCode, userUID string, userUids []string) (models.TeamDetail, error)
	removeFn   func(teamUID, facilityCode, userUID, memberUID string) (models.TeamDetail, error)
	replaceFn  func(teamUID, facilityCode, userUID string, userUids []string) (models.TeamDetail, error)
}

func (m *teamsServiceMock) GetAllTeams(facilityCode string) ([]models.TeamListItem, error) {
	if m.getAllFn != nil {
		return m.getAllFn(facilityCode)
	}
	return []models.TeamListItem{}, nil
}

func (m *teamsServiceMock) GetTeamByUID(uid, facilityCode string) (models.TeamDetail, error) {
	if m.getByUIDFn != nil {
		return m.getByUIDFn(uid, facilityCode)
	}
	return models.TeamDetail{}, nil
}

func (m *teamsServiceMock) CreateTeam(facilityCode, userUID string, req *models.TeamCreateRequest) (models.Team, error) {
	if m.createFn != nil {
		return m.createFn(facilityCode, userUID, req)
	}
	return models.Team{}, nil
}

func (m *teamsServiceMock) UpdateTeam(uid, facilityCode, userUID string, req *models.TeamUpdateRequest) (models.Team, error) {
	if m.updateFn != nil {
		return m.updateFn(uid, facilityCode, userUID, req)
	}
	return models.Team{}, nil
}

func (m *teamsServiceMock) PatchTeam(uid, facilityCode, userUID string, fields *models.PatchTeamFields) (models.Team, error) {
	if m.patchFn != nil {
		return m.patchFn(uid, facilityCode, userUID, fields)
	}
	return models.Team{}, nil
}

func (m *teamsServiceMock) DeleteTeam(uid, facilityCode, userUID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(uid, facilityCode, userUID)
	}
	return nil
}

func (m *teamsServiceMock) AddMembers(teamUID, facilityCode, userUID string, userUids []string) (models.TeamDetail, error) {
	if m.addFn != nil {
		return m.addFn(teamUID, facilityCode, userUID, userUids)
	}
	return models.TeamDetail{}, nil
}

func (m *teamsServiceMock) RemoveMember(teamUID, facilityCode, userUID, memberUID string) (models.TeamDetail, error) {
	if m.removeFn != nil {
		return m.removeFn(teamUID, facilityCode, userUID, memberUID)
	}
	return models.TeamDetail{}, nil
}

func (m *teamsServiceMock) ReplaceMembers(teamUID, facilityCode, userUID string, userUids []string) (models.TeamDetail, error) {
	if m.replaceFn != nil {
		return m.replaceFn(teamUID, facilityCode, userUID, userUids)
	}
	return models.TeamDetail{}, nil
}

// newContext builds an echo context with the facility/user values the Authorization
// middleware would normally inject, plus an optional JSON body.
func newContext(method, target, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("facilityCode", "B")
	c.Set("userUID", "actor-uid")
	return c, rec
}

func assertHTTPErrorCode(t *testing.T, err error, code int) {
	t.Helper()
	var he *echo.HTTPError
	if assert.True(t, errors.As(err, &he), "expected an echo.HTTPError, got %v", err) {
		assert.Equal(t, code, he.Code)
	}
}

func TestCreateTeam_Success(t *testing.T) {
	c, rec := newContext(http.MethodPost, "/v1/teams", `{"name":"Vacuum","code":"VAC"}`)
	svc := &teamsServiceMock{
		createFn: func(fc, uUID string, req *models.TeamCreateRequest) (models.Team, error) {
			assert.Equal(t, "B", fc)
			assert.Equal(t, "actor-uid", uUID)
			assert.Equal(t, "Vacuum", req.Name)
			return models.Team{UID: "team-1", Name: req.Name, Code: req.Code}, nil
		},
	}
	h := NewTeamsHandlers(svc)

	err := h.CreateTeam()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"uid":"team-1","name":"Vacuum","code":"VAC","description":""}`, rec.Body.String())
}

func TestCreateTeam_DuplicateCode_BadRequest(t *testing.T) {
	c, _ := newContext(http.MethodPost, "/v1/teams", `{"name":"Vacuum","code":"VAC"}`)
	svc := &teamsServiceMock{
		createFn: func(fc, uUID string, req *models.TeamCreateRequest) (models.Team, error) {
			return models.Team{}, ErrDuplicateCode
		},
	}
	h := NewTeamsHandlers(svc)

	err := h.CreateTeam()(c)

	assertHTTPErrorCode(t, err, http.StatusBadRequest)
}

func TestGetTeamByUID_NotFound(t *testing.T) {
	c, _ := newContext(http.MethodGet, "/v1/teams/missing", "")
	c.SetParamNames("uid")
	c.SetParamValues("missing")
	svc := &teamsServiceMock{
		getByUIDFn: func(uid, fc string) (models.TeamDetail, error) {
			return models.TeamDetail{}, ErrNotFound
		},
	}
	h := NewTeamsHandlers(svc)

	err := h.GetTeamByUID()(c)

	assertHTTPErrorCode(t, err, http.StatusNotFound)
}

func TestDeleteTeam_Referenced_Conflict(t *testing.T) {
	c, rec := newContext(http.MethodDelete, "/v1/teams/team-1", "")
	c.SetParamNames("uid")
	c.SetParamValues("team-1")
	svc := &teamsServiceMock{
		deleteFn: func(uid, fc, uUID string) error {
			return &ReferencedError{SystemCount: 3, RoomCardCount: 1}
		},
	}
	h := NewTeamsHandlers(svc)

	err := h.DeleteTeam()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), `"label":"System"`)
	assert.Contains(t, rec.Body.String(), `"count":3`)
	assert.Contains(t, rec.Body.String(), `"label":"RoomCard"`)
}

func TestDeleteTeam_Success(t *testing.T) {
	c, rec := newContext(http.MethodDelete, "/v1/teams/team-1", "")
	c.SetParamNames("uid")
	c.SetParamValues("team-1")
	h := NewTeamsHandlers(&teamsServiceMock{})

	err := h.DeleteTeam()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRemoveTeamMember_NonMember_Idempotent200(t *testing.T) {
	c, rec := newContext(http.MethodDelete, "/v1/teams/team-1/members/user-x", "")
	c.SetParamNames("uid", "userUid")
	c.SetParamValues("team-1", "user-x")
	svc := &teamsServiceMock{
		removeFn: func(teamUID, fc, uUID, memberUID string) (models.TeamDetail, error) {
			// user not a member -> service returns the unchanged detail, no error
			return models.TeamDetail{UID: "team-1", Name: "Vacuum", Members: []models.TeamMember{}}, nil
		},
	}
	h := NewTeamsHandlers(svc)

	err := h.RemoveTeamMember()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRemoveTeamMember_TeamMissing_404(t *testing.T) {
	c, _ := newContext(http.MethodDelete, "/v1/teams/missing/members/user-x", "")
	c.SetParamNames("uid", "userUid")
	c.SetParamValues("missing", "user-x")
	svc := &teamsServiceMock{
		removeFn: func(teamUID, fc, uUID, memberUID string) (models.TeamDetail, error) {
			return models.TeamDetail{}, ErrNotFound
		},
	}
	h := NewTeamsHandlers(svc)

	err := h.RemoveTeamMember()(c)

	assertHTTPErrorCode(t, err, http.StatusNotFound)
}

func TestAddTeamMembers_InvalidUids_BadRequest(t *testing.T) {
	c, _ := newContext(http.MethodPost, "/v1/teams/team-1/members", `{"userUids":["good","bad"]}`)
	c.SetParamNames("uid")
	c.SetParamValues("team-1")
	svc := &teamsServiceMock{
		addFn: func(teamUID, fc, uUID string, userUids []string) (models.TeamDetail, error) {
			return models.TeamDetail{}, &InvalidMembersError{Uids: []string{"bad"}}
		},
	}
	h := NewTeamsHandlers(svc)

	err := h.AddTeamMembers()(c)

	assertHTTPErrorCode(t, err, http.StatusBadRequest)
}

func TestPatchTeam_ExplicitNullClearsDescription(t *testing.T) {
	c, rec := newContext(http.MethodPatch, "/v1/teams/team-1", `{"description":null}`)
	c.SetParamNames("uid")
	c.SetParamValues("team-1")

	var captured *models.PatchTeamFields
	svc := &teamsServiceMock{
		patchFn: func(uid, fc, uUID string, fields *models.PatchTeamFields) (models.Team, error) {
			captured = fields
			return models.Team{UID: "team-1", Name: "Vacuum"}, nil
		},
	}
	h := NewTeamsHandlers(svc)

	err := h.PatchTeam()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// name absent -> nil; description present but explicit null -> Optional set with Value nil
	assert.Nil(t, captured.Name)
	assert.NotNil(t, captured.Description)
	assert.Nil(t, captured.Description.Value)
	assert.Nil(t, captured.Code)
}

func TestPatchTeam_SetsNameAndDescription(t *testing.T) {
	c, rec := newContext(http.MethodPatch, "/v1/teams/team-1", `{"name":"Cryo","description":"cryogenics team"}`)
	c.SetParamNames("uid")
	c.SetParamValues("team-1")

	var captured *models.PatchTeamFields
	svc := &teamsServiceMock{
		patchFn: func(uid, fc, uUID string, fields *models.PatchTeamFields) (models.Team, error) {
			captured = fields
			return models.Team{UID: "team-1", Name: "Cryo", Description: "cryogenics team"}, nil
		},
	}
	h := NewTeamsHandlers(svc)

	err := h.PatchTeam()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotNil(t, captured.Name)
	assert.Equal(t, "Cryo", *captured.Name)
	assert.NotNil(t, captured.Description)
	assert.NotNil(t, captured.Description.Value)
	assert.Equal(t, "cryogenics team", *captured.Description.Value)
}
