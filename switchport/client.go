// Package switchport provides a Go client library for Switchport - Prompt management and A/B testing platform.
package switchport

import (
	"fmt"
	"os"
	"strings"
)

const (
	// DefaultBaseURL is the production API URL
	DefaultBaseURL = "https://switchport-api.vercel.app"
	// Version is the SDK version
	Version = "0.1.0"
)

// Client is the main Switchport SDK client.
//
// Example:
//
//	import "github.com/switchport/switchport-go/switchport"
//
//	client, err := switchport.NewClient("")  // Reads from SWITCHPORT_API_KEY env var
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	response, err := client.Prompts.Execute("welcome-message", nil, map[string]interface{}{
//	    "name": "Alice",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Println(response.Text)
type Client struct {
	// APIKey is the Switchport API key (starts with 'sp_')
	APIKey string
	// BaseURL is the Switchport API base URL
	BaseURL string
	// Prompts is the client for executing prompts
	Prompts *PromptsClient
	// Metrics is the client for recording metrics
	Metrics *MetricsClient
}

// NewClient creates a new Switchport client.
//
// Args:
//   - apiKey: Switchport API key (starts with 'sp_').
//     If empty string is provided, reads from SWITCHPORT_API_KEY environment variable.
//
// Returns:
//   - *Client: The initialized Switchport client
//   - error: Error if no API key is provided or found in environment
//
// Example:
//
//	// Using API key from environment variable
//	client, err := switchport.NewClient("")
//
//	// Using API key directly
//	client, err := switchport.NewClient("sp_xxx")
func NewClient(apiKey string) (*Client, error) {
	// Get API key from parameter or environment
	if apiKey == "" {
		apiKey = os.Getenv("SWITCHPORT_API_KEY")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key is required. Provide it as a parameter or set SWITCHPORT_API_KEY environment variable")
	}

	if !strings.HasPrefix(apiKey, "sp_") {
		return nil, fmt.Errorf("invalid API key format. Switchport API keys start with 'sp_'")
	}

	// Allow base URL override for development/testing
	baseURL := os.Getenv("SWITCHPORT_API_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	// Create client
	client := &Client{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}

	// Initialize sub-clients
	client.Prompts = newPromptsClient(apiKey, baseURL)
	client.Metrics = newMetricsClient(apiKey, baseURL)

	return client, nil
}
