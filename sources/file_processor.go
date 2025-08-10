package sources

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Options contains configuration for file operations
type Options struct {
	FilePath string `json:"file_path"`
}

// EnvVar represents a single environment variable
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	File  string `json:"file"`
}

// EnvFile represents a parsed environment file
type EnvFile struct {
	Filename  string   `json:"filename"`
	Variables []EnvVar `json:"variables"`
}

// ProcessFileWithMerge takes existing key-value pairs and options,
// then outputs merged key-value pairs with file values taking precedence
func ProcessFileWithMerge(existingKVs map[string]string, options Options) (map[string]string, error) {
	// Parse the environment file from options
	envFile, err := parseEnvFile(options.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%s': %w", options.FilePath, err)
	}

	// Merge variables (file values take precedence over existing values)
	mergedVars := make(map[string]string)

	// First, add existing variables
	for key, value := range existingKVs {
		mergedVars[key] = value
	}

	// Then, add file variables (overriding existing ones)
	for _, variable := range envFile.Variables {
		mergedVars[variable.Key] = variable.Value
	}

	return mergedVars, nil
}

// parseEnvFile reads and parses an environment variable file
func parseEnvFile(filePath string) (EnvFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return EnvFile{}, fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	var envFile EnvFile
	envFile.Filename = filePath
	envFile.Variables = []EnvVar{}

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	variables := make(map[string]string) // For variable reference resolution

	// First pass: collect all variables
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := ""
			if len(parts) > 1 {
				value = strings.TrimSpace(parts[1])
			}

			if key != "" {
				// Unquote the value
				value = unquoteValue(value)
				variables[key] = value
			}
		}
	}

	// Second pass: resolve variable references and create EnvVar structs
	file.Seek(0, 0) // Reset file pointer
	scanner = bufio.NewScanner(file)
	lineNumber = 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := ""
			if len(parts) > 1 {
				value = strings.TrimSpace(parts[1])
			}

			if key != "" {
				// Unquote and resolve variable references
				value = unquoteValue(value)
				value = resolveVariableReferences(value, variables)

				envVar := EnvVar{
					Key:   key,
					Value: value,
					File:  filePath,
				}
				envFile.Variables = append(envFile.Variables, envVar)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return EnvFile{}, fmt.Errorf("error reading file '%s': %w", filePath, err)
	}

	return envFile, nil
}

// unquoteValue removes quotes and handles escape sequences
func unquoteValue(value string) string {
	value = strings.TrimSpace(value)

	// Handle single quotes
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		value = strings.Trim(value, "'")
		// Replace escaped single quotes
		value = strings.ReplaceAll(value, "\\'", "'")
		return value
	}

	// Handle double quotes
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = strings.Trim(value, "\"")
		// Replace escaped double quotes
		value = strings.ReplaceAll(value, "\\\"", "\"")
		return value
	}

	return value
}

// resolveVariableReferences replaces ${VAR_NAME} with actual values
func resolveVariableReferences(value string, variables map[string]string) string {
	// Use regex to find and replace variable references
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(value, func(match string) string {
		// Extract variable name from ${VAR_NAME}
		varName := match[2 : len(match)-1]
		if val, exists := variables[varName]; exists {
			return val
		}
		// If variable not found, return the original match
		return match
	})
}

// parseOptionsFile reads and parses a JSON options file
func parseOptionsFile(filePath string) (Options, error) {
	var options Options

	file, err := os.Open(filePath)
	if err != nil {
		return Options{}, fmt.Errorf("failed to open options file '%s': %w", filePath, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&options); err != nil {
		return Options{}, fmt.Errorf("failed to decode options file '%s': %w", filePath, err)
	}

	return options, nil
}
