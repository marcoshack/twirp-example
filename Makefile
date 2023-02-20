all: twirp-generate build-server build-client

twirp-generate:
	protoc --go_out=. --twirp_out=. rpc/helloworld/helloworld.proto

build-server: twirp-generate
	go build -o ./bin/ ./cmd/HelloServer/

build-client: twirp-generate
	go build -o ./bin/ ./cmd/HelloClient/

docker:
	docker build -t marcoshack/twirp-example:latest .

docker-run: docker
	docker run -d --rm -p 8080:8080 marcoshack/twirp-example

clean:
	rm -rf ./bin/
