package sources

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// YAMLProcessor handles processing of YAML files
type YAMLProcessor struct{}

// CreateYAMLProcessor creates a new YAML processor instance
func CreateYAMLProcessor() *YAMLProcessor {
	return &YAMLProcessor{}
}

// isValidKey checks if a key matches the required regex pattern
func (yp *YAMLProcessor) isValidKey(key string) bool {
	matched, _ := regexp.MatchString(`^[A-Za-z_][A-Za-z0-9_]*$`, key)
	return matched
}

// ProcessFile reads a YAML file and extracts key-value pairs
func (yp *YAMLProcessor) ProcessFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open YAML file '%s': %w", filePath, err)
	}
	defer file.Close()

	// Try to parse as a flat key-value object first
	var flatMap map[string]interface{}
	if err := yaml.NewDecoder(file).Decode(&flatMap); err != nil {
		return nil, fmt.Errorf("failed to parse YAML file '%s': %w", filePath, err)
	}

	// Convert to string key-value pairs, filtering invalid keys
	result := make(map[string]string)
	for key, value := range flatMap {
		if yp.isValidKey(key) {
			result[key] = fmt.Sprintf("%v", value)
		}
	}

	return result, nil
}

// ProcessFileWithMerge merges existing key-value pairs with those from a YAML file
func (yp *YAMLProcessor) ProcessFileWithMerge(existingKVs map[string]string, options Options) (map[string]string, error) {
	// Process the YAML file
	fileVars, err := yp.ProcessFile(options.FilePath)
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
