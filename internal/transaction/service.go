package transaction

import (
	"backend/internal/pagination"
	"context"
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

	repoResult, err := s.repository.Search(ctx, req)
	if err != nil {
		return SearchResult{}, err
	}

	page := 1
	if req.Limit > 0 {
		page = (req.Offset / req.Limit) + 1
	}

	return SearchResult{
		Items: repoResult.Items,
		Pagination: pagination.New(
			page,
			req.Limit,
			repoResult.Total,
		),
	}, nil
}
