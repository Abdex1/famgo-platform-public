// services/payment-service/interfaces/grpc/payment_handler.go
package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Abdex1/FamGo-platform/services/payment-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/payment-service/proto/payment"
)

type PaymentHandler struct {
	payment.UnimplementedPaymentServiceServer
	useCases *usecases.PaymentUseCases
}

func NewPaymentHandler(useCases *usecases.PaymentUseCases) *PaymentHandler {
	return &PaymentHandler{useCases: useCases}
}

func (h *PaymentHandler) InitiatePayment(ctx context.Context, req *payment.InitiatePaymentRequest) (*payment.InitiatePaymentResponse, error) {
	if req == nil || req.RideId == "" {
		return nil, status.Error(codes.InvalidArgument, "ride ID required")
	}

	input := &usecases.InitiatePaymentInput{
		RideID:   req.RideId,
		RiderID:  req.RiderId,
		DriverID: req.DriverId,
		Amount:   req.Amount,
		Method:   req.Method,
	}

	output, err := h.useCases.InitiatePayment(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to initiate payment: %v", err))
	}

	return &payment.InitiatePaymentResponse{
		PaymentId: output.PaymentID,
		Status:    output.Status,
		Amount:    output.Amount,
	}, nil
}

func (h *PaymentHandler) CompletePayment(ctx context.Context, req *payment.CompletePaymentRequest) (*emptypb.Empty, error) {
	if req == nil || req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment ID required")
	}

	err := h.useCases.CompletePayment(ctx, &usecases.CompletePaymentInput{
		PaymentID:     req.PaymentId,
		TransactionID: req.TransactionId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to complete payment: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error) {
	if req == nil || req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment ID required")
	}

	output, err := h.useCases.GetPayment(ctx, &usecases.GetPaymentInput{PaymentID: req.PaymentId})
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("payment not found: %v", err))
	}

	return &payment.GetPaymentResponse{
		PaymentId: output.PaymentID,
		Status:    output.Status,
		Amount:    output.Amount,
	}, nil
}

func (h *PaymentHandler) RefundPayment(ctx context.Context, req *payment.RefundPaymentRequest) (*emptypb.Empty, error) {
	if req == nil || req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment ID required")
	}

	err := h.useCases.RefundPayment(ctx, &usecases.RefundPaymentInput{
		PaymentID: req.PaymentId,
		Amount:    req.Amount,
		Reason:    req.Reason,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to refund payment: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (h *PaymentHandler) HandleWebhook(ctx context.Context, req *payment.HandleWebhookRequest) (*emptypb.Empty, error) {
	if req == nil || req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment ID required")
	}

	err := h.useCases.HandleWebhook(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to handle webhook: %v", err))
	}

	return &emptypb.Empty{}, nil
}
