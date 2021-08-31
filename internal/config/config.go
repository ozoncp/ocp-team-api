package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var cfg *Config
var fileCfg = "config.yml"

type Config struct {
	Project  *Project  `yaml:"project"`
	Database *Database `yaml:"database"`
	Server   *Server   `yaml:"server"`
	Status   *Status   `yaml:"status"`
	Jaeger   *Jaeger   `yaml:"jaeger"`
	Metrics  *Metrics  `yaml:"metrics"`
	Kafka    *Kafka    `yaml:"kafka"`
	Common   *Common   `yaml:"common"`
}

var cfgInitOnce sync.Once

func GetInstance() *Config {
	cfgInitOnce.Do(func() {
		cfg = readCfg()
	})

	return cfg
}

func readCfg() *Config {
	file, err := os.Open(fileCfg)
	if err != nil {
		log.Error().Err(err)
	}
	defer file.Close()

	var config *Config
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		log.Error().Err(err)
	}

	return config
}

type Project struct {
	Name string `yaml:"name"`
}

type Database struct {
	DSN string `yaml:"dsn"`
}

type Server struct {
	Host         string `yaml:"host"`
	HttpPort     string `yaml:"http_port"`
	GrpcPort     string `yaml:"grpc_port"`
	StartupTime  uint64 `yaml:"startup_time"`
	ShutdownTime uint64 `yaml:"shutdown_time"`
}

type Status struct {
	Port          string `yaml:"port"`
	HealthHandler string `yaml:"health_handler"`
	ReadyHandler  string `yaml:"ready_handler"`
}

type Jaeger struct {
	ServiceName string `yaml:"service_name"`
}

type Metrics struct {
	Port    string `yaml:"port"`
	Handler string `yaml:"handler"`
}

type Kafka struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}

type Common struct {
	BatchSize int `yaml:"batch_size"`
}
