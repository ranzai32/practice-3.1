package users

import (
	"context"
	"database/sql"
	"practice4/practice-4/internal/repository/_postgres"
	"practice4/practice-4/pkg/modules"
	"time"
	"fmt"
	"practice4/practice-4/pkg/apperrors"
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

func (r *Repository) GetAll(ctx context.Context) ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.SelectContext(ctx, &users, "SELECT id, name, email, created_at FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*modules.User, error) {
	user := &modules.User{}
	err := r.db.DB.GetContext(ctx, user, "SELECT id, name, email, created_at FROM users WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrNotFound // really not found
	}

	if err != nil {
		return nil, err 
	}

	return user, nil
}

func (r *Repository) Create(ctx context.Context, user *modules.User) (int64, error) {
	var id int64
	err := r.db.DB.QueryRowContext(ctx, "INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) Update(ctx context.Context, user *modules.User) error {
	result, err := r.db.DB.ExecContext(ctx, "UPDATE users SET name = $1, email = $2 WHERE id = $3",
		user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return apperrors.ErrNotFound // no rows updated, user not found
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	result, err := r.db.DB.ExecContext(ctx, "DELETE FROM users where id = $1", id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return apperrors.ErrNotFound // no rows deleted, user not found
	}

	return nil
}