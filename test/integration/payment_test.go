// test/integration/payment_test.go
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPaymentInitiation tests the payment initiation flow
func TestPaymentInitiation(t *testing.T) {
	client := setupPaymentClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.InitiatePayment(ctx, &InitiatePaymentRequest{
		RideId:   "ride_test_123",
		RiderId:  "rider_test_456",
		DriverId: "driver_test_789",
		Amount:   150.50,
		Method:   "chapa",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.PaymentId)
	assert.Equal(t, "initiated", resp.Status)
	assert.Equal(t, 150.50, resp.Amount)
}

// TestPaymentCompletion tests completing a payment
func TestPaymentCompletion(t *testing.T) {
	client := setupPaymentClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initiate
	initResp, err := client.InitiatePayment(ctx, &InitiatePaymentRequest{
		RideId:   "ride_test_124",
		RiderId:  "rider_test_457",
		DriverId: "driver_test_790",
		Amount:   200.00,
		Method:   "telebirr",
	})
	require.NoError(t, err)

	// Complete
	_, err = client.CompletePayment(ctx, &CompletePaymentRequest{
		PaymentId:     initResp.PaymentId,
		TransactionId: "txn_12345_ext",
	})
	require.NoError(t, err)

	// Get and verify
	getResp, err := client.GetPayment(ctx, &GetPaymentRequest{PaymentId: initResp.PaymentId})
	require.NoError(t, err)
	assert.Equal(t, "completed", getResp.Status)
}

// TestPaymentRefund tests refunding a completed payment
func TestPaymentRefund(t *testing.T) {
	client := setupPaymentClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initiate and complete
	initResp, err := client.InitiatePayment(ctx, &InitiatePaymentRequest{
		RideId:   "ride_test_125",
		RiderId:  "rider_test_458",
		DriverId: "driver_test_791",
		Amount:   175.75,
		Method:   "cbe_birr",
	})
	require.NoError(t, err)

	_, err = client.CompletePayment(ctx, &CompletePaymentRequest{
		PaymentId:     initResp.PaymentId,
		TransactionId: "txn_12346_ext",
	})
	require.NoError(t, err)

	// Refund
	_, err = client.RefundPayment(ctx, &RefundPaymentRequest{
		PaymentId: initResp.PaymentId,
		Amount:    175.75,
		Reason:    "customer_request",
	})
	require.NoError(t, err)

	// Verify
	getResp, err := client.GetPayment(ctx, &GetPaymentRequest{PaymentId: initResp.PaymentId})
	require.NoError(t, err)
	assert.Equal(t, "refunded", getResp.Status)
}

// TestInvalidPaymentAmount tests validation of payment amounts
func TestInvalidPaymentAmount(t *testing.T) {
	client := setupPaymentClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Below minimum
	_, err := client.InitiatePayment(ctx, &InitiatePaymentRequest{
		RideId:   "ride_test_126",
		RiderId:  "rider_test_459",
		DriverId: "driver_test_792",
		Amount:   5.00, // Below minimum
		Method:   "chapa",
	})
	assert.Error(t, err)

	// Above maximum
	_, err = client.InitiatePayment(ctx, &InitiatePaymentRequest{
		RideId:   "ride_test_127",
		RiderId:  "rider_test_460",
		DriverId: "driver_test_793",
		Amount:   200000.00, // Above maximum
		Method:   "chapa",
	})
	assert.Error(t, err)
}

func setupPaymentClient(t *testing.T) interface{} {
	// Implementation would connect to actual payment service
	return nil
}
