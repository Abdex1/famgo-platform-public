// services/user-service/internal/infrastructure/postgres_user_repo.go
// PostgreSQL User Repository

package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
)

// PostgresUserRepository implements UserRepository
type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, phone, email, first_name, last_name, status, created_at, updated_at
         FROM users WHERE id = $1`,
		userID)

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Phone,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, phone, email, first_name, last_name, status, created_at, updated_at
         FROM users WHERE phone = $1`,
		phone)

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Phone,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, phone, email, first_name, last_name, status, created_at, updated_at
         FROM users WHERE email = $1`,
		email)

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Phone,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, phone, email, first_name, last_name, status, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID,
		user.Phone,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET phone = $1, email = $2, first_name = $3, last_name = $4, status = $5, updated_at = $6
         WHERE id = $7`,
		user.Phone,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Status,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

func (r *PostgresUserRepository) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, phone, email, first_name, last_name, status, created_at, updated_at
         FROM users LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Phone,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID)
	return err
}
