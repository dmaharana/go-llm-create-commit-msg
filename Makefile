PROJECT_NAME := ccm
OUTPUT_DIR := bin

# Get the Git commit SHA1
GIT_COMMIT_SHA1 := $(shell git rev-parse --short HEAD)

# Build flags to include the Git commit SHA1
LDFLAGS := -X main.gitCommitSHA1=$(GIT_COMMIT_SHA1) -s -w

GO_BUILD_CMD := go build -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME)

OS := $(shell go env GOOS)


.PHONY: all linux mac windows clean current

all: linux mac windows

.DEFAULT_GOAL := current

current:
	@echo "Building for $(OS)..."
ifeq ($(OS),linux)
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD)-linux-amd64 cmd/main.go
	@echo "Built $(OUTPUT_DIR)/$(PROJECT_NAME)-linux-amd64"
else ifeq ($(OS),darwin)
	GOOS=darwin GOARCH=amd64 $(GO_BUILD_CMD)-mac-amd64 cmd/main.go
	@echo "Built $(OUTPUT_DIR)/$(PROJECT_NAME)-mac-amd64"
else ifeq ($(OS),windows)
	GOOS=windows GOARCH=amd64 $(GO_BUILD_CMD)-windows-amd64.exe cmd/main.go
	@echo "Built $(OUTPUT_DIR)/$(PROJECT_NAME)-windows-amd64.exe"
else
	@echo "Unsupported OS: $(OS)"
	exit 1
endif

linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD)-linux-amd64 cmd/main.go
	@echo "Built $(OUTPUT_DIR)/$(PROJECT_NAME)-linux-amd64"

mac:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GO_BUILD_CMD)-mac-amd64 cmd/main.go
	@echo "Built $(OUTPUT_DIR)/$(PROJECT_NAME)-mac-amd64"

windows:
	GOOS=windows GOARCH=amd64 $(GO_BUILD_CMD)-windows-amd64.exe cmd/main.go
	@echo "Built $(OUTPUT_DIR)/$(PROJECT_NAME)-windows-amd64.exe"

clean:
	@rm -rf $(OUTPUT_DIR)
