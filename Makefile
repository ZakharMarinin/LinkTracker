SERVICES := scrapper tgBot

.PHONY: help lint test-unit test-integration run-with-migrations run-no-migrations stop

help:
	@echo "Commands:"
	@echo "  make lint                 - Run linters for all services"
	@echo "  make test-unit            - Run unit tests"
	@echo "  make run-with-migrations  - Run app WITH migrations enabled"
	@echo "  make run-no-migrations    - Run app WITHOUT migrations"

lint:
	@echo "Running linters..."
	@for service in $(SERVICES); do \
		echo "Checking $$service..."; \
		$(MAKE) -C $$service lint; \
	done

test-unit:
	@echo "Running unit tests..."
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		$(MAKE) -C $$service test; \
	done

test-integration:
	@echo "Running integration test.."
	$(MAKE) -C scrapper test-integration;

run-with-migrations:
	@echo "Starting WITH migrations..."
	docker-compose --profile tools up --build -d

run-no-migrations:
	@echo "Starting WITHOUT migrations..."
	docker-compose up --build -d

stop:
	docker-compose down