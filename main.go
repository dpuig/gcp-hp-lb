package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/oauth2/google"
)

// Simplified representation of your backend instances
type Backend struct {
	URL          *url.URL
	Alive        bool
	HealthStatus string
}

func main() {
	ctx := context.Background()

	// Use google.FindDefaultCredentials. This will use the application default
	// credentials, which are set by the GOOGLE_APPLICATION_CREDENTIALS environment variable.
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatalf("Failed to find default credentials: %v", err)
	}

	// The project ID is part of the credentials.
	projectID := creds.ProjectID
	fmt.Printf("Project ID: %s\n", projectID)

	// The region is not part of the credentials, but can be retrieved from the
	// metadata server if running on GCP.
	region, err := metadata.Get("instance/zone")
	if err != nil {
		log.Fatalf("Failed to get region: %v", err)
	}

	// The region is the last part of the zone name.
	region = region[:len(region)-2]
	fmt.Printf("Region: %s\n", region)
}
