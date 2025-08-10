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
    env processor        Process environment variable files and output as JSON

OPTIONS:
    -e, --env <file>     Read and parse environment variable files (can be specified multiple times)

EXAMPLES:
    # Parse a single environment file
    envvars-cli --env config.env

    # Parse multiple environment files
    envvars-cli --env dev.env --env prod.env

    # Show help
    envvars-cli --help

    # Show version
    envvars-cli --version

DESCRIPTION:
    envvars-cli is a command-line tool for parsing and processing environment variable files.
    It supports parsing .env files with comments, quoted values, and variable references.
    Multiple files can be processed, with later files taking precedence over earlier ones.
`)
}
