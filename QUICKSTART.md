# Switchport Go SDK - Quick Start Guide

Get started with the Switchport Go SDK in 5 minutes!

## Installation

```bash
go get github.com/switchport/switchport-go
```

## Setup

### 1. Get Your API Key

1. Sign up at [switchport.ai](https://switchport.ai)
2. Create an organization and project
3. Navigate to Settings
4. Copy your API key (starts with `sp_`)

### 2. Set Environment Variable

```bash
export SWITCHPORT_API_KEY=sp_your_key_here
```

Or create a `.env` file:

```
SWITCHPORT_API_KEY=sp_your_key_here
```

## Basic Usage

### Initialize the Client

```go
package main

import (
    "log"
    "github.com/switchport/switchport-go/switchport"
)

func main() {
    // Reads API key from SWITCHPORT_API_KEY environment variable
    client, err := switchport.NewClient("")
    if err != nil {
        log.Fatal(err)
    }

    // Or pass the API key directly
    // client, err := switchport.NewClient("sp_your_key_here")
}
```

### Execute a Prompt

```go
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

fmt.Println(response.Text)
```

### Record a Metric

```go
result, err := client.Metrics.Record(
    "user_satisfaction",
    4.5,
    map[string]interface{}{
        "user_id": "user_123",
    },
    nil,
)

if err != nil {
    log.Fatal(err)
}

fmt.Println("Metric recorded:", result.MetricEventID)
```

## A/B Testing

Use context to enable deterministic A/B testing:

```go
// Same user_id always gets the same version
response, err := client.Prompts.Execute(
    "product-pitch",
    map[string]interface{}{
        "user_id": "user_123",
    },
    nil,
)

// Track metrics for this version
client.Metrics.Record(
    "conversion_rate",
    true,
    map[string]interface{}{
        "user_id": "user_123",
    },
    nil,
)
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/switchport/switchport-go/switchport"
)

func main() {
    // Initialize client
    client, err := switchport.NewClient("")
    if err != nil {
        log.Fatal(err)
    }

    // Execute prompt
    response, err := client.Prompts.Execute(
        "welcome-message",
        map[string]interface{}{
            "user_id": "user_123",
        },
        map[string]interface{}{
            "customer_name": "Alice",
        },
    )

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", response.Text)
    fmt.Println("Version:", response.VersionName)

    // Record metric
    result, err := client.Metrics.Record(
        "user_satisfaction",
        5.0,
        map[string]interface{}{
            "user_id": "user_123",
        },
        nil,
    )

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Metric recorded:", result.MetricEventID)
}
```

## Next Steps

- Check out the [examples](examples/) directory for more examples
- Read the [README](README.md) for full API reference
- Visit [docs.switchport.ai](https://docs.switchport.ai) for detailed documentation

## Local Development

For local development with the Switchport API running on localhost:

```bash
export SWITCHPORT_API_KEY=sp_your_key_here
export SWITCHPORT_API_URL=http://localhost:8001
```

## Need Help?

- Documentation: [https://docs.switchport.ai](https://docs.switchport.ai)
- Issues: [GitHub Issues](https://github.com/switchport/switchport-go/issues)
- Email: support@switchport.ai
