package api

import (
	"net/http"
)

// MiddlewareChain represents a chain of middlewares.
// Use Add to add middlewares (pure functions or from the factory).
// Example:
//
//	factory := NewMiddlewareFactory(logger.Log, RateLimitConfig{Window: time.Second})
//	chain := NewMiddlewareChain().
//	  Add(SecurityMiddleware).
//	  Add(CORSMiddleware).
//	  Add(factory.RateLimitMiddleware).
//	  Add(factory.LoggingMiddleware)
//	handler := chain.Apply(finalHandler)
type MiddlewareChain struct {
	middlewares []MiddlewareFunc
}

// MiddlewareFunc defines the signature for middleware functions
type MiddlewareFunc func(http.Handler) http.Handler

// NewMiddlewareChain creates a new middleware chain
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]MiddlewareFunc, 0),
	}
}

// Add adds a middleware to the chain
func (mc *MiddlewareChain) Add(middleware MiddlewareFunc) *MiddlewareChain {
	mc.middlewares = append(mc.middlewares, middleware)
	return mc
}

// Apply applies all middleware in the chain to the given handler
func (mc *MiddlewareChain) Apply(handler http.Handler) http.Handler {
	// Apply middleware in reverse order so they execute in correct order
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		handler = mc.middlewares[i](handler)
	}
	return handler
}
