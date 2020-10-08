package errors

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestFromDynamo(t *testing.T) {
	t.Run("when given a conditional check fail", func(t *testing.T) {
		dynamoErr := awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "err", assert.AnError)

		t.Run("returns a missing parameter error", func(t *testing.T) {
			err := FromDynamo(dynamoErr)
			assert.IsType(t, err, &MissingParameterError{})
		})
	})

	t.Run("when given a resource not found error", func(t *testing.T) {
		dynamoErr := awserr.New(dynamodb.ErrCodeResourceNotFoundException, "err", assert.AnError)

		t.Run("returns a not found error", func(t *testing.T) {
			err := FromDynamo(dynamoErr)
			assert.IsType(t, err, &NotFoundError{})
		})
	})

	t.Run("when given a collection size limit exceeded error", func(t *testing.T) {
		dynamoErr := awserr.New(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, "err", assert.AnError)

		t.Run("returns a resource exhausted error", func(t *testing.T) {
			err := FromDynamo(dynamoErr)
			assert.IsType(t, err, &ResourceExhausted{})
		})
	})

	t.Run("when given a provisioned throughput exceeded error", func(t *testing.T) {
		dynamoErr := awserr.New(dynamodb.ErrCodeProvisionedThroughputExceededException, "err", assert.AnError)

		t.Run("returns aresource exhausted error", func(t *testing.T) {
			err := FromDynamo(dynamoErr)
			assert.IsType(t, err, &ResourceExhausted{})
		})
	})

	t.Run("when given an internal server error", func(t *testing.T) {
		dynamoErr := awserr.New(dynamodb.ErrCodeInternalServerError, "err", assert.AnError)

		t.Run("returns a network error", func(t *testing.T) {
			err := FromDynamo(dynamoErr)
			assert.IsType(t, err, &NetworkError{})
		})
	})

	t.Run("when given an internal server error", func(t *testing.T) {
		dynamoErr := awserr.New(dynamodb.ErrCodeInternalServerError, "err", assert.AnError)

		t.Run("returns a network error", func(t *testing.T) {
			err := FromDynamo(dynamoErr)
			assert.IsType(t, err, &NetworkError{})
		})
	})

	t.Run("when given an unknown error", func(t *testing.T) {
		t.Run("returns a network error", func(t *testing.T) {
			err := FromDynamo(assert.AnError)
			assert.IsType(t, err, &UnknownError{})
		})
	})
}
