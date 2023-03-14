package middlewares

import (
	"panda/apigateway/helpers"

	"github.com/labstack/echo/v4"
)

func Authorization(handlerFunc echo.HandlerFunc, rolesToCheck ...string) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := helpers.GetUserFromJWT(c)

		if user != nil {

			//set user's facility
			c.Set("facilityCode", user.FacilityCode)
			c.Set("userUID", user.Subject)
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
