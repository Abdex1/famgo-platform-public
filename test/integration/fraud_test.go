// test/integration/fraud_test.go
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFraudCheckLowRisk tests a low-risk ride check
func TestFraudCheckLowRisk(t *testing.T) {
	client := setupFraudClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CheckRide(ctx, &CheckRideRequest{
		RideId:    "ride_low_risk_001",
		UserId:    "user_normal_001",
		UserType:  "rider",
		Amount:    150.00,
		Latitude:  9.0365,  // Addis Ababa
		Longitude: 38.7469,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.CheckId)
	assert.Equal(t, "low", resp.RiskLevel)
	assert.Equal(t, "allow", resp.Action)
}

// TestFraudCheckHighRisk tests a high-risk ride check
func TestFraudCheckHighRisk(t *testing.T) {
	client := setupFraudClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Very high amount + location anomaly
	resp, err := client.CheckRide(ctx, &CheckRideRequest{
		RideId:    "ride_high_risk_001",
		UserId:    "user_suspicious_001",
		UserType:  "rider",
		Amount:    50000.00, // Unusually high
		Latitude:  -33.8688, // Sydney
		Longitude: 151.2093,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.CheckId)
	assert.Equal(t, "high", resp.RiskLevel)
	assert.Equal(t, "block", resp.Action)
	assert.Len(t, resp.Flags, 1) // Should have payment anomaly flag
}

// TestFraudCheckReview tests reviewing a fraud check
func TestFraudCheckReview(t *testing.T) {
	client := setupFraudClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a check
	checkResp, err := client.CheckRide(ctx, &CheckRideRequest{
		RideId:    "ride_review_001",
		UserId:    "user_review_001",
		UserType:  "rider",
		Amount:    5000.00, // Medium-high amount
		Latitude:  9.0365,
		Longitude: 38.7469,
	})
	require.NoError(t, err)

	// Review it
	_, err = client.ReviewCheck(ctx, &ReviewCheckRequest{
		CheckId:    checkResp.CheckId,
		ReviewedBy: "support_agent_001",
		Reason:     "user_verified_legitimate_business",
	})
	require.NoError(t, err)

	// Get and verify
	getResp, err := client.GetFraudCheck(ctx, &GetFraudCheckRequest{CheckId: checkResp.CheckId})
	require.NoError(t, err)
	assert.True(t, getResp.IsReview)
}

func setupFraudClient(t *testing.T) interface{} {
	// Implementation would connect to actual fraud service
	return nil
}
