package dicontainer

import (
	"context"
	"linktracker/internal/clients/kafka"
	"linktracker/internal/clients/scrapper"
	"linktracker/internal/config"
	"linktracker/internal/http-server/handlers"
	"linktracker/internal/http-server/router"
	"linktracker/internal/metrics"
	"linktracker/internal/storage"
	tgHandlers "linktracker/internal/telegramBot/handlers"
	tgrouter "linktracker/internal/telegramBot/router"
	"linktracker/internal/usecase"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"gopkg.in/telebot.v4"
)

type Container struct {
	Cfg          *config.Config
	Log          *slog.Logger
	DB           *storage.RedisCom
	HTTPRouter   *chi.Mux
	UseCase      *usecase.UseCase
	HTTPHandlers *handlers.HTTP
	Bot          *tgHandlers.BotHandler
	Scrapper     *scrapper.Client
	Kafka        *kafka.Consumer
	Metric       *metrics.Metric
}

func NewContainer(log *slog.Logger, cfg *config.Config) *Container {
	return &Container{
		Cfg: cfg,
		Log: log,
	}
}

func (c *Container) Init(ctx context.Context) error {
	var err error

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Cfg.Redis.Addr,
		Password: c.Cfg.Redis.Password,
		DB:       c.Cfg.Redis.DB,
	})

	c.DB = storage.NewRedisCom(rdb, c.Log)

	c.HTTPRouter = chi.NewRouter()

	c.Scrapper = scrapper.NewScrapperClient(c.Log, c.Cfg)

	c.UseCase = usecase.New(c.Log, c.Scrapper, c.DB)

	c.Kafka = kafka.NewConsumer(c.Cfg)

	pref := telebot.Settings{
		Token:  c.Cfg.TgBot.TgToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	newBot, err := telebot.NewBot(pref)
	if err != nil {
		c.Log.Error("BotRun: failed creating a bot with error", "error", err)
		return err
	}

	c.Metric = metrics.NewMetric()

	c.Bot = tgHandlers.NewBotHandler(newBot, c.UseCase, c.Log, c.Metric)

	c.HTTPHandlers = handlers.NewURLUpdate(c.Bot, c.Log)

	router.Router(c.HTTPRouter, c.HTTPHandlers, c.Log)

	tgrouter.Router(c.Bot, ctx)

	return nil
}
