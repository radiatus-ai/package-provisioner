package deployer

import (
	"testing"

	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/pkg/models"
)

func TestDeployer_DeployPackage(t *testing.T) {
	cfg := &config.Config{
		ProjectID:      "test-project",
		SubscriptionID: "test-subscription",
		BucketName:     "test-bucket",
	}

	deployer := NewDeployer(cfg)

	msg := models.DeploymentMessage{
		ProjectID: "test-project",
		PackageID: "test-package",
		Package: models.Package{
			Type:          "test-type",
			ParameterData: map[string]interface{}{"param1": "value1"},
			Outputs:       map[string]interface{}{"output1": "value1"},
		},
		ConnectedInputData: map[string]interface{}{"input1": "value1"},
	}

	err := deployer.DeployPackage(msg)
	if err != nil {
		t.Errorf("DeployPackage() error = %v", err)
	}

	// Add more assertions here to check if the deployment was successful
	// For example, check if files were created, if Terraform commands were executed, etc.
}
