package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	dao, err := createDAO(ddbEndpoint, tableName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create DAO")
	}

	server, err := createServer(dao, bindAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Twirp handler")
	}

	captureStopSignals()

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func createDAO(ddbEndpoint string, tableName string) (*storage.HelloDAO, error) {
	ddbClient, err := storage.CreateDynamoDBLocalClient(context.TODO(), ddbEndpoint)
	if err != nil {
		panic(fmt.Sprintf("failed to create DynamoDB client: %v", err))
	}

	tables, err := ddbClient.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %v", err)
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
	return dao, nil
}

func createServer(dao *storage.HelloDAO, bindAddr string) (*http.Server, error) {
	helloServer := server.NewHelloWorldServer(dao)
	twirpHandler := service.NewHelloWorldServer(helloServer)
	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)
	server := &http.Server{
		Addr:              bindAddr,
		Handler:           twirpHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return server, nil
}

func captureStopSignals() {
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Info().Msg("HelloServer shutting down...")
		// TODO release resources, wait go routines, etc...
		log.Info().Msg("Done.")
		os.Exit(0)
	}()
}
