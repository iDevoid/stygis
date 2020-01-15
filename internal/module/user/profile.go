package user

import (
	"context"

	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/constants/state"
	"github.com/sirupsen/logrus"
)

// Profile check if the user is signed in and the user data is saved inside the caching
func (s *service) Profile(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepository.DataProfile(ctx, userID)
	if err != nil {
		// log formatting
		log := logrus.WithFields(logrus.Fields{
			"domain":  "user",
			"action":  "get user data",
			"usecase": "Profile",
		})
		log.WithField(state.LogType, "repository").Errorln(err)
	}
	return user, err
}
