package switchport

import "time"

// Subject can be either a map or a string for subject identification
type Subject interface{}

// PromptResponse represents the response from executing a prompt.
type PromptResponse struct {
	// Text is the generated prompt text from the LLM
	Text string `json:"text"`
	// Model is the model used (e.g., 'gpt-4', 'claude-3-5-sonnet-20241022', 'gemini-1.5-pro')
	Model string `json:"model"`
	// VersionID is the UUID of the prompt version used
	VersionID string `json:"version_id"`
	// VersionName is the name of the version (e.g., 'v1', 'v2')
	VersionName string `json:"version_name"`
	// PromptConfigID is the UUID of the prompt config
	PromptConfigID string `json:"prompt_config_id"`
	// RequestID is the unique request ID for tracking
	RequestID string `json:"request_id"`
	// Metadata contains additional metadata (temperature, max_tokens, etc.)
	Metadata map[string]interface{} `json:"metadata"`
}

// MetricRecordResponse represents the response from recording a metric.
type MetricRecordResponse struct {
	// Success indicates whether the metric was recorded successfully
	Success bool `json:"success"`
	// MetricEventID is the UUID of the metric event
	MetricEventID string `json:"metric_event_id"`
}

// promptExecuteRequest represents the request payload for executing a prompt.
type promptExecuteRequest struct {
	PromptKey string                 `json:"prompt_key"`
	Subject   map[string]interface{} `json:"subject"`
	Variables map[string]interface{} `json:"variables"`
}

// metricRecordRequest represents the request payload for recording a metric.
type metricRecordRequest struct {
	MetricKey string                 `json:"metric_key"`
	Value     interface{}            `json:"value"`
	Subject   map[string]interface{} `json:"subject"`
	Timestamp *string                `json:"timestamp,omitempty"`
}

// normalizeSubject converts Subject to a map for JSON serialization.
func normalizeSubject(subject Subject) map[string]interface{} {
	if subject == nil {
		return make(map[string]interface{})
	}

	switch v := subject.(type) {
	case map[string]interface{}:
		return v
	case string:
		return map[string]interface{}{"_subject": v}
	default:
		return make(map[string]interface{})
	}
}

// formatTimestamp formats a time.Time to ISO 8601 format.
func formatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}
