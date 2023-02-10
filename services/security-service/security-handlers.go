package securityService

import (
	"net/http"
	"panda/apigateway/services/security-service/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type SecurityHandlers struct {
	securityService ISecurityService
}

type ISecurityHandlers interface {
	AuthenticateByUsernameAndPassword() echo.HandlerFunc
	GetUserByJWT() echo.HandlerFunc
	// RefreshToken() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewSecurityHandlers(securitySvc ISecurityService) ISecurityHandlers {
	return &SecurityHandlers{securityService: securitySvc}
}

// Login with username and password and get jwt token to play with rest of API
func (h *SecurityHandlers) AuthenticateByUsernameAndPassword() echo.HandlerFunc {

	return func(c echo.Context) error {

		cred := new(models.UserCredentials)
		if err := c.Bind(cred); err == nil {
			// authenticate and Generate encoded token and send it as response.
			authUser, err := h.securityService.AuthenticateByUsernameAndPassword(cred.Username, cred.Password)
			if err != nil {
				if err.Error() == "Unauthorized" {
					return echo.ErrUnauthorized
				} else {
					return err
				}
			}

			return c.JSON(http.StatusOK, authUser)
		} else {
			return echo.ErrUnauthorized
		}
	}
}

func (h *SecurityHandlers) RefreshToken() echo.HandlerFunc {

	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		// authenticate and Generate encoded token and send it as response.
		t, err := h.securityService.RefreshToken(claims)
		if err != nil {
			if err.Error() == "Unauthorized" {
				return echo.ErrUnauthorized
			} else {
				return err
			}
		}
		return c.JSON(http.StatusOK, echo.Map{
			"accessToken": t,
		})
	}
}

func (h *SecurityHandlers) GetUserByJWT() echo.HandlerFunc {

	return func(c echo.Context) error {
		//user := c.Get("user").(*jwt.Token)
		//claims := user.Claims.(*models.JwtCustomClaims)

		authUser := models.UserAuthInfo{}

		return c.JSON(http.StatusOK, authUser)
	}
}
