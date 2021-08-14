package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "./configs/config.yaml", "path to config file")
}

func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	config, err := NewConfig(configPath)
	if err != nil {
		return fmt.Errorf("error during config reading: %w", err)
	}
	if err := InitLogger(config.Logger); err != nil {
		return fmt.Errorf("error during logger init: %w", err)
	}

	zap.L().Info("Banner rotation service successfully started")
	return nil
}
