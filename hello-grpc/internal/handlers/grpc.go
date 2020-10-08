package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/benbjohnson/clock"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	hellogrpcv1 "github.com/pgnedoy/protos/gen/go/hellogrpc/v2"
)

type GrpcServer struct {
	clock clock.Clock
}

type GrpcServerConfig struct {
	Clock clock.Clock
}

func NewGrpcServer(cfg *GrpcServerConfig) (*GrpcServer, error) {
	if cfg == nil {
		return nil, errors.New("invalid: cfg")
	}

	cfgClock := cfg.Clock

	if cfgClock == nil {
		cfgClock = clock.New()
	}

	return &GrpcServer{
		clock: cfgClock,
	}, nil
}

func (s *GrpcServer) Hello(ctx context.Context, req *hellogrpcv1.HelloRequest) (*hellogrpcv1.HelloResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid: name")
	}

	greetingTime, err := ptypes.TimestampProto(s.clock.Now())

	// this should never happen, and if it does, let's return an internal server error
	if err != nil {
		log.Println("error with ptypes.TimestampProto", err)

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &hellogrpcv1.HelloResponse{
		Greeting:  fmt.Sprintf("Hello, %s!", req.Name),
		GreetTime: greetingTime,
	}, nil
}
