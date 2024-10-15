package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gatepoint/gatepoint/pkg/interceptor"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"

	"github.com/gatepoint/gatepoint/internal/route"
	gatewayServer "github.com/gatepoint/gatepoint/internal/route/gateway"
	grpcServer "github.com/gatepoint/gatepoint/internal/route/grpc"
	"github.com/gatepoint/gatepoint/pkg/config"
	"github.com/gatepoint/gatepoint/pkg/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// cfgPath is the path to the EnvoyGateway configuration file.
	cfgPath string
)

func getServerCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "Run both grpc and http server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.LoadConfig(cfgPath); err != nil {
				return err
			}
			log.Init()
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Infof("FLAG: --%s=%q", flag.Name, flag.Value)
			})

			rootCtx, cancel := context.WithCancel(context.Background())

			defer func() {
				log.Info("wait for exit signal")
				sig := WaitSignal()
				log.Infof("gatepoint is exiting now, because of the signal: %s", sig)
				cancel()
				log.Flush()
			}()

			//options := utils.Options{
			//	HTTPAddr:   viper.GetString("server.http.addr"),
			//	GRPCAddr:   viper.GetString("server.grpc.addr"),
			//	Network:    "tcp",
			//	OpenAPIDir: "api/v1",
			//	KubeConfig: viper.GetString("kubernetes.kubeconfig"),
			//}
			newGrpcServer := grpcServer.NewGrpcServer(rootCtx, route.RegisterGRPCRoutes, grpcServerOption)
			go func() {
				if err := newGrpcServer.Run(); err != nil {
					log.Fatalf("grpc server run error: %s", err)
				}
			}()

			newGatewayServer := gatewayServer.NewGatewayServer(rootCtx, route.RegisterHTTPRoutes, serverMuxOption)
			return newGatewayServer.Run()
			//if err := gateway.Run(rootCtx, options, serverMuxOption); err != nil {
			//	log.Fatalf("grpc gateway start error: %s", err)
			//}
			//

		},
	}
	//cobra.OnInitialize(InitConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config/config.yaml", "config file (default is $HOME/.gatepoint.yaml)")
	//rootCmd.PersistentFlags().StringVarP(&httpAddr, "http-addr", "", ":8080", "HTTP listen address.")
	//rootCmd.PersistentFlags().StringVarP(&grpcAddr, "grpc-addr", "", ":8081", "GRPC listen address.")
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config-path", "g", "./config/config.yaml",
		"The path to the configuration file.")

	return rootCmd
}

func grpcServerOption() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(interceptor.Auth), selector.MatchFunc(interceptor.AllButHealthZ)),
		),
		grpc.ChainStreamInterceptor(
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(interceptor.Auth), selector.MatchFunc(interceptor.AllButHealthZ)),
		),
	}
}

func serverMuxOption() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{runtime.WithErrorHandler(route.HttpErrorHandler)}
}

func WaitSignal() string {
	sigsCh := make(chan os.Signal, 1)
	signal.Notify(sigsCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigsCh
	return sig.String()
}

//func InitConfig() {
//	if cfgFile != "" {
//		viper.SetConfigType("yaml")
//		viper.SetConfigFile(cfgFile)
//		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
//	} else {
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".gatepoint")
//	}
//
//	viper.AutomaticEnv()
//
//	if err := viper.ReadInConfig(); err != nil {
//		panic(fmt.Sprintf("Failed to read in the config file %s: %v", viper.ConfigFileUsed(), err))
//	}
//	fmt.Println("Using config file:", viper.ConfigFileUsed())
//
//	if err := viper.Unmarshal(&config.Shared, func(decoderConfig *mapstructure.DecoderConfig) {
//		decoderConfig.TagName = "json"
//	}); err != nil {
//		panic(fmt.Sprintf("Failed to unmarshal configuration from disk: %v", err))
//	}
//
//}
