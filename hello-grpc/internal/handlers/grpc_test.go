package handlers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	hellogrpcv1 "github.com/pgnedoy/protos/gen/go/hello-proto/v1"
)

func TestNewGrpcServer(t *testing.T) {
	t.Run("when not given a config object", func(t *testing.T) {
		_, err := NewGrpcServer(nil)

		t.Run("returns an error", func(t *testing.T) {
			assert.EqualValues(t, errors.New("invalid: cfg"), err)
		})
	})

	t.Run("when given a valid config object", func(t *testing.T) {
		server, err := NewGrpcServer(&GrpcServerConfig{})

		t.Run("returns a server", func(t *testing.T) {
			assert.NotNil(t, server)
			assert.Nil(t, err)
		})
	})
}

func setupGrpcServerFixture() (*GrpcServer, error) {
	return NewGrpcServer(&GrpcServerConfig{
		Clock: clock.NewMock(),
	})
}

func TestGrpcServer_Hello(t *testing.T) {
	t.Run("when request Name is empty", func(t *testing.T) {
		s, err := setupGrpcServerFixture()

		if err != nil {
			t.Error(err)
			return
		}

		t.Run("returns an error", func(t *testing.T) {
			_, reqErr := s.Hello(context.Background(), &hellogrpcv1.HelloRequest{})

			assert.EqualValues(t, status.Error(codes.InvalidArgument, "invalid: name"), reqErr)
		})
	})

	t.Run("when there is an error with generating the GreetTime", func(t *testing.T) {
		s, err := setupGrpcServerFixture()

		if err != nil {
			t.Error(err)
			return
		}

		clockMock := s.clock.(*clock.Mock)

		// invalid date should trigger the greetingTime error
		clockMock.Set(time.Time{}.Add(-1000000 * time.Hour))

		_, reqErr := s.Hello(context.Background(), &hellogrpcv1.HelloRequest{Name: "Foo"})

		t.Run("returns an error", func(t *testing.T) {
			assert.EqualValues(t, status.Error(codes.Unknown, "timestamp: seconds:-65735596800  before 0001-01-01"), reqErr)
		})
	})

	t.Run("when there is an error with generating the GreetTime", func(t *testing.T) {
		s, err := setupGrpcServerFixture()

		if err != nil {
			t.Error(err)
			return
		}

		clockMock := s.clock.(*clock.Mock)

		expName := "Foo"
		expGreetingTime := time.Now()

		greetingTime, err := ptypes.TimestampProto(expGreetingTime)

		if err != nil {
			t.Error(err)
			return
		}

		expResp := &hellogrpcv1.HelloResponse{
			Greeting:  fmt.Sprintf("Hello, %s!", expName),
			GreetTime: greetingTime,
		}

		clockMock.Set(expGreetingTime)

		resp, reqErr := s.Hello(context.Background(), &hellogrpcv1.HelloRequest{Name: expName})

		t.Run("returns the correct response", func(t *testing.T) {
			assert.Nil(t, reqErr)
			assert.EqualValues(t, expResp, resp)
		})
	})
}
