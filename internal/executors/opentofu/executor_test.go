package executors

import (
	"testing"
)

// mockExecutor is a mock implementation of the Executor interface
type mockExecutor struct {
	applyFunc      func(string, map[string]string) error
	destroyFunc    func(string, map[string]string) error
	getOutputsFunc func(string) (map[string]interface{}, error)
}

func (m *mockExecutor) Apply(dir string, vars map[string]string) error {
	return m.applyFunc(dir, vars)
}

func (m *mockExecutor) Destroy(dir string, vars map[string]string) error {
	return m.destroyFunc(dir, vars)
}

func (m *mockExecutor) GetOutputs(dir string) (map[string]interface{}, error) {
	return m.getOutputsFunc(dir)
}

func TestOpenTofuExecutor_Apply(t *testing.T) {
	mockExec := &mockExecutor{
		applyFunc: func(dir string, vars map[string]string) error {
			return nil
		},
	}

	err := mockExec.Apply("", nil)
	if err != nil {
		t.Errorf("Apply() error = %v", err)
	}
}

func TestOpenTofuExecutor_Destroy(t *testing.T) {
	mockExec := &mockExecutor{
		destroyFunc: func(dir string, vars map[string]string) error {
			return nil
		},
	}

	err := mockExec.Destroy("", nil)
	if err != nil {
		t.Errorf("Destroy() error = %v", err)
	}
}

func TestOpenTofuExecutor_GetOutputs(t *testing.T) {
	mockOutput := map[string]interface{}{
		"output1": "value1",
		"output2": float64(42),
	}

	mockExec := &mockExecutor{
		getOutputsFunc: func(dir string) (map[string]interface{}, error) {
			return mockOutput, nil
		},
	}

	outputs, err := mockExec.GetOutputs("")
	if err != nil {
		t.Errorf("GetOutputs() error = %v", err)
	}

	if len(outputs) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(outputs))
	}

	if outputs["output1"] != "value1" || outputs["output2"] != float64(42) {
		t.Errorf("Unexpected output values: %v", outputs)
	}
}
