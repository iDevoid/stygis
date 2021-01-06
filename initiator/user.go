package initiator

import (
	"fmt"
	"os"

	"github.com/iDevoid/cptx"
	"github.com/iDevoid/stygis/internal/glue/routing"
	"github.com/iDevoid/stygis/internal/handler/rest"
	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/iDevoid/stygis/internal/repository"
	"github.com/iDevoid/stygis/internal/storage/persistence"
	"github.com/iDevoid/stygis/platform/routers"
	"github.com/sirupsen/logrus"
)

const (
	postgresURL = "postgresql://%s:%s@%s/iDevoid-db?sslmode=disable"

	domain = "user"
)

// User initializes the domain user
func User(testInit bool) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbURL := fmt.Sprintf(postgresURL, dbUser, dbPass, dbHost)
	psql := cptx.Initialize(dbURL, dbURL, domain)
	postgresDB, postgresTX := psql.Open()

	databaseUser := persistence.UserInit(postgresDB)
	databaseProfile := persistence.ProfileInit(postgresDB)

	encryptKey := os.Getenv("ENCRYPTION_KEY")
	repo := repository.UserInit(encryptKey)

	usecase := user.Initialize(postgresTX, repo, databaseUser, databaseProfile)

	handler := rest.UserInit(usecase)
	router := routing.UserRouting(handler)

	port := os.Getenv("HOST_PORT")
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	server := routers.Initialize(port, allowedOrigins, router, domain)

	if testInit {
		logrus.Info("Initialize test mode Finished!")
		os.Exit(0)
	}

	server.Serve()
}
