package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gatepoint/gatepoint/pkg/errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server interface {
	Run() error
	Stop() error
}

type HTTPRoute func() []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error

type GrpcRoute func(*grpc.Server)

func HttpErrorHandler(ctx context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, writer http.ResponseWriter, _ *http.Request, err error) {
	e := errors.ToGatepointError(err)
	s, ee := json.Marshal(&e)
	fmt.Println(ee)
	writer.WriteHeader(e.HTTPCode())
	writer.Write(s)
}
