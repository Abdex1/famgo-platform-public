package observability

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(ContextUnaryServerInterceptor(), LoggingUnaryServerInterceptor()),
	}
}

func ContextUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		ctx = ContextFromGRPC(ctx)
		traceID, correlationID, requestID := IDsFromContext(ctx)

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(
			attribute.String("trace_id", traceID),
			attribute.String("correlation_id", correlationID),
			attribute.String("request_id", requestID),
			attribute.String("rpc.method", info.FullMethod),
		)

		resp, err := handler(ctx, req)
		if err != nil {
			span.RecordError(err)
		}
		return resp, err
	}
}

func LoggingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil && status.Code(err) == codes.Internal {
			traceID, _, _ := IDsFromContext(ctx)
			_ = traceID
		}
		return resp, err
	}
}
