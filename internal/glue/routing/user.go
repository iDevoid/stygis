package routing

import (
	"github.com/iDevoid/stygis/platform/routers"
	"net/http"

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
			Method:  http.MethodGet,
			URL:     "/test",
			Handler: uh.handler.Test,
		},
		&routers.Router{
			Method:  http.MethodGet,
			URL:     "/account/profile",
			Handler: uh.handler.ShowProfile,
		},

		&routers.Router{
			Method:  http.MethodPost,
			URL:     "/account/register",
			Handler: uh.handler.CreateNewAccount,
		},
		&routers.Router{
			Method:  http.MethodPost,
			URL:     "/account/login",
			Handler: uh.handler.SignIn,
		},
	}
}
