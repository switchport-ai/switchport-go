package switchport

import (
	"testing"
)

func TestSwitchportError_Error(t *testing.T) {
	err := &SwitchportError{Message: "test error"}

	if err.Error() != "test error" {
		t.Errorf("Expected error message 'test error', got '%s'", err.Error())
	}
}

func TestNewAuthenticationError(t *testing.T) {
	message := "Invalid API key"
	err := NewAuthenticationError(message)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	if err.Error() != message {
		t.Errorf("Expected error message '%s', got '%s'", message, err.Error())
	}

	// Verify it's the correct type
	if _, ok := interface{}(err).(*AuthenticationError); !ok {
		t.Errorf("Expected *AuthenticationError, got %T", err)
	}

	// Verify it implements error interface
	var _ error = err
}

func TestNewPromptNotFoundError(t *testing.T) {
	promptKey := "test-prompt"
	err := NewPromptNotFoundError(promptKey)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMessage := "Prompt config with key 'test-prompt' not found"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}

	// Verify it's the correct type
	if _, ok := interface{}(err).(*PromptNotFoundError); !ok {
		t.Errorf("Expected *PromptNotFoundError, got %T", err)
	}

	// Verify it implements error interface
	var _ error = err
}

func TestNewMetricNotFoundError(t *testing.T) {
	metricKey := "test-metric"
	err := NewMetricNotFoundError(metricKey)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMessage := "Metric definition with key 'test-metric' not found"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}

	// Verify it's the correct type
	if _, ok := interface{}(err).(*MetricNotFoundError); !ok {
		t.Errorf("Expected *MetricNotFoundError, got %T", err)
	}

	// Verify it implements error interface
	var _ error = err
}

func TestNewAPIError(t *testing.T) {
	message := "API request failed"
	statusCode := 500
	responseData := map[string]interface{}{
		"detail": "Internal server error",
		"code":   "SERVER_ERROR",
	}

	err := NewAPIError(message, statusCode, responseData)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	if err.Error() != message {
		t.Errorf("Expected error message '%s', got '%s'", message, err.Error())
	}

	if err.StatusCode != statusCode {
		t.Errorf("Expected status code %d, got %d", statusCode, err.StatusCode)
	}

	if err.ResponseData["detail"] != "Internal server error" {
		t.Errorf("Expected ResponseData detail 'Internal server error', got %v", err.ResponseData["detail"])
	}

	if err.ResponseData["code"] != "SERVER_ERROR" {
		t.Errorf("Expected ResponseData code 'SERVER_ERROR', got %v", err.ResponseData["code"])
	}

	// Verify it's the correct type
	if _, ok := interface{}(err).(*APIError); !ok {
		t.Errorf("Expected *APIError, got %T", err)
	}

	// Verify it implements error interface
	var _ error = err
}

func TestNewAPIError_NilResponseData(t *testing.T) {
	err := NewAPIError("Error message", 400, nil)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	if err.ResponseData != nil {
		t.Errorf("Expected ResponseData to be nil, got %v", err.ResponseData)
	}
}

func TestNewAPIError_EmptyResponseData(t *testing.T) {
	responseData := map[string]interface{}{}
	err := NewAPIError("Error message", 400, responseData)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	if len(err.ResponseData) != 0 {
		t.Errorf("Expected empty ResponseData, got %d keys", len(err.ResponseData))
	}
}

func TestErrorTypeAssertion_AuthenticationError(t *testing.T) {
	err := NewAuthenticationError("Invalid API key")

	// Test type assertion
	authErr, ok := interface{}(err).(*AuthenticationError)
	if !ok {
		t.Fatal("Failed to assert as *AuthenticationError")
	}

	if authErr.Error() != "Invalid API key" {
		t.Errorf("Expected error message 'Invalid API key', got '%s'", authErr.Error())
	}

	// Verify it's also a SwitchportError
	baseErr := &authErr.SwitchportError
	if baseErr.Message != "Invalid API key" {
		t.Errorf("Expected base error message 'Invalid API key', got '%s'", baseErr.Message)
	}
}

func TestErrorTypeAssertion_PromptNotFoundError(t *testing.T) {
	err := NewPromptNotFoundError("my-prompt")

	promptErr, ok := interface{}(err).(*PromptNotFoundError)
	if !ok {
		t.Fatal("Failed to assert as *PromptNotFoundError")
	}

	if promptErr.Error() != "Prompt config with key 'my-prompt' not found" {
		t.Errorf("Unexpected error message: %s", promptErr.Error())
	}
}

func TestErrorTypeAssertion_MetricNotFoundError(t *testing.T) {
	err := NewMetricNotFoundError("my-metric")

	metricErr, ok := interface{}(err).(*MetricNotFoundError)
	if !ok {
		t.Fatal("Failed to assert as *MetricNotFoundError")
	}

	if metricErr.Error() != "Metric definition with key 'my-metric' not found" {
		t.Errorf("Unexpected error message: %s", metricErr.Error())
	}
}

func TestErrorTypeAssertion_APIError(t *testing.T) {
	responseData := map[string]interface{}{
		"detail": "Bad request",
	}
	err := NewAPIError("Request failed", 400, responseData)

	apiErr, ok := interface{}(err).(*APIError)
	if !ok {
		t.Fatal("Failed to assert as *APIError")
	}

	if apiErr.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", apiErr.StatusCode)
	}
}

func TestErrorMessages_WithSpecialCharacters(t *testing.T) {
	// Test with special characters in keys
	testCases := []struct {
		key      string
		expected string
	}{
		{"my-prompt-123", "Prompt config with key 'my-prompt-123' not found"},
		{"prompt_with_underscore", "Prompt config with key 'prompt_with_underscore' not found"},
		{"prompt.with.dots", "Prompt config with key 'prompt.with.dots' not found"},
		{"prompt with spaces", "Prompt config with key 'prompt with spaces' not found"},
	}

	for _, tc := range testCases {
		err := NewPromptNotFoundError(tc.key)
		if err.Error() != tc.expected {
			t.Errorf("For key '%s', expected '%s', got '%s'", tc.key, tc.expected, err.Error())
		}
	}
}

func TestAPIError_ComplexResponseData(t *testing.T) {
	responseData := map[string]interface{}{
		"detail": "Validation error",
		"errors": []interface{}{
			map[string]interface{}{
				"field":   "prompt_key",
				"message": "Required field",
			},
			map[string]interface{}{
				"field":   "variables",
				"message": "Invalid format",
			},
		},
		"code":       "VALIDATION_ERROR",
		"request_id": "req-123",
	}

	err := NewAPIError("Validation failed", 422, responseData)

	if err.ResponseData["code"] != "VALIDATION_ERROR" {
		t.Errorf("Expected code 'VALIDATION_ERROR', got %v", err.ResponseData["code"])
	}

	errors, ok := err.ResponseData["errors"].([]interface{})
	if !ok {
		t.Fatal("Expected errors to be a slice")
	}

	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}
}
