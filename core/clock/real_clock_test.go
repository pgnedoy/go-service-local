package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("returns a valid object", func(t *testing.T) {
		c := New()

		assert.NotNil(t, c)
	})
}

func TestRealClock_Now(t *testing.T) {
	t.Run("returns a time.Time", func(t *testing.T) {
		now := New().Now()

		assert.NotNil(t, now)
	})
}

func TestRealClock_Sleep(t *testing.T) {
	t.Run("does a sleep", func(t *testing.T) {
		assert.NotPanics(t, func() { New().Sleep(1 * time.Millisecond) })
	})
}

func TestRealClock_LoadLocation(t *testing.T) {
	t.Run("returns a time.Location", func(t *testing.T) {
		loc, err := New().LoadLocation("Asia/Hong_Kong")

		assert.NotNil(t, loc)
		assert.Nil(t, err)
	})
}
