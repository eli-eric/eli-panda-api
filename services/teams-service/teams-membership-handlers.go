package teamsservice

import (
	"errors"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/teams-service/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// AddTeamMembers Add team members godoc
// @Summary Add team members
// @Description Add one or more users to a team (idempotent). All uids must be users in the facility.
// @Tags Teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "team uid"
// @Param members body models.TeamMembersRequest true "User uids"
// @Success 200 {object} models.TeamDetail
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/teams/{uid}/members [post]
func (h *TeamsHandlers) AddTeamMembers() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")

		req := new(models.TeamMembersRequest)
		if err := c.Bind(req); err != nil {
			return helpers.BadRequest(err.Error())
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		detail, err := h.teamsService.AddMembers(uid, facilityCode, userUID, req.UserUids)
		if err != nil {
			return membershipError(err, "Error adding team members")
		}

		return c.JSON(http.StatusOK, detail)
	}
}

// RemoveTeamMember Remove team member godoc
// @Summary Remove a team member
// @Description Remove a single user from a team. Idempotent (200 even if the user was not a member).
// @Tags Teams
// @Security BearerAuth
// @Produce json
// @Param uid path string true "team uid"
// @Param userUid path string true "user uid"
// @Success 200 {object} models.TeamDetail
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/teams/{uid}/members/{userUid} [delete]
func (h *TeamsHandlers) RemoveTeamMember() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		memberUID := c.Param("userUid")

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		detail, err := h.teamsService.RemoveMember(uid, facilityCode, userUID, memberUID)
		if err != nil {
			return membershipError(err, "Error removing team member")
		}

		return c.JSON(http.StatusOK, detail)
	}
}

// ReplaceTeamMembers Replace team members godoc
// @Summary Replace team members
// @Description Replace the entire member set of a team with the supplied user uids (add + remove diff).
// @Tags Teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "team uid"
// @Param members body models.TeamMembersRequest true "User uids"
// @Success 200 {object} models.TeamDetail
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/teams/{uid}/members [put]
func (h *TeamsHandlers) ReplaceTeamMembers() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")

		req := new(models.TeamMembersRequest)
		if err := c.Bind(req); err != nil {
			return helpers.BadRequest(err.Error())
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		detail, err := h.teamsService.ReplaceMembers(uid, facilityCode, userUID, req.UserUids)
		if err != nil {
			return membershipError(err, "Error replacing team members")
		}

		return c.JSON(http.StatusOK, detail)
	}
}

// membershipError maps service errors to HTTP responses for the membership endpoints.
func membershipError(err error, logMsg string) error {
	if errors.Is(err, ErrNotFound) {
		return echo.ErrNotFound
	}
	if isClientError(err) {
		return helpers.BadRequest(err.Error())
	}
	log.Error().Err(err).Msg(logMsg)
	return echo.ErrInternalServerError
}
