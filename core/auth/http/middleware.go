package http

import (
	"net/http"
	"strings"

	"github.com/9count/go-services/core/errors"
	"github.com/9count/go-services/core/log"
)

type RequestAuthenticator struct {
	logErrors bool
}

type RequestAuthenticatorConfig struct {
	LogErrors bool
}

func NewRequestAuthenticator(cfg *RequestAuthenticatorConfig) (*RequestAuthenticator, error) {
	if cfg == nil {
		return nil, errors.NewMissingParameterError("cfg")
	}

	return &RequestAuthenticator{
		logErrors: cfg.LogErrors,
	}, nil
}

func tokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], "bearer") {
		return "", errors.NewInvalidEntityError("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func (a *RequestAuthenticator) failRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func (a *RequestAuthenticator) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := tokenFromHeader(r)

		if err != nil {
			if a.logErrors {
				log.Warn(r.Context(), "error in AuthMiddleware", log.WithError(err))
			}

			a.failRequest(w)
			return
		}

		if token == "" {
			a.failRequest(w)
			return
		}

		// TODO: for now hardcode token 1234-5678
		if token == "1234-5678" {
			next.ServeHTTP(w, r)
		} else {
			a.failRequest(w)
		}
	})
}
