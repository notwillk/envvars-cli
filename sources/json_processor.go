package sources

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONProcessor handles processing of JSON files
type JSONProcessor struct{}

// CreateJSONProcessor creates a new JSON processor instance
func CreateJSONProcessor() *JSONProcessor {
	return &JSONProcessor{}
}

// ProcessFile reads a JSON file and extracts key-value pairs
func (jp *JSONProcessor) ProcessFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file '%s': %w", filePath, err)
	}
	defer file.Close()

	// Try to parse as a flat key-value object first
	var flatMap map[string]interface{}
	if err := json.NewDecoder(file).Decode(&flatMap); err != nil {
		return nil, fmt.Errorf("failed to parse JSON file '%s': %w", filePath, err)
	}

	// Convert to string key-value pairs
	result := make(map[string]string)
	for key, value := range flatMap {
		result[key] = fmt.Sprintf("%v", value)
	}

	return result, nil
}

// ProcessFileWithMerge merges existing key-value pairs with those from a JSON file
func (jp *JSONProcessor) ProcessFileWithMerge(existingKVs map[string]string, options Options) (map[string]string, error) {
	// Process the JSON file
	fileVars, err := jp.ProcessFile(options.FilePath)
	if err != nil {
		return nil, err
	}

	// Merge: file values take precedence
	mergedVars := make(map[string]string)

	// First, add existing variables
	for key, value := range existingKVs {
		mergedVars[key] = value
	}

	// Then, add file variables (overriding existing ones)
	for key, value := range fileVars {
		mergedVars[key] = value
	}

	return mergedVars, nil
}
