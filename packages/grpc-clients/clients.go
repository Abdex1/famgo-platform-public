// packages/grpc-clients/clients.go
// gRPC Service Clients for Cross-Service Communication

package grpcclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"go.uber.org/zap"
)

// ClientConfig holds configuration for gRPC client connections
type ClientConfig struct {
	Host    string
	Port    int
	Timeout time.Duration
}

// GRPCClientPool manages multiple service connections
type GRPCClientPool struct {
	pricing  *grpc.ClientConn
	gps      *grpc.ClientConn
	dispatch *grpc.ClientConn
	logger   *zap.Logger
}

// NewGRPCClientPool creates a new client pool
func NewGRPCClientPool(
	pricingCfg ClientConfig,
	gpsCfg ClientConfig,
	dispatchCfg ClientConfig,
	logger *zap.Logger,
) (*GRPCClientPool, error) {
	pool := &GRPCClientPool{logger: logger}

	var err error

	// Connect to Pricing Service
	pool.pricing, err = grpc.Dial(
		fmt.Sprintf("%s:%d", pricingCfg.Host, pricingCfg.Port),
		grpc.WithInsecure(), // In production, use TLS
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to pricing service: %w", err)
	}

	// Connect to GPS Service
	pool.gps, err = grpc.Dial(
		fmt.Sprintf("%s:%d", gpsCfg.Host, gpsCfg.Port),
		grpc.WithInsecure(), // In production, use TLS
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)),
	)
	if err != nil {
		pool.pricing.Close()
		return nil, fmt.Errorf("failed to connect to gps service: %w", err)
	}

	// Connect to Dispatch Service
	pool.dispatch, err = grpc.Dial(
		fmt.Sprintf("%s:%d", dispatchCfg.Host, dispatchCfg.Port),
		grpc.WithInsecure(), // In production, use TLS
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)),
	)
	if err != nil {
		pool.pricing.Close()
		pool.gps.Close()
		return nil, fmt.Errorf("failed to connect to dispatch service: %w", err)
	}

	logger.Info("gRPC client pool initialized",
		zap.String("pricing", fmt.Sprintf("%s:%d", pricingCfg.Host, pricingCfg.Port)),
		zap.String("gps", fmt.Sprintf("%s:%d", gpsCfg.Host, gpsCfg.Port)),
		zap.String("dispatch", fmt.Sprintf("%s:%d", dispatchCfg.Host, dispatchCfg.Port)))

	return pool, nil
}

// GetPricingClient returns the pricing service client connection
func (p *GRPCClientPool) GetPricingClient() *grpc.ClientConn {
	return p.pricing
}

// GetGPSClient returns the GPS service client connection
func (p *GRPCClientPool) GetGPSClient() *grpc.ClientConn {
	return p.gps
}

// GetDispatchClient returns the dispatch service client connection
func (p *GRPCClientPool) GetDispatchClient() *grpc.ClientConn {
	return p.dispatch
}

// Close closes all client connections
func (p *GRPCClientPool) Close() error {
	var errs []error

	if p.pricing != nil {
		if err := p.pricing.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if p.gps != nil {
		if err := p.gps.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if p.dispatch != nil {
		if err := p.dispatch.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	p.logger.Info("gRPC client pool closed")
	return nil
}

// CallWithRetry executes a gRPC call with retry logic
func CallWithRetry(
	ctx context.Context,
	maxRetries int,
	backoff time.Duration,
	fn func(context.Context) error,
) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(backoff * time.Duration(attempt)):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on non-retryable errors
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument, codes.NotFound, codes.PermissionDenied:
				return err // Don't retry these
			}
		}
	}

	return lastErr
}

// WithTimeout adds a timeout to a context
func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}
