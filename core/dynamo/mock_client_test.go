package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	conn dynamodbiface.DynamoDBAPI
}

func (c *mockClient) TestDescribeTable() (*dynamodb.DescribeTableOutput, error) {
	in := &dynamodb.DescribeTableInput{}

	return c.conn.DescribeTable(in)
}

func (c *mockClient) TestListTables() (*dynamodb.ListTablesOutput, error) {
	in := &dynamodb.ListTablesInput{}

	return c.conn.ListTables(in)
}

func (c *mockClient) TestPutItem() (*dynamodb.PutItemOutput, error) {
	in := &dynamodb.PutItemInput{}

	return c.conn.PutItem(in)
}

func (c *mockClient) TestGetItem() (*dynamodb.GetItemOutput, error) {
	in := &dynamodb.GetItemInput{}

	return c.conn.GetItem(in)
}

func (c *mockClient) TestQuery() (*dynamodb.QueryOutput, error) {
	in := &dynamodb.QueryInput{}

	return c.conn.Query(in)
}

func (c *mockClient) TestScan() (*dynamodb.ScanOutput, error) {
	in := &dynamodb.ScanInput{}

	return c.conn.Scan(in)
}

func (c *mockClient) TestDeleteItem() (*dynamodb.DeleteItemOutput, error) {
	in := &dynamodb.DeleteItemInput{}

	return c.conn.DeleteItem(in)
}

func (c *mockClient) TestUpdateItem() (*dynamodb.UpdateItemOutput, error) {
	in := &dynamodb.UpdateItemInput{}

	return c.conn.UpdateItem(in)
}

func newMockClient(conn dynamodbiface.DynamoDBAPI) *mockClient {
	return &mockClient{
		conn: conn,
	}
}

func TestMockDynamoDBClient_DescribeTable(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockDynamoDBClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &dynamodb.DescribeTableOutput{}
			testClient.On("DescribeTable", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestDescribeTable() })

				subject, _ := tester.TestDescribeTable()

				assert.Equal(t, resp, subject)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("DescribeTable", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestDescribeTable() })

				_, subject := tester.TestDescribeTable()

				assert.Equal(t, resp, subject)
			})
		})
	})
}

func TestMockDynamoDBClient_ListTables(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockDynamoDBClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &dynamodb.ListTablesOutput{}
			testClient.On("ListTables", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestListTables() })

				subject, _ := tester.TestListTables()

				assert.Equal(t, resp, subject)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("ListTables", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestListTables() })

				_, subject := tester.TestListTables()

				assert.Equal(t, resp, subject)
			})
		})
	})
}

func TestMockDynamoDBClient_GetItem(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockDynamoDBClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &dynamodb.GetItemOutput{}
			testClient.On("GetItem", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestGetItem() })

				subject, _ := tester.TestGetItem()

				assert.Equal(t, resp, subject)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("GetItem", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestGetItem() })

				_, subject := tester.TestGetItem()

				assert.Equal(t, resp, subject)
			})
		})
	})
}

func TestMockDynamoDBClient_Query(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockDynamoDBClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &dynamodb.QueryOutput{}
			testClient.On("Query", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestQuery() })

				subject, _ := tester.TestQuery()

				assert.Equal(t, resp, subject)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("Query", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestQuery() })

				_, subject := tester.TestQuery()

				assert.Equal(t, resp, subject)
			})
		})
	})
}

func TestMockDynamoDBClient_PutItem(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockDynamoDBClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &dynamodb.PutItemOutput{}
			testClient.On("PutItem", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestPutItem() })

				subject, _ := tester.TestPutItem()

				assert.Equal(t, resp, subject)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("PutItem", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestPutItem() })

				_, subject := tester.TestPutItem()

				assert.Equal(t, resp, subject)
			})
		})
	})
}

func TestMockDynamoDBClient_Scan(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockDynamoDBClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &dynamodb.ScanOutput{}
			testClient.On("Scan", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestScan() })

				subject, _ := tester.TestScan()

				assert.Equal(t, resp, subject)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("Scan", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestScan() })

				_, subject := tester.TestScan()

				assert.Equal(t, resp, subject)
			})
		})
	})
}

func TestMockDynamoDBClient_DeleteItem(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockDynamoDBClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &dynamodb.DeleteItemOutput{}
			testClient.On("DeleteItem", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestDeleteItem() })

				subject, _ := tester.TestDeleteItem()

				assert.Equal(t, resp, subject)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("DeleteItem", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestDeleteItem() })

				_, subject := tester.TestDeleteItem()

				assert.Equal(t, resp, subject)
			})
		})
	})
}

func TestMockDynamoDBClient_UpdateItem(t *testing.T) {
	t.Run("when given a mock client", func(t *testing.T) {
		testClient := new(MockDynamoDBClient)
		tester := newMockClient(testClient)

		t.Run("when not given an error", func(t *testing.T) {
			resp := &dynamodb.UpdateItemOutput{}
			testClient.On("UpdateItem", mock.Anything).Twice().Return(resp, nil)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestUpdateItem() })

				subject, _ := tester.TestUpdateItem()

				assert.Equal(t, resp, subject)
			})
		})

		t.Run("when not an error", func(t *testing.T) {
			resp := assert.AnError
			testClient.On("UpdateItem", mock.Anything).Twice().Return(nil, resp)

			t.Run("it doesn't panic", func(t *testing.T) {
				assert.NotPanics(t, func() { _, _ = tester.TestUpdateItem() })

				_, subject := tester.TestUpdateItem()

				assert.Equal(t, resp, subject)
			})
		})
	})
}
