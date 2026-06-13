package transaction

type SearchRequest struct {
	CustomerCode string `json:"customer_code" validate:"required"`
	AccountNo    string `json:"account_no" validate:"required"`
	FromDate     string `json:"from_date" validate:"required"`
	ToDate       string `json:"to_date" validate:"required"`
	Symbol       string `json:"symbol"`
	Filter       string `json:"filter"`
	Sort         string `json:"sort"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type TransactionItem struct {
	DataDT      string `json:"data_dt"`
	Symbol      string `json:"symbol"`
	BuySellType string `json:"buy_sell_type"`
	Volume      int64  `json:"volume"`
	Price       string `json:"price"`
	Amount      string `json:"amount"`
}

type SearchResult struct {
	Items      []TransactionItem
	Total      int
	Page       int
	PerPage    int
	TotalPages int
}
