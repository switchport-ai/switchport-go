package switchport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetricsRecord_Success(t *testing.T) {
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
		var reqBody metricRecordRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if reqBody.MetricKey != "test-metric" {
			t.Errorf("Expected metric_key 'test-metric', got %s", reqBody.MetricKey)
		}

		// Send mock response
		response := MetricRecordResponse{
			Success:       true,
			MetricEventID: "event-123",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)

	response, err := client.Record("test-metric", 42.5, nil, nil)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.MetricEventID != "event-123" {
		t.Errorf("Expected metric_event_id 'event-123', got %s", response.MetricEventID)
	}
}

func TestMetricsRecord_FloatValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody metricRecordRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		if reqBody.Value != 42.5 {
			t.Errorf("Expected value 42.5, got %v", reqBody.Value)
		}

		response := MetricRecordResponse{Success: true, MetricEventID: "event-1"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)
	_, err := client.Record("test-metric", 42.5, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMetricsRecord_IntValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody metricRecordRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		// JSON unmarshaling converts numbers to float64, so we need to check accordingly
		if int(reqBody.Value.(float64)) != 100 {
			t.Errorf("Expected value 100, got %v", reqBody.Value)
		}

		response := MetricRecordResponse{Success: true, MetricEventID: "event-1"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)
	_, err := client.Record("test-metric", 100, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMetricsRecord_BoolValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody metricRecordRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		if reqBody.Value != true {
			t.Errorf("Expected value true, got %v", reqBody.Value)
		}

		response := MetricRecordResponse{Success: true, MetricEventID: "event-1"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)
	_, err := client.Record("test-metric", true, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMetricsRecord_StringValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody metricRecordRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		if reqBody.Value != "success" {
			t.Errorf("Expected value 'success', got %v", reqBody.Value)
		}

		response := MetricRecordResponse{Success: true, MetricEventID: "event-1"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)
	_, err := client.Record("test-metric", "success", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMetricsRecord_WithMapContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody metricRecordRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		if reqBody.Context["user_id"] != "user123" {
			t.Errorf("Expected context user_id 'user123', got %v", reqBody.Context["user_id"])
		}
		if reqBody.Context["region"] != "us-west" {
			t.Errorf("Expected context region 'us-west', got %v", reqBody.Context["region"])
		}

		response := MetricRecordResponse{Success: true, MetricEventID: "event-1"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)

	context := map[string]interface{}{
		"user_id": "user123",
		"region":  "us-west",
	}

	_, err := client.Record("test-metric", 100, context, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMetricsRecord_WithStringContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody metricRecordRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		if reqBody.Context["_context"] != "user123" {
			t.Errorf("Expected context _context 'user123', got %v", reqBody.Context["_context"])
		}

		response := MetricRecordResponse{Success: true, MetricEventID: "event-1"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)
	_, err := client.Record("test-metric", 100, "user123", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMetricsRecord_WithTimestamp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody metricRecordRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		if reqBody.Timestamp == nil {
			t.Error("Expected timestamp to be set")
		}

		// Verify timestamp format (ISO 8601 / RFC3339)
		_, err := time.Parse(time.RFC3339, *reqBody.Timestamp)
		if err != nil {
			t.Errorf("Expected valid RFC3339 timestamp, got parse error: %v", err)
		}

		response := MetricRecordResponse{Success: true, MetricEventID: "event-1"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)

	timestamp := time.Now()
	_, err := client.Record("test-metric", 100, nil, &timestamp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMetricsRecord_WithoutTimestamp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody metricRecordRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		if reqBody.Timestamp != nil {
			t.Errorf("Expected timestamp to be nil, got %v", reqBody.Timestamp)
		}

		response := MetricRecordResponse{Success: true, MetricEventID: "event-1"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)
	_, err := client.Record("test-metric", 100, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMetricsRecord_MetricNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"detail": "Metric not found",
		})
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)
	_, err := client.Record("nonexistent-metric", 100, nil, nil)

	if err == nil {
		t.Fatal("Expected error for nonexistent metric")
	}

	if _, ok := err.(*MetricNotFoundError); !ok {
		t.Errorf("Expected MetricNotFoundError, got %T", err)
	}
}

func TestMetricsRecord_AuthenticationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"detail": "Invalid API key",
		})
	}))
	defer server.Close()

	client := newMetricsClient("sp_invalid", server.URL)
	_, err := client.Record("test-metric", 100, nil, nil)

	if err == nil {
		t.Fatal("Expected authentication error")
	}

	if _, ok := err.(*AuthenticationError); !ok {
		t.Errorf("Expected AuthenticationError, got %T", err)
	}
}

func TestMetricsRecord_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"detail": "Invalid value type",
		})
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)
	_, err := client.Record("test-metric", 100, nil, nil)

	if err == nil {
		t.Fatal("Expected API error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if apiErr.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", apiErr.StatusCode)
	}
}

func TestMetricsRecord_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)
	_, err := client.Record("test-metric", 100, nil, nil)

	if err == nil {
		t.Fatal("Expected error for invalid JSON response")
	}

	if _, ok := err.(*APIError); !ok {
		t.Errorf("Expected APIError for invalid JSON, got %T", err)
	}
}

func TestMetricsRecord_ComplexContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody metricRecordRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		// Verify complex context
		user := reqBody.Context["user"].(map[string]interface{})
		if user["name"] != "Alice" {
			t.Error("Expected nested context user.name to be 'Alice'")
		}

		response := MetricRecordResponse{Success: true, MetricEventID: "event-1"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newMetricsClient("sp_test123", server.URL)

	context := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
			"id":   123,
		},
		"session_id": "sess-456",
	}

	_, err := client.Record("test-metric", 100, context, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
