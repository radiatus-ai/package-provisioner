package executors

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type OpenTofuExecutor struct{}

func NewOpenTofuExecutor() *OpenTofuExecutor {
	return &OpenTofuExecutor{}
}

func (o *OpenTofuExecutor) Apply(deployDir string, params map[string]interface{}) error {
	cmd := exec.Command("tofu", "apply", "-auto-approve")
	cmd.Dir = deployDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("OpenTofu apply failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func (o *OpenTofuExecutor) Destroy(deployDir string, params map[string]interface{}) error {
	cmd := exec.Command("tofu", "destroy", "-auto-approve")
	cmd.Dir = deployDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("OpenTofu destroy failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func (o *OpenTofuExecutor) GetOutputs(deployDir string) (map[string]interface{}, error) {
	cmd := exec.Command("tofu", "output", "-json")
	cmd.Dir = deployDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenTofu outputs: %v", err)
	}

	var outputs map[string]interface{}
	if err := json.Unmarshal(output, &outputs); err != nil {
		return nil, fmt.Errorf("failed to parse OpenTofu outputs: %v", err)
	}

	return outputs, nil
}
