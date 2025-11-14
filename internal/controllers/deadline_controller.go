package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadyrbayev2005/studysync/internal/models"
	"github.com/kadyrbayev2005/studysync/internal/repository"
)

type DeadlineController struct {
	Repo     *repository.DeadlineRepository
	TaskRepo *repository.TaskRepository
}

func NewDeadlineController(repo *repository.DeadlineRepository, taskRepo *repository.TaskRepository) *DeadlineController {
	return &DeadlineController{Repo: repo, TaskRepo: taskRepo}
}

func (c *DeadlineController) CreateDeadline(ctx *gin.Context) {
	var payload struct {
		TaskID  uint      `json:"task_id" binding:"required"`
		DueDate time.Time `json:"due_date" binding:"required"`
	}
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
	ctx.JSON(http.StatusCreated, d)
}

func (c *DeadlineController) GetAllDeadlines(ctx *gin.Context) {
	deadlines, err := c.Repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch deadlines"})
		return
	}
	ctx.JSON(http.StatusOK, deadlines)
}

func (c *DeadlineController) GetDeadlineByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	d, err := c.Repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "deadline not found"})
		return
	}
	ctx.JSON(http.StatusOK, d)
}

func (c *DeadlineController) DeleteDeadline(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.Repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete deadline"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deadline deleted"})
}
