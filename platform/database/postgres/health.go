//platform/database/postgres/health.go

package postgres

import (
	"context"
	"time"
)

func (db *DB) Health(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return db.Primary.Ping(ctxTimeout)
}
