package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Raschudesny/otus_project/v1/internal/config"
	"github.com/Raschudesny/otus_project/v1/internal/logger"
	"github.com/Raschudesny/otus_project/v1/internal/server"
	"github.com/Raschudesny/otus_project/v1/internal/services"
	"github.com/Raschudesny/otus_project/v1/internal/storage/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

var TerminalSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGHUP}

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
	cnf, err := config.NewConfig(configPath)
	if err != nil {
		return fmt.Errorf("error during config reading: %w", err)
	}
	if err := logger.InitLogger(cnf.Logger); err != nil {
		return fmt.Errorf("error during logger init: %w", err)
	}
	zap.L().Info("Banner rotation service starting...")
	ctx, stop := signal.NotifyContext(context.Background(), TerminalSignals...)
	defer stop()

	zap.L().Info("rotation service storage starting...")
	dbStorage := sql.NewStorage("pgx", cnf.DB)
	if err := dbStorage.Connect(ctx); err != nil {
		return fmt.Errorf("failed to init db storage: %w", err)
	}
	defer func() {
		if err := dbStorage.Close(); err != nil {
			zap.L().Error("failed to close db storage", zap.Error(err))
		}
	}()
	zap.L().Info("rotation service storage started")

	app := services.NewRotationService(dbStorage)
	grpcServer := server.InitServer(app, cnf.Server)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer.Start(stop)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		grpcServer.Stop()
	}()
	wg.Wait()
	zap.L().Info("rotation service stopped")
	return nil
}
