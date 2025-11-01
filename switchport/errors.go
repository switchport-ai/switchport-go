// Package switchport provides a Go client library for Switchport - Prompt management and A/B testing platform.
package switchport

import "fmt"

// SwitchportError is the base error type for all Switchport SDK errors.
type SwitchportError struct {
	Message string
}

func (e *SwitchportError) Error() string {
	return e.Message
}

// AuthenticationError is raised when API key authentication fails.
type AuthenticationError struct {
	SwitchportError
}

// NewAuthenticationError creates a new AuthenticationError.
func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{
		SwitchportError: SwitchportError{Message: message},
	}
}

// PromptNotFoundError is raised when a prompt config with the given key is not found.
type PromptNotFoundError struct {
	SwitchportError
}

// NewPromptNotFoundError creates a new PromptNotFoundError.
func NewPromptNotFoundError(promptKey string) *PromptNotFoundError {
	return &PromptNotFoundError{
		SwitchportError: SwitchportError{
			Message: fmt.Sprintf("Prompt config with key '%s' not found", promptKey),
		},
	}
}

// MetricNotFoundError is raised when a metric with the given key is not found.
type MetricNotFoundError struct {
	SwitchportError
}

// NewMetricNotFoundError creates a new MetricNotFoundError.
func NewMetricNotFoundError(metricKey string) *MetricNotFoundError {
	return &MetricNotFoundError{
		SwitchportError: SwitchportError{
			Message: fmt.Sprintf("Metric definition with key '%s' not found", metricKey),
		},
	}
}

// APIError is raised when the Switchport API returns an error.
type APIError struct {
	SwitchportError
	StatusCode   int
	ResponseData map[string]interface{}
}

// NewAPIError creates a new APIError.
func NewAPIError(message string, statusCode int, responseData map[string]interface{}) *APIError {
	return &APIError{
		SwitchportError: SwitchportError{Message: message},
		StatusCode:      statusCode,
		ResponseData:    responseData,
	}
}
