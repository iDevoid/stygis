package user

import (
	"context"

	"github.com/iDevoid/stygis/internal/constants/state"
	"github.com/sirupsen/logrus"
)

// Login checks the email and password inside the persistence and save it to cache
func (s *service) Login(ctx context.Context, email, password string) (int64, error) {
	var err error
	var stringOperation string
	defer func() {
		if err == nil {
			return
		}

		// log formatting
		log := logrus.WithFields(logrus.Fields{
			"domain":  "user",
			"action":  "login user",
			"usecase": "Login",
		})
		log.WithField(state.LogType, stringOperation).Errorln(err)
	}()

	user, err := s.userPersistence.Find(ctx, email, password)
	if err != nil {
		stringOperation = "persistence"
		return 0, err
	}

	// save to caching
	err = s.userCaching.Save(ctx, user)
	if err != nil {
		stringOperation = "cache"
	}
	return user.ID, err
}
