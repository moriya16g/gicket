package main

import (
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/cmd"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
