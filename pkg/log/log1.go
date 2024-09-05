package log

import (
	"github.com/gatepoint/gatepoint/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type infoLogger struct {
	level zapcore.Level
	log   *zap.Logger
}

// zapLogger is a logr.Logger that uses Zap to log.
type zapLogger struct {
	// NB: this looks very similar to zap.SugaredLogger, but
	// deals with our desire to have multiple verbosity levels.
	zapLogger *zap.Logger
	infoLogger
}

var logger *zapLogger

var loggerConfig *zap.Config

func initLoggerConfig() {
	opts := DefaultOptions()
	opts.ApplyFlags(config.GetLog())

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// when output to local path, with color is forbidden
	if !opts.DisableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	loggerConfig = &zap.Config{
		Level:             zap.NewAtomicLevelAt(opts.Level),
		Development:       opts.Deployment,
		DisableCaller:     !opts.EnableCaller,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         string(opts.Format),
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
		InitialFields:    map[string]interface{}{},
	}
}

// Init initializes logger by opts which can custmoized by command arguments.
func Init() {
	initLoggerConfig()
	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger = &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
	//klog.InitLogger(l)
	zap.RedirectStdLog(l)
}

func SetLevel(level zapcore.Level) {
	loggerConfig.Level.SetLevel(level)
	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger = &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
}

//type Logger struct {
//
//}
//
//func (l *zapLogger) WithComponent(component string) (Logger,error) {
//	l.log.Named()
//}

// Debug method output debug level log.
func Debug(msg string, fields ...Field) {
	logger.zapLogger.Debug(msg, fields...)
}

// Debugf method output debug level log.
func Debugf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Debugf(format, v...)
}

// Info method output info level log.
func Info(msg string, fields ...Field) {
	logger.zapLogger.Info(msg, fields...)
}

// Infof method output info level log.
func Infof(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Infof(format, v...)
}

func Infow(msg string, keysAndVals ...interface{}) {
	logger.zapLogger.Sugar().Infow(msg, keysAndVals...)
}

// Warn method output warning level log.
func Warn(msg string, fields ...Field) {
	logger.zapLogger.Warn(msg, fields...)
}

// Warnf method output warning level log.
func Warnf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Warnf(format, v...)
}

// Error method output error level log.
func Error(msg string, fields ...Field) {
	logger.zapLogger.Error(msg, fields...)
}

// Errorf method output error level log.
func Errorf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Errorf(format, v...)
}

// Panic method output panic level log and shutdown application.
func Panic(msg string, fields ...Field) {
	logger.zapLogger.Panic(msg, fields...)
}

// Panicf method output panic level log and shutdown application.
func Panicf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Panicf(format, v...)
}

// Fatal method output fatal level log.
func Fatal(msg string, fields ...Field) {
	logger.zapLogger.Fatal(msg, fields...)
}

// Fatalf method output fatal level log.
func Fatalf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Fatalf(format, v...)
}

// Flush calls the underlying Core's Sync method, flushing any buffered
// log entries. Applications should take care to call Sync before exiting.
func Flush() { logger.Flush() }

func (l *zapLogger) Flush() {
	_ = l.zapLogger.Sync()
}
