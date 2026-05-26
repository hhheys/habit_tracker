package streak

type DailyConfirmationInput struct {
	UserID  uint
	HabitID uint
}

type HeatmapInput struct {
	UserID    uint
	StartDate string
	EndDate   string
}
