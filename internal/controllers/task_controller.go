package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kadyrbayev2005/studysync/internal/models"
	"github.com/kadyrbayev2005/studysync/internal/repository"
	"github.com/kadyrbayev2005/studysync/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	Repo *repository.TaskRepository
}

func NewTaskController(repo *repository.TaskRepository) *TaskController {
	return &TaskController{Repo: repo}
}

// CreateTask godoc
// @Summary      Create a new task
// @Description  Creates a new task. Requires authentication.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        task body models.Task true "Task payload"
// @Success      201 {object} models.Task
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks [post]
// @Security     BearerAuth
func (c *TaskController) CreateTask(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.Repo.Create(&task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	services.RedisClient.Del(services.Ctx, "tasks:all")

	ctx.JSON(http.StatusCreated, task)
}

// GetAllTasks godoc
// @Summary      List tasks with pagination, filtering and sorting
// @Description  Returns a paginated list of tasks with optional filters: status, subject_id, search, date range and sorting.
// @Tags         tasks
// @Produce      json
// @Param        Authorization   header   string  true   "Bearer token"
// @Param        page            query    int     false  "Page number (default: 1)"
// @Param        limit           query    int     false  "Items per page (default: 10)"
// @Param        status          query    string  false  "Filter by status (todo | in-progress | done)"
// @Param        subject_id      query    int     false  "Filter by subject ID"
// @Param        search          query    string  false  "Search text in title or description"
// @Param        sort            query    string  false  "Sort by field (created_at, deadline, title) with optional 'desc'. Example: 'deadline desc'"
// @Param        deadline_before query    string  false  "Return tasks with deadline before this timestamp (RFC3339 format)"
// @Param        deadline_after  query    string  false  "Return tasks with deadline after this timestamp (RFC3339 format)"
// @Success      200 {object} map[string]interface{} "Paginated response: data + meta"
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks [get]
// @Security     BearerAuth
func (c *TaskController) GetAllTasks(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	status := strings.TrimSpace(ctx.Query("status"))
	subjectIDStr := strings.TrimSpace(ctx.Query("subject_id"))
	search := strings.TrimSpace(ctx.Query("search"))
	sort := strings.TrimSpace(ctx.Query("sort"))
	deadlineBeforeStr := strings.TrimSpace(ctx.Query("deadline_before"))
	deadlineAfterStr := strings.TrimSpace(ctx.Query("deadline_after"))

	// Detect if request has ANY filters
	hasFilters := status != "" ||
		subjectIDStr != "" ||
		search != "" ||
		sort != "" ||
		deadlineBeforeStr != "" ||
		deadlineAfterStr != "" ||
		page != 1 ||
		limit != 10

	// Only cache if no filters at all
	if !hasFilters {
		cached, _ := services.RedisClient.Get(services.Ctx, "tasks:all").Result()
		if cached != "" {
			ctx.Data(200, "application/json", []byte(cached))
			return
		}
	}

	// --- Parse advanced filter options ---
	var subjectID *uint
	if subjectIDStr != "" {
		if num, err := strconv.Atoi(subjectIDStr); err == nil && num > 0 {
			tmp := uint(num)
			subjectID = &tmp
		} else {
			ctx.JSON(400, gin.H{"error": "invalid subject_id"})
			return
		}
	}

	var deadlineBefore *time.Time
	if deadlineBeforeStr != "" {
		t, err := time.Parse(time.RFC3339, deadlineBeforeStr)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "deadline_before must be RFC3339"})
			return
		}
		deadlineBefore = &t
	}

	var deadlineAfter *time.Time
	if deadlineAfterStr != "" {
		t, err := time.Parse(time.RFC3339, deadlineAfterStr)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "deadline_after must be RFC3339"})
			return
		}
		deadlineAfter = &t
	}

	filter := &repository.TaskFilter{
		Page:           page,
		Limit:          limit,
		Status:         status,
		SubjectID:      subjectID,
		Search:         search,
		Sort:           sort,
		DeadlineBefore: deadlineBefore,
		DeadlineAfter:  deadlineAfter,
	}

	tasks, total, err := c.Repo.GetTasks(filter)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch tasks"})
		return
	}

	resp := gin.H{
		"data": tasks,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	// Cache only unfiltered result
	if !hasFilters {
		jsonData, _ := json.Marshal(resp)
		services.RedisClient.Set(services.Ctx, "tasks:all", jsonData, 30*time.Second)
	}

	ctx.JSON(200, resp)
}

// GetTaskByID godoc
// @Summary      Get a task by ID
// @Description  Retrieves a task by its ID. Requires authentication.
// @Tags         tasks
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path int true "Task ID"
// @Success      200 {object} models.Task
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /tasks/{id} [get]
// @Security     BearerAuth
func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	task, err := c.Repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(404, gin.H{"error": "task not found"})
		return
	}

	services.RedisClient.Del(services.Ctx, "tasks:all")

	ctx.JSON(200, task)
}

// UpdateTask godoc
// @Summary      Update a task
// @Description  Updates fields of a task. Requires authentication.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path int true "Task ID"
// @Param        data body map[string]interface{} true "Task fields to update"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [put]
// @Security     BearerAuth
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var data map[string]interface{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := c.Repo.Update(uint(id), data); err != nil {
		ctx.JSON(500, gin.H{"error": "failed to update task"})
		return
	}

	services.RedisClient.Del(services.Ctx, "tasks:all")

	ctx.JSON(200, gin.H{"message": "task updated"})
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Deletes a task by ID. Requires authentication.
// @Tags         tasks
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path int true "Task ID"
// @Success      200 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [delete]
// @Security     BearerAuth
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.Repo.Delete(uint(id)); err != nil {
		ctx.JSON(500, gin.H{"error": "delete failed"})
		return
	}

	services.RedisClient.Del(services.Ctx, "tasks:all")

	ctx.JSON(200, gin.H{"message": "deleted"})
}
