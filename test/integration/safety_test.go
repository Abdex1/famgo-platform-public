// test/integration/safety_test.go
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSOSInitiation tests initiating an SOS incident
func TestSOSInitiation(t *testing.T) {
	client := setupSafetyClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.InitiateSOS(ctx, &InitiateSOSRequest{
		RideId:       "ride_sos_001",
		UserId:       "user_sos_001",
		UserType:     "rider",
		Latitude:     9.0365,
		Longitude:    38.7469,
		IncidentType: "threat",
		Description:  "feeling_unsafe_with_driver",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.IncidentId)
	assert.Equal(t, "active", resp.Status)
	assert.Equal(t, "emergency_contact", resp.EscalationLevel)
}

// TestSOSEscalation tests escalating an SOS incident
func TestSOSEscalation(t *testing.T) {
	client := setupSafetyClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initiate SOS
	initResp, err := client.InitiateSOS(ctx, &InitiateSOSRequest{
		RideId:       "ride_sos_002",
		UserId:       "user_sos_002",
		UserType:     "rider",
		Latitude:     9.0365,
		Longitude:    38.7469,
		IncidentType: "medical_emergency",
		Description:  "need_ambulance",
	})
	require.NoError(t, err)

	// Escalate to ambulance
	_, err = client.EscalateIncident(ctx, &EscalateIncidentRequest{
		IncidentId:      initResp.IncidentId,
		EscalationLevel: "ambulance",
	})
	require.NoError(t, err)

	// Verify escalation
	getResp, err := client.GetIncident(ctx, &GetIncidentRequest{IncidentId: initResp.IncidentId})
	require.NoError(t, err)
	assert.Equal(t, "escalated", getResp.Status)
	assert.Equal(t, "ambulance", getResp.EscalationLevel)
}

// TestSOSResolution tests resolving an SOS incident
func TestSOSResolution(t *testing.T) {
	client := setupSafetyClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initiate SOS
	initResp, err := client.InitiateSOS(ctx, &InitiateSOSRequest{
		RideId:       "ride_sos_003",
		UserId:       "user_sos_003",
		UserType:     "rider",
		Latitude:     9.0365,
		Longitude:    38.7469,
		IncidentType: "accident",
		Description:  "minor_collision",
	})
	require.NoError(t, err)

	// Resolve
	_, err = client.ResolveIncident(ctx, &ResolveIncidentRequest{
		IncidentId:    initResp.IncidentId,
		ResolvedBy:    "support_agent_safety_001",
		ResolutionNotes: "all_parties_safe_accident_minor",
	})
	require.NoError(t, err)

	// Verify resolution
	getResp, err := client.GetIncident(ctx, &GetIncidentRequest{IncidentId: initResp.IncidentId})
	require.NoError(t, err)
	assert.Equal(t, "resolved", getResp.Status)
}

func setupSafetyClient(t *testing.T) interface{} {
	// Implementation would connect to actual safety service
	return nil
}
