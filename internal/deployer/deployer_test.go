package deployer

import (
	"testing"

	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/internal/terraform"
	"github.com/radiatus-ai/package-provisioner/pkg/models"
	"github.com/spf13/afero"
)

var _ terraform.ExecutorInterface = (*MockExecutor)(nil)

// MockExecutor is a mock implementation of the terraform.Executor
type MockExecutor struct {
	Fs afero.Fs
}

func (m *MockExecutor) CopyTerraformModules(packageType, deployDir string) error {
	return nil // Mock implementation
}

func (m *MockExecutor) CreateParameterFile(msg models.DeploymentMessage, deployDir string) error {
	return afero.WriteFile(m.Fs, "deployments/test-package/parameters.tfvars", []byte("mocked parameters"), 0644)
}

func (m *MockExecutor) CreateBackendFile(msg models.DeploymentMessage, deployDir string) error {
	return afero.WriteFile(m.Fs, "deployments/test-package/backend.tf", []byte("mocked backend"), 0644)
}

func (m *MockExecutor) RunTerraformCommands(deployDir string) error {
	return nil // Mock implementation
}

func (m *MockExecutor) ProcessTerraformOutputs(msg models.DeploymentMessage, deployDir string) (map[string]interface{}, error) {
	return map[string]interface{}{"output1": "value1"}, nil
}

func (m *MockExecutor) WriteOutputFile(packageID, deployDir string, outputData map[string]interface{}) error {
	return afero.WriteFile(m.Fs, "deployments/test-package/output.json", []byte("mocked output"), 0644)
}

func TestDeployer_DeployPackage(t *testing.T) {
	// Create a mock filesystem
	mockFs := afero.NewMemMapFs()

	cfg := &config.Config{
		ProjectID:      "test-project",
		SubscriptionID: "test-subscription",
		BucketName:     "test-bucket",
	}

	// Create a mock executor that uses the mock filesystem
	mockExecutor := &MockExecutor{
		Fs: mockFs,
	}

	deployer := &Deployer{
		cfg:      cfg,
		executor: mockExecutor,
	}

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

	// Add assertions to check if the deployment was successful
	// For example, check if files were created in the mock filesystem
	exists, _ := afero.Exists(mockFs, "deployments/test-package/parameters.tfvars")
	if !exists {
		t.Errorf("Expected parameters.tfvars file to be created")
	}

	exists, _ = afero.Exists(mockFs, "deployments/test-package/backend.tf")
	if !exists {
		t.Errorf("Expected backend.tf file to be created")
	}

	exists, _ = afero.Exists(mockFs, "deployments/test-package/output.json")
	if !exists {
		t.Errorf("Expected output.json file to be created")
	}

	// Add more assertions as needed
}
