package terraform

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/pkg/models"
)

func TestExecutor_CreateParameterFile(t *testing.T) {
	cfg := &config.Config{BucketName: "test-bucket"}
	executor := NewExecutor(cfg)

	tempDir, err := os.MkdirTemp("", "terraform-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	msg := models.DeploymentMessage{
		PackageID: "test-package",
		Package: models.Package{
			ParameterData: map[string]interface{}{"param1": "value1"},
		},
		ConnectedInputData: map[string]interface{}{"input1": "value1"},
	}

	err = executor.CreateParameterFile(msg, tempDir)
	if err != nil {
		t.Errorf("CreateParameterFile() error = %v", err)
	}

	// Check if the file was created
	paramFile := filepath.Join(tempDir, "test-package_inputs.auto.tfvars.json")
	if _, err := os.Stat(paramFile); os.IsNotExist(err) {
		t.Errorf("Parameter file was not created")
	}

	// Add more checks here to verify the content of the file
}
