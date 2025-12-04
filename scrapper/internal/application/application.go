package application

import (
	"context"
	"log/slog"
	"net/http"
	"scrapper/internal/config"
	"scrapper/internal/cron"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	ctx    context.Context
	cfg    *config.Config
	log    *slog.Logger
	server *http.Server
	cron   *cronModule.Cron
}

func NewApplication(ctx context.Context, cfg *config.Config, log *slog.Logger, router *chi.Mux, cron *cronModule.Cron) *Application {
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
		cron:   cron,
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
		a.log.Info("Run: server started")

		err := a.server.ListenAndServe()
		if err != nil {
			a.log.Error("ListenAndServe: ", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.cron.Cron.StartBlocking()
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
}
