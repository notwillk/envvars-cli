package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	// Define flags
	var help bool
	var version bool
	var filePaths []string

	// Set up flags
	pflag.BoolVarP(&help, "help", "h", false, "Show this help message")
	pflag.BoolVarP(&version, "version", "v", false, "Show version information")
	pflag.StringSliceVarP(&filePaths, "file", "f", []string{}, "Read and output the contents of files (can be specified multiple times)")

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

	// Handle file flags
	if len(filePaths) > 0 {
		if err := readAndOutputFiles(filePaths); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading files: %v\n", err)
			os.Exit(1)
		}
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

func readAndOutputFiles(filePaths []string) error {
	for i, filePath := range filePaths {
		// Add separator between files if multiple files
		if i > 0 {
			fmt.Println("---")
		}

		// Show filename as header if multiple files
		if len(filePaths) > 1 {
			fmt.Printf("File: %s\n", filePath)
			fmt.Println("---")
		}

		if err := readAndOutputFile(filePath); err != nil {
			return fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}

		// Add newline after file content
		fmt.Println()
	}

	return nil
}

func readAndOutputFile(filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	// Copy file contents to stdout
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return fmt.Errorf("failed to read file '%s': %w", filePath, err)
	}

	return nil
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
	fmt.Println("  -f, --file      Read and output the contents of files (can be specified multiple times)")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  TODO: Add your CLI commands here")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  envvars-cli --help")
	fmt.Println("  envvars-cli --version")
	fmt.Println("  envvars-cli --file example.txt")
	fmt.Println("  envvars-cli -f /path/to/file.txt")
	fmt.Println("  envvars-cli --file file1.txt --file file2.txt")
	fmt.Println("  envvars-cli -f file1.txt -f file2.txt -f file3.txt")
}

func showVersion() {
	fmt.Println("envvars-cli v0.1.0")
}
