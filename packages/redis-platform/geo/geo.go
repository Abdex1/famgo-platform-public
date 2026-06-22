/*
# PHASE 9 — REDIS GEO ENGINE

Extracted SAFELY from Uber clone concepts.

---

# STEP 12 — CREATE GEO ENGINE

packages/redis-platform/geo/geo.go
*/
package geo

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const DriversKey = "drivers:geo"

type Engine struct {
	rdb *redis.Client
}

func NewEngine(rdb *redis.Client) *Engine {
	return &Engine{
		rdb: rdb,
	}
}

func (g *Engine) SetDriverLocation(
	ctx context.Context,
	driverID string,
	lat float64,
	lng float64,
) error {

	return g.rdb.GeoAdd(
		ctx,
		DriversKey,
		&redis.GeoLocation{
			Name:      driverID,
			Latitude:  lat,
			Longitude: lng,
		},
	).Err()
}

func (g *Engine) NearbyDrivers(
	ctx context.Context,
	lat float64,
	lng float64,
	radiusKm float64,
) ([]redis.GeoLocation, error) {

	return g.rdb.GeoSearchLocation(
		ctx,
		DriversKey,
		&redis.GeoSearchLocationQuery{
			Latitude:  lat,
			Longitude: lng,
			Radius:    radiusKm,
			RadiusUnit: "km",
			Count: 50,
			Sort: "ASC",
		},
	).Result()
}
