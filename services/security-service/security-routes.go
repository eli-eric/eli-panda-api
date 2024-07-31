package securityService

import (
	"panda/apigateway/services/security-service/models"

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

	e.GET("/v1/getuserbyazureidtoken", h.GetUserByAzureIdToken())
}
