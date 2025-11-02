package switchport

import (
	"os"
	"testing"
)

func TestNewClient_WithAPIKey(t *testing.T) {
	apiKey := "sp_test123"
	client, err := NewClient(apiKey)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, got %s", apiKey, client.APIKey)
	}

	if client.BaseURL != DefaultBaseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", DefaultBaseURL, client.BaseURL)
	}

	if client.Prompts == nil {
		t.Error("Expected Prompts client to be initialized")
	}

	if client.Metrics == nil {
		t.Error("Expected Metrics client to be initialized")
	}
}

func TestNewClient_WithEnvVar(t *testing.T) {
	apiKey := "sp_env_test123"
	os.Setenv("SWITCHPORT_API_KEY", apiKey)
	defer os.Unsetenv("SWITCHPORT_API_KEY")

	client, err := NewClient("")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, got %s", apiKey, client.APIKey)
	}
}

func TestNewClient_WithCustomBaseURL(t *testing.T) {
	apiKey := "sp_test123"
	customURL := "http://localhost:8000"
	os.Setenv("SWITCHPORT_API_URL", customURL)
	defer os.Unsetenv("SWITCHPORT_API_URL")

	client, err := NewClient(apiKey)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client.BaseURL != customURL {
		t.Errorf("Expected BaseURL to be %s, got %s", customURL, client.BaseURL)
	}
}

func TestNewClient_NoAPIKey(t *testing.T) {
	os.Unsetenv("SWITCHPORT_API_KEY")

	client, err := NewClient("")

	if err == nil {
		t.Fatal("Expected error when no API key provided")
	}

	if client != nil {
		t.Error("Expected client to be nil when no API key provided")
	}

	expectedMsg := "API key is required"
	if err.Error() != "API key is required. Provide it as a parameter or set SWITCHPORT_API_KEY environment variable" {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestNewClient_InvalidAPIKeyFormat(t *testing.T) {
	invalidKeys := []string{
		"invalid_key",
		"sk_test123",
		"api_key_123",
		"sp",
		"",
	}

	for _, key := range invalidKeys {
		// Skip empty string as it's tested separately
		if key == "" {
			continue
		}

		client, err := NewClient(key)

		if err == nil {
			t.Errorf("Expected error for invalid API key '%s', got none", key)
		}

		if client != nil {
			t.Errorf("Expected client to be nil for invalid API key '%s'", key)
		}

		expectedMsg := "invalid API key format"
		if err != nil && err.Error() != "invalid API key format. Switchport API keys start with 'sp_'" {
			t.Errorf("Expected error message to contain '%s' for key '%s', got '%s'", expectedMsg, key, err.Error())
		}
	}
}

func TestNewClient_ValidAPIKeyFormats(t *testing.T) {
	validKeys := []string{
		"sp_test123",
		"sp_prod_abc123",
		"sp_dev_xyz789",
		"sp_1234567890",
	}

	for _, key := range validKeys {
		client, err := NewClient(key)

		if err != nil {
			t.Errorf("Expected no error for valid API key '%s', got %v", key, err)
		}

		if client == nil {
			t.Errorf("Expected client to be created for valid API key '%s'", key)
		}
	}
}

func TestNewClient_SubClientsInitialized(t *testing.T) {
	apiKey := "sp_test123"
	client, err := NewClient(apiKey)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test Prompts client
	if client.Prompts == nil {
		t.Fatal("Prompts client should not be nil")
	}
	if client.Prompts.apiKey != apiKey {
		t.Errorf("Expected Prompts.apiKey to be %s, got %s", apiKey, client.Prompts.apiKey)
	}
	if client.Prompts.baseURL != DefaultBaseURL {
		t.Errorf("Expected Prompts.baseURL to be %s, got %s", DefaultBaseURL, client.Prompts.baseURL)
	}
	if client.Prompts.client == nil {
		t.Error("Prompts HTTP client should not be nil")
	}

	// Test Metrics client
	if client.Metrics == nil {
		t.Fatal("Metrics client should not be nil")
	}
	if client.Metrics.apiKey != apiKey {
		t.Errorf("Expected Metrics.apiKey to be %s, got %s", apiKey, client.Metrics.apiKey)
	}
	if client.Metrics.baseURL != DefaultBaseURL {
		t.Errorf("Expected Metrics.baseURL to be %s, got %s", DefaultBaseURL, client.Metrics.baseURL)
	}
	if client.Metrics.client == nil {
		t.Error("Metrics HTTP client should not be nil")
	}
}
