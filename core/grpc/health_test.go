package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthCheckServer_Check(t *testing.T) {
	s := HealthCheckServer{}

	resp, err := s.Check(context.Background(), &healthpb.HealthCheckRequest{})

	assert.Nil(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
}

type mockHealthServer struct {
	healthpb.Health_WatchServer
}

func TestHealthCheckServer_Watch(t *testing.T) {
	s := HealthCheckServer{}
	w := mockHealthServer{}

	err := s.Watch(&healthpb.HealthCheckRequest{}, w)

	assert.Error(t, err, "Watching is not supported")
}
