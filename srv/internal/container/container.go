package container

import (
	"time"

	"github.com/LissaiDev/Delphos/internal/api"
	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/internal/monitor"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

// Container manages application dependencies
type Container struct {
	logger       logger.Logger
	config       *config.Environment
	broker       *api.Broker
	statsService *monitor.StatsService
	server       *Server
}

// NewContainer creates a new dependency container
func NewContainer() *Container {
	return &Container{
		logger: logger.New(),
		config: &config.Env,
	}
}

// Logger returns the logger instance
func (c *Container) Logger() logger.Logger {
	return c.logger
}

// Config returns the configuration
func (c *Container) Config() *config.Environment {
	return c.config
}

// Broker returns the broker instance, creating it if needed
func (c *Container) Broker() *api.Broker {
	if c.broker == nil {
		c.broker = api.NewBroker()
	}
	return c.broker
}

// StatsService returns the stats service, creating it if needed
func (c *Container) StatsService() *monitor.StatsService {
	if c.statsService == nil {
		c.statsService = monitor.NewStatsService(c.logger)
	}
	return c.statsService
}

// Server returns the server instance, creating it if needed
func (c *Container) Server() *Server {
	if c.server == nil {
		c.server = NewServer(c)
	}
	return c.server
}

// Server manages HTTP server setup and routing
type Server struct {
	container *Container
}

// NewServer creates a new server instance
func NewServer(container *Container) *Server {
	return &Server{
		container: container,
	}
}

// Start starts the HTTP server with all configured handlers
func (s *Server) Start() error {
	broker := s.container.Broker()
	broker.Start()

	// Start stats broadcasting
	go s.startStatsBackgroundProcess()

	// Setup routes
	s.setupRoutes()

	return s.listenAndServe()
}

// startStatsBackgroundProcess handles periodic stats broadcasting
func (s *Server) startStatsBackgroundProcess() {
	ticker := time.NewTicker(time.Duration(s.container.Config().Interval) * time.Second)
	defer ticker.Stop()

	statsService := s.container.StatsService()
	broker := s.container.Broker()

	for range ticker.C {
		if len(broker.Clients) == 0 {
			continue
		}

		data, err := statsService.GetStatsJSON()
		if err != nil {
			s.container.Logger().Error("Failed to get stats", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		broker.Broadcast(string(data))
	}
}

// setupRoutes configures HTTP routes
func (s *Server) setupRoutes() {
	// This would be implemented based on your routing needs
	// For now, keeping it simple to maintain existing functionality
}

// listenAndServe starts the HTTP server
func (s *Server) listenAndServe() error {
	// Implementation would go here
	return nil
}
