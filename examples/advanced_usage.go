package main

import (
	"fmt"
	"log"
	"os"

	"github.com/switchport/switchport-go/switchport"
)

// Advanced usage example for Switchport Go SDK.
//
// This example demonstrates:
// 1. Using context for deterministic A/B testing
// 2. Testing different context values get different versions
// 3. Recording metrics linked to specific versions via context

func main() {
	// For local development
	os.Setenv("SWITCHPORT_API_URL", "http://localhost:8001")

	client, err := switchport.NewClient("")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	abTestExample(client)
	metricsAggregationExample(client)
	stringContextExample(client)

	fmt.Println("\n============================================================")
	fmt.Println("Advanced examples completed!")
	fmt.Println("============================================================")
}

// abTestExample demonstrates deterministic A/B testing where same context always
// gets the same version.
func abTestExample(client *switchport.Client) {
	fmt.Println("============================================================")
	fmt.Println("A/B Testing Example")
	fmt.Println("============================================================")

	// Test with multiple users
	users := []string{"user_001", "user_002", "user_003", "user_004", "user_005"}

	fmt.Println("\nExecuting prompt for 5 different users...")
	fmt.Println("(Same user should always get the same version)\n")

	userVersions := make(map[string]string)

	for _, userID := range users {
		response, err := client.Prompts.Execute(
			"product-pitch",
			map[string]interface{}{
				"user_id": userID,
			},
			map[string]interface{}{
				"product": "Premium Plan",
			},
		)
		if err != nil {
			fmt.Printf("✗ %s: Error - %v\n", userID, err)
		} else {
			userVersions[userID] = response.VersionName
			fmt.Printf("✓ %s: Version %s\n", userID, response.VersionName)
		}
	}

	// Test consistency: run again and verify same versions
	fmt.Println("\nTesting consistency (running again for same users)...")
	consistent := true

	for _, userID := range users {
		response, err := client.Prompts.Execute(
			"product-pitch",
			map[string]interface{}{
				"user_id": userID,
			},
			nil,
		)
		if err != nil {
			fmt.Printf("✗ %s: Error - %v\n", userID, err)
		} else {
			if response.VersionName != userVersions[userID] {
				fmt.Printf("✗ %s: Got different version! (was %s, now %s)\n", userID, userVersions[userID], response.VersionName)
				consistent = false
			} else {
				fmt.Printf("✓ %s: Same version (%s)\n", userID, response.VersionName)
			}
		}
	}

	if consistent {
		fmt.Println("\n✓ All users got consistent versions!")
	} else {
		fmt.Println("\n✗ Version assignment was inconsistent")
	}
}

// metricsAggregationExample demonstrates recording metrics for different users and how
// they're aggregated per version.
func metricsAggregationExample(client *switchport.Client) {
	fmt.Println("\n============================================================")
	fmt.Println("Metrics Aggregation Example")
	fmt.Println("============================================================")

	fmt.Println("\nRecording satisfaction scores for multiple users...")

	// Simulate user feedback for different users
	type FeedbackData struct {
		UserID string
		Score  float64
	}

	feedbackData := []FeedbackData{
		{"user_001", 5.0},
		{"user_002", 4.5},
		{"user_003", 4.0},
		{"user_004", 4.8},
		{"user_005", 3.5},
		{"user_001", 4.7}, // Same user, different session
		{"user_002", 4.9},
	}

	for _, feedback := range feedbackData {
		result, err := client.Metrics.Record(
			"user_satisfaction",
			feedback.Score,
			map[string]interface{}{
				"user_id": feedback.UserID,
			},
			nil,
		)
		if err != nil {
			fmt.Printf("✗ Error recording metric for %s: %v\n", feedback.UserID, err)
		} else {
			fmt.Printf("✓ Recorded score %.1f for %s (event: %s)\n", feedback.Score, feedback.UserID, result.MetricEventID)
		}
	}

	fmt.Println("\n✓ Metrics recorded! They will be aggregated per version based on context.")
	fmt.Println("  Check the Switchport dashboard to see aggregated results per version.")
}

// stringContextExample demonstrates using string context instead of map.
func stringContextExample(client *switchport.Client) {
	fmt.Println("\n============================================================")
	fmt.Println("String Context Example")
	fmt.Println("============================================================")

	fmt.Println("\nUsing simple string as context...")

	// Context as string
	response, err := client.Prompts.Execute(
		"greeting",
		"guest_user", // Simple string context
		nil,
	)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	} else {
		fmt.Println("✓ Prompt executed with string context")
		fmt.Printf("  Version: %s\n", response.VersionName)
	}

	// Context as map
	response, err = client.Prompts.Execute(
		"greeting",
		map[string]interface{}{
			"user_type": "guest",
		},
		nil,
	)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	} else {
		fmt.Println("✓ Prompt executed with map context")
		fmt.Printf("  Version: %s\n", response.VersionName)
	}
}
