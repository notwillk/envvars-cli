// Example usage of the ProcessFileWithMerge function
// This file demonstrates how to use the merge functionality programmatically

package merge_example

import (
	"fmt"
	"log"

	formatters "github.com/notwillk/envvars-cli/formatters"
	"github.com/notwillk/envvars-cli/sources"
)

func main() {
	// Existing key-value pairs
	existingKVs := map[string]string{
		"APP_NAME":  "oldapp",
		"DEBUG":     "false",
		"EXTRA_VAR": "test",
	}

	// Options containing file path
	options := sources.Options{
		FilePath: "testdata/basic.env",
	}

	fmt.Println("Existing key-value pairs:")
	for k, v := range existingKVs {
		fmt.Printf("  %s: %s\n", k, v)
	}

	fmt.Println("\nMerging with file:", options.FilePath)
	fmt.Println("Result (file values take precedence):")

	// Call the merge function from the sources package
	mergedVars, err := sources.ProcessFileWithMerge(existingKVs, options)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Output the result using the formatters package
	if err := formatters.OutputAsJSON(mergedVars); err != nil {
		log.Fatalf("Error outputting JSON: %v", err)
	}
}
