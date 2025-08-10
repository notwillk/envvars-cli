// Example usage of the formats package
// This file demonstrates how to use the JSON output functionality

package format_example

import (
	formatters "github.com/notwillk/envvars-cli/formatters"
)

func main() {
	// Sample key-value pairs
	sampleKVs := map[string]string{
		"NAME":     "John Doe",
		"AGE":      "30",
		"CITY":     "New York",
		"COUNTRY":  "USA",
		"LANGUAGE": "Go",
	}

	// Output as formatted JSON
	println("Formatted JSON output:")
	if err := formatters.OutputAsJSON(sampleKVs); err != nil {
		panic(err)
	}

	println("\nCompact JSON output:")
	if err := formatters.OutputAsJSONCompact(sampleKVs); err != nil {
		panic(err)
	}
}
