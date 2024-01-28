package main

import (
	"log/slog"
	"os"

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
