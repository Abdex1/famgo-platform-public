/*
# STEP 12 — CREATE REDIS SESSION PLATFORM

services/auth-service/internal/infrastructure/redis/session_store.go


# FILE: session_store.go
*/
package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type SessionStore struct {
	client *goredis.Client
}

func NewSessionStore(addr string) *SessionStore {
	client := goredis.NewClient(
		&goredis.Options{
			Addr: addr,
		},
	)

	return &SessionStore{
		client: client,
	}
}

func (s *SessionStore) IsBlacklisted(
	ctx context.Context,
	tokenID string,
) bool {

	_, err := s.client.Get(
		ctx,
		"blacklist:"+tokenID,
	).Result()

	return err == nil
}
