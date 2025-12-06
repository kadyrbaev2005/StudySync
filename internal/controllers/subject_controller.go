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

type SubjectController struct {
	Repo *repository.SubjectRepository
}

func NewSubjectController(repo *repository.SubjectRepository) *SubjectController {
	return &SubjectController{Repo: repo}
}

// CreateSubject godoc
// @Summary Create a subject
// @Description Create a new subject. Requires authentication.
// @Tags subjects
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param subject body models.Subject true "Subject payload"
// @Success 201 {object} models.Subject
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /subjects [post]
// @Security BearerAuth
func (c *SubjectController) CreateSubject(ctx *gin.Context) {
	var subject models.Subject
	if err := ctx.ShouldBindJSON(&subject); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.Repo.Create(&subject); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subject"})
		return
	}

	services.RedisClient.Del(services.Ctx, "subjects:all")

	ctx.JSON(http.StatusCreated, subject)
}

// GetAllSubjects godoc
// @Summary List subjects
// @Description Get all subjects
// @Tags subjects
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Subject
// @Failure 401 {object} map[string]string
// @Router /subjects [get]
// @Security BearerAuth
func (c *SubjectController) GetAllSubjects(ctx *gin.Context) {
	cached, _ := services.RedisClient.Get(services.Ctx, "subjects:all").Result()
	if cached != "" {
		ctx.Data(200, "application/json", []byte(cached))
		return
	}

	subjects, err := c.Repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subjects"})
		return
	}

	jsonData, _ := json.Marshal(subjects)
	services.RedisClient.Set(services.Ctx, "subjects:all", jsonData, 30*time.Second)

	ctx.JSON(http.StatusOK, subjects)
}

// GetSubjectByID godoc
// @Summary Get subject by ID
// @Description Get a subject by its ID
// @Tags subjects
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Subject ID"
// @Success 200 {object} models.Subject
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subjects/{id} [get]
// @Security BearerAuth
func (c *SubjectController) GetSubjectByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	subject, err := c.Repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
		return
	}

	services.RedisClient.Del(services.Ctx, "subjects:all")

	ctx.JSON(http.StatusOK, subject)
}

// UpdateSubject godoc
// @Summary Update a subject
// @Description Update a subject by its ID
// @Tags subjects
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Subject ID"
// @Param subject body map[string]interface{} true "Update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subjects/{id} [put]
// @Security BearerAuth
func (c *SubjectController) UpdateSubject(ctx *gin.Context) {
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

	services.RedisClient.Del(services.Ctx, "subjects:all")

	ctx.JSON(http.StatusOK, gin.H{"message": "subject updated"})
}

// DeleteSubject godoc
// @Summary Delete a subject
// @Description Delete a subject by ID
// @Tags subjects
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Subject ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subjects/{id} [delete]
// @Security BearerAuth
func (c *SubjectController) DeleteSubject(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.Repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	services.RedisClient.Del(services.Ctx, "subjects:all")

	ctx.JSON(http.StatusOK, gin.H{"message": "subject deleted"})
}
