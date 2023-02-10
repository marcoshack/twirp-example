all: twirp-generate build-server build-client

twirp-generate:
	rm -rf helloworld/
	protoc --go_out=. --twirp_out=. schema/service.proto
	mv github.com/marcoshack/twirp-example/helloworld .
	rm -rf github.com

build-server: twirp-generate
	go build -o ./bin/ ./cmd/HelloServer/

build-client: twirp-generate
	go build -o ./bin/ ./cmd/HelloClient/

clean:
	rm -rf ./bin/
	rm -rf ./helloworld/
