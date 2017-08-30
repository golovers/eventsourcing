package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/Sirupsen/logrus"
)

// Config contains configurations to spins up the HTTP server.
type Config struct {
	// Address where HTTP server bind to to run
	ServerAddr string `envconfig:"SERVER_ADDRESS" default:"127.0.0.1:10000"`
}

// LoadEnvConfig returns a Config object populated from environment variables.
func LoadEnvConfig() *Config {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		logrus.Fatalf("config: Unable to load config for %T: %s", cfg, err)
	}
	return cfg
}
