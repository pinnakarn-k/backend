package transaction

import (
	"context"
)

type Repository interface {
	Search(ctx context.Context, req SearchRequest) (SearchRepoResult, error)
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
) (SearchRepoResult, error) {
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

	return SearchRepoResult{
		Items: items,
		Total: 2,
	}, nil
}
