package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

// HashString returns a hashed string
func (ph *Hasher) HashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CompareStrings compares a hashed string with a non-hashed string
func (ph *Hasher) CompareStrings(hashedString, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(str))
	return err == nil
}

func (ph *Hasher) HashPassword(password string) (string, error) {
	return ph.HashString(password)
}

func (ph *Hasher) ComparePasswordHash(hashedPassword, password string) bool {
	return ph.CompareStrings(hashedPassword, password)
}

func (ph *Hasher) HashToken(token string) (string, error) {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:]), nil
}

func (ph *Hasher) CompareToken(hashedToken, token string) bool {
	actual, _ := ph.HashToken(token)
	return hashedToken == actual
}
