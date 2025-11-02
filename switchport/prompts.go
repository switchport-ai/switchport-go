package switchport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PromptsClient is a client for executing prompts via Switchport.
type PromptsClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// newPromptsClient creates a new PromptsClient.
func newPromptsClient(apiKey, baseURL string) *PromptsClient {
	return &PromptsClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute executes a prompt and gets the LLM response.
//
// Args:
//   - promptKey: The unique key of the prompt config
//   - user: User identification for routing (map or string)
//   - variables: Variables to substitute in the prompt template
//
// Returns:
//   - *PromptResponse: Response containing the generated text and metadata
//   - error: PromptNotFoundError if the prompt config doesn't exist, APIError if the API request fails
func (p *PromptsClient) Execute(promptKey string, user User, variables map[string]interface{}) (*PromptResponse, error) {
	url := fmt.Sprintf("%s/api/sdk/v1/prompts/execute", p.baseURL)

	// Normalize user and variables
	normalizedUser := normalizeUser(user)
	if variables == nil {
		variables = make(map[string]interface{})
	}

	// Create request payload
	payload := promptExecuteRequest{
		PromptKey: promptKey,
		User:      normalizedUser,
		Variables: variables,
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := p.client.Do(req)
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
		return nil, NewPromptNotFoundError(promptKey)
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
	var promptResponse PromptResponse
	if err := json.Unmarshal(body, &promptResponse); err != nil {
		return nil, NewAPIError(fmt.Sprintf("Failed to parse response: %v", err), resp.StatusCode, nil)
	}

	return &promptResponse, nil
}
