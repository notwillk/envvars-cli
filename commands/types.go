package commands

// Source represents a single source file with its metadata
type Source struct {
	FilePath string
	Type     string // "env", "json", "yaml", "sops"
	Priority int    // Higher priority sources override lower ones
	// For SOPS sources, additional metadata
	DecryptionKey string // The key to use for decryption (only for SOPS type)
}

// Options represents global options for the merge command
type Options struct {
	Verbose bool
	Format  string // "json", "yaml", "env"
}
