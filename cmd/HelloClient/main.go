package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	service "twirp-example/rpc/helloworld"
)

func main() {
	// parse CLI options
	var serviceEndpoint string
	flag.StringVar(&serviceEndpoint, "s", "http://localhost:8080", "service endpoint")

	// initialize service client and send Hello request
	client := service.NewHelloWorldProtobufClient(serviceEndpoint, &http.Client{})
	resp, err := client.Hello(context.Background(), &service.HelloReq{Subject: fmt.Sprintf("there, it's %s", time.Now().Format(time.RFC3339))})

	if err != nil {
		fmt.Printf("[ERROR] Failed calling HelloServer: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("[INFO] Response: %s\n", resp.Text)
}
