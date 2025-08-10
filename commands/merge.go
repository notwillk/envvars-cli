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

// MergeCommand handles the environment variable merging functionality
type MergeCommand struct {
	sources []Source
	options Options
}

// CreateMergeCommand creates a new merge command instance
func CreateMergeCommand(sources []Source, options Options) *MergeCommand {
	return &MergeCommand{
		sources: sources,
		options: options,
	}
}

// Execute runs the merge command
func (cmd *MergeCommand) Execute() error {
	// Check if any sources are specified
	if len(cmd.sources) == 0 {
		return fmt.Errorf("no sources specified")
	}

	if cmd.options.Verbose {
		fmt.Fprintf(os.Stderr, "Processing %d sources...\n", len(cmd.sources))
	}

	// Process each source and merge the results
	var allVariables []sources.EnvVar

	// Process sources in priority order (higher priority first)
	for _, source := range cmd.sources {
		if cmd.options.Verbose {
			fmt.Fprintf(os.Stderr, "Processing %s file: %s (priority: %d)\n", source.Type, source.FilePath, source.Priority)
		}

		var envFile sources.EnvFile
		var err error

		switch source.Type {
		case "json":
			envFile, err = cmd.parseJSONFile(source.FilePath)
		case "yaml":
			envFile, err = cmd.parseYAMLFile(source.FilePath)
		case "env":
			envFile, err = cmd.parseENVFile(source.FilePath)
		default:
			return fmt.Errorf("unsupported source type: %s", source.Type)
		}

		if err != nil {
			return fmt.Errorf("failed to parse %s file '%s': %w", source.Type, source.FilePath, err)
		}

		allVariables = append(allVariables, envFile.Variables...)
	}

	// Convert to map for output formatting (later sources override earlier ones)
	variablesMap := make(map[string]string)
	for _, envVar := range allVariables {
		variablesMap[envVar.Key] = envVar.Value
	}

	if cmd.options.Verbose {
		fmt.Fprintf(os.Stderr, "Merged %d variables\n", len(variablesMap))
	}

	// Output in the specified format
	switch cmd.options.Format {
	case "json":
		return formatters.OutputAsJSON(variablesMap)
	case "yaml":
		return formatters.OutputAsYAML(variablesMap)
	case "env":
		return formatters.OutputAsENV(variablesMap)
	default:
		return fmt.Errorf("unsupported output format: %s", cmd.options.Format)
	}
}

// parseENVFile reads and parses an environment variable file
func (cmd *MergeCommand) parseENVFile(filePath string) (sources.EnvFile, error) {
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
func (cmd *MergeCommand) parseJSONFile(filePath string) (sources.EnvFile, error) {
	processor := sources.CreateJSONProcessor()
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
func (cmd *MergeCommand) parseYAMLFile(filePath string) (sources.EnvFile, error) {
	processor := sources.CreateYAMLProcessor()
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
func (cmd *MergeCommand) unquoteValue(value string) string {
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
func (cmd *MergeCommand) resolveVariableReferences(value string, variables map[string]string) string {
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
