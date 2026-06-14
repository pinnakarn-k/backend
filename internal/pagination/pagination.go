package pagination

import "math"

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func New(page, perPage, total int) Pagination {
	if page <= 0 {
		page = 1
	}

	perPage = NormalizeLimit(perPage)

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(perPage)))
	}

	return Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
