package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/iDevoid/stygis/internal/constant/model"
	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/sirupsen/logrus"
)

// UserHandler contains the function of handler for domain user
type UserHandler interface {
	Test(c *fiber.Ctx) error
	RegistrationHandler(ctx *fiber.Ctx) error
}

type userHandler struct {
	userCase user.Usecase
}

// UserInit is to initialize the rest handler for domain user
func UserInit(userCase user.Usecase) UserHandler {
	return &userHandler{
		userCase,
	}
}

// Test is handler testing
func (uh *userHandler) Test(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello, World!")
}

func (uh *userHandler) RegistrationHandler(ctx *fiber.Ctx) error {
	var body model.User
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil || body.Username == "" || body.Email == "" || body.Password == "" {
		ctx.Status(http.StatusBadRequest)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	ctx.Response().SetBodyString("Registration Success!")

	err = uh.userCase.Registration(ctx.Context(), &body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"domain":  "user",
			"action":  "create new user",
			"usecase": "Register",
			"email":   body.Email,
		}).Errorln(err)

		ctx.Status(http.StatusInternalServerError)
		ctx.Response().SetBodyString(err.Error())
	}

	return nil
}
