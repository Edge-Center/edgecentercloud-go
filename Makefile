PROJECT_DIR = $(shell pwd)
BIN_DIR = $(PROJECT_DIR)/bin
CLOUD_TEST_DIR = $(PROJECT_DIR)/e2e/cloud
CLOUD_ENV_TESTS_FILE = $(CLOUD_TEST_DIR)/.env

.PHONY: lint
lint:
	test -f $(BIN_DIR)/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.64.8
	$(BIN_DIR)/golangci-lint run

.PHONY: test
test:
	go test -v -timeout=2m

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@v2.11.4

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=cover.out -covermode=atomic
	@go-test-coverage --config=coverage.yml

.PHONY: download-local-env
download-local-env:
	@if [ -z "${VAULT_TOKEN}" ] || [ -z "${VAULT_ADDR}" ]; then \
		echo "ERROR: Vault environment is not set, please setup VAULT_ADDR and VAULT_TOKEN environment variables" && exit 1;\
	fi
	vault kv get -field e2e.env cloud/edgecentercloud-go/e2e > $(CLOUD_ENV_TESTS_FILE)