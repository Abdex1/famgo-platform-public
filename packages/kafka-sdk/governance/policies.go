
// STEP 13 — CREATE ENTERPRISE GOVERNANCE RULES

//packages/kafka-sdk/governance/policies.go

// REQUIRED POLICIES

package governance

const (
	DefaultReplicationFactor = 3
	DefaultPartitions        = 12
	DefaultRetentionHours    = 168
	DLQRetentionHours        = 720
)
