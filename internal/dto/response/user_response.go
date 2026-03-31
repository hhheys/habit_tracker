package response

import (
	"habit-tracker/internal/models"
	"time"
)

type UserRegisterResponse struct {
	Auth *AuthResponse `json:"auth"`
	User *UserResponse `json:"user"`
}

// UserResponse is the response for the user.
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
type UserLoginResponse struct {
	User *UserResponse `json:"user"`
	Auth *AuthResponse `json:"auth"`
}

// AuthResponse is the response for the auth endpoint.
type AuthResponse struct {
	Token string `json:"token"`
}

func NewUserRegisterResponse(user *models.User, authToken string) *UserRegisterResponse {
	return &UserRegisterResponse{
		Auth: NewAuthResponse(authToken),
		User: &UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}
}

func NewUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func NewUserLoginResponse(user *models.User, authToken string) *UserLoginResponse {
	return &UserLoginResponse{
		Auth: NewAuthResponse(authToken),
		User: NewUserResponse(user),
	}
}

func NewAuthResponse(token string) *AuthResponse {
	return &AuthResponse{
		Token: token,
	}
}
