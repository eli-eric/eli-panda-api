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

// if the user has a given role then we return userInfo object
// if the user has not a given role then we return nil
func IsUserInRole(c echo.Context, roleToCheck string) (userInfo *models.JwtCustomClaims) {

	userContext := c.Get("user")
	if userContext != nil {
		u := userContext.(*jwt.Token)
		userInfo = u.Claims.(*models.JwtCustomClaims)
		hasRole := false
		for _, role := range userInfo.Roles {
			if role == roleToCheck {
				hasRole = true
				break
			}
		}
		if !hasRole {
			userInfo = nil
		}
	}
	return userInfo
}

const ROLE_SYSTEMS_VIEW string = "systems-view"
