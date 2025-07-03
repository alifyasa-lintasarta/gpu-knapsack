# Variables
BINARY_NAME=app
OUT_DIR=out
SRC_DIR=src

# Default target
build:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BINARY_NAME) $(SRC_DIR)/*.go

format:
	gofmt -w $(SRC_DIR)

# Clean build artifacts
clean:
	rm -rf $(OUT_DIR)

.PHONY: build clean format
