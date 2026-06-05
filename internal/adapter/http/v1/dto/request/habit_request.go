package request

type CreateHabitRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
	TagIDs      []uint `form:"tags"`
}

type UpdateHabitRequest struct {
	Title        string `form:"title" binding:"required"`
	Description  string `form:"description" binding:"required"`
	AddTagIDs    []uint `form:"add_tags"`
	RemoveTagIDs []uint `form:"remove_tags"`
}

type GetAllHabitsRequest struct {
	TagIDs []uint `form:"tag_ids"`
	SortRequest
	SearchRequest
	PageRequest
}

type GetUserHabitsRequest struct {
	SortRequest
	SearchRequest
	PageRequest
}

type GetAchievementsRequest struct {
	PageRequest
}

type CreateTagRequest struct {
	Name string `json:"title" binding:"required"`
}

type EditTagRequest struct {
	Name string `json:"title" binding:"required"`
}
