package main

import (
	"context"
	"linktracker/internal/application"
	"linktracker/internal/config"
	"log/slog"
	"os"
	"os/signal"
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

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
