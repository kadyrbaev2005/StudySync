package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kadyrbayev2005/studysync/internal/models"
	"github.com/kadyrbayev2005/studysync/internal/repository"
)

type SubjectController struct {
	Repo *repository.SubjectRepository
}

func NewSubjectController(repo *repository.SubjectRepository) *SubjectController {
	return &SubjectController{Repo: repo}
}

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
	ctx.JSON(http.StatusCreated, subject)
}

func (c *SubjectController) GetAllSubjects(ctx *gin.Context) {
	subjects, err := c.Repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subjects"})
		return
	}
	ctx.JSON(http.StatusOK, subjects)
}

func (c *SubjectController) GetSubjectByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	subject, err := c.Repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
		return
	}
	ctx.JSON(http.StatusOK, subject)
}

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
	ctx.JSON(http.StatusOK, gin.H{"message": "subject updated"})
}

func (c *SubjectController) DeleteSubject(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.Repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "subject deleted"})
}
