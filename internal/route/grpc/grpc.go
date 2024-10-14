package grpc

import (
	"context"
	"net"

	"github.com/gatepoint/gatepoint/internal/route"
	"github.com/gatepoint/gatepoint/pkg/config"
	"github.com/gatepoint/gatepoint/pkg/log"

	"google.golang.org/grpc"
)

type grpcServer struct {
	routes route.GrpcRoute
	server *grpc.Server
	ctx    context.Context
}

func (g grpcServer) Run() error {
	g.routes(g.server)
	l, err := net.Listen("tcp", config.GetGrpcAddr())
	if err != nil {
		<-g.ctx.Done()
		return err
	}
	defer func() {
		_ = g.Stop()
		if err := l.Close(); err != nil {
			log.Errorf("grpc server close error:%v\n", err)
		}
	}()

	log.Infof("Starting listening grpc at %s", config.GetGrpcAddr())

	return g.server.Serve(l)
}

func (g grpcServer) Stop() error {
	g.server.GracefulStop()

	<-g.ctx.Done()
	return nil
}

func NewGrpcServer(ctx context.Context, grpcRoute route.GrpcRoute, opt func() []grpc.ServerOption) route.Server {
	return grpcServer{
		routes: grpcRoute,
		server: grpc.NewServer(opt()...),
		ctx:    ctx,
	}
}
