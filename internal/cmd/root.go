package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	httpAddr string
	grpcAddr string
)

func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gatepoint",
		Short: "Standard HTTP and GRPC Go Project Layout with Protobuf and GORM",
	}

	cmd.AddCommand(getServerCommand())
	cmd.AddCommand(getVersionCommand())

	return cmd
}
