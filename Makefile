all: test build

help:
	@sed -ne '/@sed/!s/## /-- /p' $(MAKEFILE_LIST)

test: ## Run unit tests
	go test -v ./...

build: twirp-generate build-server build-client ## Build server and client

twirp-generate: ## Generate Twirp and Protobuf Go code
	protoc --go_out=. --twirp_out=. rpc/helloworld/helloworld.proto

build-server: ## Build HelloServer
	go build -o ./bin/ ./cmd/HelloServer/

build-client: ## Build HelloClient
	go build -o ./bin/ ./cmd/HelloClient/

docker-build: ## Build docker image for HelloServer
	docker build -t marcoshack/twirp-example:latest .

docker-run: ## Run HelloServer inside a containers
	docker run --rm -p 8080:8080 marcoshack/twirp-example

clean: ## Clean workspace
	rm -rf ./bin/

.PHONY: all test build twirp-generate build-server build-client docker docker-run clean
