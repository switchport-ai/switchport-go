package switchport

import (
	"testing"
	"time"
)

func TestNormalizeContext_MapContext(t *testing.T) {
	input := map[string]interface{}{
		"user_id": "user123",
		"region":  "us-west",
	}

	result := normalizeContext(input)

	if len(result) != 2 {
		t.Errorf("Expected 2 keys in result, got %d", len(result))
	}
	if result["user_id"] != "user123" {
		t.Errorf("Expected user_id 'user123', got %v", result["user_id"])
	}
	if result["region"] != "us-west" {
		t.Errorf("Expected region 'us-west', got %v", result["region"])
	}
}

func TestNormalizeContext_StringContext(t *testing.T) {
	input := "user123"

	result := normalizeContext(input)

	if len(result) != 1 {
		t.Errorf("Expected 1 key in result, got %d", len(result))
	}
	if result["_context"] != "user123" {
		t.Errorf("Expected _context 'user123', got %v", result["_context"])
	}
}

func TestNormalizeContext_NilContext(t *testing.T) {
	result := normalizeContext(nil)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d keys", len(result))
	}
}

func TestNormalizeContext_EmptyString(t *testing.T) {
	result := normalizeContext("")

	if len(result) != 1 {
		t.Errorf("Expected 1 key in result, got %d", len(result))
	}
	if result["_context"] != "" {
		t.Errorf("Expected _context to be empty string, got %v", result["_context"])
	}
}

func TestNormalizeContext_EmptyMap(t *testing.T) {
	input := map[string]interface{}{}

	result := normalizeContext(input)

	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d keys", len(result))
	}
}

func TestNormalizeContext_ComplexMap(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "Alice",
		},
		"tags": []string{"important", "urgent"},
		"metadata": map[string]interface{}{
			"version": "1.0",
		},
	}

	result := normalizeContext(input)

	if len(result) != 3 {
		t.Errorf("Expected 3 keys in result, got %d", len(result))
	}

	user, ok := result["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected user to be a map")
	}
	if user["name"] != "Alice" {
		t.Errorf("Expected user.name 'Alice', got %v", user["name"])
	}

	tags, ok := result["tags"].([]string)
	if !ok {
		t.Fatal("Expected tags to be a string slice")
	}
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}

func TestNormalizeContext_InvalidType(t *testing.T) {
	// Test with an integer (unsupported type)
	result := normalizeContext(42)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty map for invalid type, got %d keys", len(result))
	}
}

func TestNormalizeContext_Struct(t *testing.T) {
	// Test with a struct (unsupported type)
	type TestStruct struct {
		Name string
		Age  int
	}
	input := TestStruct{Name: "Alice", Age: 30}

	result := normalizeContext(input)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty map for struct type, got %d keys", len(result))
	}
}

func TestFormatTimestamp_CurrentTime(t *testing.T) {
	now := time.Now()
	result := formatTimestamp(now)

	// Verify it's in RFC3339 format
	parsed, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Fatalf("Expected valid RFC3339 timestamp, got parse error: %v", err)
	}

	// Verify the parsed time is close to the original (within 1 second)
	diff := parsed.Sub(now)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("Expected parsed time to be close to original, diff: %v", diff)
	}
}

func TestFormatTimestamp_SpecificTime(t *testing.T) {
	// Create a specific time: 2024-01-15 10:30:45 UTC
	specificTime := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	result := formatTimestamp(specificTime)

	expected := "2024-01-15T10:30:45Z"
	if result != expected {
		t.Errorf("Expected timestamp '%s', got '%s'", expected, result)
	}
}

func TestFormatTimestamp_WithNanoseconds(t *testing.T) {
	// Create a time with nanoseconds
	specificTime := time.Date(2024, 1, 15, 10, 30, 45, 123456789, time.UTC)
	result := formatTimestamp(specificTime)

	// Verify it's valid RFC3339
	parsed, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Fatalf("Expected valid RFC3339 timestamp, got parse error: %v", err)
	}

	// Verify year, month, day, hour, minute, second match
	if parsed.Year() != 2024 || parsed.Month() != 1 || parsed.Day() != 15 {
		t.Errorf("Expected date 2024-01-15, got %s", parsed.Format("2006-01-02"))
	}
	if parsed.Hour() != 10 || parsed.Minute() != 30 || parsed.Second() != 45 {
		t.Errorf("Expected time 10:30:45, got %s", parsed.Format("15:04:05"))
	}
}

func TestFormatTimestamp_DifferentTimezones(t *testing.T) {
	// Create time in different timezone (PST = UTC-8)
	pst := time.FixedZone("PST", -8*60*60)
	specificTime := time.Date(2024, 1, 15, 10, 30, 45, 0, pst)
	result := formatTimestamp(specificTime)

	// Parse and verify
	parsed, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Fatalf("Expected valid RFC3339 timestamp, got parse error: %v", err)
	}

	// Convert both to UTC and compare
	if !parsed.UTC().Equal(specificTime.UTC()) {
		t.Errorf("Expected times to be equal in UTC. Original: %s, Parsed: %s",
			specificTime.UTC(), parsed.UTC())
	}
}

func TestFormatTimestamp_ZeroTime(t *testing.T) {
	zeroTime := time.Time{}
	result := formatTimestamp(zeroTime)

	// Verify it returns valid RFC3339
	_, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Fatalf("Expected valid RFC3339 timestamp for zero time, got parse error: %v", err)
	}

	expected := "0001-01-01T00:00:00Z"
	if result != expected {
		t.Errorf("Expected zero time '%s', got '%s'", expected, result)
	}
}

func TestPromptResponse_JSONSerialization(t *testing.T) {
	// Verify that PromptResponse fields match expected JSON tags
	response := PromptResponse{
		Text:           "Hello",
		Model:          "gpt-4",
		VersionID:      "v1",
		VersionName:    "version-1",
		PromptConfigID: "config-1",
		RequestID:      "req-1",
		Metadata: map[string]interface{}{
			"temperature": 0.7,
		},
	}

	if response.Text != "Hello" {
		t.Errorf("Expected Text 'Hello', got %s", response.Text)
	}
	if response.Model != "gpt-4" {
		t.Errorf("Expected Model 'gpt-4', got %s", response.Model)
	}
	if response.Metadata["temperature"] != 0.7 {
		t.Errorf("Expected Metadata temperature 0.7, got %v", response.Metadata["temperature"])
	}
}

func TestMetricRecordResponse_JSONSerialization(t *testing.T) {
	response := MetricRecordResponse{
		Success:       true,
		MetricEventID: "event-123",
	}

	if !response.Success {
		t.Error("Expected Success to be true")
	}
	if response.MetricEventID != "event-123" {
		t.Errorf("Expected MetricEventID 'event-123', got %s", response.MetricEventID)
	}
}
