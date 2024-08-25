package executors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBashExecutor_Apply(t *testing.T) {
	executor := NewBashExecutor()

	tempDir, err := os.MkdirTemp("", "bash-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	params := map[string]interface{}{
		"applyScript": "echo 'Hello, World!' > test.txt",
	}

	err = executor.Apply(tempDir, params)
	if err != nil {
		t.Errorf("Apply() error = %v", err)
	}

	// Check if the file was created
	testFile := filepath.Join(tempDir, "test.txt")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Expected file was not created: %s", testFile)
	}

	// Check the content of the file
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read test file: %v", err)
	}
	if string(content) != "Hello, World!\n" {
		t.Errorf("File content mismatch. Got: %s, Want: Hello, World!", string(content))
	}
}

func TestBashExecutor_Destroy(t *testing.T) {
	executor := NewBashExecutor()

	tempDir, err := os.MkdirTemp("", "bash-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("Test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	params := map[string]interface{}{
		"destroyScript": "rm test.txt",
	}

	err = executor.Destroy(tempDir, params)
	if err != nil {
		t.Errorf("Destroy() error = %v", err)
	}

	// Check if the file was removed
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("Expected file to be removed: %s", testFile)
	}
}
