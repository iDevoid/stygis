package rest

import (
	"errors"
	"fmt"
	"hash"
	"net/http"

	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/savsgio/atreugo/v10"
)

type userService struct {
	usecase user.Usecase
	hash    hash.Hash
}

// HandleUser is to initialize the rest handler for domain user
func HandleUser(usecase user.Usecase, hash hash.Hash) user.Handler {
	return &userService{
		usecase: usecase,
		hash:    hash,
	}
}

// Test is the test handler function
func (us *userService) Test(ctx *atreugo.RequestCtx) error {
	return ctx.TextResponse("Hello World")
}

// CreateNewAccount handles user registration
func (us *userService) CreateNewAccount(ctx *atreugo.RequestCtx) error {
	username := string(ctx.FormValue("username"))
	email := string(ctx.FormValue("email"))
	password := string(ctx.FormValue("password"))

	if username == "" || email == "" || password == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		return errors.New("bad payload")
	}

	us.hash.Write([]byte(password))
	password = fmt.Sprintf("%x", us.hash.Sum(nil))

	err := us.usecase.NewAccount(ctx.RequestCtx, &model.User{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return err
	}
	return nil
}
