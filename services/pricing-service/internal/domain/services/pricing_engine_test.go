package services

import (
	"context"
	"testing"

	"github.com/Abdex1/FamGo-platform/services/pricing-service/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

// MockPricingRepository is a mock for testing
type MockPricingRepository struct{}

func (m *MockPricingRepository) CreatePricingRule(ctx context.Context, rule *entities.PricingRule) error {
	return nil
}

func (m *MockPricingRepository) GetPricingRule(ctx context.Context, ruleID string) (*entities.PricingRule, error) {
	return nil, nil
}

func (m *MockPricingRepository) GetActiveRuleForRideType(ctx context.Context, rideType, city string) (*entities.PricingRule, error) {
	return nil, nil
}

func (m *MockPricingRepository) SaveFareCalculation(ctx context.Context, fare *entities.FareCalculation) error {
	return nil
}

func (m *MockPricingRepository) SaveSurgeMultiplier(ctx context.Context, surge *entities.SurgeHistory) error {
	return nil
}

func (m *MockPricingRepository) ValidateDiscountCode(ctx context.Context, code string) (*entities.DiscountCode, error) {
	return nil, nil
}

func (m *MockPricingRepository) DecrementDiscountCodeUsage(ctx context.Context, codeID string) error {
	return nil
}

func (m *MockPricingRepository) GetFareCalculationHistory(ctx context.Context, rideID string) ([]entities.FareCalculation, error) {
	return []entities.FareCalculation{}, nil
}

func (m *MockPricingRepository) GetSurgeHistory(ctx context.Context, city string, hours int) ([]entities.SurgeHistory, error) {
	return []entities.SurgeHistory{}, nil
}

func (m *MockPricingRepository) GetAverageFareByRideType(ctx context.Context, city string, days int) (map[string]float64, error) {
	return map[string]float64{}, nil
}

// Tests

func TestBaseFareCalculation(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	fare := engine.CalculateFare(
		"ride_123", "ECONOMY",
		0, 0, // No distance or time
		9.0, 38.7, 9.0, 38.7,
		false,
		0, 100,
		nil,
	)

	assert.NotNil(t, fare)
	// Base fare for ECONOMY is 20 ETB
	assert.True(t, fare.FinalFare >= 20.0)
}

func TestDistanceFareCalculation(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	// 5 km distance
	fare := engine.CalculateFare(
		"ride_123", "ECONOMY",
		5000, 0,
		9.0, 38.7, 9.0, 38.7,
		false,
		0, 100,
		nil,
	)

	assert.NotNil(t, fare)
	// Distance fare: 5 km * 10 ETB/km = 50 ETB
	assert.Greater(t, fare.DistanceFare, 0.0)
	assert.True(t, fare.FinalFare > 50.0)
}

func TestTimeFareCalculation(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	// 15 minutes (900 seconds)
	fare := engine.CalculateFare(
		"ride_123", "ECONOMY",
		0, 900,
		9.0, 38.7, 9.0, 38.7,
		false,
		0, 100,
		nil,
	)

	assert.NotNil(t, fare)
	// Time fare: 15 min * 0.33 ETB/min = ~4.95 ETB
	assert.Greater(t, fare.TimeFare, 0.0)
}

func TestCompleteRideCalculation(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	// 5 km, 15 minutes
	fare := engine.CalculateFare(
		"ride_123", "ECONOMY",
		5000, 900,
		9.0, 38.7, 9.0, 38.7,
		false,
		50, 20, // 50 active rides, 20 available drivers
		nil,
	)

	assert.NotNil(t, fare)
	assert.Greater(t, fare.FinalFare, 0.0)
	assert.Equal(t, fare.RideID, "ride_123")
	assert.Equal(t, fare.RideType, "ECONOMY")
}

func TestSurgeMultiplierBasic(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	// Test surge calculation
	surge := engine.GetSurgeMultiplierAtLocation(9.0, 38.7, 100, 10)

	assert.NotNil(t, surge)
	assert.Greater(t, surge, 1.0)
	assert.LessOrEqual(t, surge, 5.0)
}

func TestSurgeMultiplierHighDemand(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	// High demand (100 rides, 5 drivers)
	surge := engine.GetSurgeMultiplierAtLocation(9.0, 38.7, 100, 5)

	assert.NotNil(t, surge)
	assert.Greater(t, surge, 1.0)
}

func TestSurgeMultiplierCapped(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	// Very high demand (1000 rides, 1 driver)
	surge := engine.GetSurgeMultiplierAtLocation(9.0, 38.7, 1000, 1)

	assert.NotNil(t, surge)
	assert.LessOrEqual(t, surge, 5.0) // Should be capped
}

func TestPoolDiscount(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	fareRegular := engine.CalculateFare(
		"ride_123", "ECONOMY",
		5000, 900,
		9.0, 38.7, 9.0, 38.7,
		false, 0, 100, nil,
	)

	farePool := engine.CalculateFare(
		"ride_124", "ECONOMY",
		5000, 900,
		9.0, 38.7, 9.0, 38.7,
		true, 0, 100, nil,
	)

	// Pool should be cheaper due to 25% discount
	assert.Less(t, farePool.FinalFare, fareRegular.FinalFare)
}

func TestEstimateFare(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	estimate := engine.EstimateFare("ECONOMY", 5000, 50, 20, false)

	assert.NotNil(t, estimate)
	assert.Greater(t, estimate["final_fare"].(float64), 0.0)
}

func TestMinimumFareEnforced(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	// Very short ride (100m, 1 minute)
	fare := engine.CalculateFare(
		"ride_123", "ECONOMY",
		100, 60,
		9.0, 38.7, 9.0, 38.7,
		false, 0, 100, nil,
	)

	assert.NotNil(t, fare)
	// Should be at least minimum fare (15 ETB for ECONOMY)
	assert.GreaterOrEqual(t, fare.FinalFare, 15.0)
}

func TestDifferentRideTypes(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	fares := map[string]float64{}
	rideTypes := []string{"ECONOMY", "COMFORT", "BUSINESS", "POOL"}

	for _, rideType := range rideTypes {
		fare := engine.CalculateFare(
			"ride_123", rideType,
			5000, 900,
			9.0, 38.7, 9.0, 38.7,
			false, 0, 100, nil,
		)
		fares[rideType] = fare.FinalFare
	}

	// Different ride types should have different fares
	assert.Less(t, fares["ECONOMY"], fares["COMFORT"])
	assert.Less(t, fares["COMFORT"], fares["BUSINESS"])
	assert.Less(t, fares["POOL"], fares["ECONOMY"])
}

func TestTaxCalculation(t *testing.T) {
	engine := NewPricingEngine(&MockPricingRepository{})

	fare := engine.CalculateFare(
		"ride_123", "ECONOMY",
		5000, 900,
		9.0, 38.7, 9.0, 38.7,
		false, 0, 100, nil,
	)

	assert.NotNil(t, fare)
	// Tax should be 2% of subtotal + surge
	expectedTaxBase := fare.SubtotalBeforeSurge + fare.SurgeAmount
	expectedTax := expectedTaxBase * 0.02
	assert.InDelta(t, fare.Taxes, expectedTax, 0.1)
}

// Benchmark tests

func BenchmarkFareCalculation(b *testing.B) {
	engine := NewPricingEngine(&MockPricingRepository{})

	for i := 0; i < b.N; i++ {
		engine.CalculateFare(
			"ride_123", "ECONOMY",
			5000, 900,
			9.0, 38.7, 9.0, 38.7,
			false, 50, 20, nil,
		)
	}
}

func BenchmarkSurgeCalculation(b *testing.B) {
	engine := NewPricingEngine(&MockPricingRepository{})

	for i := 0; i < b.N; i++ {
		engine.GetSurgeMultiplierAtLocation(9.0, 38.7, 100, 10)
	}
}

func BenchmarkEstimate(b *testing.B) {
	engine := NewPricingEngine(&MockPricingRepository{})

	for i := 0; i < b.N; i++ {
		engine.EstimateFare("ECONOMY", 5000, 50, 20, false)
	}
}
