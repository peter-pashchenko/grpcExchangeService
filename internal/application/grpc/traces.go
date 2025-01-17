package grpc

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func TracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	tracer := otel.Tracer("grpcExchangeRate")

	ctx, span := tracer.Start(ctx, info.FullMethod,
		trace.WithAttributes(
			attribute.String("rpc.system", "grpcs"),
			attribute.String("rpc.service", info.FullMethod),
		),

	)
	defer span.End()

	return handler(ctx, req)
}
