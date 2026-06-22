package transaction

import (
	"backend/internal/response"
	"backend/internal/validator"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SendEmail(c *fiber.Ctx) error {
	var req SearchRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := validator.Validate(req); err != nil {
		return response.Error(c, err)
	}

	result, err := h.service.SendEmail(
		c.Context(),
		req, // TODO
	)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (h *Handler) Download(c *fiber.Ctx) error {
	var req SearchRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := validator.Validate(req); err != nil {
		return response.Error(c, err)
	}

	file, err := h.service.Download(
		c.Context(),
		req, // TODO
	)
	if err != nil {
		return err
	}

	c.Set("Content-Type", file.ContentType)
	c.Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, file.FileName),
	)

	return c.Send(file.Bytes)
}

func (h *Handler) Search(c *fiber.Ctx) error {
	var req SearchRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	fmt.Printf("%+v\n", req)

	if err := validator.Validate(req); err != nil {
		return response.Error(c, err)
	}

	result, err := h.service.Search(c.Context(), req)
	if err != nil {
		return err
	}

	return response.SuccessWithPagination(
		c,
		result.Items,
		result.Pagination,
	)
}

// func (r *repository) Search(
// 	ctx context.Context,
// 	req SearchRequest,
// ) (SearchResult, error) {

// 	if req.Limit <= 0 {
// 		req.Limit = 20
// 	}

// 	if req.Limit > 100 {
// 		req.Limit = 100
// 	}

// 	baseWhere := `
// 	FROM transaction_history
// 	WHERE customer_code = @p1
// 	  AND account_no = @p2
// 	  AND data_dt BETWEEN @p3 AND @p4
// 	`

// 	args := []any{
// 		req.CustomerCode,
// 		req.AccountNo,
// 		req.FromDate,
// 		req.ToDate,
// 	}

// 	// =========================
// 	// COUNT
// 	// =========================

// 	countQuery := `
// 		SELECT COUNT(1)
// 	` + baseWhere

// 	var total int

// 	err := r.db.QueryRowContext(
// 		ctx,
// 		countQuery,
// 		args...,
// 	).Scan(&total)
// 	if err != nil {
// 		return SearchResult{}, err
// 	}

// fmt.Errorf("query transaction search: %w", err)

// 	// =========================
// 	// SORT
// 	// =========================

// 	sortColumn := "data_dt"

// 	switch req.Sort {
// 	case "symbol":
// 		sortColumn = "symbol"

// 	case "amount":
// 		sortColumn = "amount"

// 	case "price":
// 		sortColumn = "price"

// 	case "volume":
// 		sortColumn = "volume"
// 	}

// 	// =========================
// 	// SELECT
// 	// =========================

// 	query := fmt.Sprintf(`
// 		SELECT
// 			data_dt,
// 			symbol,
// 			buy_sell_type,
// 			volume,
// 			price,
// 			amount
// 		%s
// 		ORDER BY %s DESC
// 		OFFSET @p5 ROWS
// 		FETCH NEXT @p6 ROWS ONLY
// 	`, baseWhere, sortColumn)

// 	rows, err := r.db.QueryContext(
// 		ctx,
// 		query,
// 		req.CustomerCode,
// 		req.AccountNo,
// 		req.FromDate,
// 		req.ToDate,
// 		req.Offset,
// 		req.Limit,
// 	)
// 	if err != nil {
// 		return SearchResult{}, err
// 	}
// 	defer rows.Close()

// 	items := make([]TransactionItem, 0)

// 	for rows.Next() {

// 		var item TransactionItem

// 		err := rows.Scan(
// 			&item.DataDT,
// 			&item.Symbol,
// 			&item.BuySellType,
// 			&item.Volume,
// 			&item.Price,
// 			&item.Amount,
// 		)
// 		if err != nil {
// 			return SearchResult{}, err
// 		}

// 		items = append(items, item)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return SearchResult{}, err
// 	}

// page, totalPages := buildPagination(
// 	req.Offset,
// 	req.Limit,
// 	total,
// )

// 	return SearchResult{
// 		Items:      items,
// 		Page:       page,
// 		PerPage:    req.Limit,
// 		Total:      total,
// 		TotalPages: totalPages,
// 	}, nil
// }
