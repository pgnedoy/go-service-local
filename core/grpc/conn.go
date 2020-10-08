package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

const (
	grpcKeepAliveSeconds      = 30
	grpcMaxConcurrentStreams  = 1000000
	metadataHeaderServiceName = "x-service-cluster"
)

var (
	validHeaders = []string{
		"x-request-id",
	}
)

func toGrpcEndpoint(endpoint string) string {
	if endpoint == "" {
		log.Fatal("invalid: endpoint")
	}

	u, err := url.Parse(endpoint)

	if err != nil {
		log.Fatal("error parsing endpoint", err)
	}

	return fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
}

func NewServiceMeshConnection(endpoint, serviceName string, interceptors ...grpc.UnaryClientInterceptor) (*grpc.ClientConn, error) {
	if strings.HasPrefix(strings.ToLower(endpoint), "http") {
		endpoint = toGrpcEndpoint(endpoint)
	}

	dialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepAliveSeconds * time.Second,
			PermitWithoutStream: true, // send pings even without active streams
		}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			append([]grpc.UnaryClientInterceptor{
				NewClientServiceMeshInterceptor(serviceName),
			}, interceptors...)...,
		)),
		WithStreamingServiceNameHeader(serviceName),
	}

	return grpc.Dial(
		endpoint,
		dialOptions...,
	)
}

func WithServiceNameHeader(serviceName string) grpc.DialOption {
	return grpc.WithUnaryInterceptor(NewClientServiceMeshInterceptor(serviceName))
}

func WithStreamingServiceNameHeader(serviceName string) grpc.DialOption {
	return grpc.WithStreamInterceptor(NewClientStreamingServiceMeshInterceptor(serviceName))
}

// NewOutgoingContext reads the incoming metadata and appends appropriate headers to a new outgoing context
func NewOutgoingContext(in context.Context) context.Context {
	inMeta, _ := metadata.FromIncomingContext(in)

	outMeta := metadata.MD{}

	for _, h := range validHeaders {
		v := inMeta.Get(h)

		if len(v) > 0 {
			outMeta.Append(h, v...)
		}
	}

	return metadata.NewOutgoingContext(context.Background(), outMeta)
}

type ServerConfig struct {
	Interceptors        []grpc.UnaryServerInterceptor
	Options             []grpc.ServerOption
	ErrorResponsePrefix string
	HealthCheck         healthpb.HealthServer
	Environment         string
	ServiceName         string
}

// ServerOptions currently there are no custom options, but leaving this to use as variadic options in the future
type ServerOptions struct{}

type ServerOption func(*ServerOptions)

// func newServerOptions(options ...ServerOption) *ServerOptions {
// 	configuration := &ServerOptions{}

// 	for _, option := range options {
// 		option(configuration)
// 	}

// 	return configuration
// }

// NewGrpcServer creates a new gRPC server with preset options: MaxConcurrentStreams, KeepaliveParams
func NewServer(cfg *ServerConfig, options ...ServerOption) (*grpc.Server, error) {
	if cfg == nil {
		return nil, errors.New("invalid: cfg")
	}

	if cfg.ServiceName == "" {
		return nil, errors.New("invalid: ServiceName")
	}

	if cfg.Environment == "" {
		return nil, errors.New("invalid: Environment")
	}

	// placeholder for when there are custom options
	// serverOptions := newServerOptions(options...)

	grpcServerOpts := append(
		cfg.Options,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{Time: grpcKeepAliveSeconds * time.Second}),
	)

	interceptorChain := grpc_middleware.WithUnaryServerChain(grpc_middleware.ChainUnaryServer(cfg.Interceptors...))

	grpcServerOpts = append(grpcServerOpts, interceptorChain)

	s := grpc.NewServer(grpcServerOpts...)

	// stubbed health check to satisfy the GRPC health check API
	hc := cfg.HealthCheck

	if hc == nil {
		hc = &HealthCheckServer{}
	}

	healthpb.RegisterHealthServer(s, hc)

	return s, nil
}

func RunServer(ctx context.Context, srv *grpc.Server, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	runCtx, cancel := context.WithCancel(ctx)

	go watchSignals(runCtx, cancel)

	go func() {
		if err = srv.Serve(lis); err != nil {
			log.Fatal("failed to serve", err)
		}
	}()

	<-runCtx.Done()

	log.Println(ctx, "!!! shutting down grpc server")

	srv.GracefulStop()

	return nil
}

func watchSignals(ctx context.Context, fn context.CancelFunc) {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	sig := <-ch

	log.Println(ctx, "!!! received signal, shutting down", sig)

	fn()
}
