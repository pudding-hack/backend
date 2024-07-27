package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string         `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Address   string         `json:"address" db:"address"`
	Username  string         `json:"username" db:"username"`
	Email     string         `json:"email" db:"email"`
	Password  string         `json:"password" db:"password"`
	RoleID    int            `json:"role_id" db:"role_id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	CreatedBy string         `json:"created_by" db:"created_by"`
	UpdatedBy string         `json:"updated_by" db:"updated_by"`
	DeletedBy sql.NullString `json:"deleted_by" db:"deleted_by"`
}

type Role struct {
	ID        int        `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	UpdatedBy string     `json:"updated_by" db:"updated_by"`
	DeletedBy *string    `json:"deleted_by" db:"deleted_by"`
}
