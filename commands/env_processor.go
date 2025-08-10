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
}

// NewEnvProcessorCommand creates a new environment processor command instance
func NewEnvProcessorCommand(filePaths []string) *EnvProcessorCommand {
	return &EnvProcessorCommand{
		filePaths: filePaths,
	}
}

// Execute runs the environment processor command
func (cmd *EnvProcessorCommand) Execute() error {
	if len(cmd.filePaths) == 0 {
		return fmt.Errorf("no files specified for merge command")
	}

	// Process all files and output as JSON
	return cmd.parseAndOutputEnvFiles(cmd.filePaths)
}

// parseAndOutputEnvFiles processes environment files and outputs the result as JSON
func (cmd *EnvProcessorCommand) parseAndOutputEnvFiles(filePaths []string) error {
	// Collect all variables from all files into a single map
	allVariables := make(map[string]string)

	// Process all files
	for _, filePath := range filePaths {
		envFile, err := cmd.parseEnvFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file '%s': %w", filePath, err)
		}

		// Add variables to the combined map (later files take precedence)
		for _, variable := range envFile.Variables {
			allVariables[variable.Key] = variable.Value
		}
	}

	// Output as simple key-value JSON using the formatters package
	return formatters.OutputAsJSON(allVariables)
}

// parseEnvFile reads and parses an environment variable file
func (cmd *EnvProcessorCommand) parseEnvFile(filePath string) (sources.EnvFile, error) {
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
