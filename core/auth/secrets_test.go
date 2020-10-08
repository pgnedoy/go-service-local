package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHMACSecret(t *testing.T) {
	t.Run("returns error if input < 1", func(t *testing.T) {
		_, err := GenerateHMACSecret(0)

		assert.Error(t, err)
	})

	t.Run("returns string with length greater than input length", func(t *testing.T) {
		n := 32

		res, err := GenerateHMACSecret(n)

		assert.Nil(t, err)
		assert.True(t, n <= len(res))
	})
}
