package main

import (
	"flag"
	"net/http"

	"twirp-example/internal/server"

	service "twirp-example/rpc/helloworld"

	"github.com/rs/zerolog/log"
)

func main() {
	// parse CLI options
	var bindAddr string
	flag.StringVar(&bindAddr, "b", ":8080", "server listening address (e.g. :8080 or localhost:8080)")
	flag.Parse()

	// initialize server
	twirpHandler := service.NewHelloWorldServer(&server.HelloWorldServer{})
	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)

	// start server
	log.Info().Str("port", bindAddr).Msg("HelloServer started")
	http.ListenAndServe(bindAddr, mux)
}
