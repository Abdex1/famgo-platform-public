package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gpsv1 "github.com/Abdex1/FamGo-platform/services/dispatch-service/api/clients/gps/v1"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/ports"
)

// GPSDriverDiscovery implements DriverDiscovery via gps-service gRPC.
type GPSDriverDiscovery struct {
	client gpsv1.GPSServiceClient
	conn   *grpc.ClientConn
}

func NewGPSDriverDiscovery(gpsServiceURL string) (*GPSDriverDiscovery, error) {
	conn, err := grpc.NewClient(
		gpsServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GPS service: %w", err)
	}
	return &GPSDriverDiscovery{
		client: gpsv1.NewGPSServiceClient(conn),
		conn:   conn,
	}, nil
}

func (d *GPSDriverDiscovery) FindDriversWithinRadius(
	ctx context.Context,
	latitude, longitude, radiusKm float64,
	limit int,
) ([]ports.DriverCandidate, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 50
	}

	resp, err := d.client.FindNearbyDrivers(ctx, &gpsv1.FindNearbyDriversRequest{
		Latitude:    latitude,
		Longitude:   longitude,
		RadiusKm:    radiusKm,
		Limit:       int32(limit),
		OnlyOnline:  true,
		BaseSpeedKmh: 30,
	})
	if err != nil {
		return nil, fmt.Errorf("gps FindNearbyDrivers failed: %w", err)
	}

	candidates := make([]ports.DriverCandidate, 0, len(resp.GetDrivers()))
	for _, driver := range resp.GetDrivers() {
		candidates = append(candidates, ports.DriverCandidate{
			DriverID:       driver.GetDriverId(),
			Latitude:       driver.GetLatitude(),
			Longitude:      driver.GetLongitude(),
			DistanceKm:     driver.GetDistance(),
			ETAMinutes:     driver.GetEtaMinutes(),
			IsOnline:       driver.GetIsOnline(),
			AcceptanceRate: driver.GetAcceptanceRate(),
			Rating:         driver.GetRating(),
		})
	}
	return candidates, nil
}

func (d *GPSDriverDiscovery) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
