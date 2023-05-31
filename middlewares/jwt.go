package middlewares

import (
	"encoding/base64"
	"panda/apigateway/services/security-service/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func JwtMiddleware(jwtSecret string) echo.MiddlewareFunc {

	// JWT middleware - Configure middleware with the custom claims type

	jwtSecretDecoded, _ := base64.RawStdEncoding.DecodeString(jwtSecret)
	stringjwt := string(jwtSecretDecoded)
	log.Info().Msg(stringjwt)

	jwtMiddlewareConfig := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: jwtSecretDecoded,
		ErrorHandler: func(err error) error {
			if err != nil {
				log.Error().Msg(err.Error())
				return echo.ErrUnauthorized
			} else {
				return nil
			}
		},
	}
	jwtMiddleware := middleware.JWTWithConfig(jwtMiddlewareConfig)

	return jwtMiddleware
}
