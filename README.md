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
- [Testing](#testing)


## Usage

### Application launch

Create a `.env` file in the root directory with the following variables:

```env
# Telegram Bot Configuration
TG_TOKEN=your_telegram_bot_token_here
CONFIG_PATH="./your_config_file"

# GitHub Configuration
GITHUB_TOKEN=your_github_personal_access_token_here

# Database Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=linktracker

# Migration Configuration (Goose)
GOOSE_DRIVER=postgres
GOOSE_DBSTRING="host=postgres user=postgres password=postgres dbname=linktracker sslmode=disable"
GOOSE_MIGRATION_DIR=./migrations

```

### Deployment

The project supports two launch modes via **Makefile** to test application behavior with and without migrations:

**1. Launch WITH Migrations (Standard)**
Starts the database, the applications, and a separate `migrator` service that automatically applies SQL schemas.

```bash
make run-with-migrations
```

**2. Launch WITHOUT Migrations**
Starts only the database and the services. Use this to verify how the application handles an empty database or missing tables.

```bash
make run-no-migrations
```

**3. Stop**

```bash
make stop   # Stop all containers\
```

## Features

* Telegram bot interface with interactive buttons
* Track multiple GitHub repositories (Pull Requests & Issues)
* Automatic periodic checks for repository updates
* Real-time notifications about changes
* **Separate Migration Service** using Goose (Containerized)
* Fully containerized with Docker Compose
* Structured JSON logging (slog)
* Graceful shutdown support
* Secure configuration with environment variables
* Integration testing with **Testcontainers**

## TelegramBot

Telegram Bot service handles user interactions and manages communication between users and the Scrapper service.

### Commands

| Command | Description |
| --- | --- |
| `/start` | Register user and start using the bot |
| `/help` | Display all available commands and usage guide |
| `/track` | Add a GitHub repository link to track |
| `/untrack` | Remove a repository from tracking list |
| `/list` | Show all currently tracked repositories |
| `/filter` | Show tracked repositories with tag filter |
| `/cancel` | Interrupt current operation |

### API Endpoints

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/updates` | Receive notifications about repository updates |

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
  kafka:
    addr: "kafka:9092"
    topic: link-update
    group_id: updates
    timeout: 5s
    retry: 5
redis:
  addr: "valkey:6379"
  password: ""
  DB: 0

```

### Database Schema (Redis)

The service uses Redis for caching temporary user information and state:

* `UserTempState` - Stores user FSM state with 2-hour TTL
* `UserTempLinks` - Caches user links for faster interaction; falls back to Scrapper API if empty.

## Scrapper

Scrapper service manages repository tracking, performs periodic checks, and notifies users about changes.

### API Endpoints

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/tg-chat/{id}` | Register new Telegram chat |
| `DELETE` | `/tg-chat/{id}` | Delete chat and all associated links |
| `GET` | `/links/{id}` | Get all tracked links for user |
| `GET` | `/links/{id}/{tag}` | Get all tracked links with tag filter |
| `POST` | `/links` | Add new repository link to track |
| `DELETE` | `/links` | Remove repository link from tracking |

### Update Check

The service uses **gocron** to periodically check for updates:

* Runs at configurable intervals (default: 5 minutes)
* Scrapes GitHub data:
* Pull Requests (new, status changes, merges)
* Issues (new, comments, status changes)


* Notifies the Bot service via `/updates` endpoint

### Configuration

**config.yaml** - Application settings:

```yaml
env: "local"
transport_type: "kafka"
http_server:
  address: 0.0.0.0:8081
  timeout: 4s
  idle_timeout: 60s
tgbot:
  addr: "http://tgBot:8080"
  timeout: 5s
kafka:
  addr: "kafka:9092"
  topic: link-update
  timeout: 5s
  retry: 5
```

### Database Schema (PostgreSQL)

The service uses PostgreSQL with the following main tables, managed by a separate **Goose** migrator container:

* `users` - Telegram chat information
* `links` - Global tracked repository links
* `user_links` - Many-to-many relationship between users and links

## Technologies

### Telegram Bot Service

* **Language:** Go
* **Telegram API:** telebot
* **HTTP Router:** chi
* **Database:** Redis
* **Logger:** slog
* **Container:** Docker

### Scrapper Service

* **Language:** Go
* **HTTP Router:** chi
* **Database:** PostgreSQL
* **Query Builder:** Squirrel
* **Migrations:** Goose (Independent container)
* **Update Check:** gocron
* **Logger:** slog
* **Container:** Docker

## Testing

The project implements a full testing lifecycle:

* **Unit Tests:** Business logic validation using `testify/mock`.
```bash
make test-unit
```


* **Integration Tests:** Database layer testing using **Testcontainers**. It spins up an isolated PostgreSQL instance, applies migrations via Goose, and verifies repository logic.
```bash
make test-integration
```


* **Linting:**
```bash
make lint
```
