package service

import (
	"fmt"
	"habit-tracker/internal/dto/request"
	"habit-tracker/internal/dto/response"
	"habit-tracker/internal/models"
	"habit-tracker/internal/repository/postgres"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

var (
	HabitImagePath = "./images/habits/"
)

type HabitService interface {
	GetAllHabits(r *request.GetAllHabitsRequest, requestUserID uint) (*response.AllHabitsResponse, error)
	GetHabitByID(id uint) (*response.HabitResponse, error)
	CreateHabit(r *request.CreateHabitRequest, image *multipart.FileHeader) (*response.HabitResponse, error)
	DeleteHabit(id uint) error
	UpdateHabit(habitID uint, r *request.UpdateHabitRequest, image *multipart.FileHeader) (*response.HabitResponse, error)

	SaveImage(file *multipart.FileHeader, dir string) (string, error)

	GetAllUserHabits(userID uint, query *request.GetUserHabitsRequest) (*response.UserHabitsResponse, error)
	GetUserHabitByID(userID uint, habitID uint) (*response.UserHabitResponse, error)
	AddUserHabit(userID, habitID uint) error
}

type habitService struct {
	log        *zap.Logger
	repository postgres.Repository
}

func NewHabitService(logger *zap.Logger, repository postgres.Repository) HabitService {
	return &habitService{
		log:        logger,
		repository: repository,
	}
}

func (s *habitService) GetAllHabits(r *request.GetAllHabitsRequest, requestUserID uint) (*response.AllHabitsResponse, error) {
	allowedSorts := map[string]string{
		"new":       "h.created_at DESC",
		"title":     "h.title",
		"added":     "is_added DESC",
		"not_added": "is_added ASC",
	}

	sortColumn, ok := allowedSorts[r.Sort]
	if !ok {
		sortColumn = "h.created_at" // дефолт
	}

	r.Sort = sortColumn

	total, habits, err := s.repository.GetAllHabits(r, requestUserID)
	if err != nil {
		return nil, err
	}
	return response.NewAllHabitsResponse(habits, r.Page, r.PageSize, total), nil
}

func (s *habitService) GetHabitByID(id uint) (*response.HabitResponse, error) {
	habit, err := s.repository.GetHabitByID(id)
	if err != nil {
		return nil, err
	}
	return response.NewHabitResponse(habit), nil
}

func (s *habitService) SaveImage(file *multipart.FileHeader, dir string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// создаём директорию, если её нет
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename) // сохраняем исходное расширение
	filename := fmt.Sprintf("habit_%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(dir, filename)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (s *habitService) CreateHabit(r *request.CreateHabitRequest, image *multipart.FileHeader) (*response.HabitResponse, error) {
	habitModel := models.NewHabit(r.Title, r.Description)
	imageFilename, err := s.SaveImage(image, HabitImagePath)
	if err != nil {
		return nil, err
	}
	habitModel.ImageFilename = imageFilename

	err = s.repository.CreateHabit(habitModel)
	if err != nil {
		return nil, err
	}

	return response.NewHabitResponse(habitModel), nil
}

func (s *habitService) DeleteHabit(id uint) error {
	_, err := s.repository.GetHabitByID(id)
	if err != nil {
		return err
	}
	err = s.repository.DeleteHabit(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *habitService) UpdateHabit(habitID uint, r *request.UpdateHabitRequest, image *multipart.FileHeader) (*response.HabitResponse, error) {
	habit, err := s.repository.GetHabitByID(habitID)
	if err != nil {
		return nil, err
	}
	habit.Title = r.Title
	habit.Description = r.Description

	err = s.repository.UpdateHabit(habit)
	if err != nil {
		return nil, err
	}
	if image != nil {
		newImageFilename, saveErr := s.SaveImage(image, HabitImagePath)
		if saveErr != nil {
			return nil, err
		}
		habit.ImageFilename = newImageFilename
	}
	return response.NewHabitResponse(habit), nil
}

func (s *habitService) GetAllUserHabits(userID uint, query *request.GetUserHabitsRequest) (*response.UserHabitsResponse, error) {
	habits, err := s.repository.GetAllUserHabits(userID, query)
	if err != nil {
		return nil, err
	}
	return response.NewUserHabitsResponse(habits), nil
}

func (s *habitService) GetUserHabitByID(userID uint, habitID uint) (*response.UserHabitResponse, error) {
	habit, err := s.repository.GetUserHabit(userID, habitID)
	if err != nil {
		return nil, err
	}
	return response.NewUserHabitResponse(habit), nil
}

func (s *habitService) AddUserHabit(userID, habitID uint) error {
	_, err := s.repository.AddHabit(userID, habitID)
	if err != nil {
		return err
	}
	return nil
}
