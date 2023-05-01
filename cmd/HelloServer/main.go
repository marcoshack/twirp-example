package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/marcoshack/twirp-example/internal/server"
	"github.com/marcoshack/twirp-example/internal/storage"

	service "github.com/marcoshack/twirp-example/rpc/helloworld"

	"github.com/rs/zerolog/log"
)

func main() {
	// parse CLI options
	var bindAddr, tableName, ddbEndpoint string
	flag.StringVar(&bindAddr, "b", "localhost:8080", "server listening address")
	flag.StringVar(&ddbEndpoint, "e", "http://localhost:8000", "DynamoDB endpoint URL")
	flag.StringVar(&tableName, "t", "HelloTable", "DynamoDB table name")
	flag.Parse()

	// create DAO
	ddbClient, err := storage.CreateDynamoDBLocalClient(context.TODO(), ddbEndpoint)
	if err != nil {
		panic(fmt.Sprintf("failed to create DynamoDB client: %v", err))
	}

	tables, err := ddbClient.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		panic(fmt.Sprintf("failed to list tables: %v", err))
	}
	exists := false
	for _, table := range tables.TableNames {
		if table == tableName {
			exists = true
		}
	}
	if !exists {
		log.Info().Msgf("Table %s does not exist, creating it...", tableName)
		_, err = ddbClient.CreateTable(context.TODO(), storage.CreateTableInput(tableName))
		if err != nil {
			panic(fmt.Sprintf("failed to create table: %v", err))
		}
	}
	dao := storage.NewHelloDAO(ddbClient, tableName)

	// initialize Twirp service and HTTP server
	helloServer := server.NewHelloWorldServer(dao)
	twirpHandler := service.NewHelloWorldServer(helloServer)
	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)

	// capture SIGINT and SIGTERM
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Info().Msg("HelloServer shutting down...")
		// TODO release resources, wait go routines, etc...
		log.Info().Msg("Done.")
		os.Exit(0)
	}()

	// start server
	log.Info().Str("port", bindAddr).Int("pid", os.Getpid()).Msg("Starting HelloServer")
	http.ListenAndServe(bindAddr, mux)
}
