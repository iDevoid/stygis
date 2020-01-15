package persistence

import (
	"context"
	"database/sql"

	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/constants/query"
	"github.com/iDevoid/stygis/internal/constants/state"
	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/iDevoid/stygis/platform/postgres"
	"github.com/jmoiron/sqlx"
)

type userPersistence struct {
	db *postgres.Database
}

// UserInit is to init the user persistence that contains data accounts
func UserInit(db *postgres.Database) user.Persistence {
	return &userPersistence{
		db,
	}
}

// Create is the persistence to save new user to db
func (up *userPersistence) Create(ctx context.Context, user *model.User) error {
	tx, err := up.db.Master.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	param := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
		"status":   state.UserActiveAccount,
	}
	query, args, _ := sqlx.Named(query.InsertNewUser, param)
	query = up.db.Master.Rebind(query)
	err = tx.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err == nil {
		err = tx.Commit()
	}
	return err
}

// FindByID is to find user inside db using only userID
// this returns the pointer of the selected data if it doesn't error
func (up *userPersistence) FindByID(ctx context.Context, userID int64) (*model.User, error) {
	param := map[string]interface{}{
		"user_id": userID,
		"status":  state.UserActiveAccount,
	}
	query, args, _ := sqlx.Named(query.SelectUserByID, param)
	query = up.db.Slave.Rebind(query)

	user := new(model.User)
	err := up.db.Slave.GetContext(ctx, user, query, args...)
	if err == sql.ErrNoRows {
		err = nil
	}
	return user, err
}

// Find is the function of user storage that select all data, using email and hashed password
// this function is being used for login
func (up *userPersistence) Find(ctx context.Context, email, password string) (*model.User, error) {
	param := map[string]interface{}{
		"email":    email,
		"password": password,
		"status":   state.UserActiveAccount,
	}
	query, args, _ := sqlx.Named(query.SelectUserByEmail, param)
	query = up.db.Slave.Rebind(query)

	user := new(model.User)
	err := up.db.Slave.GetContext(ctx, user, query, args...)
	if err == sql.ErrNoRows {
		err = nil
	}
	return user, err
}

// ChangePassword is to for changing old hashed password with new hashed password and user data to change the current password inside database
func (up *userPersistence) ChangePassword(ctx context.Context, newPassword string, user *model.User) error {
	tx, err := up.db.Master.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	param := map[string]interface{}{
		"new_password": newPassword,
		"user_id":      user.ID,
		"email":        user.Email,
		"old_password": user.Password,
		"status":       state.UserActiveAccount,
	}
	query, args, _ := sqlx.Named(query.UpdateUserPassword, param)
	query = up.db.Master.Rebind(query)
	_, err = tx.ExecContext(ctx, query, args...)
	if err == nil {
		err = tx.Commit()
	}
	return nil
}

// Delete basically doing the soft delete for the logged in user account
func (up *userPersistence) Delete(ctx context.Context, user *model.User) error {
	tx, err := up.db.Master.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	param := map[string]interface{}{
		"user_id":         user.ID,
		"email":           user.Email,
		"password":        user.Password,
		"active_status":   state.UserActiveAccount,
		"inactive_status": state.UserInactiveAccount,
	}
	query, args, _ := sqlx.Named(query.DeactivateUser, param)
	query = up.db.Master.Rebind(query)
	_, err = tx.ExecContext(ctx, query, args...)
	if err == nil {
		err = tx.Commit()
	}
	return err
}
