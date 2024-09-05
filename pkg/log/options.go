package log

import (
	"fmt"
	"strings"

	"github.com/gatepoint/gatepoint/pkg/config"
	"go.uber.org/zap/zapcore"
)

type Format string

const (
	consoleFormat Format = "console"
	jsonFormat    Format = "json"
)

// ParseFormat takes a string format and returns the log format constant.
func ParseFormat(f string) (Format, error) {
	switch strings.ToLower(f) {
	case "console":
		return consoleFormat, nil
	case "json":
		return jsonFormat, nil
	default:
		return "", fmt.Errorf("not a valid log format: %q", f)
	}
}

type Options struct {
	Level            zapcore.Level `json:"level"`
	Format           Format        `json:"format"`
	Deployment       bool          `json:"deployment"`
	DisableColor     bool          `json:"disable_color"`
	EnableCaller     bool          `json:"enable_caller"`
	OutputPaths      []string      `json:"output_paths"`
	ErrorOutputPaths []string      `json:"error_output_paths"`
}

func DefaultOptions() *Options {
	return &Options{
		Level:            zapcore.InfoLevel,
		Format:           consoleFormat,
		DisableColor:     false,
		EnableCaller:     false,
		Deployment:       false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *Options) ApplyFlags(loggerCfg config.Log) error {

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(loggerCfg.Level)); err != nil {
		return err
	}
	o.Level = zapLevel

	format, err := ParseFormat(loggerCfg.Format)
	if err != nil {
		return err
	}
	o.Format = format

	o.DisableColor = loggerCfg.DisableColor
	o.EnableCaller = loggerCfg.EnableCaller
	if loggerCfg.OutputPaths != nil {
		o.OutputPaths = loggerCfg.OutputPaths
	}
	if loggerCfg.ErrorOutputPaths != nil {
		o.ErrorOutputPaths = loggerCfg.ErrorOutputPaths
	}
	o.Deployment = loggerCfg.Deployment
	return nil
}
