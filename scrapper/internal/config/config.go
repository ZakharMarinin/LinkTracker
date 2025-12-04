package config

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string         `yaml:"env" env-default:"local"`
	HttpServer  HttpServer     `yaml:"http_server"`
	Postgres    PostgresConfig `yaml:"postgres"`
	GitHubToken string         `yaml:"github_token"`
	TgBot       TGBot          `yaml:"tgbot"`
}

type TGBot struct {
	Addr    string        `yaml:"addr" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type PostgresConfig struct {
	Addr string `yaml:"addr"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8081"`
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
		log.Fatal("cannot find the config: ", err)
	}

	cfg.Postgres.Addr = os.Getenv("GOOSE_DBSTRING")
	cfg.GitHubToken = os.Getenv("GITHUB_TOKEN")

	return &cfg
}

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + "scrapper" + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/scrapper` + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
