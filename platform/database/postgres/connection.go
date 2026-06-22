
//platform/database/postgres/connection.go
package postgres

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
    Primary *pgxpool.Pool
    Replica *pgxpool.Pool
    Config  Config
}

func New(ctx context.Context, cfg Config) (*DB, error) {
    poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
    if err != nil {
        return nil, fmt.Errorf("postgres parse config: %w", err)
    }

    poolCfg.MaxConns = cfg.MaxConns
    poolCfg.MinConns = cfg.MinConns
    poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
    poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
    poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

    ctxTimeout, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
    defer cancel()

    primary, err := pgxpool.NewWithConfig(ctxTimeout, poolCfg)
    if err != nil {
        return nil, fmt.Errorf("postgres create pool: %w", err)
    }

    if err := waitForDatabase(ctxTimeout, primary); err != nil {
        primary.Close()
        return nil, err
    }

    db := &DB{
        Primary: primary,
        Config:  cfg,
    }

    if cfg.ReadReplicaDSN != "" {
        replicaCfg, err := pgxpool.ParseConfig(cfg.ReadReplicaDSN)
        if err != nil {
            return nil, fmt.Errorf("replica parse config: %w", err)
        }

        replica, err := pgxpool.NewWithConfig(ctxTimeout, replicaCfg)
        if err != nil {
            return nil, fmt.Errorf("replica pool create: %w", err)
        }

        db.Replica = replica
    }

    return db, nil
}

func waitForDatabase(ctx context.Context, pool *pgxpool.Pool) error {
    var lastErr error

    for i := 0; i < 30; i++ {
        if err := pool.Ping(ctx); err == nil {
            return nil
        } else {
            lastErr = err
        }

        time.Sleep(2 * time.Second)
    }

    return fmt.Errorf("postgres unavailable: %w", lastErr)
}

func (db *DB) Close() {
    if db.Primary != nil {
        db.Primary.Close()
    }

    if db.Replica != nil {
        db.Replica.Close()
    }
}
