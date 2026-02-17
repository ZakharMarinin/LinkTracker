package application

import (
	"context"
	"encoding/json"
	"errors"
	"linktracker/internal/config"
	"linktracker/internal/dicontainer"
	"linktracker/internal/domain"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Application struct {
	ctx       context.Context
	cfg       *config.Config
	log       *slog.Logger
	container *dicontainer.Container
	wg        *sync.WaitGroup
	server    *http.Server
}

func NewApplication(ctx context.Context, cfg *config.Config, log *slog.Logger) *Application {
	return &Application{
		ctx:       ctx,
		cfg:       cfg,
		log:       log,
		container: dicontainer.NewContainer(log, cfg),
		wg:        &sync.WaitGroup{},
	}
}

func (a *Application) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *Application) Run() error {
	err := a.container.Init(a.ctx)
	if err != nil {
		a.log.Error("Failed to init container", "error", err)
		return err
	}

	a.server = &http.Server{
		Addr:         a.cfg.HTTPServer.Address,
		Handler:      a.container.HTTPRouter,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
		IdleTimeout:  a.cfg.HTTPServer.IdleTimeout,
	}

	a.wg.Add(1)

	go func() {
		defer a.wg.Done()
		a.container.Bot.Bot.Start()

		a.log.Info("Run: bot started")
	}()

	a.wg.Add(1)

	go func() {
		defer a.wg.Done()
		a.log.Info("Run: server started")

		err := a.server.ListenAndServe()
		if err != nil {
			a.log.Error("ListenAndServe", "error", err)
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

		msg, err := a.container.Kafka.ReadMessage(a.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				a.log.Info("Run: kafka consumer stopped")
				return
			}

			a.log.Error("Run: kafka consumer read error", "error", err)

			time.Sleep(2 * time.Second)

			continue
		}

		var update domain.UpdatedLink

		err = json.Unmarshal(msg.Value, &update)
		if err != nil {
			_ = a.container.Kafka.CommitMessage(a.ctx, msg)
			continue
		}

		err = a.container.Bot.Updates(&update)
		if err != nil {
			a.log.Error("failed to send updates", "error", err)
			time.Sleep(2 * time.Second)

			continue
		}

		_ = a.container.Kafka.CommitMessage(a.ctx, msg)
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

	a.container.Bot.Bot.Stop()

	a.wg.Wait()
	a.log.Info("Shutdown: completed graceful shutdown")
}
