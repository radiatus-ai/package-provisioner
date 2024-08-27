package executors

import (
	"testing"
)

// mockHelmExecutor is a mock implementation of the Executor interface
type mockHelmExecutor struct {
	applyFunc      func(string, map[string]interface{}) error
	destroyFunc    func(string, map[string]interface{}) error
	getOutputsFunc func(string) (map[string]interface{}, error)
}

func (m *mockHelmExecutor) Apply(dir string, params map[string]interface{}) error {
	return m.applyFunc(dir, params)
}

func (m *mockHelmExecutor) Destroy(dir string, params map[string]interface{}) error {
	return m.destroyFunc(dir, params)
}

func (m *mockHelmExecutor) GetOutputs(dir string) (map[string]interface{}, error) {
	return m.getOutputsFunc(dir)
}

func TestHelmExecutor_Apply(t *testing.T) {
	mockExec := &mockHelmExecutor{
		applyFunc: func(dir string, params map[string]interface{}) error {
			// Verify that the required parameters are present
			if _, ok := params["chartName"]; !ok {
				t.Errorf("chartName parameter is missing")
			}
			if _, ok := params["releaseName"]; !ok {
				t.Errorf("releaseName parameter is missing")
			}
			return nil
		},
	}

	params := map[string]interface{}{
		"chartName":   "test-chart",
		"releaseName": "test-release",
	}

	err := mockExec.Apply("", params)
	if err != nil {
		t.Errorf("Apply() error = %v", err)
	}
}

func TestHelmExecutor_Destroy(t *testing.T) {
	mockExec := &mockHelmExecutor{
		destroyFunc: func(dir string, params map[string]interface{}) error {
			// Verify that the required parameter is present
			if _, ok := params["releaseName"]; !ok {
				t.Errorf("releaseName parameter is missing")
			}
			return nil
		},
	}

	params := map[string]interface{}{
		"releaseName": "test-release",
	}

	err := mockExec.Destroy("", params)
	if err != nil {
		t.Errorf("Destroy() error = %v", err)
	}
}

func TestHelmExecutor_GetOutputs(t *testing.T) {
	mockExec := &mockHelmExecutor{
		getOutputsFunc: func(dir string) (map[string]interface{}, error) {
			// Helm doesn't typically provide outputs, so we return an empty map
			return map[string]interface{}{}, nil
		},
	}

	outputs, err := mockExec.GetOutputs("")
	if err != nil {
		t.Errorf("GetOutputs() error = %v", err)
	}

	if len(outputs) != 0 {
		t.Errorf("Expected empty outputs, got %v", outputs)
	}
}
