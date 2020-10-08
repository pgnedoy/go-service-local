package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMock(t *testing.T) {
	t.Run("returns a valid object", func(t *testing.T) {
		m := NewMock()

		assert.NotNil(t, m)
	})
}

func TestMockClock_Now(t *testing.T) {
	t.Run("returns the mocked time", func(t *testing.T) {
		m := NewMock()

		exp := time.Now()

		m.On("Now").Once().Return(exp)

		actual := m.Now()

		assert.EqualValues(t, exp, actual)
		m.AssertExpectations(t)
	})
}

func TestMockClock_Sleep(t *testing.T) {
	t.Run("registers the call to Sleep", func(t *testing.T) {
		m := NewMock()

		exp := 1 * time.Millisecond

		m.On("Sleep", exp).Once()

		m.Sleep(exp)

		m.AssertExpectations(t)
	})
}

func TestMockClock_LoadLocation(t *testing.T) {
	t.Run("returns a mocked time.Location", func(t *testing.T) {
		m := NewMock()

		name := "Asia/Hong_Kong"

		exp, _ := New().LoadLocation(name)

		m.On("LoadLocation", name).Once().Return(exp, nil)

		actual, err := m.LoadLocation(name)

		assert.EqualValues(t, exp, actual)
		assert.Nil(t, err)
		m.AssertExpectations(t)
	})
}
