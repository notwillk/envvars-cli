package commands

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	formatters "github.com/notwillk/envvars-cli/formatters"
	"github.com/notwillk/envvars-cli/sources"
)

// EnvProcessorCommand handles the environment processing functionality
type EnvProcessorCommand struct {
	filePaths []string
	format    string
	jsonFile  string
	yamlFile  string
}

// NewEnvProcessorCommand creates a new environment processor command instance
func NewEnvProcessorCommand(filePaths []string, format string, jsonFile string, yamlFile string) *EnvProcessorCommand {
	return &EnvProcessorCommand{
		filePaths: filePaths,
		format:    format,
		jsonFile:  jsonFile,
		yamlFile:  yamlFile,
	}
}

// Execute runs the environment processor command
func (cmd *EnvProcessorCommand) Execute() error {
	// Check if any files are specified
	if len(cmd.filePaths) == 0 && cmd.jsonFile == "" && cmd.yamlFile == "" {
		return fmt.Errorf("no files specified")
	}

	// Process each file and merge the results
	var allVariables []sources.EnvVar

	// Process ENV files
	for _, filePath := range cmd.filePaths {
		envFile, err := cmd.parseENVFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse ENV file '%s': %w", filePath, err)
		}
		allVariables = append(allVariables, envFile.Variables...)
	}

	// Process JSON file if specified
	if cmd.jsonFile != "" {
		envFile, err := cmd.parseJSONFile(cmd.jsonFile)
		if err != nil {
			return fmt.Errorf("failed to parse JSON file '%s': %w", cmd.jsonFile, err)
		}
		allVariables = append(allVariables, envFile.Variables...)
	}

	// Process YAML file if specified
	if cmd.yamlFile != "" {
		envFile, err := cmd.parseYAMLFile(cmd.yamlFile)
		if err != nil {
			return fmt.Errorf("failed to parse YAML file '%s': %w", cmd.yamlFile, err)
		}
		allVariables = append(allVariables, envFile.Variables...)
	}

	// Convert to map for output formatting
	variablesMap := make(map[string]string)
	for _, envVar := range allVariables {
		variablesMap[envVar.Key] = envVar.Value
	}

	// Output in the specified format
	switch cmd.format {
	case "json":
		return formatters.OutputAsJSON(variablesMap)
	case "yaml":
		return formatters.OutputAsYAML(variablesMap)
	case "env":
		return formatters.OutputAsENV(variablesMap)
	default:
		return fmt.Errorf("unsupported output format: %s", cmd.format)
	}
}

// parseEnvFile reads and parses an environment variable file
func (cmd *EnvProcessorCommand) parseEnvFile(filePath string) (sources.EnvFile, error) {
	// Use the specified format flags to determine how to parse the file
	if cmd.jsonFile != "" && filePath == cmd.jsonFile {
		return cmd.parseJSONFile(filePath)
	} else if cmd.yamlFile != "" && filePath == cmd.yamlFile {
		return cmd.parseYAMLFile(filePath)
	} else {
		return cmd.parseENVFile(filePath)
	}
}

// parseENVFile reads and parses an environment variable file
func (cmd *EnvProcessorCommand) parseENVFile(filePath string) (sources.EnvFile, error) {
	// Use the sources package to parse the file
	// Since parseEnvFile is private in sources, we'll implement the parsing here
	file, err := os.Open(filePath)
	if err != nil {
		return sources.EnvFile{}, fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	envFile := sources.EnvFile{
		Filename:  filePath,
		Variables: []sources.EnvVar{},
	}

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
				value = cmd.unquoteValue(value)
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
				// Unquote the value
				value = cmd.unquoteValue(value)
				// Resolve variable references
				resolvedValue := cmd.resolveVariableReferences(value, variables)
				envFile.Variables = append(envFile.Variables, sources.EnvVar{
					Key:   key,
					Value: resolvedValue,
					File:  filePath,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return sources.EnvFile{}, fmt.Errorf("error reading file '%s': %w", filePath, err)
	}

	return envFile, nil
}

// parseJSONFile reads and parses a JSON file
func (cmd *EnvProcessorCommand) parseJSONFile(filePath string) (sources.EnvFile, error) {
	processor := sources.NewJSONProcessor()
	variables, err := processor.ProcessFile(filePath)
	if err != nil {
		return sources.EnvFile{}, fmt.Errorf("failed to parse JSON file '%s': %w", filePath, err)
	}

	envFile := sources.EnvFile{
		Filename:  filePath,
		Variables: []sources.EnvVar{},
	}

	// Convert map[string]string to []EnvVar
	for key, value := range variables {
		envFile.Variables = append(envFile.Variables, sources.EnvVar{
			Key:   key,
			Value: value,
			File:  filePath,
		})
	}

	return envFile, nil
}

// parseYAMLFile reads and parses a YAML file
func (cmd *EnvProcessorCommand) parseYAMLFile(filePath string) (sources.EnvFile, error) {
	processor := sources.NewYAMLProcessor()
	variables, err := processor.ProcessFile(filePath)
	if err != nil {
		return sources.EnvFile{}, fmt.Errorf("failed to parse YAML file '%s': %w", filePath, err)
	}

	envFile := sources.EnvFile{
		Filename:  filePath,
		Variables: []sources.EnvVar{},
	}

	// Convert map[string]string to []EnvVar
	for key, value := range variables {
		envFile.Variables = append(envFile.Variables, sources.EnvVar{
			Key:   key,
			Value: value,
			File:  filePath,
		})
	}

	return envFile, nil
}

// unquoteValue removes quotes from a value if present
func (cmd *EnvProcessorCommand) unquoteValue(value string) string {
	value = strings.TrimSpace(value)

	// Remove single quotes
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		return strings.Trim(value, "'")
	}

	// Remove double quotes
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return strings.Trim(value, "\"")
	}

	return value
}

// resolveVariableReferences resolves ${VAR_NAME} references in values
func (cmd *EnvProcessorCommand) resolveVariableReferences(value string, variables map[string]string) string {
	// Use regex to find and replace variable references
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(value, func(match string) string {
		// Extract variable name from ${VAR_NAME}
		varName := match[2 : len(match)-1]
		if replacement, exists := variables[varName]; exists {
			return replacement
		}
		// If variable not found, keep the original reference
		return match
	})
}
