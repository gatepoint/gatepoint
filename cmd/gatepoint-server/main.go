package main

import (
	"fmt"
	"os"

	"github.com/gatepoint/gatepoint/internal/cmd"
)

func main() {
	if err := cmd.GetRootCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
