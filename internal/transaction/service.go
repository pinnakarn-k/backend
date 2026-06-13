package transaction

import (
	"backend/internal/pagination"
	"context"
	"fmt"
)

type Service interface {
	Search(ctx context.Context, req SearchRequest) (SearchResult, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Search(ctx context.Context, req SearchRequest) (SearchResult, error) {
	req.Limit = pagination.NormalizeLimit(req.Limit)
	fmt.Printf("%-v", req)
	return s.repository.Search(ctx, req)
}
