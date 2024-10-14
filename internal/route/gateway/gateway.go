package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gatepoint/gatepoint/internal/route"
	"github.com/gatepoint/gatepoint/pkg/config"
	"github.com/gatepoint/gatepoint/pkg/filter"
	"github.com/gatepoint/gatepoint/pkg/health"
	"github.com/gatepoint/gatepoint/pkg/log"
	swaggerui "github.com/gatepoint/gatepoint/swagger-ui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

const swaggerPrefix = "/api/v1/swagger/"

type gatewayServer struct {
	ctx       context.Context
	httpRoute route.HTTPRoute
	serveMux  *runtime.ServeMux
}

func (g gatewayServer) Run() error {
	conn, err := grpc.NewClient(config.GetGrpcAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		_ = g.Stop()
		if err = conn.Close(); err != nil {
			log.Errorf("grpc client close error:%v\n", err)
			return
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/grpcHealthz", grpcHealthzServer(conn))
	mux.Handle("/sys/", httpHealthzServer())

	if config.EnableSwagger() {
		mux.Handle(swaggerPrefix, http.StripPrefix(swaggerPrefix, http.FileServer(http.FS(swaggerui.Resources))))
	}

	err = g.newGateway(conn)
	if err != nil {
		return err
	}
	mux.Handle("/", g.serveMux)

	handler := filter.RecordFilter(mux)
	handler = filter.AuthFilter(handler)
	s := http.Server{
		Addr:    config.GetHttpAddr(),
		Handler: handler,
	}

	log.Infof("Starting listening at %s", config.GetHttpAddr())

	if err = s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("Failed to start http server: %v", err)
		return err
	}
	return nil
}

func (g gatewayServer) newGateway(conn *grpc.ClientConn) error {
	for _, f := range g.httpRoute() {
		if err := f(g.ctx, g.serveMux, conn); err != nil {
			return err
		}
	}
	return nil
}

func httpHealthzServer() http.Handler {
	handler := health.NewHandler()
	// Add more readiness checks
	handler.AddLivenessCheck("goroutine-threshold", health.GoroutineCountCheck(500))

	return handler
}

func grpcHealthzServer(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "ok")
	}
}

func (g gatewayServer) Stop() error {
	<-g.ctx.Done()
	return nil
}

func NewGatewayServer(ctx context.Context, httpRoute route.HTTPRoute, opt func() []runtime.ServeMuxOption) route.Server {
	return gatewayServer{
		ctx:       ctx,
		httpRoute: httpRoute,
		serveMux:  runtime.NewServeMux(opt()...),
	}
}
