package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
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
// @Summary      List all tasks
// @Description  Retrieves all tasks. Requires authentication.
// @Tags         tasks
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200 {array} models.Task
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks [get]
// @Security     BearerAuth
func (c *TaskController) GetAllTasks(ctx *gin.Context) {
	//try to get from cache
    cached, _ := services.RedisClient.Get(services.Ctx, "tasks:all").Result()
    if cached != "" {
        ctx.Data(200, "application/json", []byte(cached))
        return
    }

	//otherwise, get from db
    tasks, err := c.Repo.GetAll()
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }

    jsonData, _ := json.Marshal(tasks)

    //put in redis
    services.RedisClient.Set(services.Ctx, "tasks:all", jsonData, 30*time.Second)

    ctx.JSON(200, tasks)
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	services.RedisClient.Del(services.Ctx, "tasks:all")

	ctx.JSON(http.StatusOK, task)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Repo.Update(uint(id), data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	services.RedisClient.Del(services.Ctx, "tasks:all")

	ctx.JSON(http.StatusOK, gin.H{"message": "task updated"})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	services.RedisClient.Del(services.Ctx, "tasks:all")

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
