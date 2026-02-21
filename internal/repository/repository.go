package repository

import (
	"practice4/practice-4/internal/repository/_postgres"
	"practice4/practice-4/internal/repository/_postgres/users"
	"practice4/practice-4/pkg/modules"
	"context"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]modules.User, error)
	GetByID(ctx context.Context, id int) (*modules.User, error)
	Create(ctx context.Context, user *modules.User) (int64, error)
	Update(ctx context.Context, user *modules.User) error
	Delete(ctx context.Context, id int) error
}

type Repositories struct {
	Users UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		Users: users.NewUserRepository(db),
	}
}
