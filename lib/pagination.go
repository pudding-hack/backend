package lib

type Pagination struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

func NewPagination(page, pageSize, total int) Pagination {
	totalPage := total / pageSize
	if total%pageSize > 0 {
		totalPage++
	}
	return Pagination{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
		TotalData: total,
	}
}

type PaginationRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Keyword  string `json:"keyword"`
}

func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationRequest) Limit() int {
	return p.PageSize
}
