package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/iDevoid/stygis/internal/constant/model"
)

// Registration is basically the flow of creating a new unverified user
func (s *service) Registration(ctx context.Context, user *model.User) error {
	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("%s is bad email", user.Email)
	}
	user.CreateTime = time.Now()

	err := s.userRepo.Encrypt(user)
	if err != nil {
		return fmt.Errorf("failed to encrypt with error: %s", err.Error())
	}

	tx, err := s.transaction.Begin(&ctx)
	if err != nil {
		return fmt.Errorf("cannot start the transaction with error : %s", err.Error())
	}
	defer tx.Rollback()

	err = s.userPersist.InsertUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to save new user %s", err.Error())
	}

	err = s.profilePersist.InsertProfile(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create a profile %s", err.Error())
	}

	return tx.Commit()
}
