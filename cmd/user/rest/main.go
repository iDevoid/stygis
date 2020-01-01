package main

import (
	"flag"
	"os"

	"github.com/iDevoid/stygis/internal/glue/routing"
	"github.com/iDevoid/stygis/internal/handler/rest"
	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/iDevoid/stygis/internal/repository"
	"github.com/iDevoid/stygis/internal/storage/cache"
	"github.com/iDevoid/stygis/internal/storage/persistence"
	"github.com/iDevoid/stygis/platform/postgres"
	"github.com/iDevoid/stygis/platform/redis"
	"github.com/iDevoid/stygis/platform/routers"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	redisURL      = "127.0.0.1:6379"
	redisPassword = ""

	postgresURL = "postgresql://usut@127.0.0.1/ticket?sslmode=disable"
	domain      = "user"
)

var testInit bool

func init() {
	flag.BoolVar(&testInit, "test", false, "initialize test mode without serving")
	flag.Parse()
}

func main() {
	psql := postgres.Initialize(postgresURL, postgresURL, domain)
	postgresConn := psql.Open()

	rds := redis.Initialize(redisURL, redisPassword, domain)
	redisConn := rds.Open()

	database := persistence.UserInit(postgresConn)
	caching := cache.UserInit(redisConn)
	repo := repository.UserInit(caching, database)

	usecase := user.InitializeDomain(database, repo)
	handler := rest.HandleUser(usecase)
	router := routing.UserInit(handler).NewRouters()
	servant := routers.Initialize(":9000", router, domain)

	if testInit {
		logrus.Info("Initialize test mode Finished!")
		os.Exit(0)
	}

	servant.Serve()
}