package sources

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Directive represents a processing directive
type Directive struct {
	Name      string   `json:"name"`
	Arguments []string `json:"arguments"`
	Line      int      `json:"line"`
}

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
	Filename   string      `json:"filename"`
	Variables  []EnvVar    `json:"variables"`
	Directives []Directive `json:"directives"`
}

// ProcessFileWithMerge takes existing key-value pairs and options,
// then outputs merged key-value pairs with file values taking precedence
func ProcessFileWithMerge(existingKVs map[string]string, options Options) (map[string]string, error) {
	// Parse the environment file from options
	envFile, err := parseEnvFile(options.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%s': %w", options.FilePath, err)
	}

	// First, apply remove directives to existing key-value pairs
	processedKVs := applyRemoveDirectives(existingKVs, envFile.Directives)

	// Merge variables (file values take precedence over existing values)
	mergedVars := make(map[string]string)

	// First, add existing variables (after remove directive processing)
	for key, value := range processedKVs {
		mergedVars[key] = value
	}

	// Then, add file variables (overriding existing ones)
	for _, variable := range envFile.Variables {
		mergedVars[variable.Key] = variable.Value
	}

	// Apply filter directives to remove variables based on patterns
	mergedVars = applyFilterDirectives(mergedVars, envFile.Directives)

	// Apply filter-unless directives to keep only variables matching patterns
	mergedVars = applyFilterUnlessDirectives(mergedVars, envFile.Directives)

	// Finally, apply require directives to the final merged result
	if err := applyRequireDirectives(mergedVars, envFile.Directives); err != nil {
		return nil, err
	}

	return mergedVars, nil
}

// applyRemoveDirectives applies only remove directives to the key-value pairs
func applyRemoveDirectives(kvs map[string]string, directives []Directive) map[string]string {
	result := make(map[string]string)

	// Copy existing key-value pairs
	for key, value := range kvs {
		result[key] = value
	}

	// Apply only remove directives
	for _, directive := range directives {
		if strings.ToLower(directive.Name) == "remove" {
			applyRemoveDirective(result, directive)
		}
	}

	return result
}

// applyRequireDirectives applies only require directives to the key-value pairs
func applyRequireDirectives(kvs map[string]string, directives []Directive) error {
	// Apply only require directives
	for _, directive := range directives {
		if strings.ToLower(directive.Name) == "require" {
			if err := applyRequireDirective(kvs, directive); err != nil {
				return err
			}
		}
	}

	return nil
}

// applyRemoveDirective removes environment variables based on the directive
func applyRemoveDirective(kvs map[string]string, directive Directive) {
	for _, arg := range directive.Arguments {
		// Remove the specified key (case-insensitive)
		for key := range kvs {
			if strings.EqualFold(key, arg) {
				delete(kvs, key)
			}
		}
	}
}

// applyRequireDirective ensures required environment variables are present
func applyRequireDirective(kvs map[string]string, directive Directive) error {
	for _, arg := range directive.Arguments {
		if _, exists := kvs[arg]; !exists {
			return fmt.Errorf("required environment variable '%s' not found", arg)
		}
	}
	return nil
}

// parseDirective parses a directive line
func parseDirective(line string, lineNumber int) (Directive, error) {
	// Remove the # prefix and trim whitespace
	directiveText := strings.TrimSpace(strings.TrimPrefix(line, "#"))

	// Split by whitespace to get directive name and arguments
	parts := strings.Fields(directiveText)
	if len(parts) == 0 {
		return Directive{}, fmt.Errorf("empty directive at line %d", lineNumber)
	}

	directive := Directive{
		Name:      parts[0],
		Arguments: parts[1:],
		Line:      lineNumber,
	}

	return directive, nil
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
	envFile.Directives = []Directive{}

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	variables := make(map[string]string) // For variable reference resolution

	// First pass: collect all variables and directives
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle directives
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "# ") {
			// Check if it's a directive (not a regular comment)
			directiveText := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			if directiveText != "" && !strings.HasPrefix(directiveText, " ") {
				directive, err := parseDirective(line, lineNumber)
				if err != nil {
					return EnvFile{}, fmt.Errorf("failed to parse directive at line %d: %w", lineNumber, err)
				}
				envFile.Directives = append(envFile.Directives, directive)
				continue
			}
		}

		// Skip regular comments
		if strings.HasPrefix(line, "#") {
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

			if key != "" && isValidKey(key) {
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

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle directives (skip in second pass as they're already collected)
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "# ") {
			directiveText := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			if directiveText != "" && !strings.HasPrefix(directiveText, " ") {
				continue
			}
		}

		// Skip regular comments
		if strings.HasPrefix(line, "#") {
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

			if key != "" && isValidKey(key) {
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

// isValidKey checks if a key matches the required regex pattern
func isValidKey(key string) bool {
	matched, _ := regexp.MatchString(`^[A-Za-z_][A-Za-z0-9_]*$`, key)
	return matched
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

// applyFilterDirectives applies filter directives to remove variables based on patterns
func applyFilterDirectives(kvs map[string]string, directives []Directive) map[string]string {
	result := make(map[string]string)

	// Copy existing key-value pairs
	for key, value := range kvs {
		result[key] = value
	}

	// Apply filter directives
	for _, directive := range directives {
		if strings.ToLower(directive.Name) == "filter" {
			applyFilterDirective(result, directive)
		}
	}

	return result
}

// applyFilterUnlessDirectives applies filter-unless directives to keep only variables matching patterns
func applyFilterUnlessDirectives(kvs map[string]string, directives []Directive) map[string]string {
	result := make(map[string]string)

	// Copy existing key-value pairs
	for key, value := range kvs {
		result[key] = value
	}

	// Collect all patterns from all filter-unless directives
	var allPatterns []string
	for _, directive := range directives {
		if strings.ToLower(directive.Name) == "filter-unless" {
			allPatterns = append(allPatterns, directive.Arguments...)
		}
	}

	// If no patterns, keep all keys (no filtering)
	if len(allPatterns) == 0 {
		return result
	}

	// Find keys to keep (those matching any of the patterns)
	keysToKeep := make(map[string]bool)
	for key := range result {
		for _, pattern := range allPatterns {
			if matchesPattern(key, pattern) {
				keysToKeep[key] = true
				fmt.Fprintf(os.Stderr, "DEBUG: Keeping key %q (matches pattern %q)\n", key, pattern)
				break // Key matches at least one pattern, so keep it
			}
		}
	}

	// Remove keys that don't match any pattern
	for key := range result {
		if !keysToKeep[key] {
			fmt.Fprintf(os.Stderr, "DEBUG: Removing key %q (doesn't match any pattern)\n", key)
			delete(result, key)
		}
	}

	return result
}

// applyFilterDirective removes environment variables based on the filter directive
func applyFilterDirective(kvs map[string]string, directive Directive) {
	fmt.Fprintf(os.Stderr, "DEBUG: Applying filter directive: %+v\n", directive)
	for _, arg := range directive.Arguments {
		fmt.Fprintf(os.Stderr, "DEBUG: Filtering with pattern: %q\n", arg)
		// Remove keys matching the pattern (case-insensitive)
		for key := range kvs {
			if matchesPattern(key, arg) {
				fmt.Fprintf(os.Stderr, "DEBUG: Removing key %q (matches pattern %q)\n", key, arg)
				delete(kvs, key)
			}
		}
	}
}

// matchesPattern checks if a key matches a pattern with wildcard support
func matchesPattern(key, pattern string) bool {
	// Convert both to lowercase for case-insensitive matching
	keyLower := strings.ToLower(key)
	patternLower := strings.ToLower(pattern)

	// If pattern contains wildcard, use pattern matching
	if strings.Contains(patternLower, "*") {
		return matchesWildcardPattern(keyLower, patternLower)
	}

	// Otherwise, use exact match
	return keyLower == patternLower
}

// matchesWildcardPattern checks if a key matches a wildcard pattern
func matchesWildcardPattern(key, pattern string) bool {
	// Convert wildcard pattern to regex
	regexPattern := strings.ReplaceAll(pattern, "*", ".*")
	regexPattern = "^" + regexPattern + "$"

	matched, err := regexp.MatchString(regexPattern, key)
	if err != nil {
		// If regex compilation fails, fall back to exact match
		return key == pattern
	}

	return matched
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
