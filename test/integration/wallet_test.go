// test/integration/wallet_test.go
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWalletCreation tests creating a new wallet
func TestWalletCreation(t *testing.T) {
	client := setupWalletClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CreateWallet(ctx, &CreateWalletRequest{
		UserId:   "user_test_001",
		UserType: "rider",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.WalletId)
	assert.Equal(t, "user_test_001", resp.UserId)
	assert.Equal(t, 0.0, resp.Balance)
}

// TestWalletTransaction tests recording transactions
func TestWalletTransaction(t *testing.T) {
	client := setupWalletClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create wallet
	walletResp, err := client.CreateWallet(ctx, &CreateWalletRequest{
		UserId:   "user_test_002",
		UserType: "driver",
	})
	require.NoError(t, err)

	// Record deposit
	txnResp, err := client.RecordTransaction(ctx, &RecordTransactionRequest{
		WalletId:       walletResp.WalletId,
		TransactionType: "deposit",
		Amount:         500.00,
		Description:    "initial_deposit",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, txnResp.TransactionId)
	assert.Equal(t, "completed", txnResp.Status)

	// Get balance
	balResp, err := client.GetBalance(ctx, &GetBalanceRequest{WalletId: walletResp.WalletId})
	require.NoError(t, err)
	assert.Equal(t, 500.0, balResp.Balance)
}

// TestInsufficientFunds tests withdrawal with insufficient funds
func TestInsufficientFunds(t *testing.T) {
	client := setupWalletClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create wallet with small balance
	walletResp, err := client.CreateWallet(ctx, &CreateWalletRequest{
		UserId:   "user_test_003",
		UserType: "rider",
	})
	require.NoError(t, err)

	// Record small deposit
	_, err = client.RecordTransaction(ctx, &RecordTransactionRequest{
		WalletId:        walletResp.WalletId,
		TransactionType: "deposit",
		Amount:          50.00,
		Description:     "small_deposit",
	})
	require.NoError(t, err)

	// Try to withdraw more than balance
	_, err = client.RecordTransaction(ctx, &RecordTransactionRequest{
		WalletId:        walletResp.WalletId,
		TransactionType: "withdrawal",
		Amount:          100.00, // More than balance
		Description:     "failed_withdrawal",
	})
	assert.Error(t, err)
}

func setupWalletClient(t *testing.T) interface{} {
	// Implementation would connect to actual wallet service
	return nil
}
