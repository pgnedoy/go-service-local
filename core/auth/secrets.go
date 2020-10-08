package auth

import (
	"encoding/base64"
	"errors"

	"github.com/ory/fosite/token/hmac"
)

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

func GenerateHMACSecret(n int) (string, error) {
	if n < 1 {
		return "", errors.New("n must be greater than zero")
	}

	secretBytes, err := hmac.RandomBytes(n)

	if err != nil {
		return "", err
	}

	return b64.EncodeToString(secretBytes), nil
}
