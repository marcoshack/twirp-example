package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type GetAllInput struct {
	Limit int32
}

type GetAllOutput struct {
	Entries []HelloEntry
}

func (d *HelloDAO) GetAll(ctx context.Context, input *GetAllInput) (*GetAllOutput, error) {
	var entries []HelloEntry
	output, err := d.ddbClient.Scan(ctx, &dynamodb.ScanInput{
		TableName: &d.tableName,
		Limit:     &input.Limit,
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(output.Items, &entries)
	if err != nil {
		return nil, err
	}

	return &GetAllOutput{
		Entries: entries,
	}, nil
}
