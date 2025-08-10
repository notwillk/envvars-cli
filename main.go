package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("envvars-cli - Environment Variables CLI Tool")
	
	// TODO: Implement your CLI logic here
	
	// Example: Print current working directory
	if wd, err := os.Getwd(); err == nil {
		fmt.Printf("Working directory: %s\n", wd)
	}
}

