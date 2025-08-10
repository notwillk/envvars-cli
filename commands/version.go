package commands

import (
	"fmt"
	"os"
)

// ShowVersion displays the version information for the CLI
func ShowVersion() {
	fmt.Fprintf(os.Stdout, "envvars-cli v1.0.0\n")
	fmt.Fprintf(os.Stdout, "Environment Variable File Processor\n")
	fmt.Fprintf(os.Stdout, "Built with Go\n")
}
