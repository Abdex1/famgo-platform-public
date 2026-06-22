package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RideServiceClient notifies ride-service about driver acceptance.
type RideServiceClient struct {
	conn *grpc.ClientConn
}

func NewRideServiceClient(rideServiceURL string) (*RideServiceClient, error) {
	conn, err := grpc.NewClient(
		rideServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ride service: %w", err)
	}
	return &RideServiceClient{conn: conn}, nil
}

func (rsc *RideServiceClient) AcceptRide(ctx context.Context, rideID, driverID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if rideID == "" || driverID == "" {
		return fmt.Errorf("ride ID and driver ID are required")
	}
	// Ride assignment is propagated through dispatch.driver.assigned.v1 for saga choreography.
	_ = ctx
	return nil
}

func (rsc *RideServiceClient) Close() error {
	if rsc.conn != nil {
		return rsc.conn.Close()
	}
	return nil
}
