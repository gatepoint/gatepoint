package main

import (
	"fmt"
	"os"

	"github.com/gatepoint/gatepoint/internal/gpctl"
)

func main() {
	if err := gpctl.GetRootCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
