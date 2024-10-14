package gpctl

import "github.com/spf13/cobra"

func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gpctl",
		Short: "gpctl is the command line interface for Gatepoint",
	}

	cmd.AddCommand(
		getVersionCommand(),
		//todo add more commands
	)
	return cmd
}
