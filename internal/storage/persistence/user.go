package persistence

//go:generate mockgen -destination=../../../mocks/user/persistence_mock.go -package=user_mock -source=user.go

import (
	"context"

	"github.com/iDevoid/cptx"
	"github.com/iDevoid/stygis/internal/constant/model"
	"github.com/iDevoid/stygis/internal/constant/query"
	"github.com/iDevoid/stygis/internal/constant/state"
)

// UserPersistence contains the list of functions for database table users
type UserPersistence interface {
	InsertUser(ctx context.Context, user *model.User) error
}

type userPersistence struct {
	db cptx.Database
}

// UserInit is to init the user persistence that contains data accounts
func UserInit(db cptx.Database) UserPersistence {
	return &userPersistence{
		db,
	}
}

// InsertUser is the input the data record to database table users
func (up *userPersistence) InsertUser(ctx context.Context, user *model.User) error {
	params := map[string]interface{}{
		"username":     user.Username,
		"email":        user.Email,
		"hashed_email": user.HashedEmail,
		"password":     user.Password,
		"create_time":  user.CreateTime,
		"status":       state.UserInactiveStatus,
	}

	return up.db.Main().QueryRowMustTx(ctx, query.UserInsert, params, &user.ID)
}
