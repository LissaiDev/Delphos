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

func logRequestStart(r *http.Request) {
	logger.Log.Info("Incoming request", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"headers":     r.Header,
		"host":        r.Host,
		"proto":       r.Proto,
	})
}

func logRequestComplete(r *http.Request, rw *ResponseWriter, duration time.Duration) {
	// Ensure response writer implements our interface
	if rw == nil {
		logger.Log.Error("Response writer is nil", nil)
		return
	}

	logger.Log.Info("Request completed", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"duration_ms": duration.Milliseconds(),
		"host":        r.Host,
	})
}

func logStreamingRequestStart(r *http.Request) {
	logger.Log.Info("Incoming streaming request", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
	})
}

func logStreamingRequestComplete(r *http.Request, duration time.Duration) {
	logger.Log.Info("Streaming request completed", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"duration_ms": duration.Milliseconds(),
	})
}

func logRequestError(r *http.Request, rw *ResponseWriter) {
	// Ensure response writer implements our interface
	if rw == nil {
		logger.Log.Error("Response writer is nil", nil)
		return
	}

	logger.Log.Error("Request error", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"host":        r.Host,
		"proto":       r.Proto,
		"headers":     r.Header,
	})
}

func logRateLimitExceeded(r *http.Request, clientIP string) {
	logger.Log.Warn("Rate limit exceeded", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"client_ip":   clientIP,
		"remote_addr": r.RemoteAddr,
	})
}

func logSlowRequest(r *http.Request, rw *ResponseWriter, duration time.Duration) {
	logger.Log.Warn("Slow request detected", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"duration_ms": duration.Milliseconds(),
	})
}

func logRequestMetrics(r *http.Request, rw *ResponseWriter, duration time.Duration) {
	logger.Log.Info("Request metrics", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"duration_ms": duration.Milliseconds(),
		"remote_addr": r.RemoteAddr,
	})
}

func logCORSPreflight(r *http.Request) {
	logger.Log.Debug("CORS preflight request", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"origin":      r.Header.Get("Origin"),
		"remote_addr": r.RemoteAddr,
	})
}

func logStreamingCORSPreflight(r *http.Request) {
	logger.Log.Debug("Streaming CORS preflight request", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"origin":      r.Header.Get("Origin"),
		"remote_addr": r.RemoteAddr,
	})
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
