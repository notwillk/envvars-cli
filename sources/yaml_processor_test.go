package sources

import (
	"os"
	"reflect"
	"testing"
)

func TestCreateYAMLProcessor(t *testing.T) {
	processor := CreateYAMLProcessor()
	if processor == nil {
		t.Error("CreateYAMLProcessor returned nil")
	}
}

func TestYAMLProcessor_ProcessFile_ValidYAML(t *testing.T) {
	// Create a temporary file with valid YAML
	tempFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write valid YAML content
	yamlContent := `string_value: test
bool_value: true
int_value: 42
float_value: 3.14
null_value: null`
	_, err = tempFile.WriteString(yamlContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	processor := CreateYAMLProcessor()
	result, err := processor.ProcessFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"string_value": "test",
		"bool_value":   "true",
		"int_value":    "42",
		"float_value":  "3.14",
		"null_value":   "<nil>",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestYAMLProcessor_ProcessFile_NonExistentFile(t *testing.T) {
	processor := CreateYAMLProcessor()
	_, err := processor.ProcessFile("nonexistent.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestYAMLProcessor_ProcessFile_InvalidYAML(t *testing.T) {
	// Create a temporary file with invalid YAML
	tempFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write invalid YAML content
	_, err = tempFile.WriteString("invalid: yaml: content: [")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	processor := CreateYAMLProcessor()
	_, err = processor.ProcessFile(tempFile.Name())
	if err == nil {
		t.Error("Expected error for invalid YAML content")
	}
}

func TestYAMLProcessor_ProcessFile_EmptyYAML(t *testing.T) {
	// Create a temporary file with empty YAML
	tempFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write empty YAML
	_, err = tempFile.WriteString("")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	processor := CreateYAMLProcessor()
	_, err = processor.ProcessFile(tempFile.Name())
	if err == nil {
		t.Error("Expected error for empty YAML file")
	}
}

func TestYAMLProcessor_ProcessFile_ComplexYAML(t *testing.T) {
	// Create a temporary file with complex YAML
	tempFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write complex YAML content
	yamlContent := `# This is a comment
database:
  host: localhost
  port: 5432
  credentials:
    username: admin
    password: secret

api:
  key: abc123
  secret: xyz789
  timeout: 30

features:
  - enabled
  - disabled
  - pending`
	_, err = tempFile.WriteString(yamlContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	processor := CreateYAMLProcessor()
	result, err := processor.ProcessFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// The current implementation handles nested structures by converting them to string representations
	expected := map[string]string{
		"database": "map[credentials:map[password:secret username:admin] host:localhost port:5432]",
		"api":      "map[key:abc123 secret:xyz789 timeout:30]",
		"features": "[enabled disabled pending]",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestYAMLProcessor_ProcessFileWithMerge_ValidYAML(t *testing.T) {
	// Create a temporary file with valid YAML
	tempFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write valid YAML content
	yamlContent := `file_key: file_value
override_key: file_override`
	_, err = tempFile.WriteString(yamlContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"existing_key": "existing_value",
		"override_key": "existing_override",
	}

	options := Options{FilePath: tempFile.Name()}
	processor := CreateYAMLProcessor()
	result, err := processor.ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"existing_key": "existing_value",
		"file_key":     "file_value",
		"override_key": "file_override", // File value should override existing
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestYAMLProcessor_ProcessFileWithMerge_NonExistentFile(t *testing.T) {
	existingKVs := map[string]string{"key": "value"}
	options := Options{FilePath: "nonexistent.yaml"}
	processor := CreateYAMLProcessor()
	_, err := processor.ProcessFileWithMerge(existingKVs, options)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestYAMLProcessor_ProcessFileWithMerge_EmptyExisting(t *testing.T) {
	// Create a temporary file with valid YAML
	tempFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write valid YAML content
	yamlContent := `key: value`
	_, err = tempFile.WriteString(yamlContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{}
	options := Options{FilePath: tempFile.Name()}
	processor := CreateYAMLProcessor()
	result, err := processor.ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{"key": "value"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
