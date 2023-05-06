package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/marcoshack/twirp-example/internal/server"

	"github.com/rs/zerolog"
)

func main() {
	// parse CLI options
	var bindAddr, ddbEndpoint, ddbTableName string
	flag.StringVar(&bindAddr, "b", "localhost:8080", "server listening address")
	flag.StringVar(&ddbEndpoint, "e", "http://localhost:8000", "DynamoDB endpoint URL")
	flag.StringVar(&ddbTableName, "t", "HelloTable", "DynamoDB table name")
	flag.Parse()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	server, err := server.NewHelloWorldServer(
		server.WithBindAddr(bindAddr),
		server.WithDDBEndpoint(ddbEndpoint),
		server.WithDDBTableName(ddbTableName),
		server.WithLogger(&logger),
		server.WithServerHooks(server.NewRequestLoggingServerHooks(&logger)),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create server")
	}

	onShutdown(func() int {
		logger.Info().Msg("shutting down...")
		server.Stop()
		logger.Info().Msg("done")
		return 0
	})

	err = server.Start()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}
}

func onShutdown(teardown func() int) {
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		os.Exit(teardown())
	}()
}
