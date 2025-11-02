package main

import (
	"fmt"
	"log"
	"os"

	"github.com/switchport/switchport-go/switchport"
)

// Basic usage example for Switchport Go SDK.
//
// This example demonstrates:
// 1. Initializing the Switchport client
// 2. Executing a prompt
// 3. Recording a metric

func main() {
	// Initialize client with API key from environment variable
	// Set SWITCHPORT_API_KEY=sp_xxx in your environment
	// For local development, also set SWITCHPORT_API_URL=http://localhost:8001
	client, err := switchport.NewClient("")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Or pass API key directly
	// client, err := switchport.NewClient("sp_your_key_here")

	fmt.Println("============================================================")
	fmt.Println("Switchport SDK - Basic Usage Example")
	fmt.Println("============================================================")

	// Example 1: Execute a simple prompt
	fmt.Println("\n1. Executing prompt without user...")
	response, err := client.Prompts.Execute(
		"welcome-message",
		nil,
		map[string]interface{}{
			"customer_name": "Alice",
			"product":       "Pro Plan",
		},
	)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	} else {
		fmt.Println("✓ Prompt executed successfully!")
		fmt.Printf("  Model: %s\n", response.Model)
		fmt.Printf("  Version: %s\n", response.VersionName)
		fmt.Printf("  Request ID: %s\n", response.RequestID)
		textPreview := response.Text
		if len(textPreview) > 200 {
			textPreview = textPreview[:200] + "..."
		}
		fmt.Printf("  Generated text:\n  %s\n", textPreview)
	}

	// Example 2: Execute prompt with user for A/B testing
	fmt.Println("\n2. Executing prompt with user for A/B testing...")
	response, err = client.Prompts.Execute(
		"product-description",
		map[string]interface{}{
			"user_id": "user_123",
			"tier":    "premium",
		},
		map[string]interface{}{
			"product_name": "Enterprise Widget",
		},
	)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	} else {
		fmt.Println("✓ Prompt executed successfully!")
		fmt.Printf("  Version: %s\n", response.VersionName)
		fmt.Printf("  Request ID: %s\n", response.RequestID)
	}

	// Example 3: Record a metric
	fmt.Println("\n3. Recording a metric...")
	result, err := client.Metrics.Record(
		"user_satisfaction",
		4.5,
		map[string]interface{}{
			"user_id": "user_123",
			"tier":    "premium",
		},
		nil,
	)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	} else {
		fmt.Println("✓ Metric recorded successfully!")
		fmt.Printf("  Event ID: %s\n", result.MetricEventID)
	}

	// Example 4: Record different types of metrics
	fmt.Println("\n4. Recording different metric types...")

	// Float metric
	_, err = client.Metrics.Record(
		"response_time_ms",
		125.5,
		map[string]interface{}{
			"endpoint": "/api/users",
		},
		nil,
	)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	} else {
		fmt.Println("✓ Float metric recorded")
	}

	// Boolean metric
	_, err = client.Metrics.Record(
		"conversion_success",
		true,
		map[string]interface{}{
			"campaign_id": "summer_2025",
		},
		nil,
	)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	} else {
		fmt.Println("✓ Boolean metric recorded")
	}

	// Enum metric (string value)
	_, err = client.Metrics.Record(
		"user_sentiment",
		"positive",
		map[string]interface{}{
			"feature": "new_ui",
		},
		nil,
	)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	} else {
		fmt.Println("✓ Enum metric recorded")
	}

	fmt.Println("\n============================================================")
	fmt.Println("Example completed!")
	fmt.Println("============================================================")
}
