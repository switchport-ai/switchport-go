package switchport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MetricsClient is a client for recording metrics via Switchport.
type MetricsClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// newMetricsClient creates a new MetricsClient.
func newMetricsClient(apiKey, baseURL string) *MetricsClient {
	return &MetricsClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Record records a metric value.
//
// Args:
//   - metricKey: The unique key of the metric definition
//   - value: The metric value (float, int, bool, or string for enum)
//   - user: User identification for aggregation (map or string)
//   - timestamp: Optional timestamp (defaults to now if nil)
//
// Returns:
//   - *MetricRecordResponse: Response with success status and event ID
//   - error: MetricNotFoundError if the metric definition doesn't exist, APIError if the API request fails
func (m *MetricsClient) Record(metricKey string, value interface{}, user User, timestamp *time.Time) (*MetricRecordResponse, error) {
	url := fmt.Sprintf("%s/api/sdk/v1/metrics/record", m.baseURL)

	// Normalize user
	normalizedUser := normalizeUser(user)

	// Format timestamp if provided
	var timestampStr *string
	if timestamp != nil {
		ts := formatTimestamp(*timestamp)
		timestampStr = &ts
	}

	// Create request payload
	payload := metricRecordRequest{
		MetricKey: metricKey,
		Value:     value,
		User:      normalizedUser,
		Timestamp: timestampStr,
	}

	// Marshal payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, NewAPIError(fmt.Sprintf("Failed to marshal request: %v", err), 0, nil)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, NewAPIError(fmt.Sprintf("Failed to create request: %v", err), 0, nil)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, NewAPIError(fmt.Sprintf("Network error: %v", err), 0, nil)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewAPIError(fmt.Sprintf("Failed to read response: %v", err), resp.StatusCode, nil)
	}

	// Handle error responses
	if resp.StatusCode == 404 {
		return nil, NewMetricNotFoundError(metricKey)
	} else if resp.StatusCode == 401 {
		return nil, NewAuthenticationError("Invalid API key")
	} else if resp.StatusCode != 200 {
		var responseData map[string]interface{}
		json.Unmarshal(body, &responseData)
		return nil, NewAPIError(
			fmt.Sprintf("API request failed: %s", string(body)),
			resp.StatusCode,
			responseData,
		)
	}

	// Parse successful response
	var metricResponse MetricRecordResponse
	if err := json.Unmarshal(body, &metricResponse); err != nil {
		return nil, NewAPIError(fmt.Sprintf("Failed to parse response: %v", err), resp.StatusCode, nil)
	}

	return &metricResponse, nil
}
