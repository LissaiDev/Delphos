package api

import (
	"net/http"
)

// ResponseWriter wraps http.ResponseWriter to capture status code and response size
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode   int
	ResponseSize int
}

// NewResponseWriter creates a new response writer wrapper
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		ResponseSize:   0,
	}
}

// WriteHeader captures the status code and delegates to the wrapped writer
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the response size and delegates to the wrapped writer
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.ResponseSize += size
	return size, err
}

// IsError returns true if the status code indicates an error
func (rw *ResponseWriter) IsError() bool {
	return rw.StatusCode >= 400
}

// IsServerError returns true if the status code indicates a server error
func (rw *ResponseWriter) IsServerError() bool {
	return rw.StatusCode >= 500
}

// IsClientError returns true if the status code indicates a client error
func (rw *ResponseWriter) IsClientError() bool {
	return rw.StatusCode >= 400 && rw.StatusCode < 500
}

// IsSuccess returns true if the status code indicates success
func (rw *ResponseWriter) IsSuccess() bool {
	return rw.StatusCode >= 200 && rw.StatusCode < 300
}

// GetStatusCode returns the captured status code
func (rw *ResponseWriter) GetStatusCode() int {
	return rw.StatusCode
}

// GetResponseSize returns the captured response size
func (rw *ResponseWriter) GetResponseSize() int {
	return rw.ResponseSize
}
