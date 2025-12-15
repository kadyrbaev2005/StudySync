package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadyrbayev2005/studysync/internal/models"
	"github.com/kadyrbayev2005/studysync/internal/repository"
	"github.com/kadyrbayev2005/studysync/internal/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// поднимаем in-memory БД, репозиторий и контроллер
func setupUserTestServer(t *testing.T) (*UserController, *gin.Engine) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	ctrl := NewUserController(userRepo)

	// фейковый Redis, чтобы Del не падал
	services.Ctx = context.Background()
	services.RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	r := gin.New()
	auth := r.Group("/auth")
	{
		auth.POST("/register", ctrl.Register)
		auth.POST("/login", ctrl.Login)
	}

	return ctrl, r
}

func TestUserController_Register_Success(t *testing.T) {
	_, r := setupUserTestServer(t)

	body, _ := json.Marshal(map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "123456",
		"role":     "",
	})

	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUserController_Register_InvalidBody(t *testing.T) {
	_, r := setupUserTestServer(t)

	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUserController_Login_Success(t *testing.T) {
	ctrl, r := setupUserTestServer(t)

	// создаём пользователя напрямую через репозиторий
	hashed := services.HashPassword("123456")
	u := models.User{
		Name:         "Test User",
		Email:        "login@test.com",
		PasswordHash: hashed,
		Role:         services.RoleUser,
		CreatedAt:    time.Now(),
	}
	if err := ctrl.Repo.Create(&u); err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	body, _ := json.Marshal(map[string]string{
		"email":    "login@test.com",
		"password": "123456",
	})

	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp["token"] == "" {
		t.Fatalf("expected token in response")
	}
}

func TestUserController_Login_InvalidCredentials(t *testing.T) {
	_, r := setupUserTestServer(t)

	body, _ := json.Marshal(map[string]string{
		"email":    "wrong@example.com",
		"password": "wrong",
	})

	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}
