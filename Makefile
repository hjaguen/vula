BINARY_NAME=vula
BUILD_DIR=bin
CMD_DIR=cmd/vula
GO=go

.PHONY: all build test clean run doctor hud install fmt lint

all: build

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to ~/.local/bin..."
	@mkdir -p $(HOME)/.local/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "Installed successfully! Run 'vula doctor' or 'vula hud'"

test:
	$(GO) test -v ./...

doctor: build
	./$(BUILD_DIR)/$(BINARY_NAME) doctor

hud: build
	./$(BUILD_DIR)/$(BINARY_NAME) hud

fmt:
	$(GO) fmt ./...

clean:
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."
