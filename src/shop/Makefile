
APP_NAME=app
BIN_DIR=bin

.PHONY: build run clean swag

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/main.go

run: build
	./$(BIN_DIR)/$(APP_NAME)

clean:
	rm -rf $(BIN_DIR)

swag:
	swag init --generalInfo cmd/main.go --dir ./internal,./cmd
