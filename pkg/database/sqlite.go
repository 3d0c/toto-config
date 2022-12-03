package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	// SQLite driver
	_ "github.com/mattn/go-sqlite3"

	"github.com/3d0c/toto-config/pkg/config"
)

var (
	instance *sql.DB
	once     sync.Once
)

// Instance is a DB connection singleton
// "FailFast" Strategy - should be initialized from cmd/totocfg/run on init
func Instance() *sql.DB {
	once.Do(func() {
		var err error

		if instance, err = connectDatabase(); err != nil {
			panic(err)
		}
	})

	return instance
}

func connectDatabase() (*sql.DB, error) {
	var (
		conn *sql.DB
		dsn  = config.TheConfig().Database.DSN
		err  error
	)

	if dsn == "" {
		return nil, fmt.Errorf("error opening database - DataSource can't be empty")
	}

	if err = os.MkdirAll(filepath.Dir(dsn), os.ModePerm); err != nil {
		return nil, err
	}

	if conn, err = sql.Open("sqlite3", dsn); err != nil {
		return nil, err
	}

	if err = Migrate(conn); err != nil {
		return nil, fmt.Errorf("error migrating database - %s", err)
	}

	return conn, nil
}
