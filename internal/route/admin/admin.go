package admin

import (
	"context"
	"errors"
	"net/http"

	"github.com/gatepoint/gatepoint/internal/route"
	"github.com/gatepoint/gatepoint/pkg/config"
	"github.com/gatepoint/gatepoint/pkg/filter"
	"github.com/gatepoint/gatepoint/pkg/log"
	swaggerui "github.com/gatepoint/gatepoint/swagger-ui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const swaggerPrefix = "/api/v1/swagger/"

type adminServer struct {
	context.Context
	httpRoute route.HTTPRoute
	serveMux  *runtime.ServeMux
}

func (a adminServer) Run() error {
	mux := http.NewServeMux()
	if config.EnableSwagger() {
		mux.Handle(swaggerPrefix, http.StripPrefix(swaggerPrefix, http.FileServer(http.FS(swaggerui.Resources))))
	}

	handler := filter.RecordFilter(mux)

	s := http.Server{
		Addr:    config.GetHttpAddr(),
		Handler: handler,
	}

	log.Infof("Starting admin listening at %s", config.GetAdminAddr())

	if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("Failed to start http server: %v", err)
		return err
	}
	return nil
}

func NewAdminServer() {

}
