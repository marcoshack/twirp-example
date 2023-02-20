package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	service "twirp-example/rpc/helloworld"
)

func main() {
	// parse CLI options
	var serviceEndpoint string
	var messageToSend string
	flag.StringVar(&serviceEndpoint, "s", "http://localhost:8080", "service endpoint")
	flag.StringVar(&messageToSend, "m", "world", "message to send")
	flag.Parse()

	// initialize service client and send Hello request
	client := service.NewHelloWorldProtobufClient(serviceEndpoint, &http.Client{})
	resp, err := client.Hello(context.Background(), &service.HelloReq{Subject: messageToSend})

	if err != nil {
		fmt.Printf("[ERROR] Failed calling HelloServer: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("[INFO] Response: %s\n", resp.Text)
}
