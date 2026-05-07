package request

type CreateHabitRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
	TagIDs      []int  `form:"tags"`
}

type UpdateHabitRequest struct {
	Title       string `form:"title" validate:"required"`
	Description string `form:"description" validate:"required"`
	TagIDs      []int  `form:"tags"`
}

type GetAllHabitsRequest struct {
	TagIDs []int `form:"tag_ids"`
	SortRequest
	SearchRequest
	PageRequest
}

type GetUserHabitsRequest struct {
	SortRequest
}

type GetAllHabitTagsRequest struct {
	SearchRequest
	PageRequest
}

type CreateTagRequest struct {
	Name string `json:"title" binding:"required"`
}

type EditTagRequest struct {
	Name string `json:"title" binding:"required"`
}
