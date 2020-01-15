package user

import (
	"context"

	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/constants/state"
	"github.com/sirupsen/logrus"
)

// Register sets the new data user into persistence
func (s *service) Register(ctx context.Context, user *model.User) error {
	// log formatting
	log := logrus.WithFields(logrus.Fields{
		"domain":  "user",
		"action":  "create new user",
		"usecase": "Register",
	})

	err := s.userPersistence.Create(ctx, user)
	if err != nil {
		log.WithField(state.LogType, "persistence").Errorln(err)
	}

	// this saving to the cache is a business logic
	// user can automatically login after registration is A business logic because it is a flow of user experience
	// so, don't put your business logic inside repository
	// repository is just a data storing logic
	err = s.userCaching.Save(ctx, user)
	if err != nil {
		log.WithField(state.LogType, "caching").Errorln(err)
	}
	return err
}
