package storage

type StorageAccessorMock struct {
	UserRetriever
	DatabaseTransactionProvider
}

type StorageAccessorMockOption func(*StorageAccessorMock)

func NewStorageAccessorMock(opts ...StorageAccessorMockOption) *StorageAccessorMock {
	mock := &StorageAccessorMock{}
	for _, opt := range opts {
		opt(mock)
	}

	return mock
}

func WithDatabaseTransactionProviderMock(mock DatabaseTransactionProvider) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.DatabaseTransactionProvider = mock
	}
}
