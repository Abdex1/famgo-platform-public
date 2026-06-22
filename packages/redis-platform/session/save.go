
// STEP 5 — SAVE SESSION

//packages/redis-platform/session/save.go

package session

import (
	"context"
	"encoding/json"
	"time"
)

func (s *Store) Save(
	ctx context.Context,
	session Session,
) error {

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	ttl := time.Until(session.ExpiresAt)

	pipe := s.rdb.TxPipeline()

	pipe.Set(
		ctx,
		sessionKey(session.SessionID),
		data,
		ttl,
	)

	pipe.SAdd(
		ctx,
		userSessionsKey(session.UserID),
		session.SessionID,
	)

	pipe.SAdd(
		ctx,
		tokenFamilyKey(session.TokenFamilyID),
		session.SessionID,
	)

	_, err = pipe.Exec(ctx)

	return err
}
