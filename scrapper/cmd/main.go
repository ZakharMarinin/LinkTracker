package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"scrapper/internal/application"
	"scrapper/internal/client/githubClient"
	"scrapper/internal/client/tgBotClient"
	"scrapper/internal/config"
	cronModule "scrapper/internal/cron"
	"scrapper/internal/http/handlers"
	"scrapper/internal/http/router"
	"scrapper/internal/storage"
	"scrapper/internal/usecase"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron"
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
	defer cancel()

	log := setupLogger(cfg.Env)

	db, err := storage.New(ctx, cfg.Postgres.Addr)
	if err != nil {
		log.Error("failed to connect to storage", "error", err)
		return
	}
	defer db.Close()

	httpRouter := chi.NewRouter()

	useCase := usecase.NewUseCase(db, log, ctx, cfg)

	gitClient := githubClient.NewGithubClient(cfg.GitHubToken, log)

	tgClient := tgBotClient.NewTGClient(log, cfg)

	cron, err := setupCron(log, db, gitClient, tgClient, 100)
	if err != nil {
		log.Error("failed to setup cron", "error", err)
		return
	}

	httpHandlers := handlers.NewHTTP(useCase, log)

	router.Router(ctx, httpRouter, httpHandlers, log)

	app := application.NewApplication(ctx, cfg, log, httpRouter, cron)

	app.MustRun()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown

	app.Shutdown()
}

func setupCron(log *slog.Logger, db *storage.PostgresStorage, gitClient *githubClient.GithubClient, tgClient *tgBotClient.Client, limit uint64) (*cronModule.Cron, error) {
	startCron := gocron.NewScheduler(time.UTC)
	cron := cronModule.New(log, startCron, db, gitClient, tgClient, limit)

	_, err := cron.Cron.Every(1).Minutes().Do(cron.StartCron)
	if err != nil {
		return nil, err
	}

	return cron, nil
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
