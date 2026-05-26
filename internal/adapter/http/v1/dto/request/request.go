package request

type SortRequest struct {
	SortBy    string `form:"sort_by" binding:"omitempty"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc ASC DESC"`
}

type SearchRequest struct {
	Search string `form:"search" binding:"omitempty"`
}

type PageRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}
