package query

type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalRows  int `json:"total_rows"`
	TotalPages int `json:"total_pages"`
}

func (qp *QueryParams) NewPaginationMeta(totalRows int) *PaginationMeta {
	totalPages := totalRows / qp.PageSize
	if totalRows%qp.PageSize > 0 {
		totalPages++
	}

	return &PaginationMeta{
		Page:       qp.Page,
		PageSize:   qp.PageSize,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}
}
