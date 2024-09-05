package log

//
//import (
//	"context"
//	"io"
//	"os"
//
//	commonv1 "github.com/gatepoint/gatepoint/api/common/v1"
//	"github.com/go-logr/logr"
//	"github.com/go-logr/zapr"
//	"go.uber.org/zap"
//	"go.uber.org/zap/zapcore"
//)
//
//var LogLevelMap = map[commonv1.LogLevel]zapcore.Level{
//	commonv1.LogLevel_LOG_LEVEL_UNSPECIFIED: zapcore.InfoLevel,
//	commonv1.LogLevel_LOG_LEVEL_DEBUG:       zapcore.DebugLevel,
//	commonv1.LogLevel_LOG_LEVEL_INFO:        zapcore.InfoLevel,
//	commonv1.LogLevel_LOG_LEVEL_WARN:        zapcore.WarnLevel,
//	commonv1.LogLevel_LOG_LEVEL_ERROR:       zapcore.ErrorLevel,
//}
//
//var LogComponentMap = map[commonv1.LogComponent]string{
//	commonv1.LogComponent_LOG_COMPONENT_DEFAULT:     "default",
//	commonv1.LogComponent_LOG_COMPONENT_GATEWAY:     "gateway",
//	commonv1.LogComponent_LOG_COMPONENT_DOMAIN:      "domain",
//	commonv1.LogComponent_LOG_COMPONENT_API:         "api",
//	commonv1.LogComponent_LOG_COMPONENT_SERVICE:     "service",
//	commonv1.LogComponent_LOG_COMPONENT_PLUGIN:      "plugin",
//	commonv1.LogComponent_LOG_COMPONENT_UNSPECIFIED: "default",
//}
//
//type Logger struct {
//	logr.Logger
//	logging       map[commonv1.LogComponent]commonv1.LogLevel
//	sugaredLogger *zap.SugaredLogger
//}
//
//func NewLogger() Logger {
//	logger := initZapLogger(os.Stdout, DefaultLogging(), commonv1.LogLevel_LOG_LEVEL_INFO)
//
//	return Logger{
//		Logger:        zapr.NewLogger(logger),
//		logging:       DefaultLogging(),
//		sugaredLogger: logger.Sugar(),
//	}
//}
//
//func DefaultLogger(level commonv1.LogLevel) Logger {
//	logger := initZapLogger(os.Stdout, logging, level)
//
//	return Logger{
//		Logger:        zapr.NewLogger(logger),
//		logging:       DefaultLogging(),
//		sugaredLogger: logger.Sugar(),
//	}
//}
//
//func DefaultLogging() map[commonv1.LogComponent]commonv1.LogLevel {
//	return map[commonv1.LogComponent]commonv1.LogLevel{
//		commonv1.LogComponent_LOG_COMPONENT_DEFAULT: commonv1.LogLevel_LOG_LEVEL_INFO,
//		commonv1.LogComponent_LOG_COMPONENT_GATEWAY: commonv1.LogLevel_LOG_LEVEL_INFO,
//		commonv1.LogComponent_LOG_COMPONENT_DOMAIN:  commonv1.LogLevel_LOG_LEVEL_INFO,
//		commonv1.LogComponent_LOG_COMPONENT_API:     commonv1.LogLevel_LOG_LEVEL_INFO,
//		commonv1.LogComponent_LOG_COMPONENT_SERVICE: commonv1.LogLevel_LOG_LEVEL_INFO,
//		commonv1.LogComponent_LOG_COMPONENT_PLUGIN:  commonv1.LogLevel_LOG_LEVEL_INFO,
//	}
//}
//
//var std *zap.SugaredLogger
//
//func New() (*zap.SugaredLogger, error) {
//	config := zap.NewDevelopmentConfig()
//	config.OutputPaths = []string{"stdout"}
//	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
//	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	config.DisableStacktrace = false
//	config.DisableCaller = false
//
//	log, err := config.Build()
//	if err != nil {
//		return nil, err
//	}
//
//	std = log.Sugar()
//	return std, nil
//}
//
//func L(ctx context.Context) *zap.SugaredLogger {
//	copy := *std
//	lg := &copy
//
//	if requestID := ctx.Value("requestId"); requestID != nil {
//		lg = lg.With(zap.Any("requestId", requestID))
//	}
//
//	return lg
//}
//
//func DefaultLoggingLevel(logging map[commonv1.LogComponent]commonv1.LogLevel, level commonv1.LogLevel) zapcore.Level {
//	if level != commonv1.LogLevel_LOG_LEVEL_UNSPECIFIED {
//		return LogLevelMap[level]
//	}
//	if logging[commonv1.LogComponent_LOG_COMPONENT_DEFAULT] != commonv1.LogLevel_LOG_LEVEL_UNSPECIFIED {
//		return LogLevelMap[logging[commonv1.LogComponent_LOG_COMPONENT_DEFAULT]]
//	}
//	return LogLevelMap[commonv1.LogLevel_LOG_LEVEL_INFO]
//}
//
//func initZapLogger(w io.Writer, logging map[commonv1.LogComponent]commonv1.LogLevel, level commonv1.LogLevel) *zap.Logger {
//	core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(w), zap.NewAtomicLevelAt(DefaultLoggingLevel(logging, level)))
//
//	return zap.New(core, zap.AddCaller())
//}
//
//func (l Logger) WithComponent(component commonv1.LogComponent) Logger {
//	logLevel := l.logging[component]
//	logger := initZapLogger(os.Stdout, l.logging, logLevel)
//
//	return Logger{
//		Logger:        zapr.NewLogger(logger).WithName(LogComponentMap[component]),
//		logging:       l.logging,
//		sugaredLogger: logger.Sugar(),
//	}
//}
//
//func (l Logger) WithValues(keysAndValues ...interface{}) Logger {
//	l.Logger = l.Logger.WithValues(keysAndValues...)
//	return l
//}
