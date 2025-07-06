package config

// Environment holds the application configuration settings
// Loaded from environment variables with sensible defaults
type Environment struct {
	Name     string // Application name
	Port     string // Server port (e.g., ":8080")
	Interval int    // Monitoring interval in seconds
}
