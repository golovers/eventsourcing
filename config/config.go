package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/Sirupsen/logrus"
	"github.com/lnquy/eventsourcing/datastore"
)

type Config struct {
	ServerAddr string `envconfig:"SERVER_ADDRESS" default:"127.0.0.1:10000"`
	DataStore *datastore.Config
}

// LoadEnvConfig returns a Config object populated from environment variables.
func LoadEnvConfig() *Config {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		logrus.Fatalf("config: Unable to load config for %T: %s", cfg, err)
	}
	return cfg
}
