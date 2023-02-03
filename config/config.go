package config

// tag::import[]
import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

/**
 * ReadConfig reads the application settings from enviroment varibales - its possible to use .env file
 */

func LoadConfiguraion() (*Config, error) {
	config := Config{}

	config.JwtSecret = os.Getenv("API_JWT_SECRET")
	config.Port = os.Getenv("API_PORT")
	config.SaltRounds = parseIntWithDefaultValue(os.Getenv("BCRYPT_SALT_ROUNDS"), 12)
	config.Neo4jHost = os.Getenv("NEO4J_HOST")
	config.Neo4jPort = os.Getenv("NEO4J_PORT")
	config.Neo4jUsername = os.Getenv("NEO4J_USER")
	config.Neo4jPassword = os.Getenv("NEO4J_PASSWORD")

	return &config, nil
}

func parseIntWithDefaultValue(inputString string, defaultValue int32) int {
	result, err := strconv.ParseInt(inputString, 10, 32)

	if err != nil {
		result = int64(defaultValue)
	}

	return int(result)
}

type Config struct {
	Neo4jHost     string
	Neo4jPort     string
	Neo4jUsername string
	Neo4jPassword string

	Port       string
	JwtSecret  string
	SaltRounds int
}
