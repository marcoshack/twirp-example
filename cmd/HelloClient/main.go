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

const (
	requestTimeout = time.Second * 10
)

func main() {
	// parse CLI options
	var serviceEndpoint, messageInput string
	var count, delayInMillis int
	var getAll bool
	flag.StringVar(&serviceEndpoint, "s", "http://localhost:8080", "service endpoint")
	flag.StringVar(&messageInput, "m", "", "message to send")
	flag.IntVar(&count, "c", 1, "number of messages to send (-m) or retrieve (-g)")
	flag.IntVar(&delayInMillis, "d", 500, "delay between messages, in milliseconds")
	flag.BoolVar(&getAll, "g", false, "get all hello messages (cannot be used with -m or -c)")
	flag.Parse()

	if getAll && messageInput != "" {
		flag.Usage()
		return
	}

	logger := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})).With().Timestamp().Logger()

	httpClient := &http.Client{}
	client := service.NewHelloWorldProtobufClient(serviceEndpoint, httpClient)

	if getAll {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		resp, err := client.GetAll(ctx, &service.GetAllReq{Limit: int32(count)})
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to call GetAll")
		}

		for _, item := range resp.Items {
			logger.Info().Interface("item", item).Msg("Response")
		}
		return
	}

	nameGenerator := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())
	for i := 1; i <= count; i++ {
		messageToSend := messageInput
		if messageToSend == "" {
			messageToSend = nameGenerator.Generate()
		}

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()

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
