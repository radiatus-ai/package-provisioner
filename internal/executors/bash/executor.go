package executors

import (
	"fmt"
	"os/exec"
)

type BashExecutor struct{}

func NewBashExecutor() *BashExecutor {
	return &BashExecutor{}
}

func (b *BashExecutor) Apply(deployDir string, params map[string]interface{}) error {
	script := params["applyScript"].(string)
	return b.runScript(deployDir, script)
}

func (b *BashExecutor) Destroy(deployDir string, params map[string]interface{}) error {
	script := params["destroyScript"].(string)
	return b.runScript(deployDir, script)
}

func (b *BashExecutor) GetOutputs(deployDir string) (map[string]interface{}, error) {
	// Implement logic to retrieve outputs from a bash script
	// This might involve parsing a specific output file or format
	return map[string]interface{}{}, nil
}

func (b *BashExecutor) runScript(deployDir string, script string) error {
	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = deployDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("bash script execution failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}
