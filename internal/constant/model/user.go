package model

import (
	"database/sql"
	"time"
)

// User is the structure for table users that contains the sensitive data / private data
type User struct {
	ID          int64        `db:"id" json:"user_id"`
	Username    string       `db:"username" json:"username"`
	Email       string       `db:"email" json:"email"`
	HashedEmail string       `db:"hashed_email" json:"hashed_email"`
	Password    string       `db:"password" json:"password"`
	CreateTime  time.Time    `db:"create_time" json:"create_time"`
	Status      int8         `db:"status" json:"status"`
	UpdateTime  sql.NullTime `db:"update_time" json:"update_time"`
}
