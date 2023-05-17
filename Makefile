BUILD_DIR = ./build

all: test gosec build

help:
	@sed -ne '/@sed/!s/## /-- /p' $(MAKEFILE_LIST)

test: ## Run unit tests
	mkdir -p $(BUILD_DIR)
	go test -v -coverprofile=$(BUILD_DIR)/coverage.out ./...
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	go tool cover -func=$(BUILD_DIR)/coverage.out

integ-test: ## Run integration tests
	RUN_INTEG_TESTS=1 go test -v -count=1 -run "^TestIntegration" ./test

gosec: ## Run Go Security Checker (gosec)
	gosec ./...

build: twirp-generate build-server build-client ## Build server and client

run: build-server
	$(BUILD_DIR)/bin/HelloServer

twirp-generate: ## Generate Twirp and Protobuf Go code
	protoc --go_out=. --twirp_out=. rpc/helloworld/helloworld.proto

build-server: ## Build HelloServer
	go build -o $(BUILD_DIR)/bin/ ./cmd/HelloServer/

build-client: ## Build HelloClient
	go build -o $(BUILD_DIR)/bin ./cmd/HelloClient/

docker-build: ## Build docker image for HelloServer
	docker build -t marcoshack/twirp-example:latest .

docker-run: ## Run HelloServer inside a containers
	docker run --rm -p 8080:8080 marcoshack/twirp-example

workspace: ## Setup your local workspace to build the project
# TODO Doesn't work from Makefile
#	cat tools.go | grep _ | awk -F'\"' '{print $2}' | xargs -tI % go install %

clean: ## Clean workspace
	rm -rf $(BUILD_DIR)

.PHONY: all test gosec build twirp-generate build-server build-client docker docker-run clean
