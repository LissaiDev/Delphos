package api

import (
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
)

// MiddlewareFactory creates middleware with proper dependencies
type MiddlewareFactory struct {
	logger logger.BasicLogger
}

// NewMiddlewareFactory creates a new middleware factory
func NewMiddlewareFactory(log logger.BasicLogger) *MiddlewareFactory {
	return &MiddlewareFactory{
		logger: log,
	}
}

// LoggingMiddleware creates a logging middleware with logger dependency
func (f *MiddlewareFactory) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log incoming request
		f.logger.Info("HTTP Request Started", map[string]interface{}{
			"method":         r.Method,
			"url":            r.URL.String(),
			"remote_addr":    r.RemoteAddr,
			"user_agent":     r.UserAgent(),
			"headers":        len(r.Header),
			"content_length": r.ContentLength,
		})

		// Create a custom response writer to capture status code
		responseWriter := NewResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(responseWriter, r)

		// Log request completion
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

// ErrorLoggingMiddleware creates an error logging middleware
func (f *MiddlewareFactory) ErrorLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture errors
		errorWriter := NewResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(errorWriter, r)

		// Log errors if status code indicates an error
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

// RateLimitMiddleware creates a rate limiting middleware
func (f *MiddlewareFactory) RateLimitMiddleware(next http.Handler) http.Handler {
	// Simple in-memory rate limiter (in production, use Redis or similar)
	clients := make(map[string]time.Time)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		// Check if client has made a request recently (within 1 second)
		if lastRequest, exists := clients[clientIP]; exists {
			if time.Since(lastRequest) < time.Second {
				f.logger.Warn("Rate limit exceeded", map[string]interface{}{
					"client_ip": clientIP,
					"method":    r.Method,
					"url":       r.URL.String(),
				})
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
		}

		// Update last request time
		clients[clientIP] = time.Now()

		next.ServeHTTP(w, r)
	})
}

// MetricsMiddleware creates a metrics logging middleware
func (f *MiddlewareFactory) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Create a custom response writer to capture response size
		metricsWriter := NewResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(metricsWriter, r)

		// Calculate metrics
		duration := time.Since(startTime)

		// Log performance metrics for slow requests (>100ms)
		if duration > 100*time.Millisecond {
			f.logger.Warn("Slow request detected", map[string]interface{}{
				"method":        r.Method,
				"url":           r.URL.String(),
				"duration_ms":   duration.Milliseconds(),
				"response_size": metricsWriter.GetResponseSize(),
				"status_code":   metricsWriter.GetStatusCode(),
			})
		}

		// Log all requests with metrics
		f.logger.Debug("Request metrics", map[string]interface{}{
			"method":        r.Method,
			"url":           r.URL.String(),
			"duration_ms":   duration.Milliseconds(),
			"response_size": metricsWriter.GetResponseSize(),
			"status_code":   metricsWriter.GetStatusCode(),
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
