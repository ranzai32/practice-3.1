package usecase

import (
	"context"
	"practice4/practice-4/internal/repository"
	"practice4/practice-4/pkg/modules"
	"practice4/practice-4/pkg/apperrors"
)

type userUsecase struct {
    repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
    return &userUsecase{repo: repo}
}

var _ UserUsecase = (*userUsecase)(nil)

func (u *userUsecase) GetAll(ctx context.Context) ([]modules.User, error) {
    return u.repo.GetAll(ctx)
}

func (u *userUsecase) GetByID(ctx context.Context, id int) (*modules.User, error) {
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

func (u *userUsecase) Delete(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}