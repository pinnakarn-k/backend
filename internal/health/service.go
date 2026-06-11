package health

import (
	"backend/internal/apperror"
	"backend/internal/requestcontext"
	"fmt"
)

type Service interface {
	Ready(reqCtx requestcontext.RequestContext) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Ready(reqCtx requestcontext.RequestContext) error {
	fmt.Printf("%+v\n", reqCtx)

	if err := s.repository.Ready(); err != nil {
		return apperror.ErrServiceUnavailable
	}

	return nil
}
