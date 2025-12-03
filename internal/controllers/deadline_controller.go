package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadyrbayev2005/studysync/internal/models"
	"github.com/kadyrbayev2005/studysync/internal/repository"
	"github.com/kadyrbayev2005/studysync/internal/services"
)

type DeadlineController struct {
	Repo     *repository.DeadlineRepository
	TaskRepo *repository.TaskRepository
}

func NewDeadlineController(repo *repository.DeadlineRepository, taskRepo *repository.TaskRepository) *DeadlineController {
	return &DeadlineController{Repo: repo, TaskRepo: taskRepo}
}

// CreateDeadline godoc
// @Summary Create a deadline
// @Description Create a new deadline for a task. Requires authentication.
// @Tags deadlines
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param deadline body struct{TaskID uint "json:\"task_id\""; DueDate time.Time "json:\"due_date\""} true "Deadline payload"
// @Success 201 {object} models.Deadline
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /deadlines [post]
// @Security BearerAuth

// @Param deadline body models.DeadlineRequest true "Deadline payload"
func (c *DeadlineController) CreateDeadline(ctx *gin.Context) {
	var payload models.DeadlineRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify task exists
	if _, err := c.TaskRepo.GetByID(payload.TaskID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "task not found"})
		return
	}

	d := models.Deadline{
		TaskID:    payload.TaskID,
		DueDate:   payload.DueDate,
		CreatedAt: time.Now(),
	}
	if err := c.Repo.Create(&d); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create deadline"})
		return
	}

	services.RedisClient.Del(services.Ctx, "deadlines:all")

	ctx.JSON(http.StatusCreated, d)
}

// GetAllDeadlines godoc
// @Summary List deadlines
// @Description Get all deadlines
// @Tags deadlines
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Deadline
// @Failure 401 {object} map[string]string
// @Router /deadlines [get]
// @Security BearerAuth
func (c *DeadlineController) GetAllDeadlines(ctx *gin.Context) {
	cached, _ := services.RedisClient.Get(services.Ctx, "deadlines:all").Result()
    if cached != "" {
        ctx.Data(200, "application/json", []byte(cached))
        return
    }

	deadlines, err := c.Repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch deadlines"})
		return
	}

	jsonData, _ := json.Marshal(deadlines)
    services.RedisClient.Set(services.Ctx, "deadlines:all", jsonData, 30*time.Second)
	
	ctx.JSON(http.StatusOK, deadlines)
}

// GetDeadlineByID godoc
// @Summary Get deadline by ID
// @Description Get a deadline by its ID
// @Tags deadlines
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Deadline ID"
// @Success 200 {object} models.Deadline
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /deadlines/{id} [get]
// @Security BearerAuth
func (c *DeadlineController) GetDeadlineByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	d, err := c.Repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "deadline not found"})
		return
	}

	services.RedisClient.Del(services.Ctx, "deadlines:all")

	ctx.JSON(http.StatusOK, d)
}

// DeleteDeadline godoc
// @Summary Delete a deadline
// @Description Delete a deadline by ID
// @Tags deadlines
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Deadline ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /deadlines/{id} [delete]
// @Security BearerAuth
func (c *DeadlineController) DeleteDeadline(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.Repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete deadline"})
		return
	}

	services.RedisClient.Del(services.Ctx, "deadlines:all")

	ctx.JSON(http.StatusOK, gin.H{"message": "deadline deleted"})
}
