package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnrichContextGeneratesIDs(t *testing.T) {
	ctx := EnrichContext(context.Background(), "", "", "")
	traceID, correlationID, requestID := IDsFromContext(ctx)
	require.NotEmpty(t, traceID)
	require.Equal(t, traceID, correlationID)
	require.NotEmpty(t, requestID)
}

func TestEnrichContextPreservesProvidedIDs(t *testing.T) {
	ctx := EnrichContext(context.Background(), "trace-1", "corr-1", "req-1")
	traceID, correlationID, requestID := IDsFromContext(ctx)
	require.Equal(t, "trace-1", traceID)
	require.Equal(t, "corr-1", correlationID)
	require.Equal(t, "req-1", requestID)
}
