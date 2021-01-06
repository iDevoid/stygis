package repository

//go:generate mockgen -destination=../../mocks/user/repository_mock.go -package=user_mock -source=user.go

import (
	"crypto/sha512"
	"fmt"

	"github.com/iDevoid/stygis/internal/constant/model"
	"github.com/iDevoid/stygis/platform/encryption"
)

// UserRepository contains the functions of data logic for domain user
type UserRepository interface {
	Encrypt(user *model.User) (err error)
}

type userRepository struct {
	systemEncryptKey string
}

// UserInit initializes the data logic / repository for domain user
func UserInit(systemEncryptKey string) UserRepository {
	return &userRepository{
		systemEncryptKey,
	}
}

func (ur *userRepository) Encrypt(user *model.User) (err error) {
	if user.Email != "" {
		user.Email, err = encryption.Encrypting(user.Email, ur.systemEncryptKey)
		if err != nil {
			return
		}
		hash512 := sha512.Sum512([]byte(user.Email))
		user.Email = fmt.Sprintf("%x", hash512)
	}

	if user.Password != "" {
		hash512 := sha512.Sum512([]byte(user.Password))
		user.Password = fmt.Sprintf("%x", hash512)
	}

	return
}
