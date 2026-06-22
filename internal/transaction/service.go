package transaction

import (
	"backend/internal/pagination"
	"context"
	"fmt"
)

type Service interface {
	Search(ctx context.Context, req SearchRequest) (SearchResult, error)
	Download(ctx context.Context, req SearchRequest) (FileResult, error)
	SendEmail(ctx context.Context, req SearchRequest) (MailResult, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func prepareExportRequest(req SearchRequest) SearchRequest {
	req.Offset = 0
	req.Limit = 0
	return req
}

func (s *service) SendEmail(ctx context.Context, req SearchRequest) (MailResult, error) {
	file, err := s.generateExportFile(ctx, prepareExportRequest(req))
	if err != nil {
		return MailResult{}, err
	}

	_ = file

	// TODO call email api

	return MailResult{
		Success: true,
		Message: "email sent",
		RefNo:   "TEST001",
	}, nil
}

func (s *service) Download(ctx context.Context, req SearchRequest) (FileResult, error) {
	return s.generateExportFile(ctx, prepareExportRequest(req))
}

func (s *service) generateExportFile(ctx context.Context, req SearchRequest) (FileResult, error) {
	switch req.AccountType {
	case "F":
		return s.exportF(ctx, req)
	default:
		return FileResult{}, fmt.Errorf("invalid account type")
	}
}

func (s *service) exportF(ctx context.Context, req SearchRequest) (FileResult, error) {
	// TODO repo

	// TODO columes field

	// TODO columes data

	// TODO setup

	// TODO excel/pdf
	switch req.ExportType {
	case "excel":
		return FileResult{
			FileName:    "transaction_f.xlsx",
			ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Bytes:       []byte("dummy excel"),
		}, nil

	case "pdf":
		return FileResult{
			FileName:    "transaction_f.pdf",
			ContentType: "application/pdf",
			Bytes:       []byte("dummy pdf"),
		}, nil
	default:
		return FileResult{}, fmt.Errorf("invalid export type")
	}
}

// case "excel":
// 	return FileResult{
// 		FileName:    "transaction.xlsx",
// 		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
// 		Bytes:       []byte("excel"),
// 	}, nil

// case "pdf":
// 	return FileResult{
// 		FileName:    "transaction.pdf",
// 		ContentType: "application/pdf",
// 		Bytes:       []byte("pdf"),
// 	}, nil

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
