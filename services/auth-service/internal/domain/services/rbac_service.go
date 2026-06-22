/*
# PHASE 5 — RBAC POLICY ENGINE

# =========================================================

# STEP 14 — CREATE POLICY MATRIX


internal/domain/services/rbac_service.go
*/
package services

var RBACPolicies = map[string][]string{
	"super-admin": {
		"*",
	},

	"ops": {
		"ride.read",
		"ride.update",
		"driver.suspend",
	},

	"support": {
		"user.read",
		"ride.read",
	},

	"fraud-agent": {
		"user.read",
		"user.block",
	},
}
