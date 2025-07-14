# Variables
BINARY_NAME=gpu-knapsack
OUT_DIR=out
SRC_DIR=src

# Default target
build: $(OUT_DIR)/$(BINARY_NAME)

$(OUT_DIR)/$(BINARY_NAME): $(SRC_DIR)/*.go
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BINARY_NAME) $(SRC_DIR)/*.go

format:
	gofmt -w $(SRC_DIR)

# Clean build artifacts
clean:
	rm -rf $(OUT_DIR)

.PHONY: build clean format
