package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/marcoshack/twirp-example/internal/storage"
	"github.com/rs/zerolog"
	"github.com/twitchtv/twirp"

	service "github.com/marcoshack/twirp-example/rpc/helloworld"
)

type HelloServer struct {
	config *HelloServerConfig
	dao    *storage.HelloDAO
}

type HelloServerConfig struct {
	BindAddr     string
	DDBEndpoint  string
	DDBTableName string
	DDBTimeout   time.Duration
	Logger       *zerolog.Logger
	ServerHooks  []*twirp.ServerHooks
}

func (c *HelloServerConfig) String() string {
	return fmt.Sprintf("HelloServerConfig{BindAddr: %s, DDBEndpoint: %s, DDBTableName: %s, DDBTimeout: %s}",
		c.BindAddr, c.DDBEndpoint, c.DDBTableName, c.DDBTimeout)
}

type HelloServerOption func(*HelloServerConfig)

var (
	DefaultServerConfig = HelloServerConfig{
		BindAddr:     ":8080",
		DDBEndpoint:  "http://localhost:8000",
		DDBTableName: "HelloWorldTable",
		DDBTimeout:   10 * time.Second,
	}
)

func WithBindAddr(addr string) HelloServerOption {
	return func(c *HelloServerConfig) {
		c.BindAddr = addr
	}
}

func WithDDBEndpoint(endpoint string) HelloServerOption {
	return func(c *HelloServerConfig) {
		c.DDBEndpoint = endpoint
	}
}

func WithDDBTableName(tableName string) HelloServerOption {
	return func(c *HelloServerConfig) {
		c.DDBTableName = tableName
	}
}

func WithLogger(logger *zerolog.Logger) HelloServerOption {
	return func(c *HelloServerConfig) {
		c.Logger = logger
	}
}

func WithServerHooks(hooks ...*twirp.ServerHooks) HelloServerOption {
	return func(c *HelloServerConfig) {
		newHooks := make([]*twirp.ServerHooks, len(hooks)+len(c.ServerHooks))
		newHooks = append(newHooks, c.ServerHooks...)
		newHooks = append(newHooks, hooks...)
		c.ServerHooks = newHooks
	}
}

func NewHelloWorldServer(options ...HelloServerOption) (*HelloServer, error) {
	config := DefaultServerConfig
	for _, option := range options {
		option(&config)
	}

	dao, err := createDAO(&config)
	if err != nil {
		return nil, errors.New("failed to create DAO")
	}

	return &HelloServer{
		config: &config,
		dao:    dao,
	}, nil
}

func (s *HelloServer) Start() error {
	twirpOptions := twirp.WithServerHooks(twirp.ChainHooks(s.config.ServerHooks...))
	twirpHandler := service.NewHelloWorldServer(s, twirpOptions)

	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)
	server := &http.Server{
		Addr:              s.config.BindAddr,
		Handler:           twirpHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	s.config.Logger.Info().Stringer("config", s.config).Msg("starting server")
	return server.ListenAndServe()
}

func (s *HelloServer) Stop() {
}

func createDAO(config *HelloServerConfig) (*storage.HelloDAO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DDBTimeout)
	defer cancel()
	ddbClient, err := storage.CreateDynamoDBLocalClient(ctx, config.DDBEndpoint)
	if err != nil {
		panic(fmt.Sprintf("failed to create DynamoDB client: %v", err))
	}

	tables, err := ddbClient.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %v", err)
	}
	exists := false
	for _, table := range tables.TableNames {
		if table == config.DDBTableName {
			exists = true
		}
	}
	if !exists {
		config.Logger.Info().Msgf("table %s does not exist, creating it", config.DDBTableName)
		_, err = ddbClient.CreateTable(ctx, storage.CreateTableInput(config.DDBTableName))
		if err != nil {
			panic(fmt.Sprintf("failed to create table: %v", err))
		}
	}
	dao := storage.NewHelloDAO(ddbClient, config.DDBTableName)
	return dao, nil
}
