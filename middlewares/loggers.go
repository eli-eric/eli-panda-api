package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"panda/apigateway/services/security-service/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func RequestLoggerMiddleware() echo.MiddlewareFunc {
	logger := zerolog.New(os.Stdout)
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogMethod:   true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogError:    true,
		LogLatency:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			userID := ""
			userName := ""
			userContext := c.Get("user")
			if userContext != nil {
				u := userContext.(*jwt.Token)
				claims := u.Claims.(*models.JwtCustomClaims)
				userID = claims.Subject
				userName = claims.Id
			}
			if v.Error != nil {
				logger.Error().Timestamp().Int("status", v.Status).Str("method", v.Method).Str("uri", v.URI).Str("user-id", userID).Str("user-name", userName).Int64("latency", v.Latency.Milliseconds()).Str("remote-ip", v.RemoteIP).Msg(v.Error.Error())
			} else {
				//log.Printf("status=%v method=%v uri=%v user-id=%v user-name=%v latency=%vms remote-ip=%v \n", v.Status, v.Method, v.URI, userID, userName, v.Latency.Milliseconds(), v.RemoteIP)
				logger.Info().Timestamp().Int("status", v.Status).Str("method", v.Method).Str("uri", v.URI).Str("user-id", userID).Str("user-name", userName).Int64("latency", v.Latency.Milliseconds()).Str("remote-ip", v.RemoteIP).Send()
			}

			return nil
		},
	})
}

func LogLokiLog(logEntry string) {
	log := LokiLog{}
	values := []string{strconv.FormatInt(time.Now().UnixNano(), 10), logEntry}
	log.Streams = append(log.Streams, LokiLogStream{Stream: LokiLogLabels{App: "panda-api", Level: "request"}, Values: values})
	json_data, err := json.Marshal(log)
	if err == nil {
		resp, err := http.Post("http://localhost:3111/loki/api/v1/push", "application/json", bytes.NewBuffer(json_data))
		fmt.Println(err)
		fmt.Println(resp)
	}
}

type LokiLog struct {
	Streams []LokiLogStream `json:"streams"`
}

type LokiLogStream struct {
	Stream LokiLogLabels `json:"stream"`
	Values []string      `json:"values"`
}

type LokiLogLabels struct {
	App   string `json:"app"`   // panda-api
	Level string `json:"level"` // request, error, ???
}
