# Twirp-example

## Workspace setup

### Ubuntu

Set `GOPATH` in your shell profile and add `$GOPATH/bin` to `$PATH`. e.g.:

```
echo "export GOPATH=$HOME/work/go >> ~/.bashrc
echo "export PATH=$GOPATH/bin:$PATH >> ~/.bashrc
```

Install dependencies 

```sh
sudo apt install make protobuf-compiler
```

Install Go dependencies and tools

```sh
go mod download
go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
