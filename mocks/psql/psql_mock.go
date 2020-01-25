// Code generated manually. DO NOT EDIT.

// Package mock_sql is a generated GoMock package.
package mock_psql

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// Connection returns the mocks of sqlx DB connection
func Connection() (*sqlx.DB, sqlmock.Sqlmock, *sql.DB) {
	db, mocked, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mocked, db
}
