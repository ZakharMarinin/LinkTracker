package main

import (
	"context"
	"linktracker/internal/application"
	"linktracker/internal/clients/scrapper"
	"linktracker/internal/config"
	"linktracker/internal/http-server/handlers"
	"linktracker/internal/http-server/router"
	"linktracker/internal/storage"
	tgHandlers "linktracker/internal/telegramBot/handlers"
	tgRouter "linktracker/internal/telegramBot/router"
	"linktracker/internal/usecase"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"gopkg.in/telebot.v4"
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	rPersistence := storage.NewRedisCom(rdb, log)

	httpRouter := chi.NewRouter()

	client := scrapper.NewScrapperClient(log, cfg)

	useCase := usecase.New(log, client, rPersistence)

	bot, err := botRun(cfg, log, useCase, httpRouter)
	if err != nil {
		log.Error("Cannot start bot", err)
		panic(err)
	}

	rout := handlers.NewURLUpdate(bot, log)
	router.Router(httpRouter, rout, ctx, log)

	tgRouter.Router(bot, ctx)

	app := application.NewApplication(ctx, cfg, log, httpRouter, bot)

	app.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	app.Shutdown()
}

func botRun(cfg *config.Config, log *slog.Logger, useCase *usecase.UseCase, httpRouter *chi.Mux) (*tgHandlers.BotHandler, error) {
	pref := telebot.Settings{
		Token:  cfg.TgBot.TgToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	newBot, err := telebot.NewBot(pref)
	if err != nil {
		log.Error("BotRun: failed creating a bot with error", err.Error())
		return nil, err
	}

	botHandlers := tgHandlers.NewBotHandler(newBot, useCase, log)

	return botHandlers, nil
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
