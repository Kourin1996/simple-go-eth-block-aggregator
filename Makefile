BINARY_NAME := $(PROJECT_NAME)
SRC_DIR := ./cmd/
MAIN_SRC := $(SRC_DIR)/main.go
BUILD_FILE := ./app

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building..."
	@go build $(MAIN_SRC)

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm $(BUILD_FILE)

.PHONY: run
run:
	@echo "Running..."
	@go run $(MAIN_SRC)

.PHONY: fmt
fmt:
	@echo "Formatting..."
	@go fmt ./...

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make          - build"
	@echo "  make build    - build"
	@echo "  make run      - run"
	@echo "  make help     - display helps"
