package commands

// Source represents a single source file with its metadata
type Source struct {
	FilePath string
	Type     string // "env", "json", "yaml"
	Priority int    // Higher priority sources override lower ones
}

// Options represents global options for the merge command
type Options struct {
	Verbose bool
	Format  string // "json", "yaml", "env"
}
