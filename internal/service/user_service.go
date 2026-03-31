package service

import (
	"habit-tracker/internal/auth"
	"habit-tracker/internal/dto/request"
	"habit-tracker/internal/dto/response"
	appErrors "habit-tracker/internal/errors"
	"habit-tracker/internal/models"
	"habit-tracker/internal/repository/postgres"

	"go.uber.org/zap"
)

type UserService interface {
	RegisterUser(user *request.UserRegisterRequest) (*response.UserRegisterResponse, error)
	AuthUser(user request.UserLoginRequest) (*response.UserLoginResponse, error)
}

type userService struct {
	log        *zap.Logger
	jwtService auth.JWTService
	repository postgres.Repository
}

// NewUserService creates a new authentication service.
func NewUserService(logger *zap.Logger, jwt auth.JWTService, repository postgres.Repository) UserService {
	return &userService{
		log:        logger,
		jwtService: jwt,
		repository: repository,
	}
}

// RegisterUser registers a new user.
func (s *userService) RegisterUser(user *request.UserRegisterRequest) (*response.UserRegisterResponse, error) {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	userModel := models.NewUser(user)
	err = s.repository.CreateUser(userModel)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtService.GenerateToken(userModel.ID, userModel.Role)

	return response.NewUserRegisterResponse(userModel, token), nil
}

// AuthUser authenticates a user.
func (s *userService) AuthUser(user request.UserLoginRequest) (*response.UserLoginResponse, error) {
	userModel, err := s.repository.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	ok := auth.CheckPassword(userModel.Password, user.Password)
	if !ok {
		return nil, appErrors.ErrWrongPassword
	}

	token, err := s.jwtService.GenerateToken(userModel.ID, userModel.Role)
	if err != nil {
		return nil, err
	}
	return response.NewUserLoginResponse(userModel, token), nil
}
