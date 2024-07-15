package securityService

import (
	"net/http"
	"panda/apigateway/services/security-service/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SecurityHandlers struct {
	securityService ISecurityService
}

type ISecurityHandlers interface {
	AuthenticateByUsernameAndPassword() echo.HandlerFunc
	//ChangeUserPassword() echo.HandlerFunc

	GetUserByAzureIdToken() echo.HandlerFunc
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

func (h *SecurityHandlers) GetUserByAzureIdToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		azureIdToken := c.Request().Header.Get("Authorization")
		tenantId := c.QueryParam("tenantId")

		if azureIdToken == "" {
			return echo.ErrUnauthorized
		}

		user, err := h.securityService.GetUserByAzureIdToken(azureIdToken, tenantId)
		if err != nil {
			log.Error().Err(err).Msg("Error getting user by azure id token")
			return echo.ErrUnauthorized
		}

		return c.JSON(http.StatusOK, user)
	}
}
