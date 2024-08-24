// cmd/deployer/main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/internal/deployer"
	"github.com/radiatus-ai/package-provisioner/internal/pubsub"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded successfully: %+v", cfg)

	deployer := deployer.NewDeployer(cfg)
	log.Printf("Deployer initialized")

	subscriber := pubsub.NewSubscriber(cfg, deployer.DeployPackage)
	log.Printf("Subscriber initialized")

	// Set up HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		fmt.Fprintf(w, "Deployer is running")
	})

	// Add the push endpoint
	http.HandleFunc("/push", subscriber.HandlePush)

	// Get PORT from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	log.Printf("Using port: %s", port)

	log.Printf("Starting HTTP server on port %s", port)
	log.Printf("Listening for messages on projects/%s/subscriptions/%s", cfg.ProjectID, cfg.SubscriptionID)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
