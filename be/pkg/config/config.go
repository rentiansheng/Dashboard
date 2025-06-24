package config

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ServerConfig holds the server configuration
type configStruct struct {
	Server   Server `mapstructure:"server"`
	LogLevel string `mapstructure:"log_level"`
	Engine   Engine `mapstructure:"engine"`
	ES       ES     `mapstructure:"es"`
}

type Server struct {
	Port            string        `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Engine struct {
	EngineType          string       `mapstructure:"engine_type"`
	IsAll               bool         `mapstructure:"is_all"`
	DefaultEngineConfig EngineConfig `mapstructure:"default_engine_config"`
	DataSourceEngine    EngineConfig `mapstructure:"data_source_engine_config"`
	GroupKeyEngine      EngineConfig `mapstructure:"group_key_engine_config"`
}

type EngineConfig struct {
	Mysql string `mapstructure:"dsn"`
}

type ES struct {
	Addrs    []string `mapstructure:"addrs"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

var (
	Config *configStruct
)

// NewServerConfig creates a new server configuration with defaults
func NewServerConfig() error {
	// Set default values
	viper.SetDefault("port", "8080")
	viper.SetDefault("read_timeout", "5s")
	viper.SetDefault("write_timeout", "10s")
	viper.SetDefault("shutdown_timeout", "15s")
	viper.SetDefault("log_level", "info")

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DASHBOARD")
	viper.BindEnv("port", "PORT")

	// Read config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/dashboard/")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		log.Warn("No config file found, using defaults and environment variables")
	}

	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// ConfigureLogging sets up the logging configuration based on the server config
func (c *configStruct) ConfigureLogging() error {
	logLevel, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	log.SetLevel(logLevel)
	return nil
}
