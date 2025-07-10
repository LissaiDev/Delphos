package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
)

// MiddlewareFactory creates middleware with proper dependencies and configuration
// Stateless middlewares (CORS, Security) devem ser usados diretamente das funções puras
// Middlewares com dependências (logger, config) devem ser criados via factory

type RateLimitConfig struct {
	Window time.Duration
}

type MiddlewareFactory struct {
	logger          logger.BasicLogger
	rateLimitConfig RateLimitConfig
}

// NewMiddlewareFactory creates a new middleware factory
func NewMiddlewareFactory(log logger.BasicLogger, rateLimitConfig RateLimitConfig) *MiddlewareFactory {
	return &MiddlewareFactory{
		logger:          log,
		rateLimitConfig: rateLimitConfig,
	}
}

// LoggingMiddleware cria um middleware de logging com logger injetado
func (f *MiddlewareFactory) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		f.logger.Info("HTTP Request Started", map[string]interface{}{
			"method":         r.Method,
			"url":            r.URL.String(),
			"remote_addr":    r.RemoteAddr,
			"user_agent":     r.UserAgent(),
			"headers":        len(r.Header),
			"content_length": r.ContentLength,
		})
		responseWriter := NewResponseWriter(w)
		next.ServeHTTP(responseWriter, r)
		duration := time.Since(startTime)
		f.logger.Info("HTTP Request Completed", map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"status_code": responseWriter.GetStatusCode(),
			"duration":    duration.String(),
			"duration_ms": duration.Milliseconds(),
		})
	})
}

// ErrorLoggingMiddleware cria um middleware de logging de erro
func (f *MiddlewareFactory) ErrorLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorWriter := NewResponseWriter(w)
		next.ServeHTTP(errorWriter, r)
		if errorWriter.IsError() {
			f.logger.Error("HTTP Request Error", map[string]interface{}{
				"method":      r.Method,
				"url":         r.URL.String(),
				"status_code": errorWriter.GetStatusCode(),
				"remote_addr": r.RemoteAddr,
			})
		}
	})
}

// RateLimitMiddleware thread-safe por IP
func (f *MiddlewareFactory) RateLimitMiddleware(next http.Handler) http.Handler {
	clients := sync.Map{}
	window := f.rateLimitConfig.Window
	if window == 0 {
		window = time.Second
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		now := time.Now()
		if last, ok := clients.Load(clientIP); ok {
			if t, ok := last.(time.Time); ok && now.Sub(t) < window {
				f.logger.Warn("Rate limit exceeded", map[string]interface{}{
					"client_ip": clientIP,
					"method":    r.Method,
					"url":       r.URL.String(),
				})
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
		}
		clients.Store(clientIP, now)
		next.ServeHTTP(w, r)
	})
}

// MetricsMiddleware cria um middleware de métricas
func (f *MiddlewareFactory) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		metricsWriter := NewResponseWriter(w)
		next.ServeHTTP(metricsWriter, r)
		duration := time.Since(startTime)
		if duration > 100*time.Millisecond {
			f.logger.Warn("Slow request detected", map[string]interface{}{
				"method":        r.Method,
				"url":           r.URL.String(),
				"duration_ms":   duration.Milliseconds(),
				"response_size": metricsWriter.GetResponseSize(),
				"status_code":   metricsWriter.GetStatusCode(),
			})
		}
		f.logger.Debug("Request metrics", map[string]interface{}{
			"method":        r.Method,
			"url":           r.URL.String(),
			"duration_ms":   duration.Milliseconds(),
			"response_size": metricsWriter.GetResponseSize(),
			"status_code":   metricsWriter.GetStatusCode(),
		})
	})
}

// StreamingLoggingMiddleware creates a streaming logging middleware
func (f *MiddlewareFactory) StreamingLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log incoming streaming request
		f.logger.Info("Streaming Request Started", map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		})

		// Call the next handler directly without wrapping the response writer
		next.ServeHTTP(w, r)

		// Log streaming request completion
		duration := time.Since(startTime)
		f.logger.Info("Streaming Request Completed", map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"duration":    duration.String(),
			"duration_ms": duration.Milliseconds(),
		})
	})
}

// CORSMiddleware creates a CORS middleware
func (f *MiddlewareFactory) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		setCORSHeaders(w)

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			f.logger.Info("CORS Preflight Request", map[string]interface{}{
				"method": r.Method,
				"url":    r.URL.String(),
			})
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// StreamingCORSMiddleware creates a streaming CORS middleware
func (f *MiddlewareFactory) StreamingCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for streaming
		setStreamingCORSHeaders(w)

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			f.logger.Info("CORS Preflight Request (Streaming)", map[string]interface{}{
				"method": r.Method,
				"url":    r.URL.String(),
			})
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// NewAPIChainWithLogger creates a standard API middleware chain with logger
func (f *MiddlewareFactory) NewAPIChainWithLogger() *MiddlewareChain {
	return NewMiddlewareChain().
		Add(SecurityMiddleware).
		Add(f.CORSMiddleware).
		Add(f.RateLimitMiddleware).
		Add(f.LoggingMiddleware).
		Add(f.ErrorLoggingMiddleware).
		Add(f.MetricsMiddleware)
}

// NewStreamingChainWithLogger creates a streaming-specific middleware chain with logger
func (f *MiddlewareFactory) NewStreamingChainWithLogger() *MiddlewareChain {
	return NewMiddlewareChain().
		Add(StreamingSecurityMiddleware).
		Add(f.StreamingCORSMiddleware).
		Add(f.StreamingLoggingMiddleware)
}
