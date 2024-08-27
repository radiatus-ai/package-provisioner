package executors

import (
	"fmt"
	"os/exec"
)

type HelmExecutor struct{}

func NewHelmExecutor() *HelmExecutor {
	return &HelmExecutor{}
}

func (h *HelmExecutor) Apply(deployDir string, params map[string]interface{}) error {
	chartName := params["chartName"].(string)
	releaseName := params["releaseName"].(string)

	cmd := exec.Command("helm", "upgrade", "--install", releaseName, chartName, "--values", "values.yaml")
	cmd.Dir = deployDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("helm apply failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func (h *HelmExecutor) Destroy(deployDir string, params map[string]interface{}) error {
	releaseName := params["releaseName"].(string)

	cmd := exec.Command("helm", "uninstall", releaseName)
	cmd.Dir = deployDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("helm destroy failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func (h *HelmExecutor) GetOutputs(deployDir string) (map[string]interface{}, error) {
	// Implement logic to retrieve Helm outputs
	// This might involve parsing Helm status or custom logic
	return map[string]interface{}{}, nil
}
