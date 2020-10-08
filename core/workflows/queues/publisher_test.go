package queues

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"

	commonSQS "github.com/9count/go-services/core/sqs"
)

func setupSQPublisherFixture() (*SQSPublisher, error) {
	return NewSQSPublisher(&SQSPublisherConfig{
		SqsClient: &commonSQS.MockSQSClient{},
		QueueURL:  "https://sqs.us-east-1.amazonaws.com/1234567890/foo-queue-ABCD",
	})
}

func TestNewSQSPublisher(t *testing.T) {
	queueURL := "https://sqs.us-east-1.amazonaws.com/1234567890/foo-queue-ABCD"

	t.Run("returns an error with nil config", func(t *testing.T) {
		_, err := NewSQSPublisher(nil)

		assert.Error(t, err)
		assert.EqualValues(t, errors.New("invalid: cfg"), err)
	})

	t.Run("returns an error with no QueueURL", func(t *testing.T) {
		_, err := NewSQSPublisher(&SQSPublisherConfig{})

		assert.Error(t, err)
		assert.EqualValues(t, errors.New("invalid: QueueURL"), err)
	})

	t.Run("returns a valid object", func(t *testing.T) {
		w, err := NewSQSPublisher(&SQSPublisherConfig{
			AwsEndpoint: "http://localhost:8083",
			AwsRegion:   "us-west-1",
			QueueURL:    queueURL,
		})

		assert.Nil(t, err)
		assert.NotNil(t, w)
	})
}

func TestSQSPublisher_SendMessage(t *testing.T) {
	t.Run("returns error with nil msg", func(t *testing.T) {
		pub, err := setupSQPublisherFixture()

		if err != nil {
			t.Error(err)
			return
		}

		err = pub.SendMessage(context.Background(), nil)

		assert.Error(t, err)
		assert.EqualValues(t, errors.New("invalid: msg"), err)
	})

	t.Run("returns error with invalid msg type", func(t *testing.T) {
		pub, err := setupSQPublisherFixture()

		if err != nil {
			t.Error(err)
			return
		}

		invalidMsg := make(chan bool)

		err = pub.SendMessage(context.Background(), invalidMsg)

		assert.Error(t, err)
		assert.IsType(t, &json.UnsupportedTypeError{}, err)
	})

	t.Run("returns no error with a valid string message", func(t *testing.T) {
		pub, err := setupSQPublisherFixture()

		if err != nil {
			t.Error(err)
			return
		}

		input := "Hello, World!"

		sqsMock := pub.sqsClient.(*commonSQS.MockSQSClient)

		sqsMock.On(
			"SendMessage",
			&sqs.SendMessageInput{
				QueueUrl:    aws.String(pub.queueURL),
				MessageBody: aws.String(input),
			},
		).Once().Return(&sqs.SendMessageOutput{}, nil)

		err = pub.SendMessage(context.Background(), input)

		assert.Nil(t, err)

		sqsMock.AssertExpectations(t)
	})

	t.Run("returns no error with a valid string message", func(t *testing.T) {
		pub, err := setupSQPublisherFixture()

		if err != nil {
			t.Error(err)
			return
		}

		input := struct {
			Name string `json:"name"`
		}{Name: "Foo"}

		sqsMock := pub.sqsClient.(*commonSQS.MockSQSClient)

		sqsMock.On(
			"SendMessage",
			&sqs.SendMessageInput{
				QueueUrl:    aws.String(pub.queueURL),
				MessageBody: aws.String(`{"name":"Foo"}`),
			},
		).Once().Return(&sqs.SendMessageOutput{}, nil)

		err = pub.SendMessage(context.Background(), input)

		assert.Nil(t, err)

		sqsMock.AssertExpectations(t)
	})
}
