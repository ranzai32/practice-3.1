package modules

import "time"

type User struct {
	ID        int64      `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type UserInput struct {
	Name  string `json:"name" example:"Alice"`
	Email string `json:"email" example:"alice@example.com"`
}

type AuditLog struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Action    string    `db:"action"`
	CreatedAt time.Time `db:"created_at"`
}

type PaginatedUsers struct {
	Users  []User `json:"users"`
	Total  int64  `json:"total"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}
