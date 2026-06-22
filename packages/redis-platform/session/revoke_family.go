
/* STEP 8 — REVOKE TOKEN FAMILY

CRITICAL SECURITY FEATURE.

If refresh token theft detected:

ALL descendant sessions revoked.
*/
//packages/redis-platform/session/revoke_family.go

package session

import (
	"context"
)

func (s *Store) RevokeFamily(
	ctx context.Context,
	familyID string,
) error {

	sessions, err := s.rdb.SMembers(
		ctx,
		tokenFamilyKey(familyID),
	).Result()

	if err != nil {
		return err
	}

	for _, id := range sessions {
		_ = s.Revoke(ctx, id)
	}

	return nil
}
