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

// why is this needed? because sonarqube will say there's a duplication if you write it multiple times
var stringType = "type"
var stringConnection = "connection"

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
	// log fieds for logrus, no need to write this multiple times
	logFields := logrus.Fields{
		"platform": "postgres",
		"domain":   cs.domain,
	}
	logMasterFields := logrus.Fields{
		stringType:       "master",
		stringConnection: cs.master,
	}
	logSlaveFields := logrus.Fields{
		stringType:       "slave",
		stringConnection: cs.slave,
	}

	logrus.WithFields(logFields).Info("Connecting to PostgreSQL DB")

	logrus.WithFields(logFields).Info("Opening Connection to Master")
	dbMaster, err := sqlx.Open("postgres", cs.master)
	if err != nil {
		logrus.WithFields(logMasterFields).Fatal(err)
		panic(err)
	}
	err = dbMaster.Ping()
	if err != nil {
		logrus.WithFields(logMasterFields).Fatal(err)
		panic(err)
	}

	logrus.WithFields(logFields).Info("Opening Connection to Slave")
	dbSlave, err := sqlx.Open("postgres", cs.master)
	if err != nil {
		logrus.WithFields(logSlaveFields).Fatal(err)
		panic(err)
	}
	err = dbSlave.Ping()
	if err != nil {
		logrus.WithFields(logSlaveFields).Fatal(err)
		panic(err)
	}

	return &Database{
		Master: dbMaster,
		Slave:  dbSlave,
	}
}
