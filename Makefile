# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINST=$(GOCMD) install

#Binary Name
BINARY_NAME=main

# Build
build:
	@$(GOBUILD) -o $(BINARY_NAME) ./cmd/http
	@echo "📦 Build Done"

# Test
test:
	@$(GOTEST) -v ./...
	@echo "🧪 Test Completed"

# Run
run:
	@echo "🚀 Running App"
	@./$(BINARY_NAME)

# Generate Mocks
generate-mocks:
	@$(GOINST) github.com/golang/mock/mockgen@v1.6.0
	@./scripts/generate-mocks.sh

# Dev
dev:build
    @echo "🚀 Running App"
    @./$(BINARY_NAME)

