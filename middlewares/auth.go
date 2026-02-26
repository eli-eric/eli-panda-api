package middlewares

import (
	"panda/apigateway/helpers"

	"github.com/labstack/echo/v4"
)

// Authorization is a middleware function that checks if the user has at least one of the required roles to access the endpoint. It retrieves the user information from the JWT token and checks the user's roles against the provided rolesToCheck. If the user has at least one of the required roles, it allows access to the endpoint; otherwise, it returns an unauthorized error.
func Authorization(handlerFunc echo.HandlerFunc, rolesToCheck ...string) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := helpers.GetUserFromJWT(c)

		if user != nil {

			//set user's facility
			c.Set("facilityCode", user.FacilityCode)
			c.Set("userUID", user.Subject)
			c.Set("userName", user.Id)
			c.Set("userRoles", user.Roles)
			//check user's roles -> return unauthorized if user has no of the accepted roles
			for _, roleToCheck := range rolesToCheck {
				for _, role := range user.Roles {
					if role == roleToCheck {
						return handlerFunc(c)
					}
				}
			}
		}

		return echo.ErrUnauthorized
	}
}
