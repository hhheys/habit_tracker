package handler

import (
	appErrors "habit-tracker/internal/errors"
	"habit-tracker/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StreakHandler interface {
	CreateDailyConfirmation(c *gin.Context)
}

type streakHandler struct {
	log     *zap.Logger
	service service.Service
}

func NewStreakHandler(log *zap.Logger, service service.Service) StreakHandler {
	return &streakHandler{
		log:     log,
		service: service,
	}
}

func (s *streakHandler) CreateDailyConfirmation(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		_ = c.Error(appErrors.ErrUnauthorized)
		return
	}
	userID, ok := userIDRaw.(uint)
	if !ok {
		_ = c.Error(appErrors.ErrUnauthorized)
		return
	}

	habitID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = s.service.CreateDailyConfirmation(userID, uint(habitID))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}
