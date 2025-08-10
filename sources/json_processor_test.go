package sources

import (
	"os"
	"reflect"
	"testing"
)

func TestCreateJSONProcessor(t *testing.T) {
	processor := CreateJSONProcessor()
	if processor == nil {
		t.Error("CreateJSONProcessor returned nil")
	}
}

func TestJSONProcessor_ProcessFile_ValidJSON(t *testing.T) {
	// Create a temporary file with valid JSON
	tempFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write valid JSON content
	jsonContent := `{
		"string_value": "test",
		"bool_value": true,
		"int_value": 42,
		"float_value": 3.14,
		"null_value": null
	}`
	_, err = tempFile.WriteString(jsonContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	processor := CreateJSONProcessor()
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

func TestJSONProcessor_ProcessFile_NonExistentFile(t *testing.T) {
	processor := CreateJSONProcessor()
	_, err := processor.ProcessFile("nonexistent.json")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestJSONProcessor_ProcessFile_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tempFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write invalid JSON content
	_, err = tempFile.WriteString("invalid: json: content: [")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	processor := CreateJSONProcessor()
	_, err = processor.ProcessFile(tempFile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON content")
	}
}

func TestJSONProcessor_ProcessFile_EmptyJSON(t *testing.T) {
	// Create a temporary file with empty JSON object
	tempFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write empty JSON object
	_, err = tempFile.WriteString("{}")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	processor := CreateJSONProcessor()
	result, err := processor.ProcessFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}
}

func TestJSONProcessor_ProcessFileWithMerge_ValidJSON(t *testing.T) {
	// Create a temporary file with valid JSON
	tempFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write valid JSON content
	jsonContent := `{
		"file_key": "file_value",
		"override_key": "file_override"
	}`
	_, err = tempFile.WriteString(jsonContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"existing_key": "existing_value",
		"override_key": "existing_override",
	}

	options := Options{FilePath: tempFile.Name()}
	processor := CreateJSONProcessor()
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

func TestJSONProcessor_ProcessFileWithMerge_NonExistentFile(t *testing.T) {
	existingKVs := map[string]string{"key": "value"}
	options := Options{FilePath: "nonexistent.json"}
	processor := CreateJSONProcessor()
	_, err := processor.ProcessFileWithMerge(existingKVs, options)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestJSONProcessor_ProcessFileWithMerge_EmptyExisting(t *testing.T) {
	// Create a temporary file with valid JSON
	tempFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write valid JSON content
	jsonContent := `{"key": "value"}`
	_, err = tempFile.WriteString(jsonContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{}
	options := Options{FilePath: tempFile.Name()}
	processor := CreateJSONProcessor()
	result, err := processor.ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{"key": "value"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
