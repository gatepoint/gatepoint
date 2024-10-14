package gpctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"versions", "v"},
		Short:   "Print the version number of gpctl",
		Run: func(cmd *cobra.Command, args []string) {
			// Print the version number of gpctl
			// todo
			fmt.Println("gpctl version 0.0.1")
		},
	}

	return cmd
}
