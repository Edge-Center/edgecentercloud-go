# ENVS
PROJECT_DIR = $(shell pwd)
BUILD_DIR = $(PROJECT_DIR)/bin
ENV_TESTS_FILE = .env
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
export VAULT_ADDR = https://vault.p.ecnl.ru/

# BUILD
GOOS		?= $(shell go env GOOS)
VERSION		?= $(shell git describe --tags 2> /dev/null || \
			   git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)
GOARCH		?= $(shell go env GOARCH)
LDFLAGS		:= "-w -s -X 'main.AppVersion=${VERSION}'"
CMD_PACKAGE := ./cmd/ec_client
BINARY 		:= ./ec_client

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags $(LDFLAGS) -o $(BINARY) $(CMD_PACKAGE)

# TESTS
TESTS_LIST = $(shell go list ./... | grep -v ./client)

envs_reader:
	go install github.com/joho/godotenv/cmd/godotenv@latest

tests: envs_reader
	godotenv -f $(ENV_TESTS_FILE) go test -count=1 -timeout=2m $(TESTS_LIST) | { grep -v 'no test files'; true; }

# local test run (need to export VAULT_TOKEN env)
jq:
	if test "$(OS)" = "linux"; then \
		curl -L -o jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64; \
	else \
		curl -L -o jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-osx-amd64; \
	fi
	chmod +x jq

vault:
	curl -L -o vault.zip https://releases.hashicorp.com/vault/1.12.3/vault_1.12.3_$(OS)_$(ARCH).zip
	unzip vault.zip && rm -f vault.zip && chmod +x vault

envs:
	vault login -method=token $(VAULT_TOKEN)
	vault kv get -format=json --field data /CLOUD/edgecentercloud-go | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > $(ENV_TESTS_FILE)

local_tests: envs_reader
	godotenv -f $(ENV_TESTS_FILE) go test -count=1 $(TESTS_LIST) | { grep -v 'no test files'; true; }

# CHECKS
vet:
	go vet ./...

fmt:
	go fmt ./...

gofumpt:
	go install mvdan.cc/gofumpt@v0.4.0
	gofumpt -l -w .

checks: vet fmt gofumpt

linters:
	@test -f $(BUILD_DIR)/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.2
	@$(BUILD_DIR)/golangci-lint run

.PHONY: vet fmt linters gofumpt
