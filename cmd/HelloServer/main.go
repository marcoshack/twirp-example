package main

import (
	"net/http"

	"twirp-example/internal/server"

	service "twirp-example/rpc/helloworld"

	"github.com/rs/zerolog/log"
)

func main() {
	twirpHandler := service.NewHelloWorldServer(&server.HelloWorldServer{})
	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)
	port := ":8080"

	log.Info().Str("port", port).Msg("HelloServer started")
	http.ListenAndServe(port, mux)
}
