package terraform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/pkg/models"
)

type ExecutorInterface interface {
	CopyTerraformModules(packageType, deployDir string) error
	CreateParameterFile(msg models.DeploymentMessage, deployDir string) error
	CreateSecretsFile(msg models.DeploymentMessage, deployDir string) error
	CreateBackendFile(msg models.DeploymentMessage, deployDir string) error
	RunTerraformCommands(deployDir string, action models.DeploymentAction) error
	ProcessTerraformOutputs(msg models.DeploymentMessage, deployDir string) (map[string]interface{}, error)
	PostOutputToAPI(projectID string, packageID string, outputData map[string]interface{}, action models.DeployStatus) error
	WriteOutputFile(packageID, deployDir string, outputData map[string]interface{}) error
}

type Executor struct {
	cfg                  *config.Config
	terraformModulesPath string
}

func NewExecutor(cfg *config.Config) *Executor {
	log.Println("Creating new Terraform Executor")
	return &Executor{
		cfg:                  cfg,
		terraformModulesPath: cfg.TerraformModulesPath,
	}
}

func (e *Executor) CopyTerraformModules(packageType string, deployDir string) error {
	sourceDir := e.terraformModulesPath
	log.Printf("Copying Terraform modules from %s to %s for package type %s", sourceDir, deployDir, packageType)

	// Construct the full path to the package-specific module
	sourcePath := filepath.Join(sourceDir, packageType)

	// Check if the source directory exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourcePath)
	}

	// Use the cp command to copy the directory
	cmd := exec.Command("cp", "-R", sourcePath+"/.", deployDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy terraform modules: %v\nOutput: %s", err, output)
	}

	log.Printf("Successfully copied Terraform modules to %s", deployDir)
	return nil
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

func (e *Executor) CreateSecretsFile(msg models.DeploymentMessage, deployDir string) error {
	log.Printf("Creating secrets file for package: %s in directory: %s", msg.PackageID, deployDir)

	secretsData := make(map[string]interface{})
	for k, v := range msg.Secrets {
		var jsonValue interface{}
		err := json.Unmarshal([]byte(v), &jsonValue)
		if err != nil {
			log.Printf("Error unmarshaling secret value for key %s: %v", k, err)
			jsonValue = v
		}
		secretsData[k] = jsonValue
	}

	filePath := filepath.Join(deployDir, fmt.Sprintf("%s_secrets.auto.tfvars.json", msg.PackageID))
	err := e.writeJSONFile(filePath, secretsData)
	if err != nil {
		log.Printf("Error creating secrets file: %v", err)
	} else {
		log.Printf("Successfully created secrets file: %s", filePath)
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

func (e *Executor) RunTerraformCommands(deployDir string, action models.DeploymentAction) error {
	log.Printf("Running Terraform commands in directory: %s for action: %s", deployDir, action)
	commands := []string{
		"terraform init",
		"terraform plan",
	}

	// todo: use model enum for this and create an interface for other eexecutor types to adhere to
	if action == models.ActionDeploy {
		commands = append(commands, "terraform apply -auto-approve")
	} else if action == models.ActionDestroy {
		commands = append(commands, "terraform destroy -auto-approve")
	} else {
		return fmt.Errorf("unsupported action: %s", action)
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

// todo: move this to rad-labs
type OutputPayloadBody struct {
	DeployStatus *string                `json:"deploy_status,omitempty"`
	OutputData   map[string]interface{} `json:"output_data,omitempty"`
	// errors and logs are added to the output data, which we will add a struct for shortly
	// ErrorMessage string                 `json:"error_message,omitempty"`
}

func (e *Executor) PostOutputToAPI(projectID string, packageID string, outputData map[string]interface{}, action models.DeployStatus) error {
	url := fmt.Sprintf("%s/provisioner/projects/%s/packages/%s", e.cfg.APIURL, projectID, packageID)
	log.Printf("Posting output data for package: %s to API", url)

	apiURL := e.cfg.APIURL
	if apiURL == "" {
		return fmt.Errorf("API_URL environment variable is not set")
	}

	deployStatus := string(action)
	payload := OutputPayloadBody{
		DeployStatus: &deployStatus,
		OutputData:   outputData,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling output data: %v", err)
	}
	log.Printf("JSON payload: %s", string(jsonData))

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-canvas-token", e.cfg.CanvasToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()
	log.Printf("Response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	log.Printf("Successfully patched output data for package: %s to API", packageID)
	return nil
}

func (e *Executor) runCommand(command, dir string) (string, error) {
	log.Printf("Running command: %s in directory: %s", command, dir)
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command failed: %v", err)
		return "", err
	}

	// Clean up the output
	cleanedOutput := cleanTerraformOutput(string(output))

	log.Printf("Command executed successfully")
	return cleanedOutput, nil
}

func cleanTerraformOutput(output string) string {
	// Split output into lines
	lines := strings.Split(output, "\n")

	var cleanedLines []string
	for _, line := range lines {
		// Remove timestamp and other formatting
		cleanedLine := regexp.MustCompile(`^\[[\d:]+\]\s*`).ReplaceAllString(line, "")

		// Remove empty lines and lines with only whitespace
		if strings.TrimSpace(cleanedLine) != "" {
			cleanedLines = append(cleanedLines, cleanedLine)
		}
	}

	// Join the cleaned lines back together
	return strings.Join(cleanedLines, "\n")
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
