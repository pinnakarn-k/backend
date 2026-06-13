package pagination

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}

	if limit > MaxLimit {
		return MaxLimit
	}

	return limit
}

func buildPagination(
	offset int,
	limit int,
	total int,
) (int, int) {
	page := 1

	if limit > 0 {
		page = (offset / limit) + 1
	}

	totalPages := 0

	if total > 0 && limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return page, totalPages
}
