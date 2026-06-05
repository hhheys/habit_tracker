package auth

import (
	"habit-tracker/internal/domain"
	"habit-tracker/internal/usecase/user"
)

type RegisterInput struct {
	Username string
	Email    string
	Password string
	Timezone string
}

type SessionInfoInput struct {
	UserIP    string
	UserAgent string
}

type RegisterOutput struct {
	User          user.Output
	Authorization TokenOutput
}

type LoginInput struct {
	Username string
	Email    string
	Password string
}

type LoginOutput struct {
	User   user.Output
	Tokens TokenOutput
}

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	Tokens TokenOutput
}

type TokenSubject struct {
	UserID uint
	Role   domain.UserRole
}

type TokenOutput struct {
	AccessToken  string
	RefreshToken string
}
