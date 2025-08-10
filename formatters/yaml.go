package formatters

import (
	"fmt"
	"os"
	"sort"
)

// OutputAsYAML outputs the key-value pairs as YAML to stdout
func OutputAsYAML(variables map[string]string) error {
	// Sort keys for consistent output
	keys := make([]string, 0, len(variables))
	for k := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Output as YAML
	for _, key := range keys {
		value := variables[key]
		// Escape quotes and special characters if needed
		if needsQuoting(value) {
			fmt.Fprintf(os.Stdout, "%s: %q\n", key, value)
		} else {
			fmt.Fprintf(os.Stdout, "%s: %s\n", key, value)
		}
	}

	return nil
}

// needsQuoting determines if a value needs to be quoted in YAML
func needsQuoting(value string) bool {
	if value == "" {
		return false
	}

	// Check for special characters that require quoting
	for _, char := range value {
		if char == ':' || char == '{' || char == '}' || char == '[' || char == ']' ||
			char == ',' || char == '&' || char == '*' || char == '#' || char == '?' ||
			char == '|' || char == '>' || char == '!' || char == '%' || char == '@' ||
			char == '`' || char == '"' || char == '\'' || char == '\\' {
			return true
		}
	}

	// Check if it looks like a number or boolean
	if value == "true" || value == "false" || value == "null" || value == "yes" || value == "no" {
		return true
	}

	// Check if it starts with special characters
	if len(value) > 0 && (value[0] == '-' || value[0] == '?' || value[0] == ':' || value[0] == '[' || value[0] == '{') {
		return true
	}

	return false
}
