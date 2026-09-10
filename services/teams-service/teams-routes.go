package teamsservice

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapTeamsRoutes(e *echo.Echo, h ITeamsHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// team CRUD
	e.GET("/v1/teams", m.Authorization(h.GetAllTeams(), shared.ROLE_TEAMS_VIEW, shared.ROLE_ADMIN), jwtMiddleware)
	e.GET("/v1/teams/:uid", m.Authorization(h.GetTeamByUID(), shared.ROLE_TEAMS_VIEW, shared.ROLE_ADMIN), jwtMiddleware)
	e.POST("/v1/teams", m.Authorization(h.CreateTeam(), shared.ROLE_TEAMS_EDIT, shared.ROLE_ADMIN), jwtMiddleware)
	e.PUT("/v1/teams/:uid", m.Authorization(h.UpdateTeam(), shared.ROLE_TEAMS_EDIT, shared.ROLE_ADMIN), jwtMiddleware)
	e.PATCH("/v1/teams/:uid", m.Authorization(h.PatchTeam(), shared.ROLE_TEAMS_EDIT, shared.ROLE_ADMIN), jwtMiddleware)
	e.DELETE("/v1/teams/:uid", m.Authorization(h.DeleteTeam(), shared.ROLE_TEAMS_EDIT, shared.ROLE_ADMIN), jwtMiddleware)

	// assignable facility users for the member picker
	e.GET("/v1/teams/assignable-users", m.Authorization(h.GetAssignableUsers(), shared.ROLE_TEAMS_VIEW, shared.ROLE_TEAMS_EDIT, shared.ROLE_ADMIN), jwtMiddleware)

	// team membership
	e.POST("/v1/teams/:uid/members", m.Authorization(h.AddTeamMembers(), shared.ROLE_TEAMS_EDIT, shared.ROLE_ADMIN), jwtMiddleware)
	e.PUT("/v1/teams/:uid/members", m.Authorization(h.ReplaceTeamMembers(), shared.ROLE_TEAMS_EDIT, shared.ROLE_ADMIN), jwtMiddleware)
	e.DELETE("/v1/teams/:uid/members/:userUid", m.Authorization(h.RemoveTeamMember(), shared.ROLE_TEAMS_EDIT, shared.ROLE_ADMIN), jwtMiddleware)
}
