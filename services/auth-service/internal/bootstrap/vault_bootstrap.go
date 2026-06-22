/*
# PHASE 3 — VAULT SECRET INJECTION

# =========================================================

# STEP 8 — VAULT SECRET LOADER

internal/bootstrap/vault_bootstrap.go


# STEP 19 — LOAD SECRETS FROM VAULT
*/

package bootstrap

import (
	"context"

	vaultclient "github.com/Abdex1/FamGo-platform/packages/vault-sdk/client"
	vaultkv "github.com/Abdex1/FamGo-platform/packages/vault-sdk/kv"
)

type AuthSecrets struct {
	JWTSecret        string
	JWTRefreshSecret string
}

func LoadAuthSecrets(
	ctx context.Context,
	address string,
	token string,
) (*AuthSecrets, error) {

	client, err := vaultclient.New(
		vaultclient.Config{
			Address: address,
			Token: token,
		},
	)

	if err != nil {
		return nil, err
	}

	store := vaultkv.New(client)

	data, err := store.GetSecret(
		ctx,
		"auth-service",
	)

	if err != nil {
		return nil, err
	}

	return &AuthSecrets{
		JWTSecret: data["JWT_SECRET"].(string),
		JWTRefreshSecret: data["JWT_REFRESH_SECRET"].(string),
	}, nil
}
