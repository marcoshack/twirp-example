package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/marcoshack/twirp-example/internal/storage"
	"github.com/stretchr/testify/require"
)

const (
	testTablePrefix = "TestHelloTable"
)

var (
	testStart = time.Now()
)

func TestHelloDAO_CreateDAO(t *testing.T) {
	ctx := context.Background()
	dao, ddbClient, tableName := createDAO(t, ctx)
	defer cleanUpDAO(t, ctx, ddbClient, tableName)
	require.NotNil(t, dao)
}

func TestHelloDAO_AddHelloWorld(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dao, ddbClient, tableName := createDAO(t, ctx)
	defer cleanUpDAO(t, ctx, ddbClient, tableName)

	entry, err := dao.AddHelloWorld(ctx, &storage.HelloInput{Message: "Hello World"})
	require.NoError(t, err, "add hello world should not return an error")
	require.NotNil(t, entry)
	require.NotEmpty(t, entry.ID)
	require.GreaterOrEqual(t, entry.CreatedAt, testStart)
	require.LessOrEqual(t, entry.CreatedAt, time.Now())
}

func createDAO(t *testing.T, ctx context.Context) (*storage.HelloDAO, *dynamodb.Client, string) {
	ddbClient, err := storage.CreateDynamoDBLocalClient(ctx, "http://dynamodb-local:8000")
	require.NoError(t, err, "create dynamodb client should not return an error")
	tableName := testTablePrefix + "-" + uuid.NewString()

	_, err = ddbClient.CreateTable(ctx, storage.CreateTableInput(tableName))
	require.NoError(t, err, "create table should not return an error")
	return storage.NewHelloDAO(ddbClient, tableName), ddbClient, tableName
}

func cleanUpDAO(t *testing.T, ctx context.Context, ddbClient *dynamodb.Client, tableName string) {
	_, err := ddbClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	require.NoError(t, err, "delete table should not return an error")
}
