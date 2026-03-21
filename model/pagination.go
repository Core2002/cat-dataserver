package model

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int `form:"page" binding:"min=1"`              // 页码，从1开始
	PageSize int `form:"page_size" binding:"min=1,max=100"` // 每页数量，最大100
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Data       any   `json:"data"`        // 数据列表
	Total      int64 `json:"total"`       // 总记录数
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页数量
	TotalPages int   `json:"total_pages"` // 总页数
}

// GetOffset 计算偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetPageSize 获取每页数量（设置默认值）
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10 // 默认每页10条
	}
	if p.PageSize > 100 {
		return 100 // 最大每页100条
	}
	return p.PageSize
}

// GetPage 获取当前页码（设置默认值）
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1 // 默认第1页
	}
	return p.Page
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(data any, total int64, page, pageSize int) *PaginationResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
