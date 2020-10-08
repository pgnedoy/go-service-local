package errors

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func FromDynamo(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case dynamodb.ErrCodeConditionalCheckFailedException:
			return NewMissingParameterError("", HumanizedMessage(err.Error()))
		case dynamodb.ErrCodeResourceNotFoundException:
			return NewNotFoundError()
		case dynamodb.ErrCodeItemCollectionSizeLimitExceededException, dynamodb.ErrCodeProvisionedThroughputExceededException:
			return NewResourceExhaustedError()
		case dynamodb.ErrCodeInternalServerError:
			return NewNetworkError()
		}
	}

	return NewUnknownError(err)
}
