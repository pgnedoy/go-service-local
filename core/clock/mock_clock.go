package clock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockClock struct {
	Clock
	mock.Mock
}

func NewMock() *MockClock {
	return &MockClock{}
}

func (m *MockClock) Now() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}

func (m *MockClock) Sleep(d time.Duration) {
	_ = m.Called(d)
}

func (m *MockClock) LoadLocation(name string) (*time.Location, error) {
	args := m.Called(name)

	return args.Get(0).(*time.Location), args.Error(1)
}
