package storage

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type StorageAccessor interface {
	UserRetriever
	DatabaseTransactionProvider
	TeamHydrator
}

type Storage struct {
	db          *sql.DB
	IdGenerator utilities.CuidGenerator
}

type StorageOptions struct {
	Db          *sql.DB
	IdGenerator utilities.CuidGenerator
}

func NewDbStorage(opts StorageOptions) (*Storage, error) {
	if opts.Db == nil {
		return nil, errors.New("Needs a backing database")
	}

	if opts.IdGenerator == nil {
		opts.IdGenerator = &utilities.RandomIdGenerator{}
	}

	return &Storage{
		db:          opts.Db,
		IdGenerator: opts.IdGenerator,
	}, nil
}

// A lot of queries/updates need to be part of a transaction but not all.
// So we have the below interface that will allow the caller to either pass in *sql.DB or *sql.Tx depending on it's needs and our code will handle it without any issues.
type customDbHandler interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
