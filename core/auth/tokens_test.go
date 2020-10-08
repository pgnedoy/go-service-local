package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildTokenGenerator() (*TokenGenerator, error) {
	secret, err := GenerateHMACSecret(32)

	if err != nil {
		return nil, err
	}

	return NewTokenGenerator(&TokenGeneratorConfig{Secret: secret})
}

func TestNewTokenGenerator(t *testing.T) {
	t.Run("returns error with nil cfg", func(t *testing.T) {
		_, err := NewTokenGenerator(nil)

		assert.Error(t, err)
	})

	t.Run("returns error with empty secret", func(t *testing.T) {
		_, err := NewTokenGenerator(&TokenGeneratorConfig{})

		assert.Error(t, err)
	})

	t.Run("returns valid TokenGenerator", func(t *testing.T) {
		secret, err := GenerateHMACSecret(32)

		if err != nil {
			t.Error(err)
		}

		resp, err := NewTokenGenerator(&TokenGeneratorConfig{Secret: secret})

		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestTokenGenerator_GenerateToken(t *testing.T) {
	g, err := buildTokenGenerator()

	if err != nil {
		t.Error(err)
	}

	t.Run("returns a valid token with signature", func(t *testing.T) {
		token, sig, err := g.GenerateToken()

		if err != nil {
			t.Error(err)
		}

		assert.True(t, len(token) > 0)

		// signature is "<token>.<signature>"
		embeddedSig := strings.Split(token, ".")[1]

		assert.Equal(t, sig, embeddedSig)
	})
}

func TestTokenGenerator_ValidateToken(t *testing.T) {
	g, err := buildTokenGenerator()

	if err != nil {
		t.Error(err)
	}

	t.Run("returns an error with an invalid token", func(t *testing.T) {
		err := g.ValidateToken("invalid-token")

		assert.Error(t, err)
	})

	t.Run("returns nil with valid token", func(t *testing.T) {
		token, _, err := g.GenerateToken()

		if err != nil {
			t.Error(err)
		}

		err = g.ValidateToken(token)

		assert.Nil(t, err)
	})
}
