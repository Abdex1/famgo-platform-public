package repositories

import (
	"context"
	"time"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
)

type DispatchStats struct {
	TotalMatches       int
	SuccessfulMatches  int
	FailedMatches      int
	SuccessRate        float64
	AverageTimeToMatch float64
}

type DispatchRepository interface {
	Create(ctx context.Context, request *entities.DispatchRequest) error
	Update(ctx context.Context, request *entities.DispatchRequest) error
	GetByID(ctx context.Context, id string) (*entities.DispatchRequest, error)
	GetByRideID(ctx context.Context, rideID string) (*entities.DispatchRequest, error)
	GetPendingRequests(ctx context.Context, limit int) ([]*entities.DispatchRequest, error)
	GetDispatchStats(ctx context.Context, startDate, endDate time.Time) (*DispatchStats, error)
}

type MatchingSessionRepository interface {
	Create(ctx context.Context, session *entities.MatchingSession) error
	Update(ctx context.Context, session *entities.MatchingSession) error
	GetByDispatchRequestID(ctx context.Context, dispatchRequestID string) (*entities.MatchingSession, error)
}

type MatchResultRepository interface {
	Create(ctx context.Context, result *entities.MatchResult) error
	GetByRideID(ctx context.Context, rideID string) (*entities.MatchResult, error)
}
