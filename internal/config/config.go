package config

import (
	"os"
)

type Config struct {
	APIURL               string
	CanvasToken          string
	ProjectID            string
	SubscriptionID       string
	BucketName           string
	TerraformModulesPath string
}

func Load() (*Config, error) {
	cfg := &Config{
		APIURL:               getEnvOrDefault("API_URL", "https://canvas-api.dev.r7ai.net"),
		CanvasToken:          getEnvOrDefault("CANVAS_TOKEN", "foobar"),
		ProjectID:            getEnvOrDefault("GOOGLE_CLOUD_PROJECT", "rad-dev-dev"),
		SubscriptionID:       getEnvOrDefault("PUBSUB_SUBSCRIPTION_ID", "provisioner"),
		BucketName:           getEnvOrDefault("BUCKET_NAME", "rad-provisioner-state-1234"),
		TerraformModulesPath: getEnvOrDefault("TERRAFORM_MODULES_PATH", "/mnt/canvas-packages"),
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
