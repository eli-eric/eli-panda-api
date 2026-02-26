package securityService

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/services/security-service/models"
	"panda/apigateway/shared"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func MapSecurityRoutes(e *echo.Echo, h ISecurityHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// Login route
	e.POST("/v1/authenticate", h.AuthenticateByUsernameAndPassword())

	e.GET("/v1/getuser", func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		userInfo := models.UserAuthInfo{
			Uid:          claims.Subject,
			Username:     "",
			Email:        "",
			LastName:     "",
			FirstName:    "",
			Facility:     "",
			FacilityCode: "",
			AccessToken:  "",
			Roles:        claims.Roles,
			PasswordHash: "",
		}

		return c.JSON(200, userInfo)
	}, jwtMiddleware)

	// This endpoint is intended for internal use by the API Gateway to retrieve user information based on an Azure ID token.
	e.GET("/v1/getuserbyazureidtoken", h.GetUserByAzureIdToken())

	// NOTE: User status cache is currently process-local. In multi-instance deployments,
	// each instance maintains its own cache and may serve stale status information until
	// the local cache entry expires (default ~60 seconds, depending on configuration).
	e.GET("/v1/auth/cache", m.Authorization(h.GetUserStatusCache(), shared.ROLE_ADMIN), jwtMiddleware)

	// WARNING: Invalidating the user status cache via this endpoint only affects the
	// instance handling the request. Other instances will not be notified and will
	// continue to use their local cache entries until they expire.
	e.POST("/v1/auth/cache/invalidate/:userUID", m.Authorization(h.InvalidateUserStatusCache(), shared.ROLE_ADMIN), jwtMiddleware)

	// WARNING: Enabling/disabling a user will invalidate the user status cache only on
	// the instance that processes the request. In a multi-instance setup, a recently
	// disabled user may still be treated as enabled on other instances until their local
	// cache entries expire. For security-critical operations, consider using a shorter
	// cache TTL or a distributed cache with cross-instance invalidation.
	e.POST("/v1/users/:userUID/enable", m.Authorization(h.EnableUser(), shared.ROLE_ADMIN), jwtMiddleware)
	e.POST("/v1/users/:userUID/disable", m.Authorization(h.DisableUser(), shared.ROLE_ADMIN), jwtMiddleware)
}
