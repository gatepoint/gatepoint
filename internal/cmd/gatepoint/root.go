package gatepoint

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	namespaces string

	httpAddr string
	grpcAddr string
)

func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gatepoint",
		Short: "Standard HTTP and GRPC Go Project Layout with Protobuf and GORM",
	}

	cmd.AddCommand(getServerCommand())

	return cmd
}
