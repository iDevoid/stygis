package rest

import (
	"errors"
	"fmt"
	"hash"
	"net/http"
	"strconv"

	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/constants/state"
	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/savsgio/atreugo/v10"
	"github.com/valyala/fasthttp"
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
	username := string(ctx.FormValue(state.HandlerKeyUsername))
	email := string(ctx.FormValue(state.HandlerKeyEmail))
	password := string(ctx.FormValue(state.HandlerKeyPassword))

	if username == "" || email == "" || password == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	us.hash.Write([]byte(password))
	password = fmt.Sprintf("%x", us.hash.Sum(nil))

	err := us.usecase.Register(ctx.RequestCtx, &model.User{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
	}
	return err
}

// SignIn handles API for loggin in the user
func (us *userService) SignIn(ctx *atreugo.RequestCtx) error {
	email := string(ctx.FormValue(state.HandlerKeyEmail))
	password := string(ctx.FormValue(state.HandlerKeyPassword))

	if email == "" || password == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	us.hash.Write([]byte(password))
	password = fmt.Sprintf("%x", us.hash.Sum(nil))

	userID, err := us.usecase.Login(ctx, email, password)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
	}
	authCookie := fasthttp.Cookie{}
	authCookie.SetKey(state.HandlerUserKeyCookie)
	authCookie.SetSecure(true)
	authCookie.SetValue(strconv.FormatInt(userID, 10))
	authCookie.SetMaxAge(60 * 60 * 24)
	authCookie.SetPath("/")
	ctx.Response.Header.Cookie(&authCookie)
	return err
}

// ShowProfile handlers the request to show the user profile
func (us *userService) ShowProfile(ctx *atreugo.RequestCtx) error {
	stringUserID := string(ctx.Request.Header.Cookie(state.HandlerUserKeyCookie))
	userID, err := strconv.ParseInt(stringUserID, 10, 64)
	if err != nil || userID < 1 {
		ctx.SetStatusCode(http.StatusBadRequest)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	user, err := us.usecase.Profile(ctx, userID)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return err
	}

	return ctx.JSONResponse(user)
}
