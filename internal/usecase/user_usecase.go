package usecase

import (
	"context"
	"practice4/practice-4/internal/repository"
	"practice4/practice-4/pkg/apperrors"
	"practice4/practice-4/pkg/modules"
)

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

var _ UserUsecase = (*userUsecase)(nil)

func (u *userUsecase) GetAll(ctx context.Context, limit, offset int64) (*modules.PaginatedUsers, error) {
	users, err := u.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	total, err := u.repo.CountUsers(ctx)
	if err != nil {
		return nil, err
	}
	return &modules.PaginatedUsers{
		Users:  users,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (u *userUsecase) GetByID(ctx context.Context, id int64) (*modules.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *userUsecase) Create(ctx context.Context, user *modules.User) (int64, error) {
	if user.Name == "" || user.Email == "" {
		return 0, apperrors.ErrValidation
	}
	return u.repo.Create(ctx, user)
}

func (u *userUsecase) Update(ctx context.Context, user *modules.User) error {
	if user.Name == "" || user.Email == "" {
		return apperrors.ErrValidation
	}
	return u.repo.Update(ctx, user)
}

func (u *userUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *userUsecase) CreateUserWithAudit(ctx context.Context, user *modules.User) (int64, error) {
	if user.Name == "" || user.Email == "" {
		return 0, apperrors.ErrValidation
	}
	return u.repo.CreateUserWithAudit(ctx, user)
}
