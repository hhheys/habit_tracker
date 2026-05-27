package response

import (
	authuc "habit-tracker/internal/usecase/auth"
	useruc "habit-tracker/internal/usecase/user"
	"time"
)

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserRegisterResponse struct {
	Auth AuthResponse `json:"auth"`
	User UserResponse `json:"user"`
}

type UserLoginResponse struct {
	User UserResponse `json:"user"`
	Auth AuthResponse `json:"auth"`
}

func NewUserRegisterResponse(output *authuc.RegisterOutput) *UserRegisterResponse {
	return &UserRegisterResponse{
		Auth: NewAuthResponse(output.Authorization),
		User: NewUserResponse(output.User),
	}
}

func NewUserLoginResponse(output *authuc.LoginOutput) *UserLoginResponse {
	return &UserLoginResponse{
		Auth: NewAuthResponse(output.Tokens),
		User: NewUserResponse(output.User),
	}
}

func NewAuthResponse(tokens authuc.TokenOutput) AuthResponse {
	return AuthResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}
}

func NewUserResponse(user useruc.Output) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}
}
