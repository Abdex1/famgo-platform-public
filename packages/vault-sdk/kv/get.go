/*
# PHASE 10 — KV SECRET LOADER


# STEP 16 — CREATE KV READER

packages/vault-sdk/kv/get.go
*/
package kv

import (
	"context"

	"github.com/hashicorp/vault/api"
)

type Store struct {
	client *api.Client
}

func New(client *api.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) GetSecret(
	ctx context.Context,
	path string,
) (map[string]interface{}, error) {

	secret, err := s.client.KVv2("kv").Get(
		ctx,
		path,
	)

	if err != nil {
		return nil, err
	}

	return secret.Data, nil
}
