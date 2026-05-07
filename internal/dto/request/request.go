package request

type SortRequest struct {
	Sort string `form:"sort" binding:"omitempty"`
}

type SearchRequest struct {
	Search string `form:"search" binding:"omitempty"`
}

type PageRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100"`
}
