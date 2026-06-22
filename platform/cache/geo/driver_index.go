
// platform/cache/geo/driver_index.go

package geo

import (
    "context"

    goredis "github.com/redis/go-redis/v9"
)

const DriverGeoKey = "dispatch:drivers:geo"

type DriverLocation struct {
    DriverID string
    Lat      float64
    Lng      float64
}

type DriverGeoIndex struct {
    redis goredis.UniversalClient
}

func NewDriverGeoIndex(redis goredis.UniversalClient) *DriverGeoIndex {
    return &DriverGeoIndex{redis: redis}
}

func (d *DriverGeoIndex) Upsert(
    ctx context.Context,
    driverID string,
    lat float64,
    lng float64,
) error {
    return d.redis.GeoAdd(ctx, DriverGeoKey, &goredis.GeoLocation{
        Name:      driverID,
        Latitude:  lat,
        Longitude: lng,
    }).Err()
}

func (d *DriverGeoIndex) Nearby(
    ctx context.Context,
    lat float64,
    lng float64,
    radiusKM float64,
    limit int,
) ([]string, error) {
    return d.redis.GeoSearch(ctx, DriverGeoKey, &goredis.GeoSearchQuery{
        Longitude:  lng,
        Latitude:   lat,
        Radius:     radiusKM,
        RadiusUnit: "km",
        Count:      limit,
        Sort:       "ASC",
    }).Result()
}
