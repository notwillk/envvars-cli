package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	formatters "github.com/notwillk/envvars-cli/formatters"
	"github.com/notwillk/envvars-cli/sources"
	"github.com/spf13/pflag"
)

func main() {
	// Define flags
	var help bool
	var version bool
	var filePaths []string

	// Set up flags
	pflag.BoolVarP(&help, "help", "h", false, "Show this help message")
	pflag.BoolVarP(&version, "version", "v", false, "Show version information")
	pflag.StringSliceVarP(&filePaths, "file", "f", []string{}, "Read and parse environment variable files (can be specified multiple times)")

	// Parse flags
	pflag.Parse()

	// Handle help flag
	if help {
		showHelp()
		return
	}

	// Handle version flag
	if version {
		showVersion()
		return
	}

	// Handle file flags
	if len(filePaths) > 0 {
		if err := parseAndOutputEnvFiles(filePaths); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing files: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// If no flags provided, show help by default
	if len(os.Args) == 1 {
		showHelp()
		return
	}

	// TODO: Implement your CLI logic here

	// Example: Print current working directory
	if wd, err := os.Getwd(); err == nil {
		fmt.Printf("Working directory: %s\n", wd)
	}
}

// parseAndOutputEnvFiles processes environment files with optional merging
func parseAndOutputEnvFiles(filePaths []string) error {
	return parseAndOutputEnvFilesWithMerge(filePaths, nil, "")
}

// parseAndOutputEnvFilesWithMerge processes environment files with optional merging
// existingKVs: map of existing key-value pairs to merge with
// options: configuration options including file path
func parseAndOutputEnvFilesWithMerge(filePaths []string, existingKVs map[string]string, optionsFile string) error {
	// Collect all variables from all files into a single map
	allVariables := make(map[string]string)

	// If we have existing key-values, start with them
	if existingKVs != nil {
		for key, value := range existingKVs {
			allVariables[key] = value
		}
	}

	// Process all files
	for _, filePath := range filePaths {
		envFile, err := parseEnvFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file '%s': %w", filePath, err)
		}

		// Add variables to the combined map (file values take precedence)
		for _, variable := range envFile.Variables {
			allVariables[variable.Key] = variable.Value
		}
	}

	// If options file is provided, also process that file
	if optionsFile != "" {
		options, err := parseOptionsFile(optionsFile)
		if err != nil {
			return fmt.Errorf("failed to parse options file: %w", err)
		}

		if options.FilePath != "" {
			envFile, err := parseEnvFile(options.FilePath)
			if err != nil {
				return fmt.Errorf("failed to parse options file path '%s': %w", options.FilePath, err)
			}

			// Add variables from options file (these take precedence)
			for _, variable := range envFile.Variables {
				allVariables[variable.Key] = variable.Value
			}
		}
	}

	// Output as simple key-value JSON using the formatters package
	return formatters.OutputAsJSON(allVariables)
}

// ProcessFileWithMerge is the main function that takes existing key-value strings and options
// and outputs merged key-value pairs with file values taking precedence
func ProcessFileWithMerge(existingKVs map[string]string, options sources.Options) error {
	// Use the sources package to process the merge
	mergedVars, err := sources.ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		return err
	}

	// Output the merged result using the formatters package
	return formatters.OutputAsJSON(mergedVars)
}

// parseKeyValueString parses a string in format "key1=value1,key2=value2"
func parseKeyValueString(kvString string) (map[string]string, error) {
	result := make(map[string]string)

	pairs := strings.Split(kvString, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		if strings.Contains(pair, "=") {
			parts := strings.SplitN(pair, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := ""
			if len(parts) > 1 {
				value = strings.TrimSpace(parts[1])
			}

			if key != "" {
				result[key] = value
			}
		}
	}

	return result, nil
}

// parseOptionsFile reads and parses a JSON options file
func parseOptionsFile(filePath string) (sources.Options, error) {
	var options sources.Options

	file, err := os.Open(filePath)
	if err != nil {
		return sources.Options{}, fmt.Errorf("failed to open options file '%s': %w", filePath, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&options); err != nil {
		return sources.Options{}, fmt.Errorf("failed to decode options file '%s': %w", filePath, err)
	}

	return options, nil
}

func parseEnvFile(filePath string) (sources.EnvFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return sources.EnvFile{}, fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	envFile := sources.EnvFile{
		Filename:  filePath,
		Variables: []sources.EnvVar{},
	}

	// First pass: collect all variables
	variables := make(map[string]string)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := ""
			if len(parts) > 1 {
				value = strings.TrimSpace(parts[1])
			}

			// Remove quotes if present
			value = unquoteValue(value)

			variables[key] = value
		}
	}

	// Second pass: resolve variable references
	for key, value := range variables {
		resolvedValue := resolveVariableReferences(value, variables)
		envFile.Variables = append(envFile.Variables, sources.EnvVar{
			Key:   key,
			Value: resolvedValue,
			File:  filePath,
		})
	}

	if err := scanner.Err(); err != nil {
		return sources.EnvFile{}, fmt.Errorf("error reading file '%s': %w", filePath, err)
	}

	return envFile, nil
}

func unquoteValue(value string) string {
	value = strings.TrimSpace(value)

	// Remove double quotes
	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		value = strings.TrimPrefix(value, `"`)
		value = strings.TrimSuffix(value, `"`)
		// Unescape quotes
		value = strings.ReplaceAll(value, `\"`, `"`)
	}

	// Remove single quotes
	if strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`) {
		value = strings.TrimPrefix(value, `'`)
		value = strings.TrimSuffix(value, `'`)
		// Unescape apostrophes
		value = strings.ReplaceAll(value, `\'`, `'`)
	}

	return value
}

func resolveVariableReferences(value string, variables map[string]string) string {
	// Simple variable reference resolution: ${VAR_NAME}
	re := regexp.MustCompile(`\$\{([^}]+)\}`)

	return re.ReplaceAllStringFunc(value, func(match string) string {
		// Extract variable name from ${VAR_NAME}
		varName := match[2 : len(match)-1]
		if resolvedValue, exists := variables[varName]; exists {
			return resolvedValue
		}
		// If variable not found, return the original match
		return match
	})
}

func readAndOutputFiles(filePaths []string) error {
	for i, filePath := range filePaths {
		// Add separator between files if multiple files
		if i > 0 {
			fmt.Println("---")
		}

		// Show filename as header if multiple files
		if len(filePaths) > 1 {
			fmt.Printf("File: %s\n", filePath)
			fmt.Println("---")
		}

		if err := readAndOutputFile(filePath); err != nil {
			return fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}

		// Add newline after file content
		fmt.Println()
	}

	return nil
}

func readAndOutputFile(filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	// Copy file contents to stdout
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return fmt.Errorf("failed to read file '%s': %w", filePath, err)
	}

	return nil
}

func showHelp() {
	fmt.Println("envvars-cli - Environment Variables CLI Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  envvars-cli [flags] [command]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -h, --help      Show this help message")
	fmt.Println("  -v, --version   Show version information")
	fmt.Println("  -f, --file      Read and parse environment variable files (can be specified multiple times)")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  TODO: Add your CLI commands here")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  envvars-cli --help")
	fmt.Println("  envvars-cli --version")
	fmt.Println("  envvars-cli --file example.env")
	fmt.Println("  envvars-cli -f /path/to/file.env")
	fmt.Println("  envvars-cli --file file1.env --file file2.env")
	fmt.Println("  envvars-cli -f file1.env -f file2.env -f file3.env")
	fmt.Println()
	fmt.Println("Output:")
	fmt.Println("  Files are parsed as environment variable files and output as JSON")
	fmt.Println("  Output format: simple key-value pairs without metadata")
	fmt.Println("  Variable references (${VAR_NAME}) are resolved automatically")
	fmt.Println("  Multiple files are merged into a single key-value object")
	fmt.Println()
	fmt.Println("Programmatic Usage:")
	fmt.Println("  Use ProcessFileWithMerge(existingKVs, options) function for merging with existing values")
	fmt.Println("  File values take precedence over existing values when merging")
}

func showVersion() {
	fmt.Println("envvars-cli v0.1.0")
}
