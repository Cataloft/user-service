package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cataloft/user-service/internal/config"
	"github.com/Cataloft/user-service/internal/server"
	"github.com/Cataloft/user-service/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	db := storage.New(cfg.Database, log)
	srv := server.New(db, cfg, log)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Error("Failed to start server", "error", err)
		}
	}()
	<-done
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
