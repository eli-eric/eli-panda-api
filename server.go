package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"panda/apigateway/config"
	"panda/apigateway/ioutils"
	catalogueService "panda/apigateway/services/catalogue-service"
	codebookService "panda/apigateway/services/codebook-service"
	securityService "panda/apigateway/services/security-service"
	"panda/apigateway/services/security-service/models"
	systemsService "panda/apigateway/services/systems-service"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/neo4j"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	// configuration settings
	// application enviroment varibles described in example.env file
	settings, err := config.LoadConfiguraion()
	ioutils.PanicOnError(err)

	log.Println("PANDA REST API Starting...")
	//new http Echo instance
	e := echo.New()
	e.HideBanner = true
	// Middlewares ************************************************************************************

	//Swagger documentation served from open-api-specification
	swaggerGroup := e.Group("")
	swaggerGroup.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "open-api-specification",
		Browse: true,
	}))

	//CORS middleware to allow cross origin access
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{"*"},
	}))

	//logging and autorecover from panics middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
				fmt.Printf("%v: %v, status: %v, user-id: %v, error: %v, latency: %vms\n", v.Method, v.URI, v.Status, userID, v.Error, v.Latency.Milliseconds())
			} else {
				fmt.Printf("%v: %v, status: %v, user-id: %v, latency: %vms\n", v.Method, v.URI, v.Status, userID, v.Latency.Milliseconds())
			}

			return nil
		},
	}))

	e.Use(middleware.Recover())

	//JWT middleware - Configure middleware with the custom claims type
	jwtMiddlewareConfig := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte(settings.JwtSecret),
		ErrorHandler: func(err error) error {
			if err != nil {
				fmt.Println(err)
				return echo.ErrUnauthorized
			} else {
				return nil
			}
		},
	}
	jwtMiddleware := middleware.JWTWithConfig(jwtMiddlewareConfig)

	// Middlewares END **********************************************************************************

	// NEO4J ********************************************************************************************
	// migrations - init migration lib with neo4j settings
	m, err := migrate.New(
		"file://db/neo4j/migrations",
		"neo4j://"+settings.Neo4jUsername+":"+settings.Neo4jPassword+"@"+settings.Neo4jHost+":"+settings.Neo4jPort+"?x-multi-statement=true")
	// if there is a db error log and shut down
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Applay migrations...")
	// if there is an error in migrations log and shut down, if its successful or there are no changes we can continue
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
	log.Println("Migrations OK")

	// Lets create neo4j database driver which we want to share across the "services"
	// Create new Driver instance
	neo4jDriver, err := neo4j.NewDriver(
		"bolt://"+settings.Neo4jHost+":"+settings.Neo4jPort,
		neo4j.BasicAuth(settings.Neo4jUsername, settings.Neo4jPassword, ""),
	)

	// Check error in driver instantiation
	if err != nil {
		ioutils.PanicOnError(err)
	}

	// Verify Connectivity
	err = neo4jDriver.VerifyConnectivity()

	// If connectivity fails, handle the error
	if err != nil {
		ioutils.PanicOnError(err)
	}

	log.Println("Neo4j security database connection established successfully.")

	// NEO4J END ****************************************************************************************

	//security services used in handlers and maped in routes...
	securitySvc := securityService.NewSecurityService(settings, &neo4jDriver)
	securityHandlers := securityService.NewSecurityHandlers(securitySvc)
	securityService.MapSecurityRoutes(e, securityHandlers, jwtMiddleware)
	log.Println("Security service initialized successfully.")

	//security services used in handlers and maped in routes...
	catalogueSvc := catalogueService.NewCatalogueService(settings, &neo4jDriver)
	catalogueHandlers := catalogueService.NewCatalogueHandlers(catalogueSvc)
	catalogueService.MapCatalogueRoutes(e, catalogueHandlers, jwtMiddleware)
	log.Println("Catalogue service initialized successfully.")

	systemsSvc := systemsService.NewSystemsService(settings, &neo4jDriver)
	log.Println("Catalogue service initialized successfully.")

	//security services used in handlers and maped in routes...
	codebookSvc := codebookService.NewCodebookService(settings, catalogueSvc, securitySvc, systemsSvc)
	codebookHandlers := codebookService.NewCodebookHandlers(codebookSvc)
	codebookService.MapCodebookRoutes(e, codebookHandlers, jwtMiddleware)
	log.Println("Codebook service initialized successfully.")

	// Start server
	go func() {
		if err := e.Start(":" + settings.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server: ELI - PANDA - API Gateway: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
