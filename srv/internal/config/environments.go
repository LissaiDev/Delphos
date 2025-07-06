package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Global environment configuration instance
var Env Environment

// LoadEnvironment loads configuration from environment variables
// Falls back to default values if environment variables are not set
func LoadEnvironment() {
	defaultConfig := Environment{
		Name:     "Delphos Server API",
		Port:     ":8080",
		Interval: 5,
	}

	// Try to load .env file (optional)
	_ = godotenv.Load()

	// Load configuration from environment variables
	name, nameExists := os.LookupEnv("NAME")
	port, portExists := os.LookupEnv("PORT")
	intervalStr, intervalExists := os.LookupEnv("INTERVAL")

	// Use environment variables if all are present, otherwise use defaults
	if nameExists && portExists && intervalExists {
		interval, err := strconv.Atoi(intervalStr)
		if err != nil {
			// If interval parsing fails, use default
			Env = defaultConfig
			return
		}

		Env = Environment{
			Name:     name,
			Port:     port,
			Interval: interval,
		}
	} else {
		Env = defaultConfig
	}
}

// Initialize configuration on package import
func init() {
	LoadEnvironment()
}
