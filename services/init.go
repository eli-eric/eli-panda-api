package services

import (
	"panda/apigateway/config"
	catalogueService "panda/apigateway/services/catalogue-service"
	codebookService "panda/apigateway/services/codebook-service"
	cronservice "panda/apigateway/services/cron-service"
	filesservice "panda/apigateway/services/files-service"
	"panda/apigateway/services/general"
	ordersService "panda/apigateway/services/orders-service"
	publicationsservice "panda/apigateway/services/publications-service"
	securityService "panda/apigateway/services/security-service"
	systemsService "panda/apigateway/services/systems-service"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func InitializeServicesAndMapRoutes(e *echo.Echo, settings *config.Config, neo4jDriver *neo4j.Driver, jwtMiddleware echo.MiddlewareFunc) {

	//security services used in handlers and maped in routes...
	securitySvc := securityService.NewSecurityService(settings, neo4jDriver)
	securityHandlers := securityService.NewSecurityHandlers(securitySvc)
	securityService.MapSecurityRoutes(e, securityHandlers, jwtMiddleware)
	log.Info().Msg("Security  service initialized successfully.")

	//security services used in handlers and maped in routes...
	catalogueSvc := catalogueService.NewCatalogueService(settings, neo4jDriver)
	catalogueHandlers := catalogueService.NewCatalogueHandlers(catalogueSvc)
	catalogueService.MapCatalogueRoutes(e, catalogueHandlers, jwtMiddleware)
	log.Info().Msg("Catalogue service initialized successfully.")

	systemsSvc := systemsService.NewSystemsService(settings, neo4jDriver)
	systemsHandlers := systemsService.NewsystemsHandlers(systemsSvc)
	systemsService.MapSystemsRoutes(e, systemsHandlers, jwtMiddleware)
	log.Info().Msg("Systems   service initialized successfully.")

	ordersSvc := ordersService.NewOrdersService(neo4jDriver)
	ordersHandlers := ordersService.NewOrdersHandlers(ordersSvc)
	ordersService.MapOrdersRoutes(e, ordersHandlers, jwtMiddleware)
	log.Info().Msg("Orders    service initialized successfully.")

	//security services used in handlers and maped in routes...
	codebookSvc := codebookService.NewCodebookService(settings, neo4jDriver, catalogueSvc, securitySvc, systemsSvc, ordersSvc)
	codebookHandlers := codebookService.NewCodebookHandlers(codebookSvc)
	codebookService.MapCodebookRoutes(e, codebookHandlers, jwtMiddleware)
	log.Info().Msg("Codebook  service initialized successfully.")

	// cron service
	cronSvc := cronservice.NewCronService(settings, neo4jDriver)
	cronHandlers := cronservice.NewCronHandlers(cronSvc)
	cronservice.MapCronRoutes(e, cronHandlers, jwtMiddleware)
	log.Info().Msg("Cron      service initialized successfully.")

	// files service
	filesSvc := filesservice.NewFilesService(neo4jDriver)
	filesHandlers := filesservice.NewFilesHandlers(filesSvc)
	filesservice.MapFilesRoutes(e, filesHandlers, jwtMiddleware)
	log.Info().Msg("Files     service initialized successfully.")

	// general service
	generalSvc := general.NewGeneralService(neo4jDriver)
	generalHandlers := general.NewGeneralHandlers(generalSvc)
	general.MapGeneralRoutes(e, generalHandlers, jwtMiddleware)
	log.Info().Msg("General   service initialized successfully.")

	// publications service
	publicationsSvc := publicationsservice.NewPublicationsService(neo4jDriver, settings.ApiIntegrationBeamlinesWOSBaseUrl, settings.ApiIntegrationBeamlinesWOSBaseApiKey)
	publicationsHandlers := publicationsservice.NewPublicationsHandlers(publicationsSvc)
	publicationsservice.MapPublicationsRoutes(e, publicationsHandlers, jwtMiddleware)
	log.Info().Msg("Publications service initialized successfully.")
}
