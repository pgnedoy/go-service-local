package sqs

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	conn sqsiface.SQSAPI
}

func newMockClient(conn sqsiface.SQSAPI) *mockClient {
	return &mockClient{
		conn: conn,
	}
}

func (c *mockClient) TestListQueues() (*sqs.ListQueuesOutput, error) {
	in := &sqs.ListQueuesInput{}

	return c.conn.ListQueues(in)
}

func (c *mockClient) TestGetQueueAttributes() (*sqs.GetQueueAttributesOutput, error) {
	in := &sqs.GetQueueAttributesInput{}

	return c.conn.GetQueueAttributes(in)
}

func (c *mockClient) TestGetQueueUrl() (*sqs.GetQueueUrlOutput, error) {
	in := &sqs.GetQueueUrlInput{}

	return c.conn.GetQueueUrl(in)
}

func (c *mockClient) TestSendMessage() (*sqs.SendMessageOutput, error) {
	in := &sqs.SendMessageInput{}

	return c.conn.SendMessage(in)
}

func (c *mockClient) TestDeleteMessage() (*sqs.DeleteMessageOutput, error) {
	in := &sqs.DeleteMessageInput{}

	return c.conn.DeleteMessage(in)
}

func (c *mockClient) TestDeleteMessageBatch() (*sqs.DeleteMessageBatchOutput, error) {
	in := &sqs.DeleteMessageBatchInput{}

	return c.conn.DeleteMessageBatch(in)
}

func (c *mockClient) TestReceiveMessage() (*sqs.ReceiveMessageOutput, error) {
	in := &sqs.ReceiveMessageInput{}

	return c.conn.ReceiveMessage(in)
}

func (c *mockClient) TestPurgeQueue() (*sqs.PurgeQueueOutput, error) {
	in := &sqs.PurgeQueueInput{}

	return c.conn.PurgeQueue(in)
}

func TestMockSQSClient_ListQueues(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockSQSClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &sqs.ListQueuesOutput{}
			testClient.On("ListQueues", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestListQueues() })

				actual, _ := tester.TestListQueues()

				assert.Equal(t, resp, actual)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("ListQueues", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestListQueues() })

				_, err := tester.TestListQueues()

				assert.Equal(t, resp, err)
			})
		})
	})
}

func TestMockSQSClient_GetQueueUrl(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockSQSClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &sqs.GetQueueUrlOutput{}
			testClient.On("GetQueueUrl", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestGetQueueUrl() })

				actual, _ := tester.TestGetQueueUrl()

				assert.Equal(t, resp, actual)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("GetQueueUrl", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestGetQueueUrl() })

				_, err := tester.TestGetQueueUrl()

				assert.Equal(t, resp, err)
			})
		})
	})
}

func TestMockSQSClient_GetQueueAttributes(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockSQSClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &sqs.GetQueueAttributesOutput{}
			testClient.On("GetQueueAttributes", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestGetQueueAttributes() })

				actual, _ := tester.TestGetQueueAttributes()

				assert.Equal(t, resp, actual)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("GetQueueAttributes", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestGetQueueAttributes() })

				_, err := tester.TestGetQueueAttributes()

				assert.Equal(t, resp, err)
			})
		})
	})
}

func TestMockSQSClient_SendMessage(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockSQSClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &sqs.SendMessageOutput{}
			testClient.On("SendMessage", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestSendMessage() })

				actual, _ := tester.TestSendMessage()

				assert.Equal(t, resp, actual)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("SendMessage", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestSendMessage() })

				_, err := tester.TestSendMessage()

				assert.Equal(t, resp, err)
			})
		})
	})
}

func TestMockSQSClient_ReceiveMessage(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockSQSClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &sqs.ReceiveMessageOutput{}
			testClient.On("ReceiveMessage", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestReceiveMessage() })

				actual, _ := tester.TestReceiveMessage()

				assert.Equal(t, resp, actual)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("ReceiveMessage", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestReceiveMessage() })

				_, err := tester.TestReceiveMessage()

				assert.Equal(t, resp, err)
			})
		})
	})
}

func TestMockSQSClient_PurgeQueue(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockSQSClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &sqs.PurgeQueueOutput{}
			testClient.On("PurgeQueue", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestPurgeQueue() })

				actual, _ := tester.TestPurgeQueue()

				assert.Equal(t, resp, actual)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("PurgeQueue", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestPurgeQueue() })

				_, err := tester.TestPurgeQueue()

				assert.Equal(t, resp, err)
			})
		})
	})
}

func TestMockSQSClient_DeleteMessage(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockSQSClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &sqs.DeleteMessageOutput{}
			testClient.On("DeleteMessage", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestDeleteMessage() })

				actual, _ := tester.TestDeleteMessage()

				assert.Equal(t, resp, actual)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("DeleteMessage", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestDeleteMessage() })

				_, err := tester.TestDeleteMessage()

				assert.Equal(t, resp, err)
			})
		})
	})
}

func TestMockSQSClient_DeleteMessageBatch(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockSQSClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &sqs.DeleteMessageBatchOutput{}
			testClient.On("DeleteMessageBatch", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestDeleteMessageBatch() })

				actual, _ := tester.TestDeleteMessageBatch()

				assert.Equal(t, resp, actual)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("DeleteMessageBatch", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestDeleteMessageBatch() })

				_, err := tester.TestDeleteMessageBatch()

				assert.Equal(t, resp, err)
			})
		})
	})
}
