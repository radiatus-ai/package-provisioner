package terraform

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/pkg/models"
)

type Executor struct {
	cfg *config.Config
}

func NewExecutor(cfg *config.Config) *Executor {
	log.Println("Creating new Terraform Executor")
	return &Executor{cfg: cfg}
}

func (e *Executor) CopyTerraformModules(packageType, deployDir string) error {
	log.Printf("Copying Terraform modules for package type: %s to deploy directory: %s", packageType, deployDir)

	// Get the absolute path of the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	sourceDir := filepath.Join(currentDir, "terraform-modules", packageType)
	log.Printf("Source dir: %s", sourceDir)

	// Ensure the source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourceDir)
	}

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		destPath := filepath.Join(deployDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		return copyFile(path, destPath)
	})
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	return nil
}

func (e *Executor) CreateParameterFile(msg models.DeploymentMessage, deployDir string) error {
	log.Printf("Creating parameter file for package: %s in directory: %s", msg.PackageID, deployDir)
	combinedData := make(map[string]interface{})
	for k, v := range msg.Package.ParameterData {
		combinedData[k] = v
	}
	for k, v := range msg.ConnectedInputData {
		combinedData[k] = v
	}

	filePath := filepath.Join(deployDir, fmt.Sprintf("%s_inputs.auto.tfvars.json", msg.PackageID))
	err := e.writeJSONFile(filePath, combinedData)
	if err != nil {
		log.Printf("Error creating parameter file: %v", err)
	} else {
		log.Printf("Successfully created parameter file: %s", filePath)
	}
	return err
}

func (e *Executor) CreateBackendFile(msg models.DeploymentMessage, deployDir string) error {
	log.Printf("Creating backend file for package: %s in directory: %s", msg.PackageID, deployDir)
	prefix := fmt.Sprintf("projects/%s/packages/%s", msg.ProjectID, msg.PackageID)
	content := fmt.Sprintf(`
terraform {
  backend "gcs" {
    bucket = "%s"
    prefix = "%s"
  }
}
`, e.cfg.BucketName, prefix)

	filePath := filepath.Join(deployDir, "backend.tf")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Printf("Error creating backend file: %v", err)
	} else {
		log.Printf("Successfully created backend file: %s", filePath)
	}
	return err
}

func (e *Executor) RunTerraformCommands(deployDir string) error {
	log.Printf("Running Terraform commands in directory: %s", deployDir)
	commands := []string{
		"terraform init",
		"terraform plan",
		"terraform apply -auto-approve",
	}

	for _, cmd := range commands {
		log.Printf("Executing command: %s", cmd)
		output, err := e.runCommand(cmd, deployDir)
		if err != nil {
			log.Printf("Command '%s' failed: %v\nOutput: %s", cmd, err, output)
			return fmt.Errorf("command '%s' failed: %v\nOutput: %s", cmd, err, output)
		}
		log.Printf("Command '%s' executed successfully", cmd)
	}

	log.Println("All Terraform commands executed successfully")
	return nil
}

func (e *Executor) ProcessTerraformOutputs(msg models.DeploymentMessage, deployDir string) (map[string]interface{}, error) {
	log.Printf("Processing Terraform outputs for package: %s in directory: %s", msg.PackageID, deployDir)
	output, err := e.runCommand("terraform output -json", deployDir)
	if err != nil {
		log.Printf("Failed to get Terraform outputs: %v", err)
		return nil, fmt.Errorf("failed to get terraform outputs: %v", err)
	}

	var outputJSON map[string]interface{}
	if err := json.Unmarshal([]byte(output), &outputJSON); err != nil {
		log.Printf("Failed to parse Terraform outputs: %v", err)
		return nil, fmt.Errorf("failed to parse terraform outputs: %v", err)
	}

	processedOutput := make(map[string]interface{})
	for k, v := range outputJSON {
		if m, ok := v.(map[string]interface{}); ok {
			processedOutput[k] = m["value"]
		} else {
			processedOutput[k] = v
		}
	}

	outputData := make(map[string]interface{})
	for k := range msg.Package.Outputs {
		if v, ok := processedOutput[k]; ok {
			outputData[k] = v
		}
	}

	log.Printf("Processed %d Terraform outputs", len(outputData))
	return outputData, nil
}

func (e *Executor) WriteOutputFile(packageID, deployDir string, outputData map[string]interface{}) error {
	log.Printf("Writing output file for package: %s in directory: %s", packageID, deployDir)
	filePath := filepath.Join(deployDir, fmt.Sprintf("%s_output.json", packageID))
	err := e.writeJSONFile(filePath, outputData)
	if err != nil {
		log.Printf("Error writing output file: %v", err)
	} else {
		log.Printf("Successfully wrote output file: %s", filePath)
	}
	return err
}

func (e *Executor) runCommand(command, dir string) (string, error) {
	log.Printf("Running command: %s in directory: %s", command, dir)
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command failed: %v\nOutput: %s", err, string(output))
	} else {
		log.Printf("Command executed successfully")
	}
	return string(output), err
}

func (e *Executor) writeJSONFile(filepath string, data interface{}) error {
	log.Printf("Writing JSON file: %s", filepath)
	file, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
	} else {
		log.Printf("Successfully wrote JSON file: %s", filepath)
	}
	return err
}
