package auth

import (
	"errors"
	authuc "habit-tracker/internal/usecase/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the custom claims for JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JwtService struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTService creates a new JWTService instance
func NewJWTService(secret string, duration time.Duration) *JwtService {
	return &JwtService{
		secretKey:     []byte(secret),
		tokenDuration: duration,
	}
}

// GenerateToken creates a new JWT token
func (s *JwtService) GenerateToken(subject *authuc.TokenSubject) (string, error) {
	claims := Claims{
		UserID: subject.UserID,
		Role:   string(subject.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secretKey)
}

// ParseToken parses and validates a JWT token
func (s *JwtService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateToken validates a JWT token
func (s *JwtService) ValidateToken(tokenString string) (*Claims, error) {
	return s.ParseToken(tokenString)
}
