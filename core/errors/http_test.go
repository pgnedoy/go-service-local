package errors

import (
	"context"
	"testing"

	"errors"
	"net/http"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FakeResponse struct {
	t       *testing.T
	headers http.Header
	body    []byte
	status  int
}

func NewResp(t *testing.T) *FakeResponse {
	return &FakeResponse{
		t:       t,
		headers: make(http.Header),
	}
}

func (r *FakeResponse) Header() http.Header {
	return r.headers
}

func (r *FakeResponse) Write(body []byte) (int, error) {
	r.body = body
	return len(body), nil
}

func (r *FakeResponse) WriteHeader(st int) {
	r.status = st
}

func (r *FakeResponse) Assert(st int) {
	if r.status != st {
		r.t.Errorf("expected status %+v to equal %+v", r.status, st)
	}
}

func TestToHttp(t *testing.T) {
	t.Run("when given an invalid context error", func(t *testing.T) {
		err := NewInvalidContextError(context.Background())

		t.Run("returns a 412 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(412)
		})
	})

	t.Run("when given an invalid parameter error", func(t *testing.T) {
		err := NewInvalidParameterError("id", "id")

		t.Run("returns a 422 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(422)
		})
	})

	t.Run("when given a missing parameter error", func(t *testing.T) {
		err := NewMissingParameterError("id")

		t.Run("returns a 422 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(422)
		})
	})

	t.Run("when given an invalid entity error", func(t *testing.T) {
		err := NewInvalidEntityError("id")

		t.Run("returns a 422 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(422)
		})
	})

	t.Run("when given a resource exhausted error", func(t *testing.T) {
		err := NewResourceExhaustedError()

		t.Run("returns a 429 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(429)
		})
	})

	t.Run("when given a not implemented error", func(t *testing.T) {
		err := NewNotImplementedError()

		t.Run("returns a 501 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(501)
		})
	})

	t.Run("when given an unauthenticated error", func(t *testing.T) {
		err := NewUnauthenticatedError()

		t.Run("returns a 403 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(401)
		})
	})

	t.Run("when given an unauthorized error", func(t *testing.T) {
		err := NewUnauthorizedError()

		t.Run("returns a 403 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(403)
		})
	})

	t.Run("when given an invalid entity error", func(t *testing.T) {
		err := NewNotFoundError()

		t.Run("returns a 404 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(404)
		})
	})

	t.Run("when given an invalid entity error", func(t *testing.T) {
		err := errors.New("error")

		t.Run("returns a 400 response code", func(t *testing.T) {
			resp := NewResp(t)
			ToHttp(resp, err)

			resp.Assert(400)
		})
	})
}

func TestGrpcToHTTP(t *testing.T) {
	t.Run("when not given a grpc error", func(t *testing.T) {
		err := errors.New("not grpc")

		t.Run("returns an error", func(t *testing.T) {
			resp := NewResp(t)
			marshalErr := GrpcToHTTP(resp, err)

			assert.Error(t, marshalErr)
		})
	})

	t.Run("when given a grpc error", func(t *testing.T) {
		err := status.Error(codes.InvalidArgument, "something")

		t.Run("returns an error", func(t *testing.T) {
			resp := NewResp(t)
			resErr := GrpcToHTTP(resp, err)

			assert.Nil(t, resErr)
			resp.Assert(400)
		})
	})

	t.Run("when given a grpc error for context", func(t *testing.T) {
		err := status.Error(codes.FailedPrecondition, "something")

		t.Run("returns an error", func(t *testing.T) {
			resp := NewResp(t)
			resErr := GrpcToHTTP(resp, err)

			assert.Nil(t, resErr)
			resp.Assert(400)
		})
	})
}
