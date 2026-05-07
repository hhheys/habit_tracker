package handler

import (
	"errors"
	"habit-tracker/internal/dto/request"
	appErrors "habit-tracker/internal/errors"
	"habit-tracker/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HabitHandler interface {
	UpdateHabit(c *gin.Context)
	CreateHabit(c *gin.Context)
	DeleteHabit(c *gin.Context)

	GetAllHabits(c *gin.Context)
	GetAllUserHabits(c *gin.Context)
	AddUserHabit(c *gin.Context)

	EditTag(c *gin.Context)
	GetAllTags(c *gin.Context)
	CreateTag(c *gin.Context)
	DeleteTag(c *gin.Context)
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

func (h *habitHandler) CreateHabit(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var req request.CreateHabitRequest
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	res, err := h.service.CreateHabit(&req, file)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *habitHandler) UpdateHabit(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		// если файл не обязателен — просто игнорируем
		if errors.Is(err, http.ErrMissingFile) {
			file = nil
		} else {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	habitID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req request.UpdateHabitRequest
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	res, err := h.service.UpdateHabit(uint(habitID), &req, file)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *habitHandler) DeleteHabit(c *gin.Context) {
	habitID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.service.DeleteHabit(uint(habitID))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *habitHandler) GetAllUserHabits(c *gin.Context) {
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

	var query request.GetUserHabitsRequest
	query.Sort = ""

	err := c.ShouldBindQuery(&query)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	res, err := h.service.GetAllUserHabits(userID, &query)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *habitHandler) AddUserHabit(c *gin.Context) {
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

	err = h.service.AddUserHabit(userID, uint(habitID))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *habitHandler) GetAllTags(c *gin.Context) {
	var query request.GetAllHabitTagsRequest

	err := c.ShouldBindQuery(&query)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	resp, err := h.service.GetAllTags(query)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *habitHandler) CreateTag(c *gin.Context) {
	var body request.CreateTagRequest
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	res, err := h.service.CreateTag(&body)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *habitHandler) DeleteTag(c *gin.Context) {
	habitID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	err = h.service.DeleteTag(uint(habitID))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *habitHandler) GetTagByID(c *gin.Context) {
	habitID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	resp, err := h.service.GetTagByID(uint(habitID))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *habitHandler) EditTag(c *gin.Context) {
	habitID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	var body request.EditTagRequest
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	resp, err := h.service.EditTag(uint(habitID), &body)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(http.StatusOK, resp)
}
