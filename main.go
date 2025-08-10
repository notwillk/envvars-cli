package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/notwillk/envvars-cli/commands"
	"github.com/spf13/pflag"
)

func main() {
	// Define flags
	var help bool
	var version bool
	var filePaths []string
	var format string
	var jsonFile string
	var yamlFile string
	var verbose bool

	// Set up flags
	pflag.BoolVarP(&help, "help", "h", false, "Show this help message")
	pflag.BoolVarP(&version, "version", "v", false, "Show version information")
	pflag.StringSliceVarP(&filePaths, "env", "e", []string{}, "Read and parse environment variable files (can be specified multiple times)")
	pflag.StringVarP(&format, "format", "f", "env", "Output format: json, yaml, or env (default: env)")
	pflag.StringVarP(&jsonFile, "json", "j", "", "Process a JSON file")
	pflag.StringVarP(&yamlFile, "yaml", "y", "", "Process a YAML file")
	pflag.BoolVarP(&verbose, "verbose", "V", false, "Enable verbose output")

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
	if len(filePaths) > 0 || jsonFile != "" || yamlFile != "" {
		// Create sources array with metadata
		var sources []commands.Source
		priority := 0

		// Process flags in the order they appear in the command line
		// This preserves the user's intended priority order
		for i := 1; i < len(os.Args); i++ {
			arg := os.Args[i]

			switch arg {
			case "--env", "-e":
				// Find the corresponding file path
				if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
					sources = append(sources, commands.Source{
						FilePath: os.Args[i+1],
						Type:     "env",
						Priority: priority,
					})
					priority++
					i++ // Skip the file path in next iteration
				}
			case "--json", "-j":
				// Find the corresponding file path
				if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
					sources = append(sources, commands.Source{
						FilePath: os.Args[i+1],
						Type:     "json",
						Priority: priority,
					})
					priority++
					i++ // Skip the file path in next iteration
				}
			case "--yaml", "-y":
				// Find the corresponding file path
				if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
					sources = append(sources, commands.Source{
						FilePath: os.Args[i+1],
						Type:     "yaml",
						Priority: priority,
					})
					priority++
					i++ // Skip the file path in next iteration
				}
			}
		}

		// Create global options
		options := commands.Options{
			Verbose: verbose,
			Format:  format,
		}

		mergeCmd := commands.CreateMergeCommand(sources, options)
		if err := mergeCmd.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// If no valid command provided, show help
	commands.ShowHelp()
}
