package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"panda/apigateway/services/security-service/models"
	"strconv"
	"time"

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
				log.Printf("%v: %v, status: %v, user-id: %v, error: %v, latency: %vms\n", v.Method, v.URI, v.Status, userID, v.Error, v.Latency.Milliseconds())
			} else {
				log.Printf("%v: %v, status: %v, user-id: %v, latency: %vms\n", v.Method, v.URI, v.Status, userID, v.Latency.Milliseconds())
			}
			// go func() {
			// 	LogLokiLog(fmt.Sprintf(`{ "method": "%v", "uri": "%v", "status": "%v" }`, v.Method, v.URI, v.Status))
			// }()

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
