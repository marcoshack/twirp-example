package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	service "twirp-example/helloworld"
)

func main() {
	client := service.NewHelloWorldProtobufClient("http://localhost:8080", &http.Client{})
	resp, err := client.Hello(context.Background(), &service.HelloReq{Subject: fmt.Sprintf("there, it's %s", time.Now().Format(time.RFC3339))})
	if err != nil {
		fmt.Printf("[ERROR] Failed calling HelloServer: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("[INFO] Response: %s\n", resp.Text)
}
