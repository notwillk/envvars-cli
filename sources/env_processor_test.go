package sources

import (
	"os"
	"reflect"
	"testing"
)

func TestProcessFileWithMerge_ValidEnvFile(t *testing.T) {
	// Create a temporary file with valid env content
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write valid env content
	envContent := `# This is a comment
KEY1=value1
KEY2=value2
KEY3=value3`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY": "existing_value",
		"KEY1":         "old_value", // This should be overridden
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"EXISTING_KEY": "existing_value",
		"KEY1":         "value1", // File value should override existing
		"KEY2":         "value2",
		"KEY3":         "value3",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_NonExistentFile(t *testing.T) {
	existingKVs := map[string]string{"key": "value"}
	options := Options{FilePath: "nonexistent.env"}
	_, err := ProcessFileWithMerge(existingKVs, options)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestProcessFileWithMerge_EmptyExisting(t *testing.T) {
	// Create a temporary file with valid env content
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write valid env content
	envContent := `KEY1=value1
KEY2=value2`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{}
	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_EmptyFile(t *testing.T) {
	// Create a temporary file with empty content
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	existingKVs := map[string]string{"key": "value"}
	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should return existing values unchanged
	if !reflect.DeepEqual(result, existingKVs) {
		t.Errorf("Expected %v, got %v", existingKVs, result)
	}
}

func TestProcessFileWithMerge_FileWithCommentsAndEmptyLines(t *testing.T) {
	// Create a temporary file with comments and empty lines
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with comments and empty lines
	envContent := `# Configuration file
# Database settings

DB_HOST=localhost
DB_PORT=5432

# API settings
API_KEY=abc123
API_SECRET=xyz789`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{"EXISTING": "value"}
	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"EXISTING":   "value",
		"DB_HOST":    "localhost",
		"DB_PORT":    "5432",
		"API_KEY":    "abc123",
		"API_SECRET": "xyz789",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_QuotedValues(t *testing.T) {
	// Create a temporary file with quoted values
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with quoted values
	envContent := `KEY1="quoted value"
KEY2='single quoted'
KEY3=unquoted
KEY4="value with spaces"`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{}
	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"KEY1": "quoted value",
		"KEY2": "single quoted",
		"KEY3": "unquoted",
		"KEY4": "value with spaces",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_VariableReferences(t *testing.T) {
	// Create a temporary file with variable references
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with variable references
	envContent := `BASE_URL=https://api.example.com
API_VERSION=v1
FULL_URL=${BASE_URL}/${API_VERSION}
USERNAME=admin
PASSWORD=secret
AUTH_HEADER=${USERNAME}:${PASSWORD}`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{}
	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"BASE_URL":    "https://api.example.com",
		"API_VERSION": "v1",
		"FULL_URL":    "https://api.example.com/v1",
		"USERNAME":    "admin",
		"PASSWORD":    "secret",
		"AUTH_HEADER": "admin:secret",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_EdgeCases(t *testing.T) {
	// Create a temporary file with edge cases
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with edge cases
	envContent := `EMPTY_VALUE=
KEY_WITH_EQUALS=value=with=equals
KEY_WITH_SPACES = value with spaces
KEY_WITH_TABS	=	value	with	tabs
# Comment with = sign
KEY_AFTER_COMMENT=value`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{}
	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"EMPTY_VALUE":       "",
		"KEY_WITH_EQUALS":   "value=with=equals",
		"KEY_WITH_SPACES":   "value with spaces",
		"KEY_WITH_TABS":     "value	with	tabs",
		"KEY_AFTER_COMMENT": "value",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestUnquoteValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"quoted"`, "quoted"},
		{`'single quoted'`, "single quoted"},
		{`unquoted`, "unquoted"},
		{`"value with spaces"`, "value with spaces"},
		{`'value with "quotes"'`, "value with \"quotes\""},
		{`"value with 'quotes'"`, "value with 'quotes'"},
		{`""`, ""},
		{`''`, ""},
		{`"`, ``}, // Unmatched quote - function removes it
		{`'`, ``}, // Unmatched quote - function removes it
	}

	for _, test := range tests {
		result := unquoteValue(test.input)
		if result != test.expected {
			t.Errorf("unquoteValue(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestResolveVariableReferences(t *testing.T) {
	variables := map[string]string{
		"BASE_URL":    "https://api.example.com",
		"API_VERSION": "v1",
		"USERNAME":    "admin",
		"PASSWORD":    "secret",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"${BASE_URL}", "https://api.example.com"},
		{"${BASE_URL}/${API_VERSION}", "https://api.example.com/v1"},
		{"${USERNAME}:${PASSWORD}", "admin:secret"},
		{"no_variables", "no_variables"},
		{"${NONEXISTENT}", "${NONEXISTENT}"}, // Should not resolve
		{"${BASE_URL}${API_VERSION}", "https://api.example.comv1"},
		{"", ""},
	}

	for _, test := range tests {
		result := resolveVariableReferences(test.input, variables)
		if result != test.expected {
			t.Errorf("resolveVariableReferences(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestParseOptionsFile(t *testing.T) {
	// Create a temporary options file
	tempFile, err := os.CreateTemp("", "options-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write valid options content
	optionsContent := `{"file_path": "/path/to/file.env"}`
	_, err = tempFile.WriteString(optionsContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	result, err := parseOptionsFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := Options{FilePath: "/path/to/file.env"}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseOptionsFile_NonExistentFile(t *testing.T) {
	_, err := parseOptionsFile("nonexistent.json")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestParseOptionsFile_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tempFile, err := os.CreateTemp("", "options-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write invalid JSON content
	_, err = tempFile.WriteString("invalid json content")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	_, err = parseOptionsFile(tempFile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON content")
	}
}
