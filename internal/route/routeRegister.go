package route

import (
	"context"

	v1 "github.com/gatepoint/gatepoint/api/gatepoint/v1"
	"github.com/gatepoint/gatepoint/internal/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RegisterHTTPRoutes() []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error {
	return []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error{
		v1.RegisterDemoServiceHandler,
	}
}

func RegisterGRPCRoutes(s *grpc.Server) {
	//todo add your grpc server register here
	v1.RegisterDemoServiceServer(s, service.NewDemoService())
	reflection.Register(s)
}