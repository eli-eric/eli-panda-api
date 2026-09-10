package teamsservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/teams-service/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type TeamsHandlers struct {
	teamsService ITeamsService
}

type ITeamsHandlers interface {
	GetAllTeams() echo.HandlerFunc
	GetAssignableUsers() echo.HandlerFunc
	GetTeamByUID() echo.HandlerFunc
	CreateTeam() echo.HandlerFunc
	UpdateTeam() echo.HandlerFunc
	PatchTeam() echo.HandlerFunc
	DeleteTeam() echo.HandlerFunc
	AddTeamMembers() echo.HandlerFunc
	RemoveTeamMember() echo.HandlerFunc
	ReplaceTeamMembers() echo.HandlerFunc
}

func NewTeamsHandlers(svc ITeamsService) ITeamsHandlers {
	return &TeamsHandlers{teamsService: svc}
}

// GetAllTeams Get all teams godoc
// @Summary Get all teams
// @Description Get all teams for the current facility (flat list with member counts)
// @Tags Teams
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.TeamListItem
// @Failure 500 "Internal Server Error"
// @Router /v1/teams [get]
func (h *TeamsHandlers) GetAllTeams() echo.HandlerFunc {
	return func(c echo.Context) error {
		facilityCode := c.Get("facilityCode").(string)

		teams, err := h.teamsService.GetAllTeams(facilityCode)
		if err != nil {
			log.Error().Err(err).Msg("Error getting teams")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, teams)
	}
}

// GetAssignableUsers Get assignable facility users godoc
// @Summary Get assignable facility users
// @Description Get enabled users of the current facility for the team member picker. Optional case-insensitive search across name/username/email.
// @Tags Teams
// @Security BearerAuth
// @Produce json
// @Param search query string false "Case-insensitive substring filter across firstName/lastName/username/email"
// @Success 200 {array} models.TeamMember
// @Failure 500 "Internal Server Error"
// @Router /v1/teams/assignable-users [get]
func (h *TeamsHandlers) GetAssignableUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		facilityCode := c.Get("facilityCode").(string)
		search := c.QueryParam("search")

		users, err := h.teamsService.GetAssignableUsers(facilityCode, search)
		if err != nil {
			log.Error().Err(err).Msg("Error getting assignable users")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, users)
	}
}

// GetTeamByUID Get team by uid godoc
// @Summary Get team by uid
// @Description Get a team with its full member list
// @Tags Teams
// @Security BearerAuth
// @Produce json
// @Param uid path string true "uid"
// @Success 200 {object} models.TeamDetail
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/teams/{uid} [get]
func (h *TeamsHandlers) GetTeamByUID() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		facilityCode := c.Get("facilityCode").(string)

		team, err := h.teamsService.GetTeamByUID(uid, facilityCode)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return echo.ErrNotFound
			}
			log.Error().Err(err).Msg("Error getting team")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, team)
	}
}

// CreateTeam Create team godoc
// @Summary Create team
// @Description Create a new team (name required; code optional but unique per facility)
// @Tags Teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param team body models.TeamCreateRequest true "Team"
// @Success 201 {object} models.Team
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /v1/teams [post]
func (h *TeamsHandlers) CreateTeam() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(models.TeamCreateRequest)
		if err := c.Bind(req); err != nil {
			return helpers.BadRequest(err.Error())
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		team, err := h.teamsService.CreateTeam(facilityCode, userUID, req)
		if err != nil {
			if isClientError(err) {
				return helpers.BadRequest(err.Error())
			}
			log.Error().Err(err).Msg("Error creating team")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusCreated, team)
	}
}

// UpdateTeam Update team godoc
// @Summary Update team
// @Description Full replace of team name/code/description
// @Tags Teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "uid"
// @Param team body models.TeamUpdateRequest true "Team"
// @Success 200 {object} models.Team
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/teams/{uid} [put]
func (h *TeamsHandlers) UpdateTeam() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")

		req := new(models.TeamUpdateRequest)
		if err := c.Bind(req); err != nil {
			return helpers.BadRequest(err.Error())
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		team, err := h.teamsService.UpdateTeam(uid, facilityCode, userUID, req)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return echo.ErrNotFound
			}
			if isClientError(err) {
				return helpers.BadRequest(err.Error())
			}
			log.Error().Err(err).Msg("Error updating team")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, team)
	}
}

// PatchTeam Patch team godoc
// @Summary Patch team
// @Description Partial update of team params (absent key = unchanged; explicit null = clear)
// @Tags Teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "uid"
// @Param team body object true "Partial team fields"
// @Success 200 {object} models.Team
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/teams/{uid} [patch]
func (h *TeamsHandlers) PatchTeam() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return helpers.BadRequest("cannot read request body")
		}
		raw := map[string]json.RawMessage{}
		if err := json.Unmarshal(body, &raw); err != nil {
			return helpers.BadRequest("invalid JSON body")
		}
		fields, err := parsePatchTeamPayload(raw)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		team, err := h.teamsService.PatchTeam(uid, facilityCode, userUID, fields)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return echo.ErrNotFound
			}
			if isClientError(err) {
				return helpers.BadRequest(err.Error())
			}
			log.Error().Err(err).Msg("Error patching team")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, team)
	}
}

// DeleteTeam Delete team godoc
// @Summary Delete team
// @Description Hard delete a team. Rejected with 409 if Systems or Room Cards still reference it.
// @Tags Teams
// @Security BearerAuth
// @Param uid path string true "uid"
// @Success 200 "OK"
// @Failure 404 "Not Found"
// @Failure 409 "Conflict"
// @Failure 500 "Internal Server Error"
// @Router /v1/teams/{uid} [delete]
func (h *TeamsHandlers) DeleteTeam() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		err := h.teamsService.DeleteTeam(uid, facilityCode, userUID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return echo.ErrNotFound
			}
			var refErr *ReferencedError
			if errors.As(err, &refErr) {
				return c.JSON(http.StatusConflict, helpers.ConflictErrorResponse{
					ErrorMessage: refErr.Error(),
					RelatedNodes: referencedRelatedNodes(refErr),
				})
			}
			log.Error().Err(err).Msg("Error deleting team")
			return echo.ErrInternalServerError
		}

		return c.NoContent(http.StatusOK)
	}
}

// parsePatchTeamPayload maps a raw JSON key map into typed PatchTeamFields. Absent keys stay
// nil; explicit JSON null is captured via Optional[T] with Value=nil (clear operation).
func parsePatchTeamPayload(raw map[string]json.RawMessage) (*models.PatchTeamFields, error) {
	fields := &models.PatchTeamFields{}

	if r, ok := raw["name"]; ok {
		var v string
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, errors.New("invalid name")
		}
		fields.Name = &v
	}

	parseOpt := func(key string) (*models.Optional[string], error) {
		r, ok := raw[key]
		if !ok {
			return nil, nil
		}
		if rawMessageIsNull(r) {
			return &models.Optional[string]{Value: nil}, nil
		}
		var v string
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, errors.New("invalid " + key)
		}
		return &models.Optional[string]{Value: &v}, nil
	}

	code, err := parseOpt("code")
	if err != nil {
		return nil, err
	}
	fields.Code = code

	desc, err := parseOpt("description")
	if err != nil {
		return nil, err
	}
	fields.Description = desc

	return fields, nil
}

func rawMessageIsNull(r json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(r), []byte("null"))
}

func referencedRelatedNodes(e *ReferencedError) []helpers.RelatedNodeLabelAmount {
	nodes := []helpers.RelatedNodeLabelAmount{}
	if e.SystemCount > 0 {
		nodes = append(nodes, helpers.RelatedNodeLabelAmount{Label: "System", Count: e.SystemCount})
	}
	if e.RoomCardCount > 0 {
		nodes = append(nodes, helpers.RelatedNodeLabelAmount{Label: "RoomCard", Count: e.RoomCardCount})
	}
	return nodes
}

func isClientError(err error) bool {
	if errors.Is(err, ErrDuplicateCode) || errors.Is(err, ErrValidation) {
		return true
	}
	var invalidMembers *InvalidMembersError
	return errors.As(err, &invalidMembers)
}
