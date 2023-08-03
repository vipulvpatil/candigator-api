package storage

type StorageAccessorMock struct {
	UserRetriever
	DatabaseTransactionProvider
	TeamHydrator
	FileUploadAccessor
	CandidateAccessor
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

func WithTeamHydratorMock(mock TeamHydrator) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.TeamHydrator = mock
	}
}

func WithFileUploadAccessorMock(mock FileUploadAccessor) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.FileUploadAccessor = mock
	}
}

func WithCandidateAccessorMock(mock CandidateAccessor) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.CandidateAccessor = mock
	}
}
