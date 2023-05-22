package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/vipulvpatil/candidate-tracker-go/internal/config"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

// This function will make a connection to the database only once.
func InitDb(cfg *config.Config, logger utilities.Logger) (*sql.DB, error) {
	var err error

	connStr := cfg.DbUrl
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// this will be printed in the terminal, confirming the connection to the database
	logger.LogMessageln("The database is connected")
	return db, nil
}
