package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type habitResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

func TestCreateHabitPermissions(t *testing.T) {
	testApp := setupTestApp(t)

	userTokens := registerUser(t, testApp.Router, "regular-user", "regular@example.com", "password123")
	forbidden := createHabitRequest(t, testApp.Router, userTokens.AccessToken, "Forbidden habit", "Not allowed")
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("create habit as regular user: got status %d, body %s", forbidden.Code, forbidden.Body.String())
	}

	var count int
	if err := testApp.DB.QueryRow(`SELECT COUNT(*) FROM habit WHERE title = $1`, "Forbidden habit").Scan(&count); err != nil {
		t.Fatalf("count forbidden habit: %v", err)
	}
	if count != 0 {
		t.Fatalf("create habit as regular user: created %d records", count)
	}
}

func TestCreateHabit(t *testing.T) {
	testApp := setupTestApp(t)
	adminTokens := registerAdmin(t, testApp, "admin-user", "admin@example.com", "password123")

	created := createHabit(t, testApp.Router, adminTokens.AccessToken, "Drink water", "Drink enough water")

	if created.Title != "Drink water" || created.Description != "Drink enough water" {
		t.Fatalf("create habit as admin: unexpected response %+v", created)
	}
	if created.ID == 0 || created.ImageURL == "" {
		t.Fatalf("create habit as admin: missing generated fields %+v", created)
	}

	var title, description, imageFilename string
	if err := testApp.DB.QueryRow(
		`SELECT title, description, image_filename FROM habit WHERE id = $1`,
		created.ID,
	).Scan(&title, &description, &imageFilename); err != nil {
		t.Fatalf("get created habit from database: %v", err)
	}
	if title != created.Title || description != created.Description || imageFilename == "" {
		t.Fatalf("create habit as admin: unexpected database values %q, %q, %q", title, description, imageFilename)
	}
}

func TestGetAllHabits(t *testing.T) {
	testApp := setupTestApp(t)
	adminTokens := registerAdmin(t, testApp, "list-admin", "list-admin@example.com", "password123")
	userTokens := registerUser(t, testApp.Router, "list-user", "list-user@example.com", "password123")

	createHabit(t, testApp.Router, adminTokens.AccessToken, "Meditate", "Meditate daily")
	createHabit(t, testApp.Router, adminTokens.AccessToken, "Read", "Read daily")

	req := httptest.NewRequest(http.MethodGet, "/api/habit?sort_by=title&sort_order=asc", nil)
	req.Header.Set("Authorization", "Bearer "+userTokens.AccessToken)
	res := httptest.NewRecorder()

	testApp.Router.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("get all habits: got status %d, body %s", res.Code, res.Body.String())
	}

	var payload struct {
		Habits     []habitResponse `json:"habits"`
		Pagination struct {
			Total int `json:"total"`
		} `json:"pagination"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode all habits response: %v", err)
	}
	if payload.Pagination.Total != 2 || len(payload.Habits) != 2 {
		t.Fatalf("get all habits: got total %d and %d items", payload.Pagination.Total, len(payload.Habits))
	}
	if payload.Habits[0].Title != "Meditate" || payload.Habits[1].Title != "Read" {
		t.Fatalf("get all habits: unexpected order %q, %q", payload.Habits[0].Title, payload.Habits[1].Title)
	}
}

func registerAdmin(t *testing.T, testApp *TestApp, username, email, password string) AuthTokens {
	t.Helper()

	registerUser(t, testApp.Router, username, email, password)
	if _, err := testApp.DB.Exec(`UPDATE users SET role = 'admin' WHERE email = $1`, email); err != nil {
		t.Fatalf("promote user to admin: %v", err)
	}

	return loginUser(t, testApp.Router, email, password)
}

func createHabit(t *testing.T, router *gin.Engine, accessToken, title, description string) habitResponse {
	t.Helper()

	res := createHabitRequest(t, router, accessToken, title, description)
	if res.Code != http.StatusCreated {
		t.Fatalf("create habit: got status %d, body %s", res.Code, res.Body.String())
	}

	var habit habitResponse
	if err := json.Unmarshal(res.Body.Bytes(), &habit); err != nil {
		t.Fatalf("decode create habit response: %v", err)
	}
	t.Cleanup(func() {
		imagePath := filepath.FromSlash("." + strings.TrimPrefix(habit.ImageURL, "."))
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			t.Errorf("remove habit image: %v", err)
		}
	})
	return habit
}

func createHabitRequest(t *testing.T, router *gin.Engine, accessToken, title, description string) *httptest.ResponseRecorder {
	t.Helper()

	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	if err := form.WriteField("title", title); err != nil {
		t.Fatalf("write habit title: %v", err)
	}
	if err := form.WriteField("description", description); err != nil {
		t.Fatalf("write habit description: %v", err)
	}
	image, err := form.CreateFormFile("image", "habit.png")
	if err != nil {
		t.Fatalf("create habit image field: %v", err)
	}
	if _, err := fmt.Fprint(image, "test image content"); err != nil {
		t.Fatalf("write habit image: %v", err)
	}
	if err := form.Close(); err != nil {
		t.Fatalf("close habit form: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/habit", &body)
	req.Header.Set("Content-Type", form.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}
