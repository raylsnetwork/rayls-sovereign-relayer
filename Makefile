.PHONY: all build clean run

# Variables
BINARY_NAME_RELAYER=relayer
BINARY_NAME_PUBRELAYER=pubrelayer
BINARY_NAME_CTS=cts

SRC_DIR_RELAYER=private-relayer/cmd
SRC_DIR_PUBRELAYER=public-relayer/cmd
SRC_DIR_CTS=cts/cmd

MAIN_FILE_RELAYER=$(SRC_DIR_RELAYER)/main.go
MAIN_FILE_PUBRELAYER=${SRC_DIR_PUBRELAYER}/main.go
MAIN_FILE_CTS=$(SRC_DIR_CTS)/main.go
BUILD_DIR=build

# Targets
all: build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_RELAYER) $(MAIN_FILE_RELAYER)
	go build -o ${BUILD_DIR}/$(BINARY_NAME_PUBRELAYER) $(MAIN_FILE_PUBRELAYER)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_CTS) $(MAIN_FILE_CTS)

build-relayer:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_RELAYER) $(MAIN_FILE_RELAYER)

build-pubrelayer:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_PUBRELAYER) $(MAIN_FILE_PUBRELAYER)

build-cts:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_CTS) $(MAIN_FILE_CTS)

proto-cts:
	protoc --go_out=. --go-grpc_out=. \
		--go_opt=module=github.com/raylsnetwork/rayls-privacy-relayer-api \
		--go-grpc_opt=module=github.com/raylsnetwork/rayls-privacy-relayer-api \
		cts/proto/encrypt.proto cts/proto/keys.proto cts/proto/txops.proto \
		--go-grpc_opt=require_unimplemented_servers=false

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME) run --config config.json

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

test-coverage-xml:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	gocov convert coverage.out | gocov-xml > coverage.xml
	@echo "Coverage report generated: coverage.xml"
