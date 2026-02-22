package users

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"practice4/practice-4/internal/repository/_postgres"
	"practice4/practice-4/pkg/apperrors"
	"practice4/practice-4/pkg/modules"
)

type Repository struct {
	db            *_postgres.Dialect
	executionTime time.Duration
}

func NewUserRepository(dv *_postgres.Dialect) *Repository {
	return &Repository{
		db:            dv,
		executionTime: 5 * time.Second,
	}
}

func (r *Repository) GetAll(ctx context.Context, limit, offset int64) ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.SelectContext(ctx, &users,
		"SELECT id, name, email, created_at FROM users WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetAll: %w", err)
	}
	return users, nil
}

func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.DB.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM users WHERE deleted_at IS NULL")
	if err != nil {
		return 0, fmt.Errorf("CountUsers: %w", err)
	}
	return count, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*modules.User, error) {
	user := &modules.User{}
	err := r.db.DB.GetContext(ctx, user,
		"SELECT id, name, email, created_at FROM users WHERE id = $1 AND deleted_at IS NULL", id)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetByID: %w", err)
	}
	return user, nil
}

func (r *Repository) Create(ctx context.Context, user *modules.User) (int64, error) {
	var id int64
	err := r.db.DB.QueryRowContext(ctx,
		"INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, time.Now()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Create: %w", err)
	}
	return id, nil
}

func (r *Repository) Update(ctx context.Context, user *modules.User) error {
	result, err := r.db.DB.ExecContext(ctx,
		"UPDATE users SET name = $1, email = $2 WHERE id = $3 AND deleted_at IS NULL",
		user.Name, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Update RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.DB.ExecContext(ctx,
		"UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Delete RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *Repository) CreateUserWithAudit(ctx context.Context, user *modules.User) (int64, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("CreateUserWithAudit BeginTx: %w", err)
	}
	defer tx.Rollback()

	var id int64
	err = tx.QueryRowContext(ctx,
		"INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, time.Now()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateUserWithAudit insert user: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO audit_logs (user_id, action, created_at) VALUES ($1, $2, $3)",
		id, "create", time.Now())
	if err != nil {
		return 0, fmt.Errorf("CreateUserWithAudit insert audit: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("CreateUserWithAudit commit: %w", err)
	}
	return id, nil
}
