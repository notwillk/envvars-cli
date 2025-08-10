package main

import (
	"fmt"
	"os"

	"github.com/notwillk/envvars-cli/commands"
	"github.com/spf13/pflag"
)

func main() {
	// Define flags
	var help bool
	var version bool
	var filePaths []string
	var format string

	// Set up flags
	pflag.BoolVarP(&help, "help", "h", false, "Show this help message")
	pflag.BoolVarP(&version, "version", "v", false, "Show version information")
	pflag.StringSliceVarP(&filePaths, "env", "e", []string{}, "Read and parse environment variable files (can be specified multiple times)")
	pflag.StringVarP(&format, "format", "f", "env", "Output format: json, yaml, or env (default: env)")

	// Parse flags
	pflag.Parse()

	// Handle help flag
	if help || (len(os.Args) == 1 && pflag.NFlag() == 0) {
		commands.ShowHelp()
		return
	}

	// Handle version flag
	if version {
		commands.ShowVersion()
		return
	}

	// Handle env flags (environment processor command)
	if len(filePaths) > 0 {
		envProcessorCmd := commands.NewEnvProcessorCommand(filePaths, format)
		if err := envProcessorCmd.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// If no valid command provided, show help
	commands.ShowHelp()
}
