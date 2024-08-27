package deployer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/internal/executors/terraform"
	"github.com/radiatus-ai/package-provisioner/pkg/models"
)

type Deployer struct {
	cfg *config.Config
	// executor *terraform.Executor
	executor terraform.ExecutorInterface
}

func NewDeployer(cfg *config.Config) *Deployer {
	return &Deployer{
		cfg:      cfg,
		executor: terraform.NewExecutor(cfg),
	}
}

func (d *Deployer) DeployPackage(msg models.DeploymentMessage) error {
	log.Printf("Starting deployment for package %s in project %s", msg.PackageID, msg.ProjectID)
	var startDeployData = map[string]interface{}{}
	if err := d.executor.PostOutputToAPI(msg.ProjectID, msg.PackageID, startDeployData, models.StartDeploy); err != nil {
		return fmt.Errorf("failed to post to api: %v", err)
	}

	deployDir := filepath.Join("deployments", msg.PackageID)
	if err := os.MkdirAll(deployDir, 0755); err != nil {
		return fmt.Errorf("failed to create deployment directory: %v", err)
	}

	if err := d.executor.CopyTerraformModules(msg.Package.Type, deployDir); err != nil {
		return fmt.Errorf("failed to copy terraform modules: %v", err)
	}

	if err := d.executor.CreateParameterFile(msg, deployDir); err != nil {
		return fmt.Errorf("failed to create parameter file: %v", err)
	}

	if err := d.executor.CreateSecretsFile(msg, deployDir); err != nil {
		return fmt.Errorf("failed to create secrets file: %v", err)
	}

	if err := d.executor.CreateBackendFile(msg, deployDir); err != nil {
		return fmt.Errorf("failed to create backend file: %v", err)
	}

	if err := d.executor.RunTerraformCommands(deployDir, msg.Action); err != nil {
		return fmt.Errorf("failed to run terraform commands: %v", err)
	}

	outputData, err := d.executor.ProcessTerraformOutputs(msg, deployDir)
	if err != nil {
		return fmt.Errorf("failed to process terraform outputs: %v", err)
	}

	if err := d.executor.WriteOutputFile(msg.PackageID, deployDir, outputData); err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	if err := d.executor.PostOutputToAPI(msg.ProjectID, msg.PackageID, outputData, models.Deployed); err != nil {
		return fmt.Errorf("failed to post to api: %v", err)
	}

	log.Printf("%s completed successfully for package %s in project %s", msg.Action, msg.PackageID, msg.ProjectID)
	return nil
}

func (d *Deployer) PostOutputToAPI(projectID string, packageID string, outputData map[string]interface{}, action models.DeployStatus) error {
	return d.executor.PostOutputToAPI(projectID, packageID, outputData, action)
}
