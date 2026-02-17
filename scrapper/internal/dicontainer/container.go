package dicontainer

import (
	"context"
	"fmt"
	"log/slog"
	"scrapper/internal/client/githubclient"
	"scrapper/internal/client/sender"
	"scrapper/internal/config"
	cronModule "scrapper/internal/cron"
	"scrapper/internal/http/handlers"
	"scrapper/internal/http/router"
	"scrapper/internal/metrics"
	"scrapper/internal/storage"
	"scrapper/internal/usecase"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron"
)

type Container struct {
	Cfg          *config.Config
	Log          *slog.Logger
	DB           *storage.PostgresStorage
	HTTPRouter   *chi.Mux
	UseCase      *usecase.UseCase
	GitClient    *githubclient.GithubClient
	Sender       *sender.TypeOfSender
	Cron         *cronModule.Cron
	HTTPHandlers *handlers.HTTP
	Metrics      *metrics.Metrics
}

func NewContainer(log *slog.Logger, cfg *config.Config) *Container {
	return &Container{
		Cfg: cfg,
		Log: log,
	}
}

func (c *Container) Init(ctx context.Context) error {
	var err error

	c.DB, err = storage.New(ctx, c.Cfg.Postgres.Addr)
	if err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}

	c.HTTPRouter = chi.NewRouter()

	c.GitClient = githubclient.NewGithubClient(c.Cfg.GitHubToken, c.Log)

	c.Sender, err = sender.New(c.Log, c.Cfg)
	if err != nil {
		return fmt.Errorf("failed to init Sender: %w", err)
	}

	c.UseCase = usecase.NewUseCase(c.DB, c.Log, c.Cfg)

	c.HTTPHandlers = handlers.NewHTTP(c.UseCase, c.Log)

	router.Router(ctx, c.HTTPRouter, c.HTTPHandlers, c.Log)

	startCron := gocron.NewScheduler(time.UTC)
	c.Cron = cronModule.New(c.Log, startCron, c.DB, c.GitClient, c.Sender, 1000)

	_, err = c.Cron.Cron.Every(1).Minutes().Do(c.Cron.StartCron)
	if err != nil {
		return err
	}

	c.Metrics = metrics.NewMetrics()

	return nil
}
