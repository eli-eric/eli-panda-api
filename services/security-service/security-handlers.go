package securityService

import (
	"net/http"
	"panda/apigateway/middlewares"
	"panda/apigateway/services/security-service/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SecurityHandlers struct {
	securityService     ISecurityService
	userStatusValidator middlewares.IUserStatusValidator
}

type ISecurityHandlers interface {
	AuthenticateByUsernameAndPassword() echo.HandlerFunc
	//ChangeUserPassword() echo.HandlerFunc

	GetUserByAzureIdToken() echo.HandlerFunc
	GetUserStatusCache() echo.HandlerFunc
	InvalidateUserStatusCache() echo.HandlerFunc
	EnableUser() echo.HandlerFunc
	DisableUser() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewSecurityHandlers(securitySvc ISecurityService, userStatusValidator middlewares.IUserStatusValidator) ISecurityHandlers {
	return &SecurityHandlers{securityService: securitySvc, userStatusValidator: userStatusValidator}
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

// Get user status cache godoc
// @Summary Get user status cache
// @Description Returns current users stored in auth status cache.
// @Tags Security
// @Security BearerAuth
// @Success 200 {array} models.UserStatusCacheItem
// @Failure 500 "Internal server error"
// @Router /v1/auth/cache [get]
func (h *SecurityHandlers) GetUserStatusCache() echo.HandlerFunc {
	return func(c echo.Context) error {
		if h.userStatusValidator == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "user status validator is not configured")
		}

		cacheEntries := h.userStatusValidator.GetCacheEntries()
		response := make([]models.UserStatusCacheItem, 0, len(cacheEntries))
		for _, cacheEntry := range cacheEntries {
			response = append(response, models.UserStatusCacheItem{
				UserUID:   cacheEntry.UserUID,
				IsEnabled: cacheEntry.IsEnabled,
				ExpiresAt: cacheEntry.ExpiresAt,
			})
		}

		return c.JSON(http.StatusOK, response)
	}
}

// Invalidate user status cache godoc
// @Summary Invalidate user status cache
// @Description Invalidates cached auth status for selected user UID.
// @Tags Security
// @Security BearerAuth
// @Param userUID path string true "User UID"
// @Success 200 {string} string
// @Failure 400 "Bad Request"
// @Failure 500 "Internal server error"
// @Router /v1/auth/cache/invalidate/{userUID} [post]
func (h *SecurityHandlers) InvalidateUserStatusCache() echo.HandlerFunc {
	return func(c echo.Context) error {
		userUID := c.Param("userUID")
		if userUID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "userUID is required")
		}

		if h.userStatusValidator == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "user status validator is not configured")
		}

		h.userStatusValidator.InvalidateUser(userUID)

		return c.JSON(http.StatusOK, "OK")
	}
}

// Enable user godoc
// @Summary Enable user
// @Description Enables user account and invalidates auth status cache.
// @Tags Security
// @Security BearerAuth
// @Param userUID path string true "User UID"
// @Success 200 {object} models.UserStatusResponse
// @Failure 400 "Bad Request"
// @Failure 500 "Internal server error"
// @Router /v1/users/{userUID}/enable [post]
func (h *SecurityHandlers) EnableUser() echo.HandlerFunc {
	return h.setUserEnabledHandler(true)
}

// Disable user godoc
// @Summary Disable user
// @Description Disables user account and invalidates auth status cache.
// @Tags Security
// @Security BearerAuth
// @Param userUID path string true "User UID"
// @Success 200 {object} models.UserStatusResponse
// @Failure 400 "Bad Request"
// @Failure 500 "Internal server error"
// @Router /v1/users/{userUID}/disable [post]
func (h *SecurityHandlers) DisableUser() echo.HandlerFunc {
	return h.setUserEnabledHandler(false)
}

func (h *SecurityHandlers) setUserEnabledHandler(isEnabled bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		userUID := c.Param("userUID")
		if userUID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "userUID is required")
		}

		if h.userStatusValidator == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "user status validator is not configured")
		}

		// First invalidate any existing cached state before updating the database
		h.userStatusValidator.InvalidateUser(userUID)

		updatedUserUID, err := h.securityService.SetUserEnabled(userUID, isEnabled)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// Invalidate again after the database update to avoid caching a stale state
		h.userStatusValidator.InvalidateUser(updatedUserUID)

		return c.JSON(http.StatusOK, models.UserStatusResponse{
			UserUID:   updatedUserUID,
			IsEnabled: isEnabled,
		})
	}
}
