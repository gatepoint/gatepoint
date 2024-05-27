package grpc

import (
	"context"
	"net"

	v1 "github.com/gatepoint/gatepoint/api/gatepoint/v1"

	"github.com/gatepoint/gatepoint/internal/service"
	"github.com/gatepoint/gatepoint/pkg/log"

	"google.golang.org/grpc"
)

func Run(ctx context.Context, network, address string) error {
	//init grpc server and run
	l, err := net.Listen(network, address)
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

	demoService := service.NewDemoService()
	v1.RegisterDemoServer(s, demoService)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	go func() error {
		log.L(ctx).Infof("grpc listen on:%s\n", address)
		if err := s.Serve(l); err != nil {
			return err
		}
		return nil
	}()

	return nil
}
