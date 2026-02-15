package application

import (
	"context"
	"encoding/json"
	"errors"
	"linktracker/internal/clients/kafka"
	"linktracker/internal/config"
	"linktracker/internal/domain"
	tgHandlers "linktracker/internal/telegramBot/handlers"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	ctx    context.Context
	cfg    *config.Config
	log    *slog.Logger
	wg     *sync.WaitGroup
	server *http.Server
	bot    *tgHandlers.BotHandler
	kafka  *kafka.Consumer
}

func NewApplication(ctx context.Context, cfg *config.Config, log *slog.Logger, router *chi.Mux, bot *tgHandlers.BotHandler, consumer *kafka.Consumer) *Application {
	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	return &Application{
		ctx:    ctx,
		cfg:    cfg,
		log:    log,
		wg:     &sync.WaitGroup{},
		server: srv,
		bot:    bot,
		kafka:  consumer,
	}
}

func (a *Application) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *Application) Run() error {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.bot.Bot.Start()

		a.log.Info("Run: bot started")
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.log.Info("Run: server started")

		err := a.server.ListenAndServe()
		if err != nil {
			a.log.Error("ListenAndServe: ", err)
		}
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.StartKafkaConsumer()
	}()

	return nil
}

func (a *Application) StartKafkaConsumer() {
	a.log.Info("Run: kafka consumer started")

	for {
		select {
		case <-a.ctx.Done():
			a.log.Info("Run: kafka consumer stopped")
			return
		default:
		}

		msg, err := a.kafka.ReadMessage(a.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				a.log.Info("Run: kafka consumer stopped")
				return
			}
			a.log.Error("Run: kafka consumer read error: ", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var update domain.UpdatedLink
		err = json.Unmarshal(msg.Value, &update)
		if err != nil {
			_ = a.kafka.CommitMessage(a.ctx, msg)
			continue
		}

		err = a.bot.Updates(&update)
		if err != nil {
			a.log.Error("failed to send updates", "error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		_ = a.kafka.CommitMessage(a.ctx, msg)
	}
}

func (a *Application) Shutdown() {
	a.log.Info("Shutdown")

	srvCtx, srvCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer srvCancel()

	err := a.server.Shutdown(srvCtx)
	if err != nil {
		a.log.Error("Shutdown: failed to shutdown server", "error", err)
	}

	a.bot.Bot.Stop()

	a.wg.Wait()
	a.log.Info("Shutdown: completed graceful shutdown")
}
