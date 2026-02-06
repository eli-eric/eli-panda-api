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

// AuthenticateByUsernameAndPassword godoc
// @Summary Authenticate user
// @Description Authenticates user using username and password and returns user info with access token.
// @Tags Security
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "User credentials"
// @Success 200 {object} models.UserAuthInfo
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /v1/authenticate [post]
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

// Get user by azure id token
// GetUserByAzureADIdToken godoc
// @Summary  Get user by azure id token
// @Description  Get user by azure id token
// @Tags Security
// @Param tenantId query string true "Tenant ID"
// @Param azureIdToken query string true "Azure ID Token"
// @Success 200  {object} models.UserAuthInfo
// @Failure 401 "Unauthorized"
// @Router /v1/getuserbyazureidtoken [get]
func (h *SecurityHandlers) GetUserByAzureIdToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		azureIdToken := c.QueryParam("azureIdToken")
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
