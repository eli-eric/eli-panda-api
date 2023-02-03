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
	securityService "panda/apigateway/services/security-service"
	"panda/apigateway/services/security-service/models"
	"time"

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

	//new http Echo instance
	e := echo.New()

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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//JWT middleware - Configure middleware with the custom claims type
	config := middleware.JWTConfig{
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
	jwtMiddleware := middleware.JWTWithConfig(config)

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

	// Start server
	go func() {
		if err := e.Start(":" + settings.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server: ELI - PANDA - API Gateway")
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
