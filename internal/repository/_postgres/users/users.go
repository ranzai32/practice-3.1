package users

import (
	"context"
	"database/sql"
	"fmt"
	"practice4/practice-4/internal/repository/_postgres"
	"practice4/practice-4/pkg/modules"
	"time"
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

	fmt.Println(users)
	return users, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*modules.User, error) {
	user := &modules.User{}
	err := r.db.DB.GetContext(ctx, user, "SELECT id, name, email, created_at FROM users WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound // really not found
	}

	if err != nil {
		return nil, err 
	}

	return user, nil
}