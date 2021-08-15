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

	"github.com/Raschudesny/otus_project/v1/internal"
	"github.com/Raschudesny/otus_project/v1/server"
	"github.com/Raschudesny/otus_project/v1/storage/sql"
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
	config, err := NewConfig(configPath)
	if err != nil {
		return fmt.Errorf("error during config reading: %w", err)
	}
	if err := InitLogger(config.Logger); err != nil {
		return fmt.Errorf("error during logger init: %w", err)
	}
	zap.L().Info("Banner rotation service starting...")
	ctx, stop := signal.NotifyContext(context.Background(), TerminalSignals...)
	defer stop()

	zap.L().Info("rotation service storage starting...")
	dbStorage := sql.NewStorage("pgx", "host=localhost dbname=rotation user=danny password=danny sslmode=disable")
	if err := dbStorage.Connect(ctx); err != nil {
		return fmt.Errorf("failed to init db storage: %w", err)
	}
	defer func() {
		if err := dbStorage.Close(); err != nil {
			zap.L().Error("failed to close db storage", zap.Error(err))
		}
	}()
	zap.L().Info("rotation service storage started")

	app := internal.NewRotationService(dbStorage)
	grpcServer := server.InitServer(app)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer.Start("localhost", "56789", stop)
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
