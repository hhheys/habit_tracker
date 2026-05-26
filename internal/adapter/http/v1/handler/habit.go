package handler

import (
	"errors"
	"fmt"
	"habit-tracker/internal/adapter/http/v1/dto/request"
	"habit-tracker/internal/adapter/http/v1/dto/response"
	"habit-tracker/internal/domain"
	habituc "habit-tracker/internal/usecase/habit"
	taguc "habit-tracker/internal/usecase/tag"
	userhabituc "habit-tracker/internal/usecase/userhabit"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const habitImagePath = "./images/habits"

type HabitHandler interface {
	GetAllHabits(c *gin.Context)
	GetHabitByID(c *gin.Context)
	CreateHabit(c *gin.Context)
	UpdateHabit(c *gin.Context)
	DeleteHabit(c *gin.Context)
	GetAllUserHabits(c *gin.Context)
	AddUserHabit(c *gin.Context)
	GetAllTags(c *gin.Context)
	GetTagByID(c *gin.Context)
	CreateTag(c *gin.Context)
	EditTag(c *gin.Context)
	DeleteTag(c *gin.Context)
}

type habitHandler struct {
	habits     *habituc.Service
	userHabits *userhabituc.Service
	tags       *taguc.Service
}

func NewHabitHandler(habits *habituc.Service, userHabits *userhabituc.Service, tags *taguc.Service) HabitHandler {
	return &habitHandler{habits: habits, userHabits: userHabits, tags: tags}
}

func (h *habitHandler) GetAllHabits(c *gin.Context) {
	var req request.GetAllHabitsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	limit, offset := paging(req.Page, req.PageSize)
	output, _, err := h.habits.ListHabits(c.Request.Context(), habituc.ListHabitsParams{
		UserID: userID, Search: req.Search, SortBy: req.SortBy, SortOrder: req.SortOrder, Limit: limit, Offset: offset,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewAllHabitsResponse(output.Habits, output.Limit, output.Offset, output.Total))
}

func (h *habitHandler) GetHabitByID(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	habit, err := h.habits.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewHabitResponse(habit))
}

func (h *habitHandler) CreateHabit(c *gin.Context) {
	var req request.CreateHabitRequest
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err)
		return
	}
	filename, err := saveOptionalImage(c, true)
	if err != nil {
		_ = c.Error(err)
		return
	}
	habit, err := h.habits.Create(c.Request.Context(), &habituc.CreateHabitInput{
		Title: req.Title, Description: req.Description, Tags: req.TagIDs, ImageFilename: filename,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response.NewHabitResponse(habit))
}

func (h *habitHandler) UpdateHabit(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req request.UpdateHabitRequest
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err)
		return
	}
	filename, err := saveOptionalImage(c, false)
	if err != nil {
		_ = c.Error(err)
		return
	}
	habit, err := h.habits.Update(c.Request.Context(), &habituc.UpdateHabitInput{
		ID: id, Title: req.Title, Description: req.Description, ImageFilename: filename,
		AddTagIDs: req.AddTagIDs, RemoveTagIDs: req.RemoveTagIDs,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewHabitResponse(habit))
}

func (h *habitHandler) DeleteHabit(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	if err := h.habits.DeleteByID(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *habitHandler) GetAllUserHabits(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	var req request.GetUserHabitsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}
	limit, offset := paging(req.Page, req.PageSize)
	output, err := h.userHabits.List(c.Request.Context(), userhabituc.ListUserHabitsParams{
		UserID: userID, Search: req.Search, SortBy: req.SortBy, SortOrder: req.SortOrder, Limit: limit, Offset: offset,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewUserHabitsResponse(output))
}

func (h *habitHandler) AddUserHabit(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	habitID, ok := pathID(c)
	if !ok {
		return
	}
	habit, err := h.userHabits.Add(c.Request.Context(), userhabituc.AddUserHabitInput{UserID: userID, HabitID: habitID})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, habit)
}

func (h *habitHandler) GetAllTags(c *gin.Context) {
	tags, err := h.tags.GetAll(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewHabitTagsResponse(tags))
}

func (h *habitHandler) GetTagByID(c *gin.Context) {
	tag, err := h.tags.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewHabitTagResponse(tag))
}

func (h *habitHandler) CreateTag(c *gin.Context) {
	var req request.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	tag, err := h.tags.Create(c.Request.Context(), &taguc.CreateTagInput{Name: req.Name})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response.NewHabitTagResponse(tag))
}

func (h *habitHandler) EditTag(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req request.EditTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	tag, err := h.tags.Update(c.Request.Context(), &taguc.UpdateTagInput{TagID: id, NewName: req.Name})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewHabitTagResponse(tag))
}

func (h *habitHandler) DeleteTag(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	if err := h.tags.Delete(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func pathID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(err)
		return 0, false
	}
	return uint(id), true
}

func currentUserID(c *gin.Context) (uint, bool) {
	userID, ok := c.Get("userID")
	if !ok {
		_ = c.Error(domain.ErrUnauthorized)
		return 0, false
	}
	id, ok := userID.(uint)
	if !ok {
		_ = c.Error(domain.ErrUnauthorized)
		return 0, false
	}
	return id, true
}

func paging(page, pageSize int) (int, int) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	return pageSize, (page - 1) * pageSize
}

func saveOptionalImage(c *gin.Context, required bool) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		if !required && errors.Is(err, http.ErrMissingFile) {
			return "", nil
		}
		return "", err
	}
	return saveImage(c, file)
}

func saveImage(c *gin.Context, file *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(habitImagePath, 0o755); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("habit_%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
	return filename, c.SaveUploadedFile(file, filepath.Join(habitImagePath, filename))
}
