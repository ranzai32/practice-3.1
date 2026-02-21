package users

import (
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

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}

	fmt.Println(users)
	return users, nil
}
