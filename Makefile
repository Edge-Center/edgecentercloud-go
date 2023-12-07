PROJECT_DIR = $(shell pwd)
BIN_DIR = $(PROJECT_DIR)/bin

.PHONY: lint
lint:
	test -f $(BIN_DIR)/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.2
	$(BIN_DIR)/golangci-lint run

.PHONY: test
test:
	go test -v -timeout=2m

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=cover.out -covermode=atomic
	@go-test-coverage --config=coverage.yml
