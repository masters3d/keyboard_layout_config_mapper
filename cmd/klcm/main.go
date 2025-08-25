package main

import (
	"fmt"
	"os"

	"masters3d.com/keyboard_layout_config_mapper/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}