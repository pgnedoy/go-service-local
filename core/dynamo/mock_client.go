package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/mock"
)

type MockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func (m *MockDynamoDBClient) DescribeTable(in *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.DescribeTableOutput), args.Error(1)
}

func (m *MockDynamoDBClient) DescribeTableWithContext(
	ctx context.Context, in *dynamodb.DescribeTableInput, opts ...request.Option) (*dynamodb.DescribeTableOutput, error) {
	var args mock.Arguments

	if len(opts) > 0 {
		args = m.Called(ctx, in, opts)
	} else {
		args = m.Called(ctx, in)
	}

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.DescribeTableOutput), args.Error(1)
}

func (m *MockDynamoDBClient) ListTables(in *dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.ListTablesOutput), args.Error(1)
}

func (m *MockDynamoDBClient) ListTablesWithContext(
	ctx context.Context, in *dynamodb.ListTablesInput, opts ...request.Option) (*dynamodb.ListTablesOutput, error) {
	var args mock.Arguments

	if len(opts) > 0 {
		args = m.Called(ctx, in, opts)
	} else {
		args = m.Called(ctx, in)
	}

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.ListTablesOutput), args.Error(1)
}

func (m *MockDynamoDBClient) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) PutItemWithContext(
	ctx context.Context, in *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
	var args mock.Arguments

	if len(opts) > 0 {
		args = m.Called(ctx, in, opts)
	} else {
		args = m.Called(ctx, in)
	}

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) Query(in *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func (m *MockDynamoDBClient) QueryWithContext(
	ctx aws.Context, in *dynamodb.QueryInput, opts ...request.Option) (*dynamodb.QueryOutput, error) {
	var args mock.Arguments

	if len(opts) > 0 {
		args = m.Called(ctx, in, opts)
	} else {
		args = m.Called(ctx, in)
	}

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func (m *MockDynamoDBClient) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) GetItemWithContext(
	ctx context.Context, in *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	var args mock.Arguments

	if len(opts) > 0 {
		args = m.Called(ctx, in, opts)
	} else {
		args = m.Called(ctx, in)
	}

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	args := m.Called(input)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

func (m *MockDynamoDBClient) ScanWithContext(
	ctx context.Context, in *dynamodb.ScanInput, opts ...request.Option) (*dynamodb.ScanOutput, error) {
	var args mock.Arguments

	if len(opts) > 0 {
		args = m.Called(ctx, in, opts)
	} else {
		args = m.Called(ctx, in)
	}

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

func (m *MockDynamoDBClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(input)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) DeleteItemWithContext(
	ctx context.Context, in *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error) {
	var args mock.Arguments

	if len(opts) > 0 {
		args = m.Called(ctx, in, opts)
	} else {
		args = m.Called(ctx, in)
	}

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(input)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) UpdateItemWithContext(
	ctx context.Context, in *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
	var args mock.Arguments

	if len(opts) > 0 {
		args = m.Called(ctx, in, opts)
	} else {
		args = m.Called(ctx, in)
	}

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}
