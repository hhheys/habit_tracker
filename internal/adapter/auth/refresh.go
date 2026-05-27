package auth

import (
	"crypto/rand"
	"encoding/base64"
	authuc "habit-tracker/internal/usecase/auth"
)

type RefreshTokenGenerator struct{}

func NewRefreshTokenGenerator() *RefreshTokenGenerator {
	return &RefreshTokenGenerator{}
}

func (g *RefreshTokenGenerator) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (g *RefreshTokenGenerator) GenerateToken(_ *authuc.TokenSubject) (string, error) {
	return g.GenerateRefreshToken()
}
