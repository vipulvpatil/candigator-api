package storage

import "database/sql"

type DatabaseTransactionProvider interface {
	BeginTransaction() (DatabaseTransaction, error)
}

type databaseTransaction struct {
	*sql.Tx
}

type DatabaseTransaction interface {
	customDbHandler
	Commit() error
	Rollback() error
}

func (s *Storage) BeginTransaction() (DatabaseTransaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return &databaseTransaction{
		Tx: tx,
	}, nil
}
