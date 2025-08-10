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

func TestProcessFileWithMerge_FiltersInvalidKeys(t *testing.T) {
	// Create a temporary file with invalid keys
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write content with invalid keys
	envContent := `123INVALID=value1
VALID_KEY=value2
@INVALID=value3
VALID_KEY_2=value4`
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
		"VALID_KEY":   "value2",
		"VALID_KEY_2": "value4",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// New tests for directive functionality

func TestProcessFileWithMerge_WithRemoveDirective(t *testing.T) {
	// Create a temporary file with remove directive
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write content with remove directive
	envContent := `#remove EXISTING_KEY
KEY1=value1
KEY2=value2`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY": "existing_value",
		"OTHER_KEY":    "other_value",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"OTHER_KEY": "other_value", // EXISTING_KEY should be removed
		"KEY1":      "value1",
		"KEY2":      "value2",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_WithRemoveDirectiveCaseInsensitive(t *testing.T) {
	// Create a temporary file with remove directive (case insensitive)
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write content with remove directive in different cases
	envContent := `#REMOVE existing_key
#remove OTHER_KEY
KEY1=value1`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"existing_key": "existing_value",
		"OTHER_KEY":    "other_value",
		"KEEP_KEY":     "keep_value",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"KEEP_KEY": "keep_value", // existing_key and OTHER_KEY should be removed
		"KEY1":     "value1",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_WithMultipleRemoveDirectives(t *testing.T) {
	// Create a temporary file with multiple remove directives
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write content with multiple remove directives
	envContent := `#remove KEY1 KEY2
KEY3=value3
KEY4=value4`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
		"KEY5": "value5",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"KEY5": "value5", // KEY1 and KEY2 should be removed
		"KEY3": "value3",
		"KEY4": "value4",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_WithRemoveDirectiveAndRegularComments(t *testing.T) {
	// Create a temporary file with remove directive and regular comments
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write content with remove directive and regular comments
	envContent := `# This is a regular comment
#remove KEY1
# Another comment
KEY2=value2
# Final comment`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"KEY1": "value1",
		"KEY3": "value3",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"KEY3": "value3", // KEY1 should be removed
		"KEY2": "value2",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_WithRemoveDirectiveNoArguments(t *testing.T) {
	// Create a temporary file with remove directive but no arguments
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write content with remove directive but no arguments
	envContent := `#remove
KEY1=value1`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY": "existing_value",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"EXISTING_KEY": "existing_value", // No keys should be removed
		"KEY1":         "value1",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseDirective_ValidDirective(t *testing.T) {
	directive, err := parseDirective("#remove KEY1 KEY2", 1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := Directive{
		Name:      "remove",
		Arguments: []string{"KEY1", "KEY2"},
		Line:      1,
	}

	if !reflect.DeepEqual(directive, expected) {
		t.Errorf("Expected %v, got %v", expected, directive)
	}
}

func TestParseDirective_CaseInsensitive(t *testing.T) {
	directive, err := parseDirective("#REMOVE key1", 1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := Directive{
		Name:      "REMOVE",
		Arguments: []string{"key1"},
		Line:      1,
	}

	if !reflect.DeepEqual(directive, expected) {
		t.Errorf("Expected %v, got %v", expected, directive)
	}
}

func TestParseDirective_EmptyDirective(t *testing.T) {
	_, err := parseDirective("#", 1)
	if err == nil {
		t.Error("Expected error for empty directive")
	}
}

func TestParseDirective_WhitespaceOnly(t *testing.T) {
	_, err := parseDirective("#   ", 1)
	if err == nil {
		t.Error("Expected error for whitespace-only directive")
	}
}

func TestParseDirective_NoArguments(t *testing.T) {
	directive, err := parseDirective("#remove", 1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := Directive{
		Name:      "remove",
		Arguments: []string{},
		Line:      1,
	}

	if !reflect.DeepEqual(directive, expected) {
		t.Errorf("Expected %v, got %v", expected, directive)
	}
}

func TestApplyRemoveDirective(t *testing.T) {
	kvs := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
		"KEY3": "value3",
	}

	directive := Directive{
		Name:      "remove",
		Arguments: []string{"KEY1", "KEY3"},
		Line:      1,
	}

	applyRemoveDirective(kvs, directive)

	expected := map[string]string{
		"KEY2": "value2", // KEY1 and KEY3 should be removed
	}

	if !reflect.DeepEqual(kvs, expected) {
		t.Errorf("Expected %v, got %v", expected, kvs)
	}
}

func TestApplyRemoveDirective_CaseInsensitive(t *testing.T) {
	kvs := map[string]string{
		"Key1": "value1",
		"KEY2": "value2",
		"key3": "value3",
	}

	directive := Directive{
		Name:      "remove",
		Arguments: []string{"key1", "KEY3"},
		Line:      1,
	}

	applyRemoveDirective(kvs, directive)

	expected := map[string]string{
		"KEY2": "value2", // Key1 and key3 should be removed (case-insensitive)
	}

	if !reflect.DeepEqual(kvs, expected) {
		t.Errorf("Expected %v, got %v", expected, kvs)
	}
}

func TestApplyRemoveDirective_NonExistentKeys(t *testing.T) {
	kvs := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	directive := Directive{
		Name:      "remove",
		Arguments: []string{"NONEXISTENT_KEY"},
		Line:      1,
	}

	// This should not cause any errors
	applyRemoveDirective(kvs, directive)

	expected := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	if !reflect.DeepEqual(kvs, expected) {
		t.Errorf("Expected %v, got %v", expected, kvs)
	}
}

// Tests for #require directive

func TestProcessFileWithMerge_WithRequireDirective(t *testing.T) {
	// Create a temporary file with require directive
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with require directive
	envContent := `#require EXISTING_KEY
KEY1=value1
KEY2=value2`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY": "existing_value",
		"OTHER_KEY":    "other_value",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"EXISTING_KEY": "existing_value",
		"OTHER_KEY":    "other_value",
		"KEY1":         "value1",
		"KEY2":         "value2",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_WithRequireDirectiveFailure(t *testing.T) {
	// Create a temporary file with require directive for non-existent key
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with require directive for non-existent key
	envContent := `#require NONEXISTENT_KEY
KEY1=value1`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY": "existing_value",
	}

	options := Options{FilePath: tempFile.Name()}
	_, err = ProcessFileWithMerge(existingKVs, options)
	if err == nil {
		t.Error("Expected error for missing required environment variable")
	}

	expectedErrorMsg := "required environment variable 'NONEXISTENT_KEY' not found"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got: %s", expectedErrorMsg, err.Error())
	}
}

func TestProcessFileWithMerge_WithRequireDirectiveCaseInsensitive(t *testing.T) {
	// Create a temporary file with require directive using different case
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with require directive using different case
	envContent := `#require existing_key
KEY1=value1`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY": "existing_value",
	}

	options := Options{FilePath: tempFile.Name()}
	_, err = ProcessFileWithMerge(existingKVs, options)
	if err == nil {
		t.Error("Expected error for missing required environment variable (case-sensitive)")
	}

	expectedErrorMsg := "required environment variable 'existing_key' not found"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got: %s", expectedErrorMsg, err.Error())
	}
}

func TestProcessFileWithMerge_WithMultipleRequireDirectives(t *testing.T) {
	// Create a temporary file with multiple require directives
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with multiple require directives
	envContent := `#require EXISTING_KEY1
#require EXISTING_KEY2
KEY1=value1`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY1": "existing_value1",
		"EXISTING_KEY2": "existing_value2",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"EXISTING_KEY1": "existing_value1",
		"EXISTING_KEY2": "existing_value2",
		"KEY1":          "value1",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_WithRequireDirectiveAndRegularComments(t *testing.T) {
	// Create a temporary file with require directive and regular comments
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with require directive and regular comments
	envContent := `# This is a regular comment
#require EXISTING_KEY
# Another comment
KEY1=value1`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY": "existing_value",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"EXISTING_KEY": "existing_value",
		"KEY1":         "value1",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_WithRequireDirectiveNoArguments(t *testing.T) {
	// Create a temporary file with require directive but no arguments
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with require directive but no arguments
	envContent := `#require
KEY1=value1`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY": "existing_value",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// No arguments means no requirements to check, so it should succeed
	expected := map[string]string{
		"EXISTING_KEY": "existing_value",
		"KEY1":         "value1",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProcessFileWithMerge_WithRequireAndRemoveDirectives(t *testing.T) {
	// Create a temporary file with both require and remove directives
	tempFile, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write env content with both require and remove directives
	envContent := `#require EXISTING_KEY1
#remove EXISTING_KEY2
KEY1=value1`
	_, err = tempFile.WriteString(envContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	existingKVs := map[string]string{
		"EXISTING_KEY1": "existing_value1",
		"EXISTING_KEY2": "existing_value2",
		"OTHER_KEY":     "other_value",
	}

	options := Options{FilePath: tempFile.Name()}
	result, err := ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := map[string]string{
		"EXISTING_KEY1": "existing_value1", // Required and kept
		"OTHER_KEY":     "other_value",     // Kept
		"KEY1":          "value1",          // Added from file
		// EXISTING_KEY2 should be removed by #remove directive
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestApplyRequireDirective(t *testing.T) {
	kvs := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
		"KEY3": "value3",
	}

	directive := Directive{
		Name:      "require",
		Arguments: []string{"KEY1", "KEY2"},
		Line:      1,
	}

	// This should succeed
	err := applyRequireDirective(kvs, directive)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestApplyRequireDirective_MissingKey(t *testing.T) {
	kvs := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	directive := Directive{
		Name:      "require",
		Arguments: []string{"KEY1", "MISSING_KEY"},
		Line:      1,
	}

	// This should fail
	err := applyRequireDirective(kvs, directive)
	if err == nil {
		t.Error("Expected error for missing required key")
	}

	expectedErrorMsg := "required environment variable 'MISSING_KEY' not found"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got: %s", expectedErrorMsg, err.Error())
	}
}

func TestApplyRequireDirective_NoArguments(t *testing.T) {
	kvs := map[string]string{
		"KEY1": "value1",
	}

	directive := Directive{
		Name:      "require",
		Arguments: []string{},
		Line:      1,
	}

	// This should succeed (no requirements to check)
	err := applyRequireDirective(kvs, directive)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
