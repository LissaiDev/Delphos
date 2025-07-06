package api

import (
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
)

// LoggingMiddleware wraps HTTP handlers with comprehensive request/response logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log incoming request
		logger.Log.Info("HTTP Request Started", map[string]interface{}{
			"method":         r.Method,
			"url":            r.URL.String(),
			"remote_addr":    r.RemoteAddr,
			"user_agent":     r.UserAgent(),
			"headers":        len(r.Header),
			"content_length": r.ContentLength,
		})

		// Create a custom response writer to capture status code
		responseWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(responseWriter, r)

		// Log request completion
		duration := time.Since(startTime)
		logger.Log.Info("HTTP Request Completed", map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"status_code": responseWriter.statusCode,
			"duration":    duration.String(),
			"duration_ms": duration.Milliseconds(),
		})
	})
}

// StreamingLoggingMiddleware is a simplified logging middleware for streaming endpoints
func StreamingLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log incoming streaming request
		logger.Log.Info("Streaming Request Started", map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		})

		// Call the next handler directly without wrapping the response writer
		next.ServeHTTP(w, r)

		// Log streaming request completion
		duration := time.Since(startTime)
		logger.Log.Info("Streaming Request Completed", map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"duration":    duration.String(),
			"duration_ms": duration.Milliseconds(),
		})
	})
}

// CORSMiddleware adds CORS headers to all responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			logger.Log.Info("CORS Preflight Request", map[string]interface{}{
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

// SecurityMiddleware adds basic security headers
func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

// ErrorLoggingMiddleware logs errors that occur during request processing
func ErrorLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture errors
		errorWriter := &errorResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(errorWriter, r)

		// Log errors if status code indicates an error
		if errorWriter.statusCode >= 400 {
			logger.Log.Error("HTTP Request Error", map[string]interface{}{
				"method":      r.Method,
				"url":         r.URL.String(),
				"status_code": errorWriter.statusCode,
				"remote_addr": r.RemoteAddr,
			})
		}
	})
}

// errorResponseWriter wraps http.ResponseWriter to capture status code for error logging
type errorResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (erw *errorResponseWriter) WriteHeader(code int) {
	erw.statusCode = code
	erw.ResponseWriter.WriteHeader(code)
}

func (erw *errorResponseWriter) Write(b []byte) (int, error) {
	return erw.ResponseWriter.Write(b)
}

// RateLimitMiddleware provides basic rate limiting per IP address
func RateLimitMiddleware(next http.Handler) http.Handler {
	// Simple in-memory rate limiter (in production, use Redis or similar)
	clients := make(map[string]time.Time)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		// Check if client has made a request recently (within 1 second)
		if lastRequest, exists := clients[clientIP]; exists {
			if time.Since(lastRequest) < time.Second {
				logger.Log.Warn("Rate limit exceeded", map[string]interface{}{
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

// MetricsMiddleware logs performance metrics for requests
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Create a custom response writer to capture response size
		metricsWriter := &metricsResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(metricsWriter, r)

		// Calculate metrics
		duration := time.Since(startTime)

		// Log performance metrics for slow requests (>100ms)
		if duration > 100*time.Millisecond {
			logger.Log.Warn("Slow request detected", map[string]interface{}{
				"method":        r.Method,
				"url":           r.URL.String(),
				"duration_ms":   duration.Milliseconds(),
				"response_size": metricsWriter.responseSize,
				"status_code":   metricsWriter.statusCode,
			})
		}

		// Log all requests with metrics
		logger.Log.Debug("Request metrics", map[string]interface{}{
			"method":        r.Method,
			"url":           r.URL.String(),
			"duration_ms":   duration.Milliseconds(),
			"response_size": metricsWriter.responseSize,
			"status_code":   metricsWriter.statusCode,
		})
	})
}

// metricsResponseWriter wraps http.ResponseWriter to capture response size
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}

func (mrw *metricsResponseWriter) Write(b []byte) (int, error) {
	mrw.responseSize += len(b)
	return mrw.ResponseWriter.Write(b)
}

// StreamingCORSMiddleware adds CORS headers for streaming endpoints
func StreamingCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for streaming
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cache-Control")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			logger.Log.Info("CORS Preflight Request (Streaming)", map[string]interface{}{
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

// StreamingSecurityMiddleware adds basic security headers for streaming
func StreamingSecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add minimal security headers for streaming
		w.Header().Set("X-Content-Type-Options", "nosniff")

		next.ServeHTTP(w, r)
	})
}

// ChainMiddleware applies multiple middlewares in order
func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
