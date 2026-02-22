package usecase

import (
	"context"
	"practice4/practice-4/pkg/modules"
)

type UserUsecase interface {
	GetAll(ctx context.Context, limit, offset int64) (*modules.PaginatedUsers, error)
	GetByID(ctx context.Context, id int64) (*modules.User, error)
	Create(ctx context.Context, user *modules.User) (int64, error)
	Update(ctx context.Context, user *modules.User) error
	Delete(ctx context.Context, id int64) error
	CreateUserWithAudit(ctx context.Context, user *modules.User) (int64, error)
}
