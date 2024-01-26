package main

import (
	"log/slog"
	"os"
	"test-task/internal/config"
	"test-task/internal/server"
	"test-task/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	db := storage.New(cfg.DatabaseUrl)
	srv := server.New(db, cfg, log)

	err := srv.Start()

	if err != nil {
		log.Error("Server crashed", "error", err)
		return
	}
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
