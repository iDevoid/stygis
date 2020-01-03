package user

import (
	"context"

	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/sirupsen/logrus"
)

// NewAccount sets the new data user into persistance
func (s *service) NewAccount(ctx context.Context, user *model.User) error {
	err := s.userPersistence.Create(ctx, user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"domain":  "user",
			"action":  "create new user",
			"type":    "persistance",
			"usecase": "NewAccount",
		}).Errorln(err)
	}
	return err
}
