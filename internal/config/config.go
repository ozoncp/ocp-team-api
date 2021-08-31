package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var cfg *Config
var fileCfg = "config.yml"

// Config is the struct that represents application configuration.
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

// GetInstance is the method for getting singleton instance
// of the project configuration.
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

// Project is the struct representing project description in configuration.
type Project struct {
	Name string `yaml:"name"`
}

// Database is the struct representing database settings in configuration.
type Database struct {
	DSN string `yaml:"dsn"`
}

// Server is the struct representing main server settings in configuration.
type Server struct {
	Host         string `yaml:"host"`
	HttpPort     string `yaml:"http_port"`
	GrpcPort     string `yaml:"grpc_port"`
	StartupTime  uint64 `yaml:"startup_time"`
	ShutdownTime uint64 `yaml:"shutdown_time"`
}

// Status is the struct representing status server settings in configuration.
type Status struct {
	Port          string `yaml:"port"`
	HealthHandler string `yaml:"health_handler"`
	ReadyHandler  string `yaml:"ready_handler"`
}

// Jaeger is the struct representing jaeger settings in configuration.
type Jaeger struct {
	ServiceName string `yaml:"service_name"`
}

// Metrics is the struct representing metrics settings in configuration.
type Metrics struct {
	Port    string `yaml:"port"`
	Handler string `yaml:"handler"`
}

// Kafka is the struct representing kafka settings in configuration.
type Kafka struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}

// Common is the struct representing common settings in configuration.
type Common struct {
	BatchSize int `yaml:"batch_size"`
}
