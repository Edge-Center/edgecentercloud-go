PROJECT_DIR = $(shell pwd)
BIN_DIR = $(PROJECT_DIR)/bin

.PHONY: lint
lint:
	test -f $(BIN_DIR)/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.2
	$(BIN_DIR)/golangci-lint run
