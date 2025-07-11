package api

import (
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
)

// SecurityMiddleware adds basic security headers
func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		setSecurityHeaders(w)

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware adds CORS headers to all responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			logCORSPreflight(r)
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// StreamingSecurityMiddleware adds basic security headers for streaming
func StreamingSecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setStreamingSecurityHeaders(w)
		next.ServeHTTP(w, r)
	})
}

// StreamingCORSMiddleware adds CORS headers for streaming endpoints
func StreamingCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setStreamingCORSHeaders(w)
		if r.Method == http.MethodOptions {
			logStreamingCORSPreflight(r)
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Helper functions to eliminate DRY violations

// logRequestStart logs the start of a request
func logRequestStart(r *http.Request) {
	log := logger.GetInstance()
	log.Info("Incoming request", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"headers":     r.Header,
		"host":        r.Host,
		"proto":       r.Proto,
	})
}

// logRequestComplete logs the completion of a request
func logRequestComplete(r *http.Request, rw *ResponseWriter, duration time.Duration) {
	log := logger.GetInstance()
	// Ensure response writer implements our interface
	if rw == nil {
		log.Error("Response writer is nil", nil)
		return
	}

	log.Info("Request completed", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"duration_ms": duration.Milliseconds(),
		"host":        r.Host,
	})
}

// logStreamingRequestStart logs the start of a streaming request
func logStreamingRequestStart(r *http.Request) {
	log := logger.GetInstance()

	log.Info("Incoming streaming request", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
	})
}

// logStreamingRequestComplete logs the completion of a streaming request
func logStreamingRequestComplete(r *http.Request, duration time.Duration) {
	log := logger.GetInstance()

	log.Info("Streaming request completed", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"duration_ms": duration.Milliseconds(),
	})
}

// logRequestError logs a request error
func logRequestError(r *http.Request, rw *ResponseWriter) {
	log := logger.GetInstance()

	// Ensure response writer implements our interface
	if rw == nil {
		log.Error("Response writer is nil", nil)
		return
	}

	log.Error("Request error", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"host":        r.Host,
		"proto":       r.Proto,
		"headers":     r.Header,
	})
}

// logRateLimitExceeded logs when the rate limit is exceeded
func logRateLimitExceeded(r *http.Request, clientIP string) {
	log := logger.GetInstance()

	log.Warn("Rate limit exceeded", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"client_ip":   clientIP,
		"remote_addr": r.RemoteAddr,
	})
}

// logSlowRequest logs a slow request
func logSlowRequest(r *http.Request, rw *ResponseWriter, duration time.Duration) {
	log := logger.GetInstance()

	log.Warn("Slow request detected", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"duration_ms": duration.Milliseconds(),
	})
}

// logRequestMetrics logs request metrics
func logRequestMetrics(r *http.Request, rw *ResponseWriter, duration time.Duration) {
	log := logger.GetInstance()

	log.Info("Request metrics", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"duration_ms": duration.Milliseconds(),
		"remote_addr": r.RemoteAddr,
	})
}

// logCORSPreflight logs a CORS preflight request
func logCORSPreflight(r *http.Request) {
	log := logger.GetInstance()

	log.Debug("CORS preflight request", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"origin":      r.Header.Get("Origin"),
		"remote_addr": r.RemoteAddr,
	})
}

// logStreamingCORSPreflight logs a streaming CORS preflight request
func logStreamingCORSPreflight(r *http.Request) {
	log := logger.GetInstance()

	log.Debug("Streaming CORS preflight request", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"origin":      r.Header.Get("Origin"),
		"remote_addr": r.RemoteAddr,
	})
}

// setCORSHeaders sets CORS headers
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
}

// setStreamingCORSHeaders sets CORS headers for streaming
func setStreamingCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cache-Control")
	w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
}

// setSecurityHeaders sets security headers
func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
}

// setStreamingSecurityHeaders sets security headers for streaming
func setStreamingSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
}
