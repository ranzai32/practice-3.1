package repository

import (
	"context"
	"practice4/practice-4/internal/repository/_postgres"
	"practice4/practice-4/internal/repository/_postgres/users"
	"practice4/practice-4/pkg/modules"
)

type UserRepository interface {
	GetAll(ctx context.Context, limit, offset int64) ([]modules.User, error)
	CountUsers(ctx context.Context) (int64, error)
	GetByID(ctx context.Context, id int64) (*modules.User, error)
	Create(ctx context.Context, user *modules.User) (int64, error)
	Update(ctx context.Context, user *modules.User) error
	Delete(ctx context.Context, id int64) error
	CreateUserWithAudit(ctx context.Context, user *modules.User) (int64, error)
}

type Repositories struct {
	Users UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		Users: users.NewUserRepository(db),
	}
}
