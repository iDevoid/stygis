package routers

import (
	"github.com/savsgio/atreugo/v10"
	"github.com/sirupsen/logrus"
)

// Router data will be registered to http listener
type Router struct {
	Method  string
	URL     string
	Handler atreugo.View
}

type routing struct {
	host   string
	domain string
}

// Routers contains the functions of http handler to clean payloads and pass it the service
type Routers interface {
	Serve()
}

var handlers []*Router

// Initialize is for initialize the handler
func Initialize(host string, routers []*Router, domain string) Routers {
	handlers = routers
	return &routing{
		host:   host,
		domain: domain,
	}
}

// Serve is to start serving the HTTP Listener for every domain
func (us *routing) Serve() {
	config := &atreugo.Config{
		Addr: us.host,
	}
	server := atreugo.New(config)

	for _, router := range handlers {
		server.Path(router.Method, router.URL, router.Handler)
	}

	logrus.WithFields(logrus.Fields{
		"host":   us.host,
		"domain": us.domain,
	}).Info("Starts Serving on HTTP")
	err := server.ListenAndServe()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host":   us.host,
			"domain": us.domain,
		}).Fatal(err)
		panic(err)
	}
}
