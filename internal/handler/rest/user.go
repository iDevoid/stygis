package rest

import (
	"errors"
	"net/http"

	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/savsgio/atreugo/v10"
)

type userService struct {
	Usecase user.Usecase
}

// UserHandler contains all the functions for handling http request
type UserHandler interface {
	Test(ctx *atreugo.RequestCtx) error
	CreateNewAccount(ctx *atreugo.RequestCtx) error
}

// HandleUser is to initialize the rest handler for domain user
func HandleUser(usecase user.Usecase) UserHandler {
	return &userService{
		Usecase: usecase,
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

	// us.Usecase.NewAccount(ctx.RequestCtx, &model.User{})
	return nil
}
