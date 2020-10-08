package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewClientServiceMeshInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		newCtx := metadata.AppendToOutgoingContext(ctx, metadataHeaderServiceName, serviceName)

		// Calls the invoker to execute RPC
		err := invoker(newCtx, method, req, reply, cc, opts...)

		return err
	}
}

func NewClientStreamingServiceMeshInterceptor(serviceName string) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		newCtx := metadata.AppendToOutgoingContext(ctx, metadataHeaderServiceName, serviceName)

		return streamer(newCtx, desc, cc, method, opts...)
	}
}
