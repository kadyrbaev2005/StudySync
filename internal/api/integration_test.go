package api

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/kadyrbayev2005/studysync/internal/models"
    "github.com/kadyrbayev2005/studysync/internal/services"
    "github.com/redis/go-redis/v9"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupIntegrationDB(t *testing.T) (*gorm.DB, func()) {
    t.Helper()

    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open test db: %v", err)
    }

    if err := db.AutoMigrate(
        &models.User{},
        &models.Subject{},
        &models.Task{},
        &models.Deadline{},
        &models.Sprint{},
    ); err != nil {
        t.Fatalf("failed to migrate test db: %v", err)
    }

    // Seed test data
    user := models.User{
        Name:         "Test User",
        Email:        "test@example.com",
        PasswordHash: services.HashPassword("password123"),
        Role:         services.RoleUser,
        CreatedAt:    time.Now(),
    }
    db.Create(&user)

    subject := models.Subject{
        Name:        "Mathematics",
        Description: "Advanced Calculus",
        CreatedAt:   time.Now(),
    }
    db.Create(&subject)

    return db, func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }
}

func TestIntegration_CompleteFlow(t *testing.T) {
    db, cleanup := setupIntegrationDB(t)
    defer cleanup()

    gin.SetMode(gin.TestMode)

    // Setup Redis mock
    services.Ctx = context.Background()
    services.RedisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    router := SetupRouter(db)

    // 1. Register new user
    registerBody, _ := json.Marshal(map[string]string{
        "name":     "Integration User",
        "email":    "integration@test.com",
        "password": "securepass123",
    })

    req := httptest.NewRequest(http.MethodPost, "/auth/register",
        bytes.NewReader(registerBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201 on register, got %d: %s", w.Code, w.Body.String())
    }

    // 2. Login to get token
    loginBody, _ := json.Marshal(map[string]string{
        "email":    "integration@test.com",
        "password": "securepass123",
    })

    req = httptest.NewRequest(http.MethodPost, "/auth/login",
        bytes.NewReader(loginBody))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200 on login, got %d: %s", w.Code, w.Body.String())
    }

    var loginResp map[string]string
    json.Unmarshal(w.Body.Bytes(), &loginResp)
    token := loginResp["token"]

    // 3. Create subject
    subjectBody, _ := json.Marshal(map[string]string{
        "name":        "Physics",
        "description": "Quantum Mechanics",
    })

    req = httptest.NewRequest(http.MethodPost, "/subjects",
        bytes.NewReader(subjectBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201 on subject creation, got %d: %s", w.Code, w.Body.String())
    }

    // 4. Create task
    taskBody, _ := json.Marshal(map[string]interface{}{
        "title":       "Integration Test Task",
        "description": "Test task for integration",
        "status":      "todo",
        "subject_id":  1,
        "deadline":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
    })

    req = httptest.NewRequest(http.MethodPost, "/tasks",
        bytes.NewReader(taskBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201 on task creation, got %d: %s", w.Code, w.Body.String())
    }

    // 5. Get tasks with pagination
    req = httptest.NewRequest(http.MethodGet, "/tasks?page=1&limit=10", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200 on get tasks, got %d: %s", w.Code, w.Body.String())
    }

    // 6. Clean up (optional)
    // Verify all operations were successful
    var taskCount int64
    db.Model(&models.Task{}).Count(&taskCount)
    if taskCount == 0 {
        t.Fatalf("expected at least one task in database")
    }

    fmt.Println("✅ All integration tests passed")
}