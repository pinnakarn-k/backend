package health

import "backend/internal/apperror"

type Service interface {
	Ready() error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Ready() error {
	if err := s.repository.Ready(); err != nil {
		return apperror.ErrServiceUnavailable
	}

	return nil
}
