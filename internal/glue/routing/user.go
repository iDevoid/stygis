package routing

import (
	"net/http"

	"github.com/iDevoid/stygis/internal/handler/rest"
	"github.com/iDevoid/stygis/platform/routers"
)

// UserRouting returns the list of routers for domain user
func UserRouting(handler rest.UserHandler) []routers.Router {
	return []routers.Router{
		{
			Method:  http.MethodGet,
			Path:    "/test",
			Handler: handler.Test,
		},

		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: handler.RegistrationHandler,
		},
	}
}
