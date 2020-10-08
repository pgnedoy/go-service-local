package errors

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestToGrpc(t *testing.T) {
	t.Run("when given an invalid parameter error", func(t *testing.T) {
		err := NewInvalidParameterError("id", "lol")

		t.Run("returns an invalid argument error", func(t *testing.T) {
			grpcErr := ToGrpc(err)

			hasCode := strings.Contains(grpcErr.Error(), "code = InvalidArgument")

			assert.Equal(t, hasCode, true)
		})
	})

	t.Run("when given an invalid context error", func(t *testing.T) {
		err := NewInvalidContextError(context.Background())

		t.Run("returns an invalid context error", func(t *testing.T) {
			grpcErr := ToGrpc(err)

			hasCode := strings.Contains(grpcErr.Error(), "code = FailedPrecondition")

			assert.Equal(t, hasCode, true)
		})
	})

	t.Run("when given a not implemented error", func(t *testing.T) {
		err := NewNotImplementedError()

		t.Run("returns an unimplemented code error", func(t *testing.T) {
			grpcErr := ToGrpc(err)

			hasCode := strings.Contains(grpcErr.Error(), "code = Unimplemented")

			assert.Equal(t, hasCode, true)
		})
	})

	t.Run("when given an unauthenticated error", func(t *testing.T) {
		err := NewUnauthenticatedError()

		t.Run("returns an unauthenticated code error", func(t *testing.T) {
			grpcErr := ToGrpc(err)

			hasCode := strings.Contains(grpcErr.Error(), "code = Unauthenticated")

			assert.Equal(t, hasCode, true)
		})
	})

	t.Run("when given a not found error", func(t *testing.T) {
		err := NewNotFoundError()

		t.Run("returns an NotFound code error", func(t *testing.T) {
			grpcErr := ToGrpc(err)

			hasCode := strings.Contains(grpcErr.Error(), "code = NotFound")

			assert.Equal(t, hasCode, true)
		})
	})

	t.Run("when given an invalid entity error", func(t *testing.T) {
		err := NewInvalidEntityError("invalid")

		t.Run("returns an InvalidArgument code error", func(t *testing.T) {
			grpcErr := ToGrpc(err)

			hasCode := strings.Contains(grpcErr.Error(), "code = InvalidArgument")

			assert.Equal(t, hasCode, true)
		})
	})

	t.Run("when given a resource exhausted error", func(t *testing.T) {
		err := NewResourceExhaustedError()

		t.Run("returns an ResourceExhausted code error", func(t *testing.T) {
			grpcErr := ToGrpc(err)

			hasCode := strings.Contains(grpcErr.Error(), "code = ResourceExhausted")

			assert.Equal(t, hasCode, true)
		})
	})

	t.Run("when given a generic error", func(t *testing.T) {
		err := errors.New("generic error")

		t.Run("returns an Unavailable code error", func(t *testing.T) {
			grpcErr := ToGrpc(err)

			hasCode := strings.Contains(grpcErr.Error(), "code = Unavailable")

			assert.Equal(t, hasCode, true)
		})
	})
}
