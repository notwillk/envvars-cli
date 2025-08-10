package formatters

import (
	"encoding/json"
	"os"
)

// OutputAsJSON outputs the given key-value pairs as JSON to stdout
func OutputAsJSON(kvs map[string]string) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(kvs)
}

// OutputAsJSONCompact outputs the given key-value pairs as compact JSON to stdout
func OutputAsJSONCompact(kvs map[string]string) error {
	encoder := json.NewEncoder(os.Stdout)
	return encoder.Encode(kvs)
}
