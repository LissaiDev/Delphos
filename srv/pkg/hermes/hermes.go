package hermes

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
)

var (
	HermesInstance Fetcher
	once           sync.Once
)

type HermesClient struct {
	logger       logger.Logger
	retries      int
	retryTimeout time.Duration
}

func (h *HermesClient) Fetch(req *Request) Response {
	requestFields := map[string]interface{}{
		"service": string(req.Service),
		"method":  req.Method.String(),
		"url":     req.Url,
	}

	h.logger.Debug("Processing request", requestFields)

	if err := req.Sanitize(); err != nil {
		h.logger.Error("Request sanitization failed", map[string]interface{}{
			"error":   err.Error(),
			"service": string(req.Service),
			"method":  req.Method.String(),
			"url":     req.Url,
		})
		return Response{
			Success: false,
			Code:    500,
			Data:    nil,
		}
	}

	url := SERVICES[req.Service] + req.Url

	// Execute with retry logic
	var (
		responseData []byte
		statusCode   int
		finalErr     error
		attempts     int
	)

	h.logger.Debug("Starting request execution", map[string]interface{}{
		"max_retries":   h.retries,
		"retry_timeout": h.retryTimeout.String(),
	})

	for attempts = 0; attempts <= h.retries; attempts++ {
		if attempts > 0 {
			retryDelay := h.retryTimeout
			h.logger.Warn(fmt.Sprintf("Retrying request (attempt %d/%d)", attempts, h.retries), map[string]interface{}{
				"service":     string(req.Service),
				"method":      req.Method.String(),
				"url":         url,
				"retry_delay": retryDelay.String(),
				"prev_error":  finalErr.Error(),
			})
			time.Sleep(retryDelay)
		}

		// Need to recreate the request for each attempt since body may be consumed
		request, err := h.createRequest(req, url)
		if err != nil {
			finalErr = err
			continue
		}

		responseData, statusCode, finalErr = h.executeRequest(request, url, req)

		// If request succeeded or we got a 4xx error (which won't be fixed by retrying),
		// don't retry
		if finalErr == nil || (statusCode >= 400 && statusCode < 500) {
			break
		}
	}

	if finalErr != nil {
		h.logger.Error("All request attempts failed", map[string]interface{}{
			"error":       finalErr.Error(),
			"service":     string(req.Service),
			"method":      req.Method.String(),
			"url":         url,
			"attempts":    attempts,
			"max_retries": h.retries,
		})
		return Response{
			Success: false,
			Code:    500,
			Data:    nil,
		}
	}

	h.logger.Debug("Request completed successfully", map[string]interface{}{
		"service":       string(req.Service),
		"method":        req.Method.String(),
		"url":           url,
		"status_code":   statusCode,
		"response_size": len(responseData),
		"attempts":      attempts + 1,
		"max_retries":   h.retries,
	})

	return Response{
		Success: true,
		Code:    statusCode,
		Data:    responseData,
	}
}

// createRequest creates a new HTTP request
func (h *HermesClient) createRequest(req *Request, url string) (*http.Request, error) {
	request, err := http.NewRequest(req.Method.String(), url, req.ResolvedBody)
	if err != nil {
		h.logger.Error("Failed to create HTTP request", map[string]interface{}{
			"error":   err.Error(),
			"service": string(req.Service),
			"method":  req.Method.String(),
			"url":     url,
		})
		return nil, err
	}

	// Set request headers
	headerFields := make(map[string]interface{})
	for key, value := range *req.Headers {
		request.Header.Set(key, value)
		headerFields[key] = value
	}
	h.logger.Debug("Request headers", headerFields)

	return request, nil
}

// executeRequest performs the HTTP request and returns response data
func (h *HermesClient) executeRequest(request *http.Request, url string, req *Request) ([]byte, int, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	h.logger.Info("Sending request", map[string]interface{}{
		"service":    string(req.Service),
		"method":     req.Method.String(),
		"url":        url,
		"timeout_ms": client.Timeout.Milliseconds(),
	})

	startTime := time.Now()
	response, err := client.Do(request)
	duration := time.Since(startTime)

	if err != nil {
		h.logger.Error("Request failed", map[string]interface{}{
			"error":       err.Error(),
			"service":     string(req.Service),
			"method":      req.Method.String(),
			"url":         url,
			"duration_ms": duration.Milliseconds(),
		})
		return nil, 0, err
	}

	h.logger.Info("Response received", map[string]interface{}{
		"service":      string(req.Service),
		"method":       req.Method.String(),
		"url":          url,
		"status_code":  response.StatusCode,
		"duration_ms":  duration.Milliseconds(),
		"content_type": response.Header.Get("Content-Type"),
	})

	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)

	if err != nil {
		h.logger.Error("Failed to read response body", map[string]interface{}{
			"error":       err.Error(),
			"service":     string(req.Service),
			"method":      req.Method.String(),
			"url":         url,
			"status_code": response.StatusCode,
		})
		return nil, response.StatusCode, err
	}

	h.logger.Debug("Response completed", map[string]interface{}{
		"service":       string(req.Service),
		"method":        req.Method.String(),
		"url":           url,
		"status_code":   response.StatusCode,
		"response_size": len(data),
		"duration_ms":   duration.Milliseconds(),
	})

	return data, response.StatusCode, nil
}

func New(logger logger.Logger, retries int, retryTimeout time.Duration) Fetcher {
	// Set default values if invalid parameters are provided
	if retries < 0 {
		retries = 0
		logger.Warn("Invalid retries value provided, defaulting to 0", nil)
	}

	if retryTimeout <= 0 {
		retryTimeout = time.Second * 1
		logger.Warn("Invalid retry timeout provided, defaulting to 1 second", nil)
	}

	client := &HermesClient{
		logger:       logger,
		retries:      retries,
		retryTimeout: retryTimeout,
	}

	logger.Info("Hermes client initialized", map[string]interface{}{
		"retries":       retries,
		"retry_timeout": retryTimeout.String(),
		"services":      len(SERVICES),
		"services_list": fmt.Sprintf("%v", SERVICES),
	})

	return client
}

// IsRetryableError determines if an error should trigger a retry
func IsRetryableError(statusCode int, err error) bool {
	// Network errors or timeouts should be retried
	if err != nil {
		return true
	}

	// Server errors (5xx) should be retried
	if statusCode >= 500 && statusCode < 600 {
		return true
	}

	// Specific status codes that might be worth retrying
	switch statusCode {
	case http.StatusTooManyRequests: // 429
		return true
	case http.StatusRequestTimeout: // 408
		return true
	default:
		return false
	}
}

// Get is a convenience method for making GET requests
func (h *HermesClient) Get(service Service, url string, headers *map[string]string) Response {
	h.logger.Debug("Creating GET request", map[string]interface{}{
		"service": string(service),
		"url":     url,
	})

	req := &Request{
		Service: service,
		Url:     url,
		Method:  MethodGet,
		Headers: headers,
	}

	return h.Fetch(req)
}

// Post is a convenience method for making POST requests
func (h *HermesClient) Post(service Service, url string, body *map[string]any, headers *map[string]string) Response {
	h.logger.Debug("Creating POST request", map[string]interface{}{
		"service": string(service),
		"url":     url,
		"body":    body != nil,
	})

	req := &Request{
		Service: service,
		Url:     url,
		Method:  MethodPost,
		Headers: headers,
		Body:    body,
	}

	return h.Fetch(req)
}

// Put is a convenience method for making PUT requests
func (h *HermesClient) Put(service Service, url string, body *map[string]any, headers *map[string]string) Response {
	h.logger.Debug("Creating PUT request", map[string]interface{}{
		"service": string(service),
		"url":     url,
		"body":    body != nil,
	})

	req := &Request{
		Service: service,
		Url:     url,
		Method:  MethodPut,
		Headers: headers,
		Body:    body,
	}

	return h.Fetch(req)
}

// Delete is a convenience method for making DELETE requests
func (h *HermesClient) Delete(service Service, url string, headers *map[string]string) Response {
	h.logger.Debug("Creating DELETE request", map[string]interface{}{
		"service": string(service),
		"url":     url,
	})

	req := &Request{
		Service: service,
		Url:     url,
		Method:  MethodDelete,
		Headers: headers,
	}

	return h.Fetch(req)
}

// Patch is a convenience method for making PATCH requests
func (h *HermesClient) Patch(service Service, url string, body *map[string]any, headers *map[string]string) Response {
	h.logger.Debug("Creating PATCH request", map[string]interface{}{
		"service": string(service),
		"url":     url,
		"body":    body != nil,
	})

	req := &Request{
		Service: service,
		Url:     url,
		Method:  MethodPatch,
		Headers: headers,
		Body:    body,
	}

	return h.Fetch(req)
}

func GetInstance() Fetcher {
	log := logger.GetInstance()
	once.Do(func() {
		HermesInstance = New(log, 3, time.Second*10)
	})

	return HermesInstance
}
