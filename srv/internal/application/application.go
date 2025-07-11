package application

import (
	"net/http"
	"sync"
	"time"

	"github.com/LissaiDev/Delphos/internal/api"
	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/internal/monitor"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

// Application encapsulates the main application logic
type Application struct {
	broker            *api.Broker
	statsService      *monitor.StatsService
	logger            logger.BasicLogger
	config            *config.Environment
	middlewareFactory *api.MiddlewareFactory
}

var (
	applicationInstance *Application
	once                sync.Once
)

// NewApplication creates a new application instance
func New() *Application {
	log := logger.GetInstance()
	rateLimitConfig := api.RateLimitConfig{Window: time.Second}
	return &Application{
		broker:            api.GetInstance(),
		statsService:      monitor.GetInstance(),
		logger:            log,
		config:            &config.Env,
		middlewareFactory: api.NewMiddlewareFactory(log, rateLimitConfig),
	}
}

// Start starts the application
func (app *Application) Start() error {
	app.broker.Start()
	defer app.broker.Stop()

	// Start background stats broadcasting
	go app.startStatsBackgroundProcess()

	// Setup HTTP routes
	app.setupRoutes()

	// Start HTTP server
	return app.startHTTPServer()
}

// startStatsBackgroundProcess handles periodic stats broadcasting
func (app *Application) startStatsBackgroundProcess() {
	ticker := time.NewTicker(time.Duration(app.config.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if len(app.broker.Clients) == 0 {
			continue
		}

		data, err := app.statsService.GetStatsJSON()
		if err != nil {
			app.logger.Error("Failed to get stats JSON", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		app.broker.Broadcast(string(data))
	}
}

// setupRoutes configures HTTP routes with middleware chains
func (app *Application) setupRoutes() {
	// Create middleware chains using the factory and pure functions
	apiChain := api.NewMiddlewareChain().
		Add(api.SecurityMiddleware).
		Add(api.CORSMiddleware).
		Add(app.middlewareFactory.RateLimitMiddleware).
		Add(app.middlewareFactory.LoggingMiddleware).
		Add(app.middlewareFactory.ErrorLoggingMiddleware).
		Add(app.middlewareFactory.MetricsMiddleware)

	streamingChain := api.NewMiddlewareChain().
		Add(api.StreamingSecurityMiddleware).
		Add(api.StreamingCORSMiddleware).
		Add(app.middlewareFactory.StreamingLoggingMiddleware)

	// Create handlers
	statsHandler := apiChain.Apply(http.HandlerFunc(api.SystemStatsHandler))
	sseHandler := streamingChain.Apply(app.broker)

	// Register routes
	http.Handle("/api/stats", statsHandler)
	http.Handle("/api/stats/sse", sseHandler)
}

// startHTTPServer starts the HTTP server
func (app *Application) startHTTPServer() error {
	app.logger.Info("Starting HTTP server", map[string]interface{}{
		"port": app.config.Port,
		"name": app.config.Name,
	})

	if err := http.ListenAndServe(app.config.Port, nil); err != nil {
		app.logger.Fatal("Failed to start HTTP server", map[string]interface{}{
			"error": err.Error(),
			"port":  app.config.Port,
		})
		return err
	}
	return nil
}

func GetInstance() *Application {
	once.Do(func() {
		applicationInstance = New()
	})
	return applicationInstance
}
