package service

import (
	"habit-tracker/internal/dto/request"
	"habit-tracker/internal/dto/response"
	"habit-tracker/internal/models"
	"habit-tracker/internal/repository/postgres"
	"io"
	"mime/multipart"
	"os"

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
	UpdateHabit(r *request.UpdateHabitRequest, image *multipart.FileHeader) (*response.HabitResponse, error)
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

func (s *habitService) SaveImage(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	_, err = io.Copy(out, src)
	return err
}

func (s *habitService) CreateHabit(r *request.CreateHabitRequest, image *multipart.FileHeader) (*response.HabitResponse, error) {
	habitModel := models.NewHabit(r.Title, r.Description)
	habitModel.ImageFilename = image.Filename

	err := s.repository.CreateHabit(habitModel)
	if err != nil {
		return nil, err
	}

	err = s.SaveImage(image, HabitImagePath+habitModel.ImageFilename)
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

func (s *habitService) UpdateHabit(r *request.UpdateHabitRequest, image *multipart.FileHeader) (*response.HabitResponse, error) {
	habit, err := s.repository.GetHabitByID(r.ID)
	if err != nil {
		return nil, err
	}
	habit.Title = r.Title
	habit.Description = r.Description

	if image != nil {
		habit.ImageFilename = image.Filename
	}
	err = s.repository.UpdateHabit(habit)
	if err != nil {
		return nil, err
	}
	if image != nil {
		saveErr := s.SaveImage(image, HabitImagePath+habit.ImageFilename)
		if saveErr != nil {
			return nil, err
		}
	}
	return response.NewHabitResponse(habit), nil
}
