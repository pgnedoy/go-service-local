package dynamo

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalStreamImage(t *testing.T) {
	t.Run("when given a list of dynamo attributes", func(t *testing.T) {
		attrs := map[string]events.DynamoDBAttributeValue{
			"id": events.NewNumberAttribute("someid"),
		}

		type newThing struct {
			ID string `dynamodbav:"id"`
		}

		t.Run("it returns a new struct", func(t *testing.T) {
			newStruct := &newThing{}

			err := UnmarshalStreamImage(attrs, newStruct)

			assert.Nil(t, err)
			assert.Equal(t, newStruct.ID, "someid")
		})
	})
}

type testMarshalFoo struct {
	Foo string            `json:"foo,omitempty"`
	Baz int               `json:"baz,omitempty"`
	Bar []*testMarshalFoo `json:"bar,omitempty"`
}

func TestAttributesFromJSONObject(t *testing.T) {
	t.Run("foo", func(t *testing.T) {
		expected, err := dynamodbattribute.MarshalMap(&testMarshalFoo{
			Foo: "abc",
			Baz: 42,
			Bar: []*testMarshalFoo{
				{Foo: "xyz"},
			},
		})

		if err != nil {
			t.Error(err)
			return
		}

		actual, err := AttributesFromJSONObject(`{"foo":"abc", "baz":42, "bar": [{"foo":"xyz"}]}`)

		if err != nil {
			t.Error(err)
			return
		}

		assert.EqualValues(t, expected, actual)
	})
}
