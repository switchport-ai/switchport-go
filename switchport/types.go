package switchport

import "time"

// User can be either a map or a string for user identification
type User interface{}

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
	User      map[string]interface{} `json:"user"`
	Variables map[string]interface{} `json:"variables"`
}

// metricRecordRequest represents the request payload for recording a metric.
type metricRecordRequest struct {
	MetricKey string                 `json:"metric_key"`
	Value     interface{}            `json:"value"`
	User      map[string]interface{} `json:"user"`
	Timestamp *string                `json:"timestamp,omitempty"`
}

// normalizeUser converts User to a map for JSON serialization.
func normalizeUser(user User) map[string]interface{} {
	if user == nil {
		return make(map[string]interface{})
	}

	switch v := user.(type) {
	case map[string]interface{}:
		return v
	case string:
		return map[string]interface{}{"_user": v}
	default:
		return make(map[string]interface{})
	}
}

// formatTimestamp formats a time.Time to ISO 8601 format.
func formatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}
