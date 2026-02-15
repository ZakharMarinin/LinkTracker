GO_TEST_FLAGS ?= -v -race -cover

.PHONY: test clean

test:
	@echo "---Launching go tests---"
	go test ./internal/usecase $(GO_TEST_FLAGS)

cover-html:
	@echo "---Generating coverage report---"
	go test -coverprofile=coverage.out ./internal/usecase
	go tool cover -html=coverage.out
	@rm coverage.out

clean:
	@echo "---Cleaning---"
	@rm -f coverage.out
	@go clean -testcache