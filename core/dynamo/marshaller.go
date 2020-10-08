package dynamo

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func UnmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out interface{}) error {
	dynamoMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range attribute {
		var dbAttr dynamodb.AttributeValue

		bytes, err := v.MarshalJSON()
		if err != nil {
			return err
		}

		if marErr := json.Unmarshal(bytes, &dbAttr); marErr != nil {
			return marErr
		}

		dynamoMap[k] = &dbAttr
	}

	return dynamodbattribute.UnmarshalMap(dynamoMap, out)
}

func AttributesFromJSONObject(str string) (map[string]*dynamodb.AttributeValue, error) {
	ob := make(map[string]interface{})

	err := json.Unmarshal([]byte(str), &ob)

	if err != nil {
		return nil, err
	}

	return dynamodbattribute.MarshalMap(&ob)
}
