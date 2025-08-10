package commands

import (
	"os"
	"testing"
)

func TestCreateMergeCommand(t *testing.T) {
	sources := []Source{
		{FilePath: "test.env", Type: "env", Priority: 0},
		{FilePath: "test.json", Type: "json", Priority: 1},
	}
	options := Options{Verbose: false, Format: "env"}

	cmd := CreateMergeCommand(sources, options)
	if cmd == nil {
		t.Error("CreateMergeCommand returned nil")
	}

	if len(cmd.sources) != 2 {
		t.Errorf("Expected 2 sources, got %d", len(cmd.sources))
	}

	if !cmd.options.Verbose == false {
		t.Error("Expected verbose to be false")
	}

	if cmd.options.Format != "env" {
		t.Errorf("Expected format to be 'env', got '%s'", cmd.options.Format)
	}
}

func TestMergeCommand_Execute_NoSources(t *testing.T) {
	cmd := CreateMergeCommand([]Source{}, Options{Verbose: false, Format: "env"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error when no sources provided")
	}
}

func TestMergeCommand_Execute_UnsupportedSourceType(t *testing.T) {
	sources := []Source{
		{FilePath: "test.unsupported", Type: "unsupported", Priority: 0},
	}
	cmd := CreateMergeCommand(sources, Options{Verbose: false, Format: "env"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for unsupported source type")
	}
}

func TestMergeCommand_Execute_UnsupportedOutputFormat(t *testing.T) {
	// Create a temporary env file for testing
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.WriteString("TEST_KEY=test_value\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	sources := []Source{
		{FilePath: tempFile.Name(), Type: "env", Priority: 0},
	}
	cmd := CreateMergeCommand(sources, Options{Verbose: false, Format: "unsupported"})
	err = cmd.Execute()
	if err == nil {
		t.Error("Expected error for unsupported output format")
	}
}

func TestMergeCommand_Execute_VerboseOutput(t *testing.T) {
	// Create a temporary env file for testing
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.WriteString("TEST_KEY=test_value\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	sources := []Source{
		{FilePath: tempFile.Name(), Type: "env", Priority: 0},
	}
	cmd := CreateMergeCommand(sources, Options{Verbose: true, Format: "env"})
	err = cmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestMergeCommand_Execute_SourcePriority(t *testing.T) {
	// Create temporary env files for testing
	tempFile1, err := os.CreateTemp("", "test1-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file 1: %v", err)
	}
	defer os.Remove(tempFile1.Name())
	defer tempFile1.Close()

	tempFile2, err := os.CreateTemp("", "test2-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file 2: %v", err)
	}
	defer os.Remove(tempFile2.Name())
	defer tempFile2.Close()

	// Write different values for the same key
	_, err = tempFile1.WriteString("DUPLICATE_KEY=first_value\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file 1: %v", err)
	}

	_, err = tempFile2.WriteString("DUPLICATE_KEY=second_value\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file 2: %v", err)
	}

	// First source has lower priority (0), second has higher priority (1)
	sources := []Source{
		{FilePath: tempFile1.Name(), Type: "env", Priority: 0},
		{FilePath: tempFile2.Name(), Type: "env", Priority: 1},
	}
	cmd := CreateMergeCommand(sources, Options{Verbose: false, Format: "env"})
	err = cmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestMergeCommand_parseSOPSFile_NonExistentFile(t *testing.T) {
	cmd := CreateMergeCommand([]Source{}, Options{})
	_, err := cmd.parseSOPSFile("nonexistent.yaml", "test-key")
	if err == nil {
		t.Error("Expected error for non-existent SOPS file")
	}
}

func TestMergeCommand_parseSOPSFile_InvalidDecryptionKey(t *testing.T) {
	// Create a temporary file that's not actually encrypted
	tempFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.WriteString("test: value\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	cmd := CreateMergeCommand([]Source{}, Options{})
	_, err = cmd.parseSOPSFile(tempFile.Name(), "invalid-key")
	// This should fail because the file is not actually encrypted with SOPS
	if err == nil {
		t.Error("Expected error for invalid SOPS decryption")
	}
}

func TestMergeCommand_Execute_WithSOPSSource(t *testing.T) {
	// Create a temporary YAML file that's not encrypted (for testing purposes)
	tempFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.WriteString("test_key: test_value\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	sources := []Source{
		{FilePath: tempFile.Name(), Type: "sops", Priority: 0, DecryptionKey: "test-key"},
	}
	cmd := CreateMergeCommand(sources, Options{Verbose: false, Format: "env"})

	// This will fail because the file is not actually encrypted, but it tests the SOPS path
	err = cmd.Execute()
	// We expect an error because the file is not encrypted, but the SOPS processing path is tested
	if err == nil {
		t.Error("Expected error for non-encrypted file in SOPS processing")
	}
}

func TestMergeCommand_Execute_MixedSourceTypes(t *testing.T) {
	// Create temporary files for different types
	envFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create env temp file: %v", err)
	}
	defer os.Remove(envFile.Name())
	defer envFile.Close()

	jsonFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("Failed to create json temp file: %v", err)
	}
	defer os.Remove(jsonFile.Name())
	defer jsonFile.Close()

	// Write content to files
	_, err = envFile.WriteString("ENV_KEY=env_value\n")
	if err != nil {
		t.Fatalf("Failed to write to env temp file: %v", err)
	}

	_, err = jsonFile.WriteString(`{"json_key": "json_value"}`)
	if err != nil {
		t.Fatalf("Failed to write to json temp file: %v", err)
	}

	sources := []Source{
		{FilePath: envFile.Name(), Type: "env", Priority: 0},
		{FilePath: jsonFile.Name(), Type: "json", Priority: 1},
	}
	cmd := CreateMergeCommand(sources, Options{Verbose: false, Format: "env"})
	err = cmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
