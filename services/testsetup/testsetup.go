package testsetup

import (
	"log"
	"panda/apigateway/config"
	"panda/apigateway/db"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Shared test driver and session
var TestDriver neo4j.Driver
var TestSession neo4j.Session
var Config config.Config

func init() {
	InitTestDatabase()
}

// InitTestDatabase initializes the shared test Neo4j connection
func InitTestDatabase() {
	if TestDriver != nil {
		return
	}
	projectRoot, _ := filepath.Abs("../../") // Adjust based on your structure
	err := godotenv.Load(filepath.Join(projectRoot, ".env"))
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	// Load environment variables from a .env file (if any)
	Config, err := config.LoadConfiguraion()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize Neo4j driver
	td := db.CreateNeo4jMainInstanceOrPanics(Config.Neo4jUsername, Config.Neo4jPassword, Config.Neo4jHost, Config.Neo4jPort, Config.Neo4jSchema)
	TestDriver = *td
	TestSession = TestDriver.NewSession(neo4j.SessionConfig{})
}

// CleanTestDatabase removes all test data after each test
func CleanTestDatabase() {
	// _, err := TestSession.Run("MATCH (n) DETACH DELETE n", nil)
	// if err != nil {
	// 	log.Println("Error cleaning up database:", err)
	// }
}

// CloseTestDatabase closes the Neo4j driver after tests
func CloseTestDatabase() {
	TestSession.Close()
	TestDriver.Close()
}
