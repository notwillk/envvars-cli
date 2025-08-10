package commands

import (
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
	variablesMap := make(map[string]string)

	// Process sources in priority order (higher priority first)
	for _, source := range cmd.sources {
		if cmd.options.Verbose {
			fmt.Fprintf(os.Stderr, "Processing %s file: %s (priority: %d)\n", source.Type, source.FilePath, source.Priority)

			// Show current state of merged variables before processing this source
			if len(variablesMap) > 0 {
				fmt.Fprintf(os.Stderr, "Current merged variables (%d):\n", len(variablesMap))
				for key, value := range variablesMap {
					fmt.Fprintf(os.Stderr, "  %s=%s\n", key, value)
				}
			} else {
				fmt.Fprintf(os.Stderr, "No variables merged yet\n")
			}
			fmt.Fprintf(os.Stderr, "\n")
		}

		var err error

		switch source.Type {
		case "json":
			envFile, err := cmd.parseJSONFile(source.FilePath)
			if err != nil {
				return fmt.Errorf("failed to parse %s file '%s': %w", source.Type, source.FilePath, err)
			}
			// Merge JSON variables
			for _, envVar := range envFile.Variables {
				variablesMap[envVar.Key] = envVar.Value
			}
		case "yaml":
			envFile, err := cmd.parseYAMLFile(source.FilePath)
			if err != nil {
				return fmt.Errorf("failed to parse %s file '%s': %w", source.Type, source.FilePath, err)
			}
			// Merge YAML variables
			for _, envVar := range envFile.Variables {
				variablesMap[envVar.Key] = envVar.Value
			}
		case "env":
			// Use the directive-aware ProcessFileWithMerge function
			options := sources.Options{FilePath: source.FilePath}
			variablesMap, err = sources.ProcessFileWithMerge(variablesMap, options)
			if err != nil {
				return fmt.Errorf("failed to parse %s file '%s': %w", source.Type, source.FilePath, err)
			}
		case "sops":
			envFile, err := cmd.parseSOPSFile(source.FilePath, source.DecryptionKey)
			if err != nil {
				return fmt.Errorf("failed to parse %s file '%s': %w", source.Type, source.FilePath, err)
			}
			// Merge SOPS variables
			for _, envVar := range envFile.Variables {
				variablesMap[envVar.Key] = envVar.Value
			}
		default:
			return fmt.Errorf("unsupported source type: %s", source.Type)
		}
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

// parseSOPSFile reads and parses a SOPS-encrypted file
func (cmd *MergeCommand) parseSOPSFile(filePath string, decryptionKey string) (sources.EnvFile, error) {
	processor := sources.CreateSOPSProcessor()
	variables, err := processor.ProcessFile(filePath, decryptionKey)
	if err != nil {
		return sources.EnvFile{}, fmt.Errorf("failed to parse SOPS file '%s': %w", filePath, err)
	}

	envFile := sources.EnvFile{
		Filename:  filePath,
		Variables: variables,
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
