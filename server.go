package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"panda/apigateway/config"
	"panda/apigateway/db"
	"panda/apigateway/ioutils"
	"panda/apigateway/middlewares"
	"panda/apigateway/services"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"

	_ "github.com/golang-migrate/migrate/v4/database/neo4j"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "panda/apigateway/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title PANDA REST API - localhost
// @version 1.0
// @description This is the REST API to the PANDA database. \n This is the only place to access data from the PANDA database.

// @contact.name Jiří Švácha
// @contact.email jiri.svacha@eli-beams.eu

// @schemes http
// @host localhost:50000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token. <br> How to obtain: https://eli-eric.atlassian.net/wiki/spaces/CS/pages/948797504/How+to+get+PANDA+API+Token <br> Add word Bearer before the token here.
func main() {

	//set locale to europe/prague
	loc, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Local = loc

	// configuration settings
	// application enviroment varibles described in example.env file
	settings, err := config.LoadConfiguraion()
	ioutils.PanicOnError(err)

	fmt.Print(ioutils.GetWelcomeMessage())
	log.Info().Msg("PANDA REST API Starting...")
	//new http Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middlewares ************************************************************************************
	//Swagger documentation served from open-api-specification
	swaggerGroup := e.Group("")
	swaggerGroup.Use(middlewares.StaticMiddleware())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//CORS middleware to allow cross origin access
	e.Use(middlewares.CORSMiddleware())

	//logging middleware
	e.Use(middlewares.RequestLoggerMiddleware())

	//register recover middleware
	e.Use(middlewares.RecoverMiddleware())

	//JWT middleware - Configure middleware with the custom claims type
	jwtMiddleware := middlewares.JwtMiddleware(settings.JwtSecret)

	// Middlewares END **********************************************************************************

	// NEO4J ********************************************************************************************
	// migrations - init migration lib with neo4j settings
	db.MigrateNeo4jMainInstance(settings.Neo4jUsername, settings.Neo4jPassword, settings.Neo4jHost, settings.Neo4jPort, settings.Neo4jSchema)
	// Lets create neo4j database driver which we want to share across the "services"
	// Create new DB Driver instance
	neo4jDriver := db.CreateNeo4jMainInstanceOrPanics(settings.Neo4jUsername, settings.Neo4jPassword, settings.Neo4jHost, settings.Neo4jPort, settings.Neo4jSchema)
	// NEO4J END ****************************************************************************************

	//Init all services
	services.InitializeServicesAndMapRoutes(e, settings, neo4jDriver, jwtMiddleware)

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
