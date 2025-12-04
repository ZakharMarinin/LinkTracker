package application

import (
	"context"
	"linktracker/internal/config"
	tgHandlers "linktracker/internal/telegramBot/handlers"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	ctx    context.Context
	cfg    *config.Config
	log    *slog.Logger
	server *http.Server
	bot    *tgHandlers.BotHandler
}

func NewApplication(ctx context.Context, cfg *config.Config, log *slog.Logger, router *chi.Mux, bot *tgHandlers.BotHandler) *Application {
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
		server: srv,
		bot:    bot,
	}
}

func (a *Application) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *Application) Run() error {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.bot.Bot.Start()

		a.log.Info("Run: bot started")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.log.Info("Run: server started")

		err := a.server.ListenAndServe()
		if err != nil {
			a.log.Error("ListenAndServe: ", err)
		}
	}()

	go func() {
		wg.Wait()
		a.log.Info("Run: server stopped")
	}()

	return nil
}

func (a *Application) Shutdown() {
	a.log.Info("Shutdown")

	err := a.server.Shutdown(a.ctx)
	if err != nil {
		a.log.Error("Shutdown: failed to shutdown server", "error", err)
	}

	a.bot.Bot.Stop()
}
