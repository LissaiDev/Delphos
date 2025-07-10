package api

import (
	"net/http"
	"time"
)

// LoggingMiddleware wraps HTTP handlers with comprehensive request/response logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log incoming request
		logRequestStart(r)

		// Create a custom response writer to capture status code
		responseWriter := NewResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(responseWriter, r)

		// Log request completion
		logRequestComplete(r, responseWriter, time.Since(startTime))
	})
}

// StreamingLoggingMiddleware is a simplified logging middleware for streaming endpoints
func StreamingLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log incoming streaming request
		logStreamingRequestStart(r)

		// Call the next handler directly without wrapping the response writer
		next.ServeHTTP(w, r)

		// Log streaming request completion
		logStreamingRequestComplete(r, time.Since(startTime))
	})
}

// CORSMiddleware adds CORS headers to all responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		setCORSHeaders(w)

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			logCORSPreflight(r)
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
		setSecurityHeaders(w)

		next.ServeHTTP(w, r)
	})
}

// ErrorLoggingMiddleware logs errors that occur during request processing
func ErrorLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture errors
		errorWriter := NewResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(errorWriter, r)

		// Log errors if status code indicates an error
		if errorWriter.IsError() {
			logRequestError(r, errorWriter)
		}
	})
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
				logRateLimitExceeded(r, clientIP)
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
		metricsWriter := NewResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(metricsWriter, r)

		// Calculate metrics
		duration := time.Since(startTime)

		// Log performance metrics for slow requests (>100ms)
		if duration > 100*time.Millisecond {
			logSlowRequest(r, metricsWriter, duration)
		}

		// Log all requests with metrics
		logRequestMetrics(r, metricsWriter, duration)
	})
}

// StreamingCORSMiddleware adds CORS headers for streaming endpoints
func StreamingCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for streaming
		setStreamingCORSHeaders(w)

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			logStreamingCORSPreflight(r)
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
		setStreamingSecurityHeaders(w)

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

// Helper functions to eliminate DRY violations

func logRequestStart(r *http.Request) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func logRequestComplete(r *http.Request, rw *ResponseWriter, duration time.Duration) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func logStreamingRequestStart(r *http.Request) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func logStreamingRequestComplete(r *http.Request, duration time.Duration) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func logRequestError(r *http.Request, rw *ResponseWriter) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func logRateLimitExceeded(r *http.Request, clientIP string) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func logSlowRequest(r *http.Request, rw *ResponseWriter, duration time.Duration) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func logRequestMetrics(r *http.Request, rw *ResponseWriter, duration time.Duration) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func logCORSPreflight(r *http.Request) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func logStreamingCORSPreflight(r *http.Request) {
	// Note: Logger is not imported to avoid circular dependency
	// This is a placeholder - in production, inject logger as dependency
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
}

func setStreamingCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cache-Control")
	w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
}

func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
}

func setStreamingSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
}
