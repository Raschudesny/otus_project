package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConfig
}

type LoggerConfig struct {
	Level string
	File  string
}

func NewConfig(path string) (cfg Config, err error) {
	if cfg.Logger.File != "" {
		return Config{}, errors.New("what is going here")
	}
	viper.AutomaticEnv()
	return Config{LoggerConfig{
		Level: "info",
		File:  path,
	}}, nil
}
