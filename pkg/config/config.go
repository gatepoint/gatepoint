package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var Shared Config
var DefaultShared Config

type Config struct {
	Debug  bool    `json:"debug"`
	Server *Server `json:"server"`
	Log    *Log    `json:"log"`
	Kube   *Kube   `json:"kube"`

	Swagger bool `json:"swagger"`
	//Kubernetes
}

type Kube struct {
	Incluster  bool   `json:"incluster"`
	Kubeconfig string `json:"kubeconfig"`
	Namespace  string `json:"namespace"`
}

type Server struct {
	HTTP  Conn `json:"http"`
	Grpc  Conn `json:"grpc"`
	Admin Conn `json:"admin"`
}

type Conn struct {
	Address string `json:"addr"`
	Timeout string `json:"timeout"`
}

type Log struct {
	Level            string   `json:"level"`
	Format           string   `json:"format"`
	DisableColor     bool     `json:"disable_color"`
	EnableCaller     bool     `json:"enable_caller"`
	OutputPaths      []string `json:"output_paths"`
	ErrorOutputPaths []string `json:"error_output_paths"`
	Deployment       bool     `json:"deployment"`
}

func GetNamespace() string {
	return config.Kube.Namespace
}

func EnableSwagger() bool {
	return config.Swagger
}

func EnableDebug() bool {
	return config.Debug
}

func GetGrpcAddr() string {
	return config.Server.Grpc.Address
}

func GetHttpAddr() string {
	return config.Server.HTTP.Address

}

func GetAdminAddr() string {
	return config.Server.Admin.Address
}

func GetKubeConfig() string {
	return config.Kube.Kubeconfig
}

func GetLog() *Log {
	return config.Log
}

var config *Config

func LoadConfig(cfgFile string) error {
	cfg := &Config{}
	if cfgFile != "" {
		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfgFile)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".gatepoint")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Failed to read in the config file %s: %v", viper.ConfigFileUsed(), err))
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	if err := viper.Unmarshal(cfg, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "json"
	}); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal configuration from disk: %v", err))
	}
	cfg.setDefaultConfig()
	config = cfg
	return nil
}

func (c *Config) setDefaultConfig() {
	if c.Server == nil {
		c.Server = defaultServer()
	}
	if c.Kube != nil {
		c.Kube = defaultKube()
	}
}

func defaultKube() *Kube {
	return &Kube{
		Incluster: true,
		Namespace: "gatepoint",
	}
}

func defaultServer() *Server {
	return &Server{
		HTTP: Conn{
			Address: ":8080",
			Timeout: "10s",
		},
		Grpc: Conn{
			Address: ":8081",
			Timeout: "10s",
		},
		Admin: Conn{
			Address: ":8082",
			Timeout: "10s",
		},
	}
}
