package redis

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type connectionString struct {
	connection string
	password   string
	domain     string
}

// Connections contains the functions to handle the redis platform
type Connections interface {
	Open() *redis.Client
}

// Initialize to init the redis platform with connection string and password
// bad parameter can cause panic and stops the entire app where the initialize is being called
func Initialize(connection, password, domain string) Connections {
	return &connectionString{
		connection: connection,
		password:   password,
		domain:     domain,
	}
}

// Open is to open a connection to redis server
func (cs *connectionString) Open() *redis.Client {
	logrus.WithFields(logrus.Fields{
		"platform": "redis",
		"domain":   cs.domain,
	}).Info("Connecting to Redis Server")
	client := redis.NewClient(&redis.Options{
		Addr:     cs.connection,
		Password: cs.password,
		DB:       0,
	})
	err := client.Ping().Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"platform":   "redis",
			"connection": cs.connection,
			"password":   cs.password,
		}).Fatal(err)
		panic(err)
	}
	return client
}
