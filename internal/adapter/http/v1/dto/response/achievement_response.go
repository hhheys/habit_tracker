package response

import (
	domainachievement "habit-tracker/internal/domain/achievement"
	achievementuc "habit-tracker/internal/usecase/achievement"
	"time"

	"github.com/google/uuid"
)

type AchievementsResponse struct {
	Achievements []*AchievementResponse `json:"achievements"`
	Pagination   PaginationResponse     `json:"pagination"`
}

type AchievementResponse struct {
	ID          uuid.UUID  `json:"id"`
	Code        string     `json:"code"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Enabled     bool       `json:"enabled"`
	Unlocked    bool       `json:"unlocked"`
	UnlockedAt  *time.Time `json:"unlocked_at,omitempty"`
}

func NewAchievementsResponse(output *achievementuc.ListUserAchievementsOutput) *AchievementsResponse {
	items := make([]*AchievementResponse, len(output.Achievements))
	for i, item := range output.Achievements {
		items[i] = NewAchievementResponse(item)
	}

	return &AchievementsResponse{
		Achievements: items,
		Pagination:   NewPaginationResponse(output.Total, output.Limit, output.Offset),
	}
}

func NewAchievementResponse(item *domainachievement.UserAchievementListItem) *AchievementResponse {
	return &AchievementResponse{
		ID:          item.Achievement.ID,
		Code:        item.Achievement.Code,
		Title:       item.Achievement.Title,
		Description: item.Achievement.Description,
		Enabled:     item.Achievement.Enabled,
		Unlocked:    item.Unlocked,
		UnlockedAt:  item.UnlockedAt,
	}
}
