package server

import "github.com/marcoshack/twirp-example/internal/storage"

func NewHelloWorldServer(dao *storage.HelloDAO) *HelloWorldServer {
	return &HelloWorldServer{dao: dao}
}

type HelloWorldServer struct {
	dao *storage.HelloDAO
}
