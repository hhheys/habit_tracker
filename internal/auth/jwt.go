package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the custom claims for JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService provides JWT token generation and validation
//
//go:generate mockery --name=JWTService --output=../../mocks --outpkg=mocks
type JWTService interface {
	GenerateToken(userID uint, role string) (string, error)
	ParseToken(tokenString string) (*Claims, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type jwtServiceImpl struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTService creates a new JWTService instance
func NewJWTService(secret string, duration time.Duration) JWTService {
	return &jwtServiceImpl{
		secretKey:     []byte(secret),
		tokenDuration: duration,
	}
}

// GenerateToken creates a new JWT token
func (s *jwtServiceImpl) GenerateToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secretKey)
}

// ParseToken parses and validates a JWT token
func (s *jwtServiceImpl) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
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
func (s *jwtServiceImpl) ValidateToken(tokenString string) (*Claims, error) {
	return s.ParseToken(tokenString)
}
