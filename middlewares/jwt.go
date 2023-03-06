package middlewares

import (
	"fmt"
	"panda/apigateway/services/security-service/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JwtMiddleware(jwtSecret string) echo.MiddlewareFunc {

	// JWT middleware - Configure middleware with the custom claims type
	jwtMiddlewareConfig := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte(jwtSecret),
		ErrorHandler: func(err error) error {
			if err != nil {
				fmt.Println(err)
				return echo.ErrUnauthorized
			} else {
				return nil
			}
		},
	}
	jwtMiddleware := middleware.JWTWithConfig(jwtMiddlewareConfig)

	return jwtMiddleware
}
