package middlewares

import (
	"fmt"
	"panda/apigateway/services/security-service/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JwtMiddleware(jwtSecret string, userStatusValidator IUserStatusValidator) echo.MiddlewareFunc {

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

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return jwtMiddleware(func(c echo.Context) error {
			if userStatusValidator == nil {
				return next(c)
			}

			userContext := c.Get("user")
			if userContext == nil {
				return echo.ErrUnauthorized
			}

			token, ok := userContext.(*jwt.Token)
			if !ok {
				return echo.ErrUnauthorized
			}

			claims, ok := token.Claims.(*models.JwtCustomClaims)
			if !ok {
				return echo.ErrUnauthorized
			}

			if claims.Subject == "" {
				return echo.ErrUnauthorized
			}

			isEnabled, err := userStatusValidator.ValidateUserEnabled(claims.Subject)
			if err != nil {
				return echo.ErrUnauthorized
			}

			if !isEnabled {
				return echo.ErrUnauthorized
			}

			return next(c)
		})
	}
}
