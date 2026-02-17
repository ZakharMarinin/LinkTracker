package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"scrapper/internal/application"
	"scrapper/internal/config"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()
	ctx, cancel := context.WithCancel(context.Background())

	log := setupLogger(cfg.Env)

	app := application.NewApplication(ctx, cfg, log)

	app.MustRun()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown

	cancel()
	app.Shutdown()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
