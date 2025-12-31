package main

import (
	"fmt"
	"os"

	"github.com/ansoncodes/workshot/internal/cli"
)

var version = "dev" // set by goreleaser

func main() {
	// pass version to cli
	cli.SetVersion(version)

	// run cli and handle error
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
