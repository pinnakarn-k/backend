package transaction

import (
	"context"
)

type Repository interface {
	Search(ctx context.Context, req SearchRequest) (SearchResult, error)
}

type repository struct {
	// db *sql.DB
}

func NewRepo() Repository {
	return &repository{
		// db: db,
	}
}

func (r *repository) Search(
	ctx context.Context,
	req SearchRequest,
) (SearchResult, error) {
	items := []TransactionItem{
		{
			DataDT:      "2026-06-13",
			Symbol:      "PTT",
			BuySellType: "BUY",
			Volume:      100,
			Price:       "32.50",
			Amount:      "3250.00",
		},
		{
			DataDT:      "2026-06-13",
			Symbol:      "KBANK",
			BuySellType: "SELL",
			Volume:      50,
			Price:       "145.00",
			Amount:      "7250.00",
		},
	}

	return SearchResult{
		Items:      items,
		Page:       1,
		PerPage:    20,
		Total:      2,
		TotalPages: 1,
	}, nil
}
