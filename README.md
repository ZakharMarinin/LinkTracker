# LinkTracker
<img src="logowidth.png" alt="logo" style="width: 100%">
<br>
The LinkTracker bot is a service for tracking GitHub repositories. 
By pinning a repository link in the Telegram bot, 
you will receive notifications about updates to that repository.

## Table of Content

- [Usage](#usage)
- [Features](#features)
- [TGBot](#telegrambot)
- [Scrapper](#scrapper)
- [Technologies](#technologies)


## Usage

### Application launch

Create a `.env` file in the root directory with the following variables:

```env
# Telegram Bot Configuration
TG_TOKEN=your_telegram_bot_token_here
CONFIG_PATH="./your_config_file"

# GitHub Configuration
GITHUB_TOKEN=your_github_personal_access_token_here
CONFIG_PATH="./your_config_file"
GOOSE_DBSTRING= postgres addr
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=./migrations
```

### DataBase launch

To launch the database, you need to run the ```docker-compose up -d``` command.



## Features

- Telegram bot interface with interactive buttons
- Track multiple GitHub repositories (Pull Requests & Issues)
- Automatic periodic checks for repository updates
- Real-time notifications about changes
- Fully containerized with Docker Compose
- Structured JSON logging
- Graceful shutdown support
- Secure configuration with environment variables


## TelegramBot

Telegram Bot service handles user interactions and manages communication between users and the Scrapper service.

### Commands

| Command    | Description                                    |
|------------|------------------------------------------------|
| `/start`   | Register user and start using the bot          |
| `/help`    | Display all available commands and usage guide |
| `/track`   | Add a GitHub repository link to track          |
| `/untrack` | Remove a repository from tracking list         |
| `/list`    | Show all currently tracked repositories        |
| `/cancel`  | Interrupt current operation                    |

### Configuration

**config.yaml** - Application settings:
```yaml
env: "local"
http_server:
  address: 0.0.0.0:8080
  timeout: 4s
  idle_timeout: 60s
bot_clients:
  scrapper:
    addr: "http://scrapper:8081"
    timeout: 5s
    retry: 5
redis:
  addr: "valkey:6379"
  password: ""
  DB: 0
```

## **.env:**
- `TG_TOKEN` - Telegram Bot API token
- `CONFIG_PATH` - Config file location

### API Endpoints

- `POST /updates` - Receive notifications about repository updates from Scrapper


### Database Schema

The service uses Redis caching temporary user information (list of links, user state):

- `UserTempState` - Stores data about the user state with 2Hour TTL
- `UserTempLinks` - Stores all user links in the cache for faster interaction with them; if the cache is empty, a query is sent to PostgreSQL


### Features

- Interactive Telegram buttons for quick command input
- Graceful HTTP server shutdown
- JSON structured logging with different levels (info/warn/error)
- ScrapperClient for communication with Scrapper service



## Scrapper

Scrapper service manages repository tracking, performs periodic checks, and notifies users about changes.

### API Endpoints

| Method | Endpoint        | Description |
|--------|-----------------|-------------|
| `POST` | `/tg-chat/{id}` | Register new Telegram chat (called on `/start`) |
| `DELETE` | `/tg-chat/{id}` | Delete chat and all associated links |
| `GET` | `/links/{id}`   | Get all tracked links for authenticated user |
| `POST` | `/links`        | Add new repository link to track |
| `DELETE` | `/links`        | Remove repository link from tracking |

**Note:** User ID can be passed via `/{id}` path parameter or through request headers.

### Update Check

The service uses **gocron** to periodically check for updates in tracked repositories:

- Runs at configurable intervals (default: every 5 minutes)
- Scrapes multiple data sources for each repository:
    - Pull Requests (new PRs, status changes, merges)
    - Issues (new issues, comments, status changes)
- Sends notifications to Bot service via `/updates` endpoint

### Configuration

**config.yaml** - Application settings:
```yaml
env: "local"
http_server:
  address: 0.0.0.0:8081
  timeout: 4s
  idle_timeout: 60s
tgbot:
  addr: "http://tgBot:8080"
  timeout: 5s
```

## **.env:**
- `GITHUB_TOKEN` - GitHub Personal Access Token
- `CONFIG_PATH` - Config file location
- `GOOSE_DBSTRING` - Postgresql address
- `GOOSE-DRIVER` - postgres
- `GOOSE_MIGRATION_DIR` - ./migrations

### Database Schema

The service uses PostgreSQL with the following main tables:

- `chats` - Telegram chat information
- `links` - Tracked repository links
- `chat_links` - Many-to-many relationship between chats and links
- `updates` - History of detected updates

### Features

- RESTful API built with **chi** router
- Database queries built with **Squirrel** SQL builder
- Multiple GitHub data sources tracking (PRs + Issues)
- Graceful HTTP server shutdown
- JSON structured logging (slog)
- GitHubClient for GitHub API interactions
- TgBotClient for sending notifications


 ## Technologies

### Telegram Bot Service
- **Language:** Go
- **Telegram API:** telebot
- **HTTP Router:** chi
- **Database:** Redis
- **Logger:** slog
- **Configuration:** yaml + .env files
- **Container:** Docker

### Scrapper Service
- **Language:** Go
- **HTTP Router:** chi
- **Database:** PostgreSQL
- **Query Builder:** Squirrel
- **Update Check:** gocron
- **Logger:** slog
- **Configuration:** yaml + .env files
- **Container:** Docker
