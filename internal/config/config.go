package config

import (
	"os"
)

type Config struct {
	ProjectID      string
	SubscriptionID string
	BucketName     string
}

func Load() (*Config, error) {
	cfg := &Config{
		ProjectID:      getEnvOrDefault("GOOGLE_CLOUD_PROJECT", "default-project-id"),
		SubscriptionID: getEnvOrDefault("PUBSUB_SUBSCRIPTION_ID", "default-subscription-id"),
		BucketName:     getEnvOrDefault("BUCKET_NAME", "rad-provisioner-state-1234"),
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
