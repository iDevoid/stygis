package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type connectionString struct {
	master string
	slave  string
	domain string
}

// Connections return the wrapped actions of postgres database
type Connections interface {
	Open() *Database
}

// Database where the opened database connection is being used
type Database struct {
	Master *sqlx.DB
	Slave  *sqlx.DB
}

// Initialize is to init the postgres platform with connection string both master or slave
// use the master connection string if there's no slave database, both must be the same connection string
// never set it to empty string, it will cause the fatal and stops the entire app where the database is being initialize
func Initialize(master, slave, domain string) Connections {
	return &connectionString{
		master: master,
		slave:  slave,
		domain: domain,
	}
}

// Open is creating the database postgres connections, both master and slave
func (cs *connectionString) Open() *Database {
	logrus.WithFields(logrus.Fields{
		"platform": "postgres",
		"domain":   cs.domain,
	}).Info("Connecting to PostgreSQL DB")

	logrus.WithFields(logrus.Fields{
		"platform": "postgres",
		"domain":   cs.domain,
	}).Info("Opening Connection to Master")
	dbMaster, err := sqlx.Open("postgres", cs.master)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"platform":   "postgres",
			"type":       "master",
			"connection": cs.master,
		}).Fatal(err)
		panic(err)
	}

	logrus.WithFields(logrus.Fields{
		"platform": "postgres",
		"domain":   cs.domain,
	}).Info("Opening Connection to Slave")
	dbSlave, err := sqlx.Open("postgres", cs.master)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"platform":   "postgres",
			"type":       "slave",
			"connection": cs.slave,
		}).Fatal(err)
		panic(err)
	}
	return &Database{
		Master: dbMaster,
		Slave:  dbSlave,
	}
}
