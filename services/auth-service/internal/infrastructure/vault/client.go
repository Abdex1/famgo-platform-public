/*
# PHASE 14 — VAULT SECRET LOADING

internal/infrastructure/vault/client.go

go get github.com/hashicorp/vault/api
*/
package vault

import vault "github.com/hashicorp/vault/api"

func NewClient(addr string, token string) (*vault.Client, error) {
    config := vault.DefaultConfig()
    config.Address = addr

    client, err := vault.NewClient(config)

    if err != nil {
        return nil, err
    }

    client.SetToken(token)

    return client, nil
}
