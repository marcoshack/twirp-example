package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	service "github.com/marcoshack/twirp-example/rpc/helloworld"
	"github.com/rs/zerolog"
	"github.com/twitchtv/twirp"

	"github.com/goombaio/namegenerator"
)

func main() {
	// parse CLI options
	var serviceEndpoint, messageInput string
	var count, delayInMillis int
	flag.StringVar(&serviceEndpoint, "s", "http://localhost:8080", "service endpoint")
	flag.StringVar(&messageInput, "m", "", "message to send")
	flag.IntVar(&count, "c", 1, "number of messages to send")
	flag.IntVar(&delayInMillis, "d", 500, "delay between messages, in milliseconds")
	flag.Parse()

	logger := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})).With().Timestamp().Logger()

	nameGenerator := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())
	httpClient := &http.Client{}
	client := service.NewHelloWorldProtobufClient(serviceEndpoint, httpClient)

	for i := 1; i <= count; i++ {
		messageToSend := messageInput
		if messageToSend == "" {
			messageToSend = nameGenerator.Generate()
		}

		ctx := context.Background()

		// Set X-Request-ID header
		requestIDHeader := make(http.Header)
		requestIDHeader.Add("X-Request-ID", "test-request-id")
		ctx, err := twirp.WithHTTPRequestHeaders(ctx, requestIDHeader)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to set request ID header")
		}

		// Call service
		resp, err := client.Hello(ctx, &service.HelloReq{Subject: messageToSend})

		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to call HelloServer")
		}

		logger.Info().Str("response", resp.Text).Msg("Response")
		if i < count {
			time.Sleep(time.Duration(delayInMillis) * time.Millisecond)
		}
	}
}
