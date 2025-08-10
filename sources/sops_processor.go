package sources

import (
	"fmt"
	"os"
	"strings"

	"github.com/getsops/sops/v3/decrypt"
	"gopkg.in/yaml.v3"
)

// SOPSProcessor handles processing of SOPS-encrypted files
type SOPSProcessor struct{}

// CreateSOPSProcessor creates a new SOPS processor instance
func CreateSOPSProcessor() *SOPSProcessor {
	return &SOPSProcessor{}
}

// ProcessFile decrypts a SOPS-encrypted file and returns the key-value pairs
func (p *SOPSProcessor) ProcessFile(filePath string, decryptionKey string) ([]EnvVar, error) {
	// Read the encrypted file
	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SOPS file: %w", err)
	}

	// Decrypt the file using SOPS
	decryptedData, err := decrypt.Data(encryptedData, "yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt SOPS file: %w", err)
	}

	// Parse the decrypted YAML content
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(decryptedData, &yamlData); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted YAML: %w", err)
	}

	// Convert to key-value pairs
	var variables []EnvVar
	p.flattenMap("", yamlData, &variables)

	return variables, nil
}

// flattenMap recursively flattens a nested map into key-value pairs
func (p *SOPSProcessor) flattenMap(prefix string, data map[string]interface{}, variables *[]EnvVar) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "_" + key
		}

		switch v := value.(type) {
		case string:
			*variables = append(*variables, EnvVar{
				Key:   strings.ToUpper(fullKey),
				Value: v,
			})
		case bool:
			*variables = append(*variables, EnvVar{
				Key:   strings.ToUpper(fullKey),
				Value: fmt.Sprintf("%t", v),
			})
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			*variables = append(*variables, EnvVar{
				Key:   strings.ToUpper(fullKey),
				Value: fmt.Sprintf("%v", v),
			})
		case float32, float64:
			*variables = append(*variables, EnvVar{
				Key:   strings.ToUpper(fullKey),
				Value: fmt.Sprintf("%v", v),
			})
		case map[string]interface{}:
			p.flattenMap(fullKey, v, variables)
		case []interface{}:
			// Convert arrays to comma-separated strings
			var strValues []string
			for _, item := range v {
				strValues = append(strValues, fmt.Sprintf("%v", item))
			}
			*variables = append(*variables, EnvVar{
				Key:   strings.ToUpper(fullKey),
				Value: strings.Join(strValues, ","),
			})
		default:
			// For any other type, convert to string
			*variables = append(*variables, EnvVar{
				Key:   strings.ToUpper(fullKey),
				Value: fmt.Sprintf("%v", v),
			})
		}
	}
}
