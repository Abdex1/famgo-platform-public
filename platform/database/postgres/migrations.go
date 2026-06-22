
//platform/database/postgres/migrations.go

package postgres

import (
    "context"
    "fmt"
    "io/fs"
    "sort"
    "strings"
)

func (db *DB) RunMigrations(
    ctx context.Context,
    migrationFS fs.FS,
) error {

    _, err := db.Primary.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version TEXT PRIMARY KEY,
        applied_at TIMESTAMPTZ DEFAULT NOW()
    )
    `)
    if err != nil {
        return err
    }

    entries, err := fs.ReadDir(migrationFS, ".")
    if err != nil {
        return err
    }

    var files []string

    for _, e := range entries {
        if strings.HasSuffix(e.Name(), ".sql") {
            files = append(files, e.Name())
        }
    }

    sort.Strings(files)

    for _, file := range files {

        var count int

        err := db.Primary.QueryRow(
            ctx,
            `SELECT COUNT(*) FROM schema_migrations WHERE version=$1`,
            file,
        ).Scan(&count)

        if err != nil {
            return err
        }

        if count > 0 {
            continue
        }

        content, err := fs.ReadFile(migrationFS, file)
        if err != nil {
            return err
        }

        tx, err := db.Primary.Begin(ctx)
        if err != nil {
            return err
        }

        _, err = tx.Exec(ctx, string(content))
        if err != nil {
            tx.Rollback(ctx)
            return fmt.Errorf("migration failed %s: %w", file, err)
        }

        _, err = tx.Exec(
            ctx,
            `INSERT INTO schema_migrations(version) VALUES($1)`,
            file,
        )

        if err != nil {
            tx.Rollback(ctx)
            return err
        }

        if err := tx.Commit(ctx); err != nil {
            return err
        }
    }

    return nil
}
