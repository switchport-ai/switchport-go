# Switchport Go SDK

Official Go SDK for [Switchport](https://switchport.ai) - Prompt management and A/B testing platform.

## Features

- 🚀 **Easy Integration**: Simple, intuitive API
- 🎯 **Prompt Execution**: Call LLMs with managed prompts
- 📊 **A/B Testing**: Deterministic version routing based on context
- 📈 **Metrics Recording**: Track performance and user feedback
- 🔐 **Secure**: API key authentication
- 🎨 **Flexible Context**: Support for map or string context
- ⚡ **Zero Dependencies**: Uses only Go standard library

## Installation

```bash
go get github.com/switchport/switchport-go
```

## Quick Start

### 1. Get Your API Key

1. Sign up at [switchport.ai](https://switchport.ai)
2. Create an organization and project
3. Get your API key from Settings (starts with `sp_`)

### 2. Set Environment Variable

```bash
export SWITCHPORT_API_KEY=sp_your_key_here
```

### 3. Use the SDK

```go
package main

import (
    "fmt"
    "log"

    "github.com/switchport/switchport-go/switchport"
)

func main() {
    // Initialize client (reads API key from environment)
    client, err := switchport.NewClient("")
    if err != nil {
        log.Fatal(err)
    }

    // Execute a prompt
    response, err := client.Prompts.Execute(
        "welcome-message",
        nil,
        map[string]interface{}{
            "customer_name": "Alice",
        },
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response.Text) // Generated text from LLM
}
```

## API Reference

### Client Initialization

```go
// Using environment variable (SWITCHPORT_API_KEY)
client, err := switchport.NewClient("")

// Using API key directly
client, err := switchport.NewClient("sp_your_key_here")
```

### Executing Prompts

```go
response, err := client.Prompts.Execute(
    promptKey string,
    context switchport.Context,      // Can be map[string]interface{} or string
    variables map[string]interface{}, // Optional variables for template substitution
)
```

**Example:**

```go
response, err := client.Prompts.Execute(
    "product-description",
    map[string]interface{}{
        "user_id": "user_123",
        "tier": "premium",
    },
    map[string]interface{}{
        "product_name": "Enterprise Widget",
    },
)

if err != nil {
    log.Fatal(err)
}

fmt.Println(response.Text)
fmt.Println(response.Model)
fmt.Println(response.VersionName)
```

### Recording Metrics

```go
result, err := client.Metrics.Record(
    metricKey string,
    value interface{},              // Can be float64, int, bool, or string
    context switchport.Context,     // Can be map[string]interface{} or string
    timestamp *time.Time,           // Optional timestamp (nil = now)
)
```

**Example:**

```go
import "time"

// Record a float metric
result, err := client.Metrics.Record(
    "user_satisfaction",
    4.5,
    map[string]interface{}{
        "user_id": "user_123",
    },
    nil,
)

// Record a boolean metric
result, err := client.Metrics.Record(
    "conversion_success",
    true,
    map[string]interface{}{
        "campaign_id": "summer_2025",
    },
    nil,
)

// Record with timestamp
now := time.Now()
result, err := client.Metrics.Record(
    "response_time_ms",
    125.5,
    nil,
    &now,
)
```

## Error Handling

The SDK provides custom error types for different scenarios:

```go
import "errors"

response, err := client.Prompts.Execute("my-prompt", nil, nil)
if err != nil {
    var promptNotFound *switchport.PromptNotFoundError
    var authError *switchport.AuthenticationError
    var apiError *switchport.APIError

    if errors.As(err, &promptNotFound) {
        // Handle prompt not found
        fmt.Println("Prompt not found")
    } else if errors.As(err, &authError) {
        // Handle authentication error
        fmt.Println("Invalid API key")
    } else if errors.As(err, &apiError) {
        // Handle other API errors
        fmt.Printf("API error: %d - %s\n", apiError.StatusCode, apiError.Message)
    }
}
```

### Error Types

- `SwitchportError` - Base error type
- `AuthenticationError` - Invalid API key
- `PromptNotFoundError` - Prompt config not found
- `MetricNotFoundError` - Metric definition not found
- `APIError` - General API errors (includes `StatusCode` and `ResponseData`)

## A/B Testing

Switchport provides deterministic A/B testing where the same context always routes to the same version:

```go
// User will always get the same version based on their user_id
response, err := client.Prompts.Execute(
    "product-pitch",
    map[string]interface{}{
        "user_id": "user_123",
    },
    nil,
)

// Record metrics with the same context to track performance per version
result, err := client.Metrics.Record(
    "conversion_rate",
    true,
    map[string]interface{}{
        "user_id": "user_123",
    },
    nil,
)
```

## Examples

See the [examples](examples/) directory for complete examples:

- [basic_usage.go](examples/basic_usage.go) - Basic prompt execution and metric recording
- [advanced_usage.go](examples/advanced_usage.go) - A/B testing and metrics aggregation

### Running Examples

```bash
# Set your API key
export SWITCHPORT_API_KEY=sp_your_key_here

# For local development
export SWITCHPORT_API_URL=http://localhost:8001

# Run basic example
go run examples/basic_usage.go

# Run advanced example
go run examples/advanced_usage.go
```

## Development

### Local Development Setup

1. Set environment variables:
```bash
export SWITCHPORT_API_KEY=sp_your_key_here
export SWITCHPORT_API_URL=http://localhost:8001  # For local API
```

2. Run examples:
```bash
go run examples/basic_usage.go
```

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

## Requirements

- Go 1.21 or higher
- No external dependencies (uses only Go standard library)

## License

MIT License

## Support

- Documentation: [https://docs.switchport.ai](https://docs.switchport.ai)
- Issues: [GitHub Issues](https://github.com/switchport/switchport-go/issues)
- Email: support@switchport.ai

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request
