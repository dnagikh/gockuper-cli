APP_NAME=gockuper-cli
BUILD_DIR=build
GOOS_LINUX=linux
GOARCH_AMD64=amd64

.PHONY: build test lint release

build:
	@echo "🚀 Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .

build-linux:
	@echo "🐧 Building $(APP_NAME) for Linux x86_64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_AMD64) go build -o $(BUILD_DIR)/linux/$(APP_NAME) .

run:
	@echo "🏃 Running $(APP_NAME)..."
	@go run . backup

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

lint:
	@echo "🔍 Running golangci-lint..."
	@golangci-lint run ./...

lint-fix:
	@echo "🛠 Fixing with golangci-lint..."
	@golangci-lint run --fix ./...