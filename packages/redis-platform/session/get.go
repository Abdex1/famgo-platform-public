
// STEP 6 — GET SESSION

//packages/redis-platform/session/get.go

package session

import (
	"context"
	"encoding/json"
)

func (s *Store) Get(
	ctx context.Context,
	sessionID string,
) (*Session, error) {

	val, err := s.rdb.Get(
		ctx,
		sessionKey(sessionID),
	).Result()

	if err != nil {
		return nil, err
	}

	var sess Session

	if err := json.Unmarshal([]byte(val), &sess); err != nil {
		return nil, err
	}

	return &sess, nil
}
