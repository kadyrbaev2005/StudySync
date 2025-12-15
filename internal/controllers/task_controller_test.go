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

// поднимаем in-memory БД, TaskRepository, контроллер и роутер
func setupTaskTestServer(t *testing.T) (*repository.TaskRepository, *gin.Engine) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&models.Subject{}, &models.Task{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	taskRepo := repository.NewTaskRepository(db)
	ctrl := NewTaskController(taskRepo)

	// фейковый Redis
	services.Ctx = context.Background()
	services.RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	r := gin.New()
	tasks := r.Group("/tasks")
	{
		tasks.POST("", ctrl.CreateTask)
		tasks.GET("", ctrl.GetAllTasks)
		tasks.GET("/:id", ctrl.GetTaskByID)
		tasks.PUT("/:id", ctrl.UpdateTask)
		tasks.DELETE("/:id", ctrl.DeleteTask)
	}

	return taskRepo, r
}

func TestTaskController_CreateTask_Success(t *testing.T) {
	_, r := setupTaskTestServer(t)

	deadline := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	body, _ := json.Marshal(map[string]interface{}{
		"title":       "Test task",
		"description": "Some description",
		"status":      "todo",
		"subject_id":  0,
		"deadline":    deadline,
	})

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestTaskController_CreateTask_InvalidBody(t *testing.T) {
	_, r := setupTaskTestServer(t)

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestTaskController_GetAllTasks_Success(t *testing.T) {
	repo, r := setupTaskTestServer(t)

	// создаём несколько задач в БД
	_ = repo.Create(&models.Task{Title: "Task 1", Status: "todo"})
	_ = repo.Create(&models.Task{Title: "Task 2", Status: "in-progress"})

	req, _ := http.NewRequest(http.MethodGet, "/tasks?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data []models.Task  `json:"data"`
		Meta map[string]any `json:"meta"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(resp.Data))
	}
}

func TestTaskController_GetTaskByID_Success(t *testing.T) {
	repo, r := setupTaskTestServer(t)

	task := models.Task{
		Title:  "Single task",
		Status: "todo",
	}
	if err := repo.Create(&task); err != nil {
		t.Fatalf("failed to seed task: %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp models.Task
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Title != "Single task" {
		t.Fatalf("expected title %q, got %q", "Single task", resp.Title)
	}
}

func TestTaskController_UpdateTask_Success(t *testing.T) {
	repo, r := setupTaskTestServer(t)

	task := models.Task{Title: "Old title", Status: "todo"}
	if err := repo.Create(&task); err != nil {
		t.Fatalf("failed to seed task: %v", err)
	}

	body, _ := json.Marshal(map[string]interface{}{
		"title": "New title",
	})

	req, _ := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestTaskController_DeleteTask_Success(t *testing.T) {
	repo, r := setupTaskTestServer(t)

	task := models.Task{Title: "To delete", Status: "todo"}
	if err := repo.Create(&task); err != nil {
		t.Fatalf("failed to seed task: %v", err)
	}

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	// проверяем, что запись реально удалена
	if _, err := repo.GetByID(1); err == nil {
		t.Fatalf("expected error when getting deleted task")
	}
}
