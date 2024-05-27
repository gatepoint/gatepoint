package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gatepoint",
		Short: "Standard HTTP and GRPC Go Project Layout with Protobuf and GORM",
	}

	cmd.AddCommand(GetServerCommand())

	return cmd
}
