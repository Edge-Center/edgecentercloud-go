# ENVS
PROJECT_DIR = $(shell pwd)
BUILD_DIR = $(PROJECT_DIR)/bin

# CHECKS
vet:
	go vet ./...

fmt:
	go fmt ./...

gofumpt:
	go install mvdan.cc/gofumpt@v0.4.0
	gofumpt -l -w .

linters:
	@test -f $(BUILD_DIR)/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.2
	@$(BUILD_DIR)/golangci-lint run

.PHONY: vet fmt linters gofumpt
