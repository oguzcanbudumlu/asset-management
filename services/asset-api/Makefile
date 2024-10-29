# Project directory and binary paths
BIN_DIR=../../bin
WALLET_BIN=$(BIN_DIR)/wallet-api

# Build the wallet-api binary
.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -o $(WALLET_BIN) main.go

# Generate Swagger documentation
.PHONY: swagger
swagger:
	swag init -g main.go -o docs

# Clean up generated files
.PHONY: clean
clean:
	rm -rf $(WALLET_BIN) docs

# Run the wallet-api
.PHONY: run
run:
	go run .
