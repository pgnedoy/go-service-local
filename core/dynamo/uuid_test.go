package dynamo

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

type fooUUIDFixture struct {
	uuidGen UUIDGenerator
}

func (f *fooUUIDFixture) GetUUID() string {
	return f.uuidGen.New()
}

func TestMockUUID_New(t *testing.T) {
	m := &MockUUID{}

	foo := &fooUUIDFixture{uuidGen: m}

	expUUID := "1234-abcd"

	m.On("New").Return(expUUID)

	actual := foo.GetUUID()

	assert.Equal(t, expUUID, actual)
}
