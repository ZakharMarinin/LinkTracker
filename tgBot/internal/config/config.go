package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string      `yaml:"env" env-default:"local"`
	HttpServer HttpServer  `yaml:"http_server"`
	TgBot      TelegramBot `yaml:"tg_bot"`
	Redis      RedisConfig `yaml:"redis"`
	BotClients BotClients  `yaml:"bot_clients"`
}

type BotClients struct {
	Scrapper Scrapper `yaml:"scrapper"`
}

type Scrapper struct {
	Addr    string        `yaml:"addr" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
	Retry   int           `yaml:"retry" env-required:"true"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" env-required:"true"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
type TelegramBot struct {
	TgToken string `yaml:"tg_token"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG env variable not set")
	}

	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	cfg.TgBot.TgToken = os.Getenv("TG_TOKEN")

	return &cfg
}
