package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	// Define flags
	var help bool
	var version bool

	// Set up flags
	pflag.BoolVarP(&help, "help", "h", false, "Show this help message")
	pflag.BoolVarP(&version, "version", "v", false, "Show version information")

	// Parse flags
	pflag.Parse()

	// Handle help flag
	if help {
		showHelp()
		return
	}

	// Handle version flag
	if version {
		showVersion()
		return
	}

	// If no flags provided, show help by default
	if len(os.Args) == 1 {
		showHelp()
		return
	}

	// TODO: Implement your CLI logic here

	// Example: Print current working directory
	if wd, err := os.Getwd(); err == nil {
		fmt.Printf("Working directory: %s\n", wd)
	}
}

func showHelp() {
	fmt.Println("envvars-cli - Environment Variables CLI Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  envvars-cli [flags] [command]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -h, --help      Show this help message")
	fmt.Println("  -v, --version   Show version information")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  TODO: Add your CLI commands here")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  envvars-cli --help")
	fmt.Println("  envvars-cli --version")
}

func showVersion() {
	fmt.Println("envvars-cli v0.1.0")
}
