package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRoute(rw http.ResponseWriter, req *http.Request) {
	panic("test entered test handler, this should not happen")
}

func TestNewRequestAuthenticator(t *testing.T) {
	t.Run("errors with no cfg", func(t *testing.T) {
		_, err := NewRequestAuthenticator(nil)

		assert.Error(t, err)
	})

	t.Run("initializes correctly", func(t *testing.T) {
		m, err := NewRequestAuthenticator(&RequestAuthenticatorConfig{
			LogErrors: true,
		})

		assert.Nil(t, err)
		assert.Equal(t, true, m.logErrors)
	})

	t.Run("returns error with empty token", func(t *testing.T) {
		m, err := NewRequestAuthenticator(&RequestAuthenticatorConfig{})

		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()

		m.AuthMiddleware(http.HandlerFunc(testRoute)).ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("returns error with invalid token format", func(t *testing.T) {
		m, err := NewRequestAuthenticator(&RequestAuthenticatorConfig{})

		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		req.Header.Set("Authorization", "not-valid")

		w := httptest.NewRecorder()

		m.AuthMiddleware(http.HandlerFunc(testRoute)).ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
