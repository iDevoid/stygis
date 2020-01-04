package user

import (
	"context"

	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/sirupsen/logrus"
)

// NewAccount sets the new data user into persistence
func (s *service) NewAccount(ctx context.Context, user *model.User) error {
	// log formatting
	stringType := "type"
	log := logrus.WithFields(logrus.Fields{
		"domain":  "user",
		"action":  "create new user",
		"usecase": "NewAccount",
	})

	err := s.userPersistence.Create(ctx, user)
	if err != nil {
		log.WithField(stringType, "persistence").Errorln(err)
	}

	// this saving to the cache is a business logic
	// user can automatically login after registration is A business logic because it is a flow of user experience
	// so, don't put your business logic inside repository
	// repository is just a data storing logic
	err = s.userCaching.Save(ctx, user)
	if err != nil {
		log.WithField(stringType, "caching").Errorln(err)
	}
	return err
}
