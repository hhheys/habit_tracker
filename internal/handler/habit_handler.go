package handler

import (
	"habit-tracker/internal/dto/request"
	appErrors "habit-tracker/internal/errors"
	"habit-tracker/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HabitHandler interface {
	GetAllHabits(c *gin.Context)
}

type habitHandler struct {
	log     *zap.Logger
	service service.Service
}

func NewHabitHandler(log *zap.Logger, service service.Service) HabitHandler {
	return &habitHandler{
		log:     log,
		service: service,
	}
}

func (h *habitHandler) GetAllHabits(c *gin.Context) {
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

	var query request.GetAllHabitsRequest
	query.Page = 1
	query.PageSize = 20
	query.Search = ""

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.GetAllHabits(&query, userID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, res)
}

//func (h *userHandler) CreateHabit(c *gin.Context) {
//	file, err := c.FormFile("image")
//	if err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//}
