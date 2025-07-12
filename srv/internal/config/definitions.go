package config

import "errors"

// Environment holds the application configuration settings
// Loaded from environment variables with sensible defaults
type Environment struct {
	Name            string  // Application name
	Port            string  // Server port (e.g., ":8080")
	Interval        int     // Monitoring interval in seconds
	CPUThreshold    float64 // CPU usage threshold in percentage
	MemoryThreshold float64 // Memory usage threshold in percentage
	DiskThreshold   float64 // Disk usage threshold in percentage
	WebhookUrl      string  // Discord webhook URL for notifications
	WebhookUsername string  // Username to use for Discord webhook notifications
	Background      bool    // If true, always checks stats in background (without broadcast)
	Cooldown        int     // Cooldown period for notifications (in seconds)
}

// Configuration errors
var (
	ErrInvalidPort            = errors.New("invalid port configuration")
	ErrInvalidInterval        = errors.New("invalid interval configuration")
	ErrInvalidName            = errors.New("invalid name configuration")
	ErrInvalidWebhookUrl      = errors.New("invalid webhook url configuration")
	ErrInvalidWebhookUsername = errors.New("invalid webhook username configuration")
	ErrInvalidCPUThreshold    = errors.New("invalid cpu threshold configuration")
	ErrInvalidMemoryThreshold = errors.New("invalid memory threshold configuration")
	ErrInvalidDiskThreshold   = errors.New("invalid disk threshold configuration")
	ErrInvalidBackground      = errors.New("invalid background configuration")
	ErrInvalidCooldown        = errors.New("invalid cooldown configuration")
)
