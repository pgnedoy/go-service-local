package sqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/mock"
)

type MockSQSClient struct {
	sqsiface.SQSAPI
	mock.Mock
}

func (m *MockSQSClient) ListQueues(in *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*sqs.ListQueuesOutput), args.Error(1)
}

func (m *MockSQSClient) GetQueueUrl(in *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*sqs.GetQueueUrlOutput), args.Error(1)
}

func (m *MockSQSClient) GetQueueAttributes(in *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*sqs.GetQueueAttributesOutput), args.Error(1)
}

func (m *MockSQSClient) SendMessage(in *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*sqs.SendMessageOutput), args.Error(1)
}

func (m *MockSQSClient) DeleteMessage(in *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*sqs.DeleteMessageOutput), args.Error(1)
}

func (m *MockSQSClient) DeleteMessageBatch(in *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*sqs.DeleteMessageBatchOutput), args.Error(1)
}

func (m *MockSQSClient) ReceiveMessage(in *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*sqs.ReceiveMessageOutput), args.Error(1)
}

func (m *MockSQSClient) PurgeQueue(in *sqs.PurgeQueueInput) (*sqs.PurgeQueueOutput, error) {
	args := m.Called(in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*sqs.PurgeQueueOutput), args.Error(1)
}
