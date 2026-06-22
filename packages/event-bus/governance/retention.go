
// 17. CREATE EVENT RETENTION GOVERNANCE

// packages/event-bus/governance/retention.go

package governance

const (
	DefaultRetentionHours = 168
	DLQRetentionHours     = 720
	AuditRetentionHours   = 8760
)
