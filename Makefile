all: test build

test:
	go test -v ./...

build: twirp-generate build-server build-client

twirp-generate:
	protoc --go_out=. --twirp_out=. rpc/helloworld/helloworld.proto

build-server:
	go build -o ./bin/ ./cmd/HelloServer/

build-client: 
	go build -o ./bin/ ./cmd/HelloClient/

docker:
	docker build -t marcoshack/twirp-example:latest .

docker-run: docker
	docker run -d --rm -p 8080:8080 marcoshack/twirp-example

clean:
	rm -rf ./bin/

.PHONY: all test build twirp-generate build-server build-client docker docker-run clean
