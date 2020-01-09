package routing

import (
	"github.com/iDevoid/stygis/platform/routers"

	"github.com/iDevoid/stygis/internal/module/user"
)

type userHandlers struct {
	handler user.Handler
}

// UserInit is to initialize the routers for domain user
func UserInit(handler user.Handler) user.Router {
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
