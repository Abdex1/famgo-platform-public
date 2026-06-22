package queries

import "time"

type GetMatchesQuery struct {
	DispatchRequestID string
}

type GetMatchesResult struct {
	DispatchRequestID string
	Status            string
	MatchedDriverID   *string
	ProposedDrivers   []string
	ExpiresAt         time.Time
}

type GetStatsQuery struct {
	StartDate time.Time
	EndDate   time.Time
}

type GetStatsResult struct {
	TotalMatches       int
	SuccessfulMatches  int
	FailedMatches      int
	SuccessRate        float64
	AverageTimeToMatch float64
}
