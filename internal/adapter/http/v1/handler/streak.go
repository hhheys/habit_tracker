package handler

import (
	"habit-tracker/internal/adapter/http/v1/dto/response"
	streakuc "habit-tracker/internal/usecase/streak"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StreakHandler interface {
	CreateDailyConfirmation(c *gin.Context)
	GetHeatMap(c *gin.Context)
}

type streakHandler struct {
	service *streakuc.Service
}

func NewStreakHandler(service *streakuc.Service) StreakHandler {
	return &streakHandler{service: service}
}

func (h *streakHandler) CreateDailyConfirmation(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	habitID, ok := pathID(c)
	if !ok {
		return
	}
	streak, err := h.service.CreateDailyConfirmation(c.Request.Context(), streakuc.DailyConfirmationInput{
		UserID: userID, HabitID: habitID,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response.NewStreakResponse(streak))
}

func (h *streakHandler) GetHeatMap(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	days, err := h.service.GetHeatmap(c.Request.Context(), streakuc.HeatmapInput{
		UserID: userID, StartDate: c.Query("start_date"), EndDate: c.Query("end_date"),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, days)
}
