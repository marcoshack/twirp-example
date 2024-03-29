# Twirp-example

## Workspace setup

### Ubuntu

Set `GOPATH`, `GOPROXY` and add `$GOPATH/bin` to `$PATH`. e.g.:

> You need to set `GOPROXY=direct` if using a network that blocks proxy.golang.org, otherwise you can skip that envvar.

```sh
echo "export GOPATH=$HOME/work/go >> ~/.profile
echo "export GOPROXY=direct" >> ~/.profile
echo "export PATH=$GOPATH/bin:$PATH >> ~/.profile
```

Install dependencies 

```sh
sudo apt update
sudo apt install make protobuf-compiler
```

### Checkout, install and build

```sh
mkdir -p $HOME/work
cd $HOME/work
git clone git@github.com:marcoshack/twirp-example.git
cd twirp-example
```
```sh
go mod download
go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/lint/golint@latest
```
```sh
make
```
