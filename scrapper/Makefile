GO_TEST_FLAGS ?= -v -race -cover
INTEGRATION_TAG := integration

.PHONY: test test-integration lint clean

test:
	@echo "Launching go tests"
	go test ./internal/usecase $(GO_TEST_FLAGS)

test-integration:
	@echo "Launching integration tests"
	TESTCONTAINERS_RYUK_DISABLED=true go test -v -tags=$(INTEGRATION_TAG) -count=1 -race -timeout=5m ./internal/storage

cover-html:
	@echo "---Generating coverage report---"
	go test -coverprofile=coverage.out ./internal/usecase
	go tool cover -html=coverage.out
	@rm coverage.out

lint:
	@echo "---Linter launching---"
	golangci-lint run

clean:
	@echo "---Cleaning---"
	@rm -f coverage.out
	@go clean -testcache