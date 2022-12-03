package database

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Migrate migrates database scheme
func Migrate(conn *sql.DB) error {
	var (
		err error
	)

	goose.SetBaseFS(migrations)

	if err = goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err = goose.Up(conn, "migrations"); err != nil {
		return err
	}

	return nil
}
