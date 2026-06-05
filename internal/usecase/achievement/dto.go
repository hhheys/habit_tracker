package achievement

import (
	domainachievement "habit-tracker/internal/domain/achievement"
)

type ListUserAchievementsParams struct {
	UserID uint
	Limit  int
	Offset int
}

type ListUserAchievementsOutput struct {
	Achievements []*domainachievement.UserAchievementListItem
	Limit        int
	Offset       int
	Total        int64
}
