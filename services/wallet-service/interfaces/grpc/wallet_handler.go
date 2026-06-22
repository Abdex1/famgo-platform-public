// services/wallet-service/interfaces/grpc/wallet_handler.go
package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Abdex1/FamGo-platform/services/wallet-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/wallet-service/proto/wallet"
)

type WalletHandler struct {
	wallet.UnimplementedWalletServiceServer
	useCases *usecases.WalletUseCases
}

func NewWalletHandler(useCases *usecases.WalletUseCases) *WalletHandler {
	return &WalletHandler{useCases: useCases}
}

func (h *WalletHandler) CreateWallet(ctx context.Context, req *wallet.CreateWalletRequest) (*wallet.WalletResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID required")
	}

	output, err := h.useCases.CreateWallet(ctx, &usecases.CreateWalletInput{
		UserID:   req.UserId,
		UserType: req.UserType,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create wallet: %v", err))
	}

	return &wallet.WalletResponse{
		WalletId:  output.WalletID,
		UserId:    output.UserID,
		Balance:   output.Balance,
		Currency:  output.Currency,
		IsActive:  true,
		CreatedAt: timestamppb.Now(),
	}, nil
}

func (h *WalletHandler) GetWallet(ctx context.Context, req *wallet.GetWalletRequest) (*wallet.WalletResponse, error) {
	if req == nil || req.WalletId == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet ID required")
	}

	output, err := h.useCases.GetWallet(ctx, &usecases.GetWalletInput{WalletID: req.WalletId})
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("wallet not found: %v", err))
	}

	return &wallet.WalletResponse{
		WalletId:  output.WalletID,
		UserId:    output.UserID,
		Balance:   output.Balance,
		Currency:  output.Currency,
		IsActive:  true,
		CreatedAt: timestamppb.Now(),
	}, nil
}

func (h *WalletHandler) GetBalance(ctx context.Context, req *wallet.GetBalanceRequest) (*wallet.GetBalanceResponse, error) {
	if req == nil || req.WalletId == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet ID required")
	}

	output, err := h.useCases.GetWallet(ctx, &usecases.GetWalletInput{WalletID: req.WalletId})
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("wallet not found: %v", err))
	}

	return &wallet.GetBalanceResponse{
		Balance:       output.Balance,
		TotalDeposited: 0,
		TotalEarned:   0,
	}, nil
}

func (h *WalletHandler) Transfer(ctx context.Context, req *wallet.TransferRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "transfer request required")
	}

	return &emptypb.Empty{}, nil
}

func (h *WalletHandler) RecordTransaction(ctx context.Context, req *wallet.RecordTransactionRequest) (*wallet.TransactionResponse, error) {
	if req == nil || req.WalletId == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet ID required")
	}

	output, err := h.useCases.RecordTransaction(ctx, &usecases.RecordTransactionInput{
		WalletID:    req.WalletId,
		TxType:      req.TransactionType,
		Amount:      req.Amount,
		Description: req.Description,
		CreatedBy:   "system",
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to record transaction: %v", err))
	}

	return &wallet.TransactionResponse{
		TransactionId: "txn_1234",
		Status:        "completed",
		Amount:        req.Amount,
		CreatedAt:     timestamppb.Now(),
	}, nil
}
