package switchport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPromptsExecute_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer sp_test123" {
			t.Errorf("Expected Authorization header, got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Verify request body
		var reqBody promptExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if reqBody.PromptKey != "test-prompt" {
			t.Errorf("Expected prompt_key 'test-prompt', got %s", reqBody.PromptKey)
		}

		// Send mock response
		response := PromptResponse{
			Text:           "Hello, World!",
			Model:          "gpt-4",
			VersionID:      "version-123",
			VersionName:    "v1",
			PromptConfigID: "config-123",
			RequestID:      "request-123",
			Metadata: map[string]interface{}{
				"temperature": 0.7,
				"max_tokens":  100,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := newPromptsClient("sp_test123", server.URL)

	// Execute prompt
	response, err := client.Execute("test-prompt", nil, map[string]interface{}{
		"name": "World",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Text != "Hello, World!" {
		t.Errorf("Expected text 'Hello, World!', got %s", response.Text)
	}
	if response.Model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got %s", response.Model)
	}
	if response.VersionID != "version-123" {
		t.Errorf("Expected version_id 'version-123', got %s", response.VersionID)
	}
	if response.RequestID != "request-123" {
		t.Errorf("Expected request_id 'request-123', got %s", response.RequestID)
	}
}

func TestPromptsExecute_WithMapContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody promptExecuteRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		// Verify context
		if reqBody.Context["user_id"] != "user123" {
			t.Errorf("Expected context user_id 'user123', got %v", reqBody.Context["user_id"])
		}

		response := PromptResponse{
			Text:      "Response",
			Model:     "gpt-4",
			VersionID: "v1",
			RequestID: "req1",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newPromptsClient("sp_test123", server.URL)

	context := map[string]interface{}{
		"user_id": "user123",
		"region":  "us-west",
	}

	_, err := client.Execute("test-prompt", context, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPromptsExecute_WithStringContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody promptExecuteRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		// Verify context is normalized to map with _context key
		if reqBody.Context["_context"] != "user123" {
			t.Errorf("Expected context _context 'user123', got %v", reqBody.Context["_context"])
		}

		response := PromptResponse{
			Text:      "Response",
			Model:     "gpt-4",
			VersionID: "v1",
			RequestID: "req1",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newPromptsClient("sp_test123", server.URL)

	_, err := client.Execute("test-prompt", "user123", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPromptsExecute_PromptNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"detail": "Prompt config not found",
		})
	}))
	defer server.Close()

	client := newPromptsClient("sp_test123", server.URL)

	_, err := client.Execute("nonexistent-prompt", nil, nil)

	if err == nil {
		t.Fatal("Expected error for nonexistent prompt")
	}

	if _, ok := err.(*PromptNotFoundError); !ok {
		t.Errorf("Expected PromptNotFoundError, got %T", err)
	}
}

func TestPromptsExecute_AuthenticationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"detail": "Invalid API key",
		})
	}))
	defer server.Close()

	client := newPromptsClient("sp_invalid", server.URL)

	_, err := client.Execute("test-prompt", nil, nil)

	if err == nil {
		t.Fatal("Expected authentication error")
	}

	if _, ok := err.(*AuthenticationError); !ok {
		t.Errorf("Expected AuthenticationError, got %T", err)
	}
}

func TestPromptsExecute_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"detail": "Internal server error",
		})
	}))
	defer server.Close()

	client := newPromptsClient("sp_test123", server.URL)

	_, err := client.Execute("test-prompt", nil, nil)

	if err == nil {
		t.Fatal("Expected API error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if apiErr.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", apiErr.StatusCode)
	}
}

func TestPromptsExecute_NilVariables(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody promptExecuteRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		// Verify variables is empty map, not nil
		if reqBody.Variables == nil {
			t.Error("Expected variables to be empty map, got nil")
		}

		response := PromptResponse{
			Text:      "Response",
			Model:     "gpt-4",
			VersionID: "v1",
			RequestID: "req1",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newPromptsClient("sp_test123", server.URL)

	_, err := client.Execute("test-prompt", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPromptsExecute_ComplexVariables(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody promptExecuteRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		// Verify complex variables
		if reqBody.Variables["user"].(map[string]interface{})["name"] != "Alice" {
			t.Error("Expected nested variable user.name to be 'Alice'")
		}
		if reqBody.Variables["tags"].([]interface{})[0] != "important" {
			t.Error("Expected array variable tags[0] to be 'important'")
		}

		response := PromptResponse{
			Text:      "Response",
			Model:     "gpt-4",
			VersionID: "v1",
			RequestID: "req1",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newPromptsClient("sp_test123", server.URL)

	variables := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
		"tags": []string{"important", "urgent"},
	}

	_, err := client.Execute("test-prompt", nil, variables)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPromptsExecute_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := newPromptsClient("sp_test123", server.URL)

	_, err := client.Execute("test-prompt", nil, nil)

	if err == nil {
		t.Fatal("Expected error for invalid JSON response")
	}

	if _, ok := err.(*APIError); !ok {
		t.Errorf("Expected APIError for invalid JSON, got %T", err)
	}
}
