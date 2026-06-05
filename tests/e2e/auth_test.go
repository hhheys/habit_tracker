package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func TestLogin(t *testing.T) {
	testApp := setupTestApp(t)

	const (
		username = "test-user"
		email    = "test@example.com"
		password = "password123"
		timezone = "Europe/Moscow"
	)
	registerUser(t, testApp.Router, username, email, password, timezone)

	tokens := loginUser(t, testApp.Router, email, password)
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatal("login user: response does not contain auth tokens")
	}
}

func loginUser(t *testing.T, router *gin.Engine, email, password string) AuthTokens {
	t.Helper()

	body, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		t.Fatalf("marshal login request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("login user: got status %d, body %s", res.Code, res.Body.String())
	}

	var payload struct {
		Auth struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"auth"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode login response: %v", err)
	}

	return AuthTokens{
		AccessToken:  payload.Auth.AccessToken,
		RefreshToken: payload.Auth.RefreshToken,
	}
}

func registerUser(t *testing.T, router *gin.Engine, username, email, password, timezone string) AuthTokens {
	t.Helper()

	body, err := json.Marshal(map[string]string{
		"username": username,
		"email":    email,
		"password": password,
		"timezone": timezone,
	})
	if err != nil {
		t.Fatalf("marshal register request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("register user: got status %d, body %s", res.Code, res.Body.String())
	}

	var payload struct {
		Auth struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"auth"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode register response: %v", err)
	}

	return AuthTokens{
		AccessToken:  payload.Auth.AccessToken,
		RefreshToken: payload.Auth.RefreshToken,
	}
}
