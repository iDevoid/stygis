package model

import (
	"time"
)

// User is a data model for user, use for login, register, etc
type User struct {
	ID        int64     `db:"id" json:"user_id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"hash_password" json:"password,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
