package persistence

//go:generate mockgen -destination=../../../mocks/profile/persistence_mock.go -package=profile_mock -source=profile.go

import (
	"context"

	"github.com/iDevoid/cptx"
	"github.com/iDevoid/stygis/internal/constant/model"
	"github.com/iDevoid/stygis/internal/constant/query"
	"github.com/iDevoid/stygis/internal/constant/state"
)

// ProfilePersistence contains the list of functions for database table profiles
type ProfilePersistence interface {
	InsertProfile(ctx context.Context, user *model.User) error
}

type profilePersistence struct {
	db cptx.Database
}

// ProfileInit is to init the profile persistence that contains data accounts
func ProfileInit(db cptx.Database) ProfilePersistence {
	return &profilePersistence{
		db,
	}
}

// InsertProfile records new user profile / public data user to database table profiles based on data user
func (pp *profilePersistence) InsertProfile(ctx context.Context, user *model.User) error {
	params := map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"full_name":   user.Username,
		"status":      state.UserInactiveStatus,
		"create_time": user.CreateTime,
	}
	_, err := pp.db.Main().ExecuteMustTx(ctx, query.ProfileInsert, params)
	return err
}
