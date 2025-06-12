// Package config provides configuration management functionality for the authentication service.
package config

import (
	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

// MongoConfig represents the configuration settings required for connecting to a MongoDB instance.
type MongoConfig struct {
	URI      string `env:"MONGO_URI"`
	User     string `env:"MONGO_USER"`
	Password string `env:"MONGO_PWD"`
}

// LoadMongoConfig loads MongoDB configuration from environment variables and logs any errors encountered during parsing.
// It returns a MongoConfig object and an error if the configuration fails to load.
func LoadMongoConfig(logger *zap.Logger) (MongoConfig, error) {
	cfg := MongoConfig{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error("Failed to load MongoDB configuration", zap.Error(err))
		return MongoConfig{}, err
	}
	return cfg, nil
}
