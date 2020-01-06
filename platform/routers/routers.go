package routers

import (
	"fmt"
	"time"

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

type handlerInfo struct {
	method string
	path   string
	time   time.Time
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

	server.UseBefore(func(ctx *atreugo.RequestCtx) error {
		ctx.SetUserValue("info", handlerInfo{
			method: string(ctx.Method()),
			path:   string(ctx.Path()),
			time:   time.Now(),
		})
		return ctx.Next()
	})

	server.UseAfter(func(ctx *atreugo.RequestCtx) error {
		info := ctx.UserValue("info").(handlerInfo)
		timeHandler := time.Since(info.time)

		// you can put datadog, prometheus or any monitoring platform
		logrus.WithFields(logrus.Fields{
			"path":   info.path,
			"method": info.method,
		}).Infoln(fmt.Sprintf("execution time: %v millisecond; %v microsecond", timeHandler.Milliseconds(), timeHandler.Microseconds()))

		return ctx.Next()
	})

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
