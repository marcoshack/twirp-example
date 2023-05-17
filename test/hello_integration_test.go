package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/marcoshack/twirp-example/rpc/helloworld"
	"github.com/stretchr/testify/require"
)

var (
	testServerEndpoint = "http://localhost:8080"
	testHTTPClient     = http.DefaultClient
	testRequestTimeout = 5 * time.Second
)

// TestMain is used to run integration tests. In order to run integration tests, the
// environment variable RUN_INTEG_TESTS must be set to 1. The server endpoint is set
// by the environment variable HELLO_SERVER_ENDPOINT.
func TestMain(m *testing.M) {
	if os.Getenv("RUN_INTEG_TESTS") == "" {
		fmt.Println("=== SKIP: integration tests not enabled, use RUN_INTEG_TESTS=1")
		return
	}
	endpoint := os.Getenv("HELLO_SERVER_ENDPOINT")
	if endpoint != "" {
		testServerEndpoint = endpoint
	}
	m.Run()
}

func TestIntegration_Hello(t *testing.T) {
	message := "integ-test-subject"
	client := helloworld.NewHelloWorldJSONClient(testServerEndpoint, testHTTPClient)
	ctx, cancel := context.WithTimeout(context.Background(), testRequestTimeout)
	defer cancel()

	response, err := client.Hello(ctx, &helloworld.HelloReq{Subject: message})
	require.NoError(t, err)
	require.Contains(t, response.Text, message)
}

func TestIntegration_Hello_WithEmptyRequest(t *testing.T) {
	client := helloworld.NewHelloWorldJSONClient(testServerEndpoint, testHTTPClient)
	ctx, cancel := context.WithTimeout(context.Background(), testRequestTimeout)
	defer cancel()

	response, err := client.Hello(ctx, &helloworld.HelloReq{})
	require.NoError(t, err)
	require.NotEmpty(t, response.Text)
}
