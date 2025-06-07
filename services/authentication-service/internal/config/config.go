package config

import (
	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

type MongoConfig struct {
	URI      string `env:"MONGO_URI"`
	User     string `env:"MONGO_USER"`
	Password string `env:"MONGO_PWD"`
}

func LoadMongoConfig(logger *zap.Logger) (MongoConfig, error) {
	cfg := MongoConfig{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error("Failed to load MongoDB configuration", zap.Error(err))
		return MongoConfig{}, err
	}
	return cfg, nil
}
