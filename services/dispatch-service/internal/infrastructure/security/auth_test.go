package security

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenValidatorHasPermission(t *testing.T) {
	validator := NewTokenValidator("secret", "")
	require.True(t, validator.HasPermission(&Claims{Roles: []string{"admin"}}, "dispatch:write"))
	require.True(t, validator.HasPermission(&Claims{Scopes: []string{"dispatch:read"}}, "dispatch:read"))
	require.False(t, validator.HasPermission(&Claims{Scopes: []string{"dispatch:read"}}, "dispatch:write"))
}
