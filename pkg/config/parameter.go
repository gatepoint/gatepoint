package config

//import (
//	commonv1 "github.com/gatepoint/gatepoint/api/common/v1"
//	"github.com/gatepoint/gatepoint/pkg/log"
//	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
//)
//
//type Parameters struct {
//	Network    string
//	GRPCAddr   string
//	HTTPAddr   string
//	OpenAPIDir string
//	Mux        []runtime.ServeMuxOption
//	KubeConfig string
//	//Logging    log.Logger
//	Logging map[commonv1.LogComponent]commonv1.LogLevel
//	Logger  log.Logger
//}
//
//func Default() Parameters {
//	return Parameters{
//		Network:    "tcp",
//		GRPCAddr:   "0.0.0.0:9091",
//		HTTPAddr:   "0.0.0.0:8081",
//		OpenAPIDir: "api/v1",
//		Mux:        nil,
//		KubeConfig: "",
//		Logging:    log.DefaultLogging(),
//		Logger:     log.DefaultLogger(commonv1.LogLevel_LOG_LEVEL_INFO),
//	}
//}
//
//func DefaultLogging() map[commonv1.LogComponent]commonv1.LogLevel {
//	return map[commonv1.LogComponent]commonv1.LogLevel{
//		commonv1.LogComponent_LOG_COMPONENT_DEFAULT: commonv1.LogLevel_LOG_LEVEL_INFO,
//	}
//}
//
//func (p *Parameters) logging() log.Logger {
//	return p.Logger.WithComponent(commonv1.LogComponent_LOG_COMPONENT_DEFAULT)
//}
