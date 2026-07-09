package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project/database/migrations"
	"project/database/seeders"
	"project/internal/app"
	"project/internal/config"
	"project/internal/routes"

	"github.com/dracory/websrv"
)

func main() {
	cfg, err := config.NewFromEnv()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	fmt.Printf("Starting %s v%s\n", cfg.GetAppName(), config.GetVersion())

	a, err := app.New(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		return
	}
	defer func() {
		if err := a.Close(); err != nil {
			slog.Error("Failed to close app", "error", err)
		}
	}()

	if err := migrations.MigrateAll(a); err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		return
	}

	if err := seeders.SeedAll(a); err != nil {
		fmt.Printf("Failed to run seeders: %v\n", err)
		return
	}

	server, err := websrv.Start(websrv.Options{
		Host:    cfg.GetAppHost(),
		Port:    cfg.GetAppPort(),
		URL:     cfg.GetAppUrl(),
		Handler: routes.Router(a).ServeHTTP,
	})

	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	fmt.Println("Shutdown signal received")

	if server != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown failed", "error", err)
		}
	}
}
