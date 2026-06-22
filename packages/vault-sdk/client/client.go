/*
# PHASE 9 — ENTERPRISE VAULT CLIENT

---

# STEP 15 — CREATE CLIENT

packages/vault-sdk/client/client.go
*/
package client

import (
	"github.com/hashicorp/vault/api"
)

type Config struct {
	Address string
	Token   string
}

func New(cfg Config) (*api.Client, error) {

	conf := api.DefaultConfig()

	conf.Address = cfg.Address

	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	client.SetToken(cfg.Token)

	return client, nil
}
