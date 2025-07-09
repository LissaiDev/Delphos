package hermes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (method Method) String() string {
	switch method {
	case MethodGet:
		return "GET"
	case MethodPost:
		return "POST"
	case MethodPut:
		return "PUT"
	case MethodDelete:
		return "DELETE"
	case MethodPatch:
		return "PATCH"
	case MethodOptions:
		return "OPTIONS"
	default:
		return "UNKNOWN"
	}
}

// MethodFromString converts a string to a Method
func MethodFromString(method string) Method {
	switch strings.ToUpper(method) {
	case "GET":
		return MethodGet
	case "POST":
		return MethodPost
	case "PUT":
		return MethodPut
	case "DELETE":
		return MethodDelete
	case "PATCH":
		return MethodPatch
	case "OPTIONS":
		return MethodOptions
	default:
		return -1 // Invalid method
	}
}

func (req *Request) Sanitize() (err error) {
	// Validate service
	if _, exists := SERVICES[req.Service]; !exists {
		return errors.New("INVALID SERVICE: " + string(req.Service))
	}

	// Validate and clean URL
	if trimmedUrl := strings.TrimSpace(req.Url); trimmedUrl == "" {
		return errors.New("INVALID URL: empty URL")
	} else {
		req.Url = trimmedUrl
	}

	// Ensure URL starts with "/"
	if !strings.HasPrefix(req.Url, "/") {
		req.Url = "/" + req.Url
	}

	// Set default headers if none provided
	if req.Headers == nil {
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		req.Headers = &headers
	}

	// Validate HTTP method
	if req.Method.String() == "UNKNOWN" {
		return errors.New("INVALID METHOD: method not recognized")
	}

	// Process request body
	if req.Body != nil {
		// Validate method for request with body
		if req.Method.String() == "GET" || req.Method.String() == "DELETE" {
			return errors.New("INVALID METHOD: " + req.Method.String() + " cannot have a request body")
		}

		// Marshal body to JSON
		resolved, err := json.Marshal(req.Body)
		if err != nil {
			return errors.New("INVALID BODY: failed to marshal JSON: " + err.Error())
		}

		req.ResolvedBody = bytes.NewBuffer(resolved)
	}

	return nil
}

// LogRequest returns a string representation of the request for logging
func (req *Request) LogRequest() string {
	var bodyInfo string
	if req.Body != nil {
		bodyInfo = fmt.Sprintf(", with body of %d fields", len(*req.Body))
	} else {
		bodyInfo = ", no body"
	}

	headerCount := 0
	if req.Headers != nil {
		headerCount = len(*req.Headers)
	}

	return fmt.Sprintf("%s %s%s to %s%s (headers: %d)",
		req.Method.String(),
		SERVICES[req.Service],
		req.Url,
		string(req.Service),
		bodyInfo,
		headerCount)
}

// ValidateServiceURL checks if a service and URL combination is valid
func ValidateServiceURL(service Service, url string) error {
	if _, exists := SERVICES[service]; !exists {
		return errors.New("unknown service: " + string(service))
	}

	if strings.TrimSpace(url) == "" {
		return errors.New("empty URL")
	}

	return nil
}
