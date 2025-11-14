package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadyrbayev2005/studysync/internal/models"
	"github.com/kadyrbayev2005/studysync/internal/repository"
	"github.com/kadyrbayev2005/studysync/internal/services"
)

type UserController struct {
	Repo *repository.UserRepository
}

func NewUserController(repo *repository.UserRepository) *UserController {
	return &UserController{Repo: repo}
}

type registerPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type loginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (c *UserController) Register(ctx *gin.Context) {
	var p registerPayload
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed := services.HashPassword(p.Password)
	user := models.User{
		Name:         p.Name,
		Email:        p.Email,
		PasswordHash: hashed,
		Role:         p.Role,
		CreatedAt:    time.Now(),
	}

	if user.Role == "" {
		user.Role = services.RoleUser
	}

	if err := c.Repo.Create(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// hide sensitive
	user.PasswordHash = ""
	ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) Login(ctx *gin.Context) {
	var p loginPayload
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.Repo.GetByEmail(p.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !services.CheckPasswordHash(p.Password, user.PasswordHash) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := services.GenerateJWT(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.Repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	// strip password hashes
	for i := range users {
		users[i].PasswordHash = ""
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	user, err := c.Repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	user.PasswordHash = ""
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.Repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
