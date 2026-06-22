// STEP 4 — CREATE SESSION STORE
//packages/redis-platform/session/store.go

package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	rdb *redis.Client
}

func NewStore(rdb *redis.Client) *Store {
	return &Store{
		rdb: rdb,
	}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func userSessionsKey(userID string) string {
	return fmt.Sprintf("user_sessions:%s", userID)
}

func tokenFamilyKey(id string) string {
	return fmt.Sprintf("token_family:%s", id)
}
