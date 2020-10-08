package grpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	rpc_code "google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/9count/go-services/core/auth"
	"github.com/9count/go-services/core/errors"
	"github.com/9count/go-services/core/log"
)

const (
	DefaultUserHeaderPrefix = "x-internal"
)

var (
	defaultAuthRequestTimeout = 10 * time.Second
	healthCheckMethods        = []string{
		"/grpc.health.v1.Health/Check",
	}
)

type InterceptorOptions struct {
	AuthRequestTimeout time.Duration
	WhiteListedRpcs    []string
	OptionalAuthRpcs   []string
}

type InterceptorOption func(*InterceptorOptions)

func newInteceptorOpts(interceptorOpts ...InterceptorOption) *InterceptorOptions {
	opts := &InterceptorOptions{
		AuthRequestTimeout: defaultAuthRequestTimeout,
	}

	for _, option := range interceptorOpts {
		option(opts)
	}

	return opts
}

// RPCs given here as options will attempt to auth and attach a user to the context, but won't
// fail the call if the auth fails
func OptionalAuthHandler(handlers ...string) InterceptorOption {
	return func(opts *InterceptorOptions) {
		opts.OptionalAuthRpcs = append(opts.OptionalAuthRpcs, handlers...)
	}
}

// RPCs given here as options will completely avoid authentication or calls to the auth service
func WhiteListHandler(handlers ...string) InterceptorOption {
	return func(opts *InterceptorOptions) {
		opts.WhiteListedRpcs = append(opts.WhiteListedRpcs, handlers...)
	}
}

func AuthRequestTimeout(d time.Duration) InterceptorOption {
	return func(opts *InterceptorOptions) {
		opts.AuthRequestTimeout = d
	}
}

type EnvoyAuthorizer struct {
	authClient v2.AuthorizationClient
	options    *InterceptorOptions
}

func (c *EnvoyAuthorizer) UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	for _, hc := range healthCheckMethods {
		if info.FullMethod == hc {
			return handler(ctx, req)
		}
	}

	newCtx := c.addUserToContext(ctx)

	unparsedMessageName := strings.Split(info.FullMethod, "/")
	parsedMessageName := unparsedMessageName[len(unparsedMessageName)-1]

	for _, rpcName := range c.options.WhiteListedRpcs {
		if parsedMessageName == rpcName {
			return handler(newCtx, req)
		}
	}

	err := c.authorize(newCtx)

	if err != nil {
		return nil, errors.ToGrpc(err)
	}

	return handler(newCtx, req)
}

func (c *EnvoyAuthorizer) authorize(ctx context.Context) error {
	user := auth.GetUserFromContext(ctx)

	if user == nil {
		return errors.NewUnauthenticatedError()
	}

	return nil
}

func (c *EnvoyAuthorizer) addUserToContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return ctx
	}

	authHeader := md.Get("authorization")

	if len(authHeader) == 0 {
		return ctx
	}

	authReq := &v2.CheckRequest{
		Attributes: &v2.AttributeContext{
			Request: &v2.AttributeContext_Request{
				Http: &v2.AttributeContext_HttpRequest{
					Headers: map[string]string{
						"authorization": authHeader[0],
					},
				},
			},
		},
	}

	clientCtx, clientCancel := context.WithTimeout(ctx, c.options.AuthRequestTimeout)
	defer clientCancel()

	authResp, err := c.authClient.Check(clientCtx, authReq)

	if err != nil {
		// TODO: replace with honeycomb event
		log.Warn(ctx, "invalid auth")
		return ctx
	}

	if authResp.Status != nil && authResp.Status.Code == int32(rpc_code.Code_OK) {
		httpResponse := authResp.HttpResponse.(*v2.CheckResponse_OkResponse)

		for _, header := range httpResponse.OkResponse.Headers {
			if header.Header.Key == fmt.Sprintf("%s-user-id", DefaultUserHeaderPrefix) {
				return context.WithValue(ctx, "user", &auth.UserData{Id: header.Header.Value})
			}
		}
	}

	return ctx
}

func NewEnvoyAuthorizer(authClient v2.AuthorizationClient, opts ...InterceptorOption) (grpc.UnaryServerInterceptor, error) {
	if authClient == nil {
		return nil, errors.NewMissingParameterError("authClient")
	}

	cfg := newInteceptorOpts(opts...)

	authorizer := &EnvoyAuthorizer{
		authClient: authClient,
		options:    cfg,
	}

	return authorizer.UnaryServerInterceptor, nil
}
