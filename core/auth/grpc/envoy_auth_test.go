package grpc

import (
	"context"
	"testing"

	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	rpc_code "google.golang.org/genproto/googleapis/rpc/code"
	rpc_status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func newMockAuthContext() context.Context {
	ctx := context.Background()

	md := metadata.Pairs(
		"authorization",
		"Bearer OFl-lgZRJ_JEXFgihRCXXDsn1Ye2SmvDPUo0oztYq_8.HSf8_J8qGeOy6via6TlWfkcJKwJcB7lTOBE010wK0Ag",
	)

	ctx = metadata.NewIncomingContext(ctx, md)

	return ctx
}

type mockRequest struct{}
type mockResponse struct{}

func mockHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return &mockResponse{}, nil
}

func newMockHandler() grpc.UnaryHandler {
	return mockHandler
}

func newRequestInfo(rpcName string) *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{
		FullMethod: rpcName,
	}
}

type MockEnvoyAuthClient struct {
	mock.Mock
	v2.AuthorizationClient
}

func (m *MockEnvoyAuthClient) Check(
	ctx context.Context,
	in *v2.CheckRequest,
	opts ...grpc.CallOption) (*v2.CheckResponse, error) {
	args := m.Called(ctx, in)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*v2.CheckResponse), nil
}

func newMockEnvoyAuthResponseOK(userID string) *v2.CheckResponse {
	headers := []*envoy_api_v2_core.HeaderValueOption{
		{
			Append: &wrappers.BoolValue{
				Value: false,
			},
			Header: &envoy_api_v2_core.HeaderValue{
				Key:   "x-internal-user-id",
				Value: userID,
			},
		},
	}

	return &v2.CheckResponse{
		Status: &rpc_status.Status{
			Code: int32(rpc_code.Code_OK),
		},
		HttpResponse: &v2.CheckResponse_OkResponse{
			OkResponse: &v2.OkHttpResponse{
				Headers: headers,
			},
		},
	}
}

// func newMockEnvoyAuthResponseDenied() *v2.CheckResponse {
//	return &v2.CheckResponse{
//		Status: &rpc_status.Status{
//			Code: int32(rpc_code.Code_UNAUTHENTICATED),
//		},
//		HttpResponse: &v2.CheckResponse_DeniedResponse{
//			DeniedResponse: &v2.DeniedHttpResponse{},
//		},
//	}
// }

func TestNewEnvoyAuthorizer(t *testing.T) {
	t.Run("when not given an oauth client returns an error", func(t *testing.T) {
		interceptor, err := NewEnvoyAuthorizer(nil)

		assert.Error(t, err)
		assert.Nil(t, interceptor)
	})

	t.Run("when given an oauth client returns a client", func(t *testing.T) {
		interceptor, err := NewEnvoyAuthorizer(&MockEnvoyAuthClient{})

		assert.Nil(t, err)
		assert.NotNil(t, interceptor)
	})
}

func TestEnvoyAuthorizer_UnaryServerInterceptor(t *testing.T) {
	t.Run("when the auth client returns an error", func(t *testing.T) {
		authClient := &MockEnvoyAuthClient{}
		authClient.On("Check", mock.Anything, mock.Anything).Once().Return(nil, assert.AnError)

		interceptor, _ := NewEnvoyAuthorizer(authClient)

		_, err := interceptor(newMockAuthContext(), &mockRequest{}, newRequestInfo("some.service/SomeRpc"), newMockHandler())

		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t)
	})

	t.Run("when the auth client returns an error but the rpc is whitelisted", func(t *testing.T) {
		authClient := &MockEnvoyAuthClient{}
		authClient.On("Check", mock.Anything, mock.Anything).Once().Return(nil, assert.AnError)

		handler := newMockHandler()

		interceptor, _ := NewEnvoyAuthorizer(authClient, WhiteListHandler("SomeRpc"))

		_, err := interceptor(newMockAuthContext(), &mockRequest{}, newRequestInfo("some.service/SomeRpc"), handler)

		assert.Nil(t, err)
		mock.AssertExpectationsForObjects(t)
	})

	t.Run("when the auth client returns a user", func(t *testing.T) {
		expUserID := "1234"

		authClient := &MockEnvoyAuthClient{}
		authClient.On("Check", mock.Anything, mock.Anything).Once().Return(newMockEnvoyAuthResponseOK(expUserID), nil)

		interceptor, _ := NewEnvoyAuthorizer(authClient)

		_, err := interceptor(newMockAuthContext(), &mockRequest{}, newRequestInfo("some.service/SomeRpc"), newMockHandler())

		assert.Nil(t, err)
		mock.AssertExpectationsForObjects(t)
	})

	t.Run("when the auth client returns a user but the rpc is white listed", func(t *testing.T) {
		authClient := &MockEnvoyAuthClient{}
		authClient.On("Check", mock.Anything, mock.Anything).Once().Return(newMockEnvoyAuthResponseOK("12345"), nil)

		handler := newMockHandler()

		interceptor, _ := NewEnvoyAuthorizer(authClient, WhiteListHandler("SomeRpc"))

		_, err := interceptor(newMockAuthContext(), &mockRequest{}, newRequestInfo("some.service/SomeRpc"), handler)

		assert.Nil(t, err)
		mock.AssertExpectationsForObjects(t)
	})

	t.Run("passes through grpc health checks", func(t *testing.T) {
		authClient := &MockEnvoyAuthClient{}

		handler := newMockHandler()

		interceptor, _ := NewEnvoyAuthorizer(authClient, WhiteListHandler("SomeRpc"))

		_, err := interceptor(
			context.Background(),
			&grpc_health_v1.HealthCheckRequest{},
			newRequestInfo("/grpc.health.v1.Health/Check"),
			handler,
		)

		assert.Nil(t, err)
		mock.AssertExpectationsForObjects(t)
	})
}
