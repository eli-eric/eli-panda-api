package services

import (
	"log"
	"panda/apigateway/config"
	catalogueService "panda/apigateway/services/catalogue-service"
	codebookService "panda/apigateway/services/codebook-service"
	securityService "panda/apigateway/services/security-service"
	systemsService "panda/apigateway/services/systems-service"

	"github.com/labstack/echo/v4"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func InitializeServicesAndMapRoutes(e *echo.Echo, settings *config.Config, neo4jDriver *neo4j.Driver, jwtMiddleware echo.MiddlewareFunc) {

	//security services used in handlers and maped in routes...
	securitySvc := securityService.NewSecurityService(settings, neo4jDriver)
	securityHandlers := securityService.NewSecurityHandlers(securitySvc)
	securityService.MapSecurityRoutes(e, securityHandlers, jwtMiddleware)
	log.Println("Security  service initialized successfully.")

	//security services used in handlers and maped in routes...
	catalogueSvc := catalogueService.NewCatalogueService(settings, neo4jDriver)
	catalogueHandlers := catalogueService.NewCatalogueHandlers(catalogueSvc)
	catalogueService.MapCatalogueRoutes(e, catalogueHandlers, jwtMiddleware)
	log.Println("Catalogue service initialized successfully.")

	systemsSvc := systemsService.NewSystemsService(settings, neo4jDriver)
	systemsHandlers := systemsService.NewsystemsHandlers(systemsSvc)
	systemsService.MapSystemsRoutes(e, systemsHandlers, jwtMiddleware)
	log.Println("Systems   service initialized successfully.")

	//security services used in handlers and maped in routes...
	codebookSvc := codebookService.NewCodebookService(settings, catalogueSvc, securitySvc, systemsSvc)
	codebookHandlers := codebookService.NewCodebookHandlers(codebookSvc)
	codebookService.MapCodebookRoutes(e, codebookHandlers, jwtMiddleware)
	log.Println("Codebook  service initialized successfully.")
}
