package cmd

import (
	"github.com/gatepoint/gatepoint/pkg/build"
	"github.com/spf13/cobra"
)

func getVersionCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"versions", "v"},
		Short:   "Show versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.Print(cmd.OutOrStdout(), output)
		},
	}

	cmd.PersistentFlags().StringVarP(&output, "output", "o", "", "One of 'yaml' or 'json")

	return cmd

}
