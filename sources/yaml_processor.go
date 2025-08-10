package sources

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
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

	// First, read the entire file to check for $schema
	var rawData map[string]interface{}
	if err := yaml.NewDecoder(file).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("failed to parse YAML file '%s': %w", filePath, err)
	}

	// Check if there's a $schema field
	if schemaURL, hasSchema := rawData["$schema"]; hasSchema {
		// Validate against the schema before processing
		if err := yp.validateAgainstSchema(rawData, schemaURL.(string), filePath); err != nil {
			return nil, fmt.Errorf("JSON schema validation failed for '%s': %w", filePath, err)
		}
	}

	// Convert to string key-value pairs, filtering invalid keys and $schema
	result := make(map[string]string)
	for key, value := range rawData {
		// Skip the $schema field itself
		if key == "$schema" {
			continue
		}

		if yp.isValidKey(key) {
			result[key] = fmt.Sprintf("%v", value)
		}
	}

	return result, nil
}

// validateAgainstSchema validates the YAML data against the specified schema
func (yp *YAMLProcessor) validateAgainstSchema(data map[string]interface{}, schemaURL string, yamlFilePath string) error {
	// Handle local schema files
	if strings.HasPrefix(schemaURL, "./") || strings.HasPrefix(schemaURL, "../") || !strings.HasPrefix(schemaURL, "http") {
		// For local schemas, resolve the path relative to the YAML file being processed
		yamlDir := filepath.Dir(yamlFilePath)
		schemaPath := filepath.Join(yamlDir, schemaURL)

		// Create a new compiler and compile the schema directly from the file
		compiler := jsonschema.NewCompiler()
		schema, err := compiler.Compile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to compile local schema from '%s': %w", schemaPath, err)
		}

		// Validate the data against the schema
		if err := schema.Validate(data); err != nil {
			return fmt.Errorf("data does not match local schema: %w", err)
		}

		return nil
	}

	// For remote schemas, try to fetch and validate
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaURL, nil); err != nil {
		return fmt.Errorf("failed to add remote schema resource: %w", err)
	}

	// Compile the schema
	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		return fmt.Errorf("failed to compile remote schema: %w", err)
	}

	// Validate the data against the schema
	if err := schema.Validate(data); err != nil {
		return fmt.Errorf("data does not match remote schema: %w", err)
	}

	return nil
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
