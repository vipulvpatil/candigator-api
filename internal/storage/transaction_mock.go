package storage

import (
	"errors"
)

type DatabaseTransactionMock struct {
	customDbHandler
	Committed  bool
	Rolledback bool
}

func (d *DatabaseTransactionMock) Commit() error {
	d.Committed = true
	return nil
}

func (d *DatabaseTransactionMock) Rollback() error {
	d.Rolledback = true
	return nil
}

type DatabaseTransactionProviderMock struct {
	Transaction *DatabaseTransactionMock
}

func (s *DatabaseTransactionProviderMock) BeginTransaction() (DatabaseTransaction, error) {
	if s.Transaction == nil {
		return nil, errors.New("unable to begin a db transaction")
	}
	return s.Transaction, nil
}
