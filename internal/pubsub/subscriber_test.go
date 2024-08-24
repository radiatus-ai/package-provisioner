package pubsub

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/pkg/models"
)

func TestSubscriber_HandlePush(t *testing.T) {
	// Create a test config and subscriber
	cfg := &config.Config{
		ProjectID:      "test-project",
		SubscriptionID: "test-subscription",
	}

	deploymentCount := 0
	testDeployFn := func(msg models.DeploymentMessage) error {
		deploymentCount++
		return nil
	}

	subscriber := NewSubscriber(cfg, testDeployFn)

	// Create a test message
	testMsg := models.DeploymentMessage{
		ProjectID: "test-project",
		Package: models.Package{
			Type: "test-package",
		},
	}
	msgBytes, _ := json.Marshal(testMsg)

	// Create a push request
	pushRequest := struct {
		Message struct {
			Data []byte `json:"data,omitempty"`
			ID   string `json:"id"`
		} `json:"message"`
	}{
		Message: struct {
			Data []byte `json:"data,omitempty"`
			ID   string `json:"id"`
		}{
			Data: msgBytes,
			ID:   "test-message-id",
		},
	}

	body, _ := json.Marshal(pushRequest)

	// Create a test request
	req, err := http.NewRequest("POST", "/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandlePush method
	subscriber.HandlePush(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the deployment function was called
	if deploymentCount != 1 {
		t.Errorf("Expected 1 deployment, got %d", deploymentCount)
	}
}

func TestSubscriber_HandlePush_InvalidMethod(t *testing.T) {
	subscriber := NewSubscriber(&config.Config{}, nil)

	req, err := http.NewRequest("GET", "/push", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	subscriber.HandlePush(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestSubscriber_HandlePush_InvalidBody(t *testing.T) {
	subscriber := NewSubscriber(&config.Config{}, nil)

	req, err := http.NewRequest("POST", "/push", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	subscriber.HandlePush(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
