package main

import (
	"crypto/sha256"
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
	// change all connection strings based on your own credentials
	redisURL      = "127.0.0.1:6379"
	redisPassword = ""
	// postgresURL   = "postgresql://postgres@127.0.0.1/postgres?sslmode=disable"
	postgresURL = "postgresql://postgres:tokopedia789@127.0.0.1/usut?sslmode=disable"

	domain = "user"
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

	usecase := user.InitializeDomain(database, caching, repo)

	hash := sha256.New()

	handler := rest.HandleUser(usecase, hash)
	router := routing.UserInit(handler).Routers()
	servant := routers.Initialize(":9000", router, domain)

	if testInit {
		logrus.Info("Initialize test mode Finished!")
		os.Exit(0)
	}

	servant.Serve()
}
