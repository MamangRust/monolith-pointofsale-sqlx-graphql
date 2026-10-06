package server

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewGRPCClient creates a gRPC client connection that propagates the active
// OpenTelemetry trace context (traceparent) to the receiving service via
// otelgrpc, so a single request can be traced across service boundaries.
func NewGRPCClient(address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}

	return grpc.NewClient(address, append(defaultOpts, opts...)...)
}
