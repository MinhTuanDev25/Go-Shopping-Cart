package utils

import "strconv"

type Pagination struct {
	Page            int32 `json:"page"`
	Limit           int32 `json:"limit"`
	TotalRecords    int32 `json:"total_records"`
	TotalPages      int32 `json:"total_pages"`
	HasNextPage     bool  `json:"has_next_page"`
	HasPreviousPage bool  `json:"has_previous_page"`
}

func NewPagination(page, limit, totalRecords int32) *Pagination {

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limitStr := GetEnv("LIMIT_RECORDS_PER_PAGE", "10")
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil || limitInt <= 0 {
			limit = 10
		}
		limit = int32(limitInt)
	}

	totalPages := (totalRecords + limit - 1) / limit

	return &Pagination{
		Page:            page,
		Limit:           limit,
		TotalRecords:    totalRecords,
		TotalPages:      totalPages,
		HasNextPage:     totalPages > page,
		HasPreviousPage: page > 1,
	}
}

func NewPaginationResponse(data any, page, limit, totalRecords int32) map[string]any {
	return map[string]any{
		"data":       data,
		"pagination": NewPagination(page, limit, totalRecords),
	}
}
