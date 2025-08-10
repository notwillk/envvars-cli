package sources

import (
	"os"
	"testing"
)

func TestCreateSOPSProcessor(t *testing.T) {
	processor := CreateSOPSProcessor()
	if processor == nil {
		t.Error("CreateSOPSProcessor returned nil")
	}
}

func TestSOPSProcessor_ProcessFile_NonExistentFile(t *testing.T) {
	processor := CreateSOPSProcessor()
	_, err := processor.ProcessFile("nonexistent.yaml", "test-key")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestSOPSProcessor_ProcessFile_InvalidYAML(t *testing.T) {
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

	processor := CreateSOPSProcessor()
	_, err = processor.ProcessFile(tempFile.Name(), "test-key")
	if err == nil {
		t.Error("Expected error for invalid YAML content")
	}
}

func TestSOPSProcessor_flattenMap_SimpleTypes(t *testing.T) {
	processor := CreateSOPSProcessor()
	var variables []EnvVar

	testData := map[string]interface{}{
		"string_value": "test",
		"bool_value":   true,
		"int_value":    42,
		"float_value":  3.14,
	}

	processor.flattenMap("", testData, &variables)

	expected := map[string]string{
		"STRING_VALUE": "test",
		"BOOL_VALUE":   "true",
		"INT_VALUE":    "42",
		"FLOAT_VALUE":  "3.14",
	}

	if len(variables) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(variables))
	}

	for _, envVar := range variables {
		if expectedValue, exists := expected[envVar.Key]; !exists {
			t.Errorf("Unexpected key: %s", envVar.Key)
		} else if envVar.Value != expectedValue {
			t.Errorf("Expected value %s for key %s, got %s", expectedValue, envVar.Key, envVar.Value)
		}
	}
}

func TestSOPSProcessor_flattenMap_NestedStructures(t *testing.T) {
	processor := CreateSOPSProcessor()
	var variables []EnvVar

	testData := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
			"credentials": map[string]interface{}{
				"username": "admin",
				"password": "secret",
			},
		},
		"api": map[string]interface{}{
			"key":    "abc123",
			"secret": "xyz789",
		},
	}

	processor.flattenMap("", testData, &variables)

	expected := map[string]string{
		"DATABASE_HOST":                 "localhost",
		"DATABASE_PORT":                 "5432",
		"DATABASE_CREDENTIALS_USERNAME": "admin",
		"DATABASE_CREDENTIALS_PASSWORD": "secret",
		"API_KEY":                       "abc123",
		"API_SECRET":                    "xyz789",
	}

	if len(variables) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(variables))
	}

	for _, envVar := range variables {
		if expectedValue, exists := expected[envVar.Key]; !exists {
			t.Errorf("Unexpected key: %s", envVar.Key)
		} else if envVar.Value != expectedValue {
			t.Errorf("Expected value %s for key %s, got %s", expectedValue, envVar.Key, envVar.Value)
		}
	}
}

func TestSOPSProcessor_flattenMap_Arrays(t *testing.T) {
	processor := CreateSOPSProcessor()
	var variables []EnvVar

	testData := map[string]interface{}{
		"endpoints": []interface{}{
			"https://api1.example.com",
			"https://api2.example.com",
			"https://api3.example.com",
		},
		"ports": []interface{}{
			8080,
			8081,
			8082,
		},
		"mixed": []interface{}{
			"string",
			42,
			true,
		},
	}

	processor.flattenMap("", testData, &variables)

	expected := map[string]string{
		"ENDPOINTS": "https://api1.example.com,https://api2.example.com,https://api3.example.com",
		"PORTS":     "8080,8081,8082",
		"MIXED":     "string,42,true",
	}

	if len(variables) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(variables))
	}

	for _, envVar := range variables {
		if expectedValue, exists := expected[envVar.Key]; !exists {
			t.Errorf("Unexpected key: %s", envVar.Key)
		} else if envVar.Value != expectedValue {
			t.Errorf("Expected value %s for key %s, got %s", expectedValue, envVar.Key, envVar.Value)
		}
	}
}

func TestSOPSProcessor_flattenMap_WithPrefix(t *testing.T) {
	processor := CreateSOPSProcessor()
	var variables []EnvVar

	testData := map[string]interface{}{
		"config": map[string]interface{}{
			"debug": true,
			"level": "info",
		},
	}

	processor.flattenMap("APP", testData, &variables)

	expected := map[string]string{
		"APP_CONFIG_DEBUG": "true",
		"APP_CONFIG_LEVEL": "info",
	}

	if len(variables) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(variables))
	}

	for _, envVar := range variables {
		if expectedValue, exists := expected[envVar.Key]; !exists {
			t.Errorf("Unexpected key: %s", envVar.Key)
		} else if envVar.Value != expectedValue {
			t.Errorf("Expected value %s for key %s, got %s", expectedValue, envVar.Key, envVar.Value)
		}
	}
}

func TestSOPSProcessor_flattenMap_EmptyMap(t *testing.T) {
	processor := CreateSOPSProcessor()
	var variables []EnvVar

	processor.flattenMap("", map[string]interface{}{}, &variables)

	if len(variables) != 0 {
		t.Errorf("Expected 0 variables for empty map, got %d", len(variables))
	}
}

func TestSOPSProcessor_flattenMap_NilMap(t *testing.T) {
	processor := CreateSOPSProcessor()
	var variables []EnvVar

	processor.flattenMap("", nil, &variables)

	if len(variables) != 0 {
		t.Errorf("Expected 0 variables for nil map, got %d", len(variables))
	}
}
