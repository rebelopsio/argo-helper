package main

import (
	"fmt"
	"os"

	"github.com/rebelopsio/argo-helper/cmd"
	"github.com/rebelopsio/argo-helper/tui"
)

func main() {
	// If no args are provided, start the TUI
	if len(os.Args) == 1 {
		if err := tui.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Otherwise, use the CLI as before
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
