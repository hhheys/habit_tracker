package handler

import (
	"habit-tracker/internal/adapter/http/v1/dto/request"
	"habit-tracker/internal/adapter/http/v1/dto/response"
	achievementuc "habit-tracker/internal/usecase/achievement"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AchievementHandler interface {
	GetAchievements(c *gin.Context)
}

type achievementHandler struct {
	achievements *achievementuc.Service
}

func NewAchievementHandler(achievements *achievementuc.Service) AchievementHandler {
	return &achievementHandler{achievements: achievements}
}

func (h *achievementHandler) GetAchievements(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req request.GetAchievementsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	limit, offset := paging(req.Page, req.PageSize)
	output, err := h.achievements.ListUserAchievements(c.Request.Context(), achievementuc.ListUserAchievementsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewAchievementsResponse(output))
}
