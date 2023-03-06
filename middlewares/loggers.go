package middlewares

import (
	"fmt"
	"panda/apigateway/services/security-service/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogMethod:   true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogError:    true,
		LogLatency:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			userID := ""
			userContext := c.Get("user")
			if userContext != nil {
				u := userContext.(*jwt.Token)
				claims := u.Claims.(*models.JwtCustomClaims)
				userID = claims.Subject
			}
			if v.Error != nil {
				fmt.Printf("%v: %v, status: %v, user-id: %v, error: %v, latency: %vms\n", v.Method, v.URI, v.Status, userID, v.Error, v.Latency.Milliseconds())
			} else {
				fmt.Printf("%v: %v, status: %v, user-id: %v, latency: %vms\n", v.Method, v.URI, v.Status, userID, v.Latency.Milliseconds())
			}

			return nil
		},
	})
}
