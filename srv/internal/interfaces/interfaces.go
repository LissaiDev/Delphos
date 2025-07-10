package interfaces

import (
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/internal/monitor"
)

// StatsProvider defines the contract for system statistics collection
type StatsProvider interface {
	GetStats() (*monitor.Monitor, error)
	GetStatsJSON() ([]byte, error)
}

// MessageBroker defines the contract for message broadcasting
type MessageBroker interface {
	Start()
	Stop()
	Broadcast(message string)
	AddClient(client chan string)
	RemoveClient(client chan string)
	ClientCount() int
}

// HTTPHandler defines the contract for HTTP request handling
type HTTPHandler interface {
	http.Handler
}

// MiddlewareFunc defines the contract for HTTP middleware
type MiddlewareFunc func(http.Handler) http.Handler

// MiddlewareChain defines the contract for chaining middleware
type MiddlewareChain interface {
	Apply(handler http.Handler) http.Handler
	Add(middleware MiddlewareFunc) MiddlewareChain
}

// ConfigProvider defines the contract for configuration access
type ConfigProvider interface {
	GetPort() string
	GetName() string
	GetInterval() time.Duration
}

// BasicLogger defines minimal logging interface
type BasicLogger interface {
	Debug(message string, fields map[string]interface{})
	Info(message string, fields map[string]interface{})
	Warn(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
}

// SystemMonitor defines the contract for system monitoring
type SystemMonitor interface {
	GetHostInfo() (*monitor.Host, error)
	GetMemoryInfo() (*monitor.Memory, error)
	GetCPUInfo() ([]*monitor.CPU, error)
	GetDiskInfo() ([]*monitor.Disk, error)
	GetNetworkInfo() ([]*monitor.Network, error)
}

// ServerRunner defines the contract for server management
type ServerRunner interface {
	Start() error
	Stop() error
	IsRunning() bool
}

// RouteHandler defines the contract for route handling
type RouteHandler interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler http.HandlerFunc)
}
