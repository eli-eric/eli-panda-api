package helpers

import (
	"panda/apigateway/services/security-service/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetUserFromJWT(c echo.Context) (userInfo *models.JwtCustomClaims) {
	userContext := c.Get("user")
	if userContext != nil {
		u := userContext.(*jwt.Token)
		userInfo = u.Claims.(*models.JwtCustomClaims)
	}
	return userInfo
}
