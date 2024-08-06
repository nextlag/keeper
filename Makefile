# Include environment variables
include .env
export

# Run docker-compose
compose:
	docker-compose up --build -d postgres app

# Cross-platform build
BINARY_NAME=keeper
SRC_DIR=./cmd/client/
BUILD_DIR=.
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Default platform
OS ?= darwin/amd64

# Validate platform
ifeq ($(filter $(OS), $(PLATFORMS)),)
  $(error Invalid platform $(OS). Please choose from: $(PLATFORMS))
endif

# Build target for the specified OS and install
client:
	@echo "Building $(BINARY_NAME) for $(OS)..."
	GOOS=$(word 1, $(subst /, ,$(OS))) GOARCH=$(word 2, $(subst /, ,$(OS))) go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	ln -sf $(abspath $(BUILD_DIR)/$(BINARY_NAME)) /usr/local/bin/$(BINARY_NAME)

# Clean target (also uninstalls)
clean:
	@echo "Cleaning build directories and uninstalling $(BINARY_NAME) from /usr/local/bin..."
	rm -f /usr/local/bin/$(BINARY_NAME)
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)

# Generate swagger docs
swag:
	swag init -g cmd/server/main.go

# Linter and download settings
LINTER := golangci-lint

lint:
	@echo "=== Lint ==="
	$(LINTER) --version
	$(LINTER) cache clean && $(LINTER) run

# Run tests
test:
	go test -v -cover -race ./internal/...

# Application name and files
APP := keeper
KEY := key.pem
BASE64 := key_base64
INLINE := key_base64_inline
PUBLIC := public_key.pem
PUBLIC_BASE64 := public_key_base64
PUBLIC_INLINE := public_key_base64_inline

# Key management
collectKeys: key_pem base64 inline public_pem public_base64 public_inline clean_key

key_pem:
	openssl genpkey -algorithm RSA -out $(KEY)

base64:
	openssl base64 -in $(KEY) -out $(BASE64)

inline:
	tr -d '\n' < $(BASE64) > $(INLINE)

public_pem:
	openssl rsa -in $(KEY) -pubout -out $(PUBLIC)

public_base64:
	openssl base64 -in $(PUBLIC) -out $(PUBLIC_BASE64)

public_inline:
	tr -d '\n' < $(PUBLIC_BASE64) > $(PUBLIC_INLINE)

# Clean key files
clean_key:
	rm -f $(KEY) $(BASE64) $(PUBLIC) $(PUBLIC_BASE64)

# Phony targets
.PHONY: client clean swag lint test collectKeys key_pem base64 inline public_pem public_base64 public_inline clean_key compose-up
