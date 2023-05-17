package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// HelloEntry represents a hello world entry in the database.
type HelloEntry struct {
	ID        string    `dynamodbav:"PK,omitempty"`
	CreatedAt time.Time `dynamodbav:"SK,omitempty"`
	Message   string    `dynamodbav:"message,omitempty"`
}

// HelloInput is the input for the AddHelloWorld function.
type HelloInput struct {
	Message string `json:"message"`
}

// AddHelloWorld saves the hello world message to the database.
func (d *HelloDAO) AddHelloWorld(ctx context.Context, input *HelloInput) (*HelloEntry, error) {
	log.Ctx(ctx).Debug().Msg("adding hello world")
	entry := &HelloEntry{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		Message:   input.Message,
	}
	av, err := attributevalue.MarshalMap(entry)
	if err != nil {
		return nil, err
	}

	_, err = d.ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      av,
	})
	if err != nil {
		return nil, err
	}
	return entry, nil
}
