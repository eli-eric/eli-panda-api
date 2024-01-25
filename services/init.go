package services

import (
	"panda/apigateway/config"
	catalogueService "panda/apigateway/services/catalogue-service"
	codebookService "panda/apigateway/services/codebook-service"
	cronservice "panda/apigateway/services/cron-service"
	ordersService "panda/apigateway/services/orders-service"
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
}
