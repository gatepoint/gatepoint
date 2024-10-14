package main

import (
	"fmt"
	"os"

	"github.com/gatepoint/gatepoint/internal/cmd/gatepoint"
)

func main() {
	if err := gatepoint.GetRootCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
