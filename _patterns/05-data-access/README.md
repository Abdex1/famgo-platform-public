# 🗄️ DATA ACCESS PATTERNS
## Extracted from uber-master

**Status:** Pattern 5/8

---

## Connection Pooling

```go
import "github.com/jackc/pgx/v5/pgxpool"

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, err
    }
    
    config.MaxConns = 25
    config.MinConns = 5
    config.MaxConnLifetime = 5 * time.Minute
    config.MaxConnIdleTime = 2 * time.Minute
    
    return pgxpool.NewWithConfig(ctx, config)
}
```

## Repository Pattern

```go
type Repository struct {
    db *pgxpool.Pool
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Entity, error) {
    var entity Entity
    err := r.db.QueryRow(ctx, `
        SELECT id, name, status, created_at
        FROM entities
        WHERE id = $1 AND deleted_at IS NULL
    `, id).Scan(&entity.ID, &entity.Name, &entity.Status, &entity.CreatedAt)
    
    if err == pgx.ErrNoRows {
        return nil, ErrNotFound
    }
    return &entity, err
}

func (r *Repository) Create(ctx context.Context, entity *Entity) error {
    return r.db.QueryRow(ctx, `
        INSERT INTO entities (name, status)
        VALUES ($1, $2)
        RETURNING id, created_at
    `, entity.Name, entity.Status).Scan(&entity.ID, &entity.CreatedAt)
}

func (r *Repository) Update(ctx context.Context, entity *Entity) error {
    result, err := r.db.Exec(ctx, `
        UPDATE entities
        SET name = $1, status = $2, updated_at = NOW()
        WHERE id = $3
    `, entity.Name, entity.Status, entity.ID)
    
    if result.RowsAffected() == 0 {
        return ErrNotFound
    }
    return err
}
```

**Pattern 5 Status:** READY FOR USE

---
