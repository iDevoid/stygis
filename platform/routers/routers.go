package routers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Router data will be registered to http listener
type Router struct {
	Method  string
	Path    string
	Handler fiber.Handler
}

type routing struct {
	host           string
	domain         string
	allowedOrigins string
	routers        []Router
}

type handlerInfo struct {
	method string
	path   string
	time   time.Time
}

// Routers contains the functions of http handler to clean payloads and pass it the service
type Routers interface {
	Serve()
}

// Initialize is for initialize the handler
func Initialize(host, allowedOrigins string, routers []Router, domain string) Routers {
	return &routing{
		host,
		domain,
		allowedOrigins,
		routers,
	}
}

// Serve is to start serving the HTTP Listener for every domain
func (r *routing) Serve() {
	server := fiber.New()

	group := server.Group(fmt.Sprintf("/%s", r.domain))

	for _, router := range r.routers {
		group.Add(router.Method, router.Path, router.Handler)
	}

	logrus.WithFields(logrus.Fields{
		"host":   r.host,
		"domain": r.domain,
	}).Info("Starts Serving on HTTP")
	err := server.Listen(r.host)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host":   r.host,
			"domain": r.domain,
		}).Fatal(err)
		panic(err)
	}
}
