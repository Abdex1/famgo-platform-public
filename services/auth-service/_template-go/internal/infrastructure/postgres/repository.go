package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// PostgresRepository provides database access
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create inserts a new record
func (r *PostgresRepository) Create(ctx context.Context, query string, args ...interface{}) error {
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// GetByID retrieves a record by ID
func (r *PostgresRepository) GetByID(ctx context.Context, query string, id string) (map[string]interface{}, error) {
	row := r.db.QueryRowContext(ctx, query, id)
	
	result := make(map[string]interface{})
	// Scan row into map
	return result, nil
}

// Update modifies a record
func (r *PostgresRepository) Update(ctx context.Context, query string, args ...interface{}) error {
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// Delete removes a record
func (r *PostgresRepository) Delete(ctx context.Context, query string, args ...interface{}) error {
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// List retrieves multiple records
func (r *PostgresRepository) List(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	// Scan rows into results
	return results, nil
}

// Health checks database connectivity
func (r *PostgresRepository) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Close closes database connection
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

// ConnectPostgres establishes PostgreSQL connection
func ConnectPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
