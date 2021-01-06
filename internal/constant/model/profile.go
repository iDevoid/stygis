package model

import "time"

// Profile is for getting data of user profile
type Profile struct {
	ID             int64     `db:"id" json:"id"`
	Username       string    `db:"username" json:"username"`
	FullName       string    `db:"full_name" json:"full_name"`
	ProfilePicture string    `db:"profile_picture" json:"profile_picture"`
	CoverPicture   string    `db:"cover_picture" json:"cover_picture"`
	Bio            string    `db:"bio" json:"bio"`
	Card           string    `db:"card" json:"card"`
	Followers      int64     `db:"followers" json:"followers"`
	Following      int64     `db:"following" json:"following"`
	RegisterTime   time.Time `db:"register_time" json:"register_time"`
	Status         bool      `db:"status" json:"status"`
	CreateTime     time.Time `db:"create_time" json:"create_time"`
	UpdateTime     time.Time `db:"update_time" json:"update_time"`
}
