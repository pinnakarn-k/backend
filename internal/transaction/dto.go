package transaction

import "backend/internal/pagination"

type FileResult struct {
	FileName    string
	ContentType string
	Bytes       []byte
}

type MailResult struct {
	Success bool
	Message string
	RefNo   string
}

type SearchRequest struct {
	AccountType  string `json:"accountType" validate:"required"`
	CustomerCode string `json:"customer_code" validate:"required"`
	AccountNo    string `json:"account_no" validate:"required"`
	FromDate     string `json:"from_date" validate:"required"`
	ToDate       string `json:"to_date" validate:"required"`
	Symbol       string `json:"symbol"`
	Filter       string `json:"filter"`
	Sort         string `json:"sort"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`

	ExportType string `json:"export_type"`
}

type TransactionItem struct {
	DataDT      string `json:"data_dt"`
	Symbol      string `json:"symbol"`
	BuySellType string `json:"buy_sell_type"`
	Volume      int64  `json:"volume"`
	Price       string `json:"price"`
	Amount      string `json:"amount"`
}

type SearchRepoResult struct {
	Items []TransactionItem
	Total int
}

type SearchResult struct {
	Items      []TransactionItem
	Pagination pagination.Pagination
}
