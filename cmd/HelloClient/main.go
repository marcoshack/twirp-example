package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	service "twirp-example/rpc/helloworld"

	"github.com/goombaio/namegenerator"
)

func main() {
	// parse CLI options
	var serviceEndpoint string
	var messageInput string
	var count int
	var delayInMillis int
	flag.StringVar(&serviceEndpoint, "s", "http://localhost:8080", "service endpoint")
	flag.StringVar(&messageInput, "m", "", "message to send")
	flag.IntVar(&count, "c", 1, "number of messages to send")
	flag.IntVar(&delayInMillis, "d", 500, "delay between messages, in milliseconds")
	flag.Parse()
	nameGenerator := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())

	for i := 1; i <= count; i++ {
		messageToSend := messageInput
		if messageToSend == "" {
			messageToSend = nameGenerator.Generate()
		}

		// initialize service client and send Hello request
		client := service.NewHelloWorldProtobufClient(serviceEndpoint, &http.Client{})
		resp, err := client.Hello(context.Background(), &service.HelloReq{Subject: messageToSend})

		if err != nil {
			fmt.Printf("[ERROR] Failed calling HelloServer: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("(%d/%d) Response: %s\n", i, count, resp.Text)
		time.Sleep(time.Duration(delayInMillis) * time.Millisecond)
	}
}
