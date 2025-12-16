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

type SprintController struct {
    Repo *repository.SprintRepository
}

func NewSprintController(repo *repository.SprintRepository) *SprintController {
    return &SprintController{Repo: repo}
}

// CreateSprint godoc
// @Summary Create a sprint
// @Description Create a new sprint. Requires authentication.
// @Tags sprints
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param sprint body models.Sprint true "Sprint payload"
// @Success 201 {object} models.Sprint
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sprints [post]
// @Security BearerAuth
func (c *SprintController) CreateSprint(ctx *gin.Context) {
    var sprint models.Sprint
    if err := ctx.ShouldBindJSON(&sprint); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if sprint.Status == "" {
        sprint.Status = "planned"
    }
    
    if err := c.Repo.Create(&sprint); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create sprint"})
        return
    }

    services.RedisClient.Del(services.Ctx, "sprints:all")

    ctx.JSON(http.StatusCreated, sprint)
}

// GetAllSprints godoc
// @Summary List sprints
// @Description Get all sprints
// @Tags sprints
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Sprint
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sprints [get]
// @Security BearerAuth
func (c *SprintController) GetAllSprints(ctx *gin.Context) {
    cached, _ := services.RedisClient.Get(services.Ctx, "sprints:all").Result()
    if cached != "" {
        ctx.Data(200, "application/json", []byte(cached))
        return
    }

    sprints, err := c.Repo.GetAll()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sprints"})
        return
    }

    jsonData, _ := json.Marshal(sprints)
    services.RedisClient.Set(services.Ctx, "sprints:all", jsonData, 30*time.Second)

    ctx.JSON(http.StatusOK, sprints)
}

// GetSprintByID godoc
// @Summary Get sprint by ID
// @Description Get a sprint by its ID
// @Tags sprints
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Sprint ID"
// @Success 200 {object} models.Sprint
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /sprints/{id} [get]
// @Security BearerAuth
func (c *SprintController) GetSprintByID(ctx *gin.Context) {
    id, _ := strconv.Atoi(ctx.Param("id"))
    sprint, err := c.Repo.GetByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "sprint not found"})
        return
    }

    services.RedisClient.Del(services.Ctx, "sprints:all")

    ctx.JSON(http.StatusOK, sprint)
}

// UpdateSprint godoc
// @Summary Update a sprint
// @Description Update a sprint by its ID
// @Tags sprints
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Sprint ID"
// @Param sprint body map[string]interface{} true "Update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sprints/{id} [put]
// @Security BearerAuth
func (c *SprintController) UpdateSprint(ctx *gin.Context) {
    id, _ := strconv.Atoi(ctx.Param("id"))
    var data map[string]interface{}
    if err := ctx.ShouldBindJSON(&data); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.Repo.Update(uint(id), data); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
        return
    }

    services.RedisClient.Del(services.Ctx, "sprints:all")

    ctx.JSON(http.StatusOK, gin.H{"message": "sprint updated"})
}

// DeleteSprint godoc
// @Summary Delete a sprint
// @Description Delete a sprint by ID
// @Tags sprints
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Sprint ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sprints/{id} [delete]
// @Security BearerAuth
func (c *SprintController) DeleteSprint(ctx *gin.Context) {
    id, _ := strconv.Atoi(ctx.Param("id"))
    if err := c.Repo.Delete(uint(id)); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
        return
    }

    services.RedisClient.Del(services.Ctx, "sprints:all")

    ctx.JSON(http.StatusOK, gin.H{"message": "sprint deleted"})
}