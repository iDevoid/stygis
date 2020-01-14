package mocks

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// PSQLMock returns the mocks of sqlx DB connection
func PSQLMock() (*sqlx.DB, sqlmock.Sqlmock, *sql.DB) {
	db, mocked, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mocked, db
}
