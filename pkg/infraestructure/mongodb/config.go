package mongodb

import (
	"github.com/caarlos0/env/v10"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// Config represents the configuration settings required for connecting to a MongoDB instance.
type Config struct {
	URI      string `env:"MONGO_URI"`
	User     string `env:"MONGO_USER"`
	Password string `env:"MONGO_PWD"`
}

// LoadConfig loads MongoDB configuration from environment variables and logs any errors encountered during parsing.
// It returns a Config object and an error if the configuration fails to load.
func LoadConfig(logger log.Logger) (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error("Failed to load MongoDB configuration", err)
		return Config{}, err
	}
	return cfg, nil
}
