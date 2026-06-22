package entities

import "context"

// PricingRepository defines the interface for pricing persistence
type PricingRepository interface {
	CreatePricingRule(ctx context.Context, rule *PricingRule) error
	GetPricingRule(ctx context.Context, ruleID string) (*PricingRule, error)
	GetActiveRuleForRideType(ctx context.Context, rideType, city string) (*PricingRule, error)
	SaveFareCalculation(ctx context.Context, fare *FareCalculation) error
	SaveSurgeMultiplier(ctx context.Context, surge *SurgeHistory) error
	ValidateDiscountCode(ctx context.Context, code string) (*DiscountCode, error)
	DecrementDiscountCodeUsage(ctx context.Context, codeID string) error
	GetFareCalculationHistory(ctx context.Context, rideID string) ([]FareCalculation, error)
	GetSurgeHistory(ctx context.Context, city string, hours int) ([]SurgeHistory, error)
	GetAverageFareByRideType(ctx context.Context, city string, days int) (map[string]float64, error)
}
