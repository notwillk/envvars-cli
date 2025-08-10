package formatters

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// OutputAsENV outputs the key-value pairs in environment variable format to stdout
func OutputAsENV(variables map[string]string) error {
	// Sort keys for consistent output
	keys := make([]string, 0, len(variables))
	for k := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Output as environment variables
	for _, key := range keys {
		value := variables[key]
		// Escape the value if it contains special characters
		escapedValue := escapeEnvValue(value)
		fmt.Fprintf(os.Stdout, "%s=%s\n", key, escapedValue)
	}

	return nil
}

// escapeEnvValue escapes special characters in environment variable values
func escapeEnvValue(value string) string {
	if value == "" {
		return ""
	}

	// If the value contains spaces, quotes, or special characters, wrap it in quotes
	if strings.ContainsAny(value, " \t\n\r\"'\\$`") {
		// Escape backslashes and quotes
		escaped := strings.ReplaceAll(value, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		return "\"" + escaped + "\""
	}

	return value
}
