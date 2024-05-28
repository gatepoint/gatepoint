package grpc

import (
	"context"
	"net"

	v1 "github.com/gatepoint/gatepoint/api/gatepoint/v1"
	"github.com/gatepoint/gatepoint/internal/gateway"
	"github.com/gatepoint/gatepoint/pkg/kubernetes"

	"github.com/gatepoint/gatepoint/internal/service"
	"github.com/gatepoint/gatepoint/pkg/log"

	"google.golang.org/grpc"
)

func Run(ctx context.Context, opts gateway.Options) error {
	//init grpc server and run
	l, err := net.Listen(opts.Network, opts.GRPCAddr)
	if err != nil {
		return err
	}
	go func() {
		defer func() error {
			if err := l.Close(); err != nil {
				return err
			}
			return nil
		}()
		<-ctx.Done()
	}()
	s := grpc.NewServer()

	clientset, err := kubernetes.GetClientset(opts.KubeConfig)
	if err != nil {
		return err
	}

	// Reg services
	demoService := service.NewDemoService(clientset)
	v1.RegisterDemoServer(s, demoService)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	go func() error {
		log.L(ctx).Infof("grpc listen on:%s\n", opts.GRPCAddr)
		if err := s.Serve(l); err != nil {
			return err
		}
		return nil
	}()

	return nil
}
