package commands

import (
	"fmt"
	"os"
)

// ShowHelp displays the help message for the CLI
func ShowHelp() {
	fmt.Fprintf(os.Stdout, `envvars-cli - Environment Variable File Processor

USAGE:
    envvars-cli [COMMAND] [OPTIONS]

COMMANDS:
    help, -h, --help     Show this help message
    version, -v, --version  Show version information
    merge               Process and merge environment variable files

OPTIONS:
    -e, --env <file>     Read and parse environment variable files (can be specified multiple times)
    -f, --format <fmt>   Output format: json, yaml, or env (default: env)
    -j, --json <file>    Process a JSON file
    -y, --yaml <file>    Process a YAML file
    -V, --verbose        Enable verbose output

EXAMPLES:
    # Parse a single environment file (default ENV format)
    envvars-cli --env config.env

    # Parse multiple environment files
    envvars-cli --env dev.env --env prod.env

    # Output as JSON
    envvars-cli --env config.env --format json

    # Output as YAML
    envvars-cli --env config.env --format yaml

    # Output as ENV (default)
    envvars-cli --env config.env --format env

    # Process JSON files
    envvars-cli --json config.json
    envvars-cli --json config.json --format yaml

    # Process YAML files
    envvars-cli --yaml config.yaml
    envvars-cli --yaml config.yaml --format json

    # Mix different file types
    envvars-cli --env config.env --json config.json --yaml config.yaml

    # Show help
    envvars-cli --help

    # Show version
    envvars-cli --version

DESCRIPTION:
    envvars-cli is a command-line tool for parsing and processing environment variable files.
    It supports parsing .env, .json, and .yaml files with comments, quoted values, and variable references.
    Multiple files can be processed, with later files taking precedence over earlier ones.
`)
}
