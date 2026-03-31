package request

type CreateHabitRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
}

type UpdateHabitRequest struct {
	ID          uint   `form:"id" binding:"required"`
	Title       string `form:"title" validate:"required"`
	Description string `form:"description" validate:"required"`
}

type GetAllHabitsRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search" binding:"omitempty"`
}
