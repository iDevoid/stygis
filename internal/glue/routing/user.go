package routing

import (
	"github.com/iDevoid/stygis/platform/routers"

	"github.com/iDevoid/stygis/internal/handler/rest"
)

type userHandlers struct {
	handler rest.UserHandler
}

// UserRouter contains the functions that will be used for the routing domain user
type UserRouter interface {
	NewRouters() []*routers.Router
}

// UserInit is to initialize the routers for domain user
func UserInit(handler rest.UserHandler) UserRouter {
	return &userHandlers{
		handler: handler,
	}
}

// NewRouters returns the data router for domain user that will serve the rest API
func (uh *userHandlers) NewRouters() []*routers.Router {
	return []*routers.Router{
		&routers.Router{
			Method:  "GET",
			URL:     "/test",
			Handler: uh.handler.Test,
		},

		&routers.Router{
			Method:  "POST",
			URL:     "/account/new",
			Handler: uh.handler.CreateNewAccount,
		},
	}
}
