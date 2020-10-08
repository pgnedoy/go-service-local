package auth

import (
	"errors"

	"github.com/ory/fosite/token/hmac"
)

const (
	tokenEntropy = 32
)

type TokenGenerator struct {
	strategy hmac.HMACStrategy
}

type TokenGeneratorConfig struct {
	Secret string
}

func NewTokenGenerator(cfg *TokenGeneratorConfig) (*TokenGenerator, error) {
	if cfg == nil {
		return nil, errors.New("cfg must be defined")
	}

	if cfg.Secret == "" {
		return nil, errors.New("secret cannot be empty")
	}

	return &TokenGenerator{
		strategy: hmac.HMACStrategy{
			GlobalSecret: []byte(cfg.Secret),
			TokenEntropy: tokenEntropy,
		},
	}, nil
}

func (g *TokenGenerator) GenerateToken() (token, signature string, err error) {
	return g.strategy.Generate()
}

func (g *TokenGenerator) ValidateToken(t string) error {
	return g.strategy.Validate(t)
}
