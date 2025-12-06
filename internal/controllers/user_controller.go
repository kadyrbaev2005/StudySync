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

type UserController struct {
	Repo *repository.UserRepository
}

func NewUserController(repo *repository.UserRepository) *UserController {
	return &UserController{Repo: repo}
}

type registerPayload struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type loginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user. Role is optional; defaults to "user".
// @Tags users
// @Accept json
// @Produce json
// @Param user body registerPayload true "User payload"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var p registerPayload
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed := services.HashPassword(p.Password)
	user := models.User{
		Name:         p.Name,
		Username:     p.Username,
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

	user.PasswordHash = ""

	services.RedisClient.Del(services.Ctx, "users:all")

	ctx.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body loginPayload true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var p loginPayload
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.Repo.GetByUsername(p.Username)
	if err != nil || !services.CheckPasswordHash(p.Password, user.PasswordHash) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := services.GenerateJWT(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	services.RedisClient.Del(services.Ctx, "users:all")

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// GetAll godoc
// @Summary List all users
// @Description Returns all users (admin only)
// @Tags users
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [get]
// @Security BearerAuth
func (c *UserController) GetAll(ctx *gin.Context) {
	cached, _ := services.RedisClient.Get(services.Ctx, "users:all").Result()
	if cached != "" {
		ctx.Data(200, "application/json", []byte(cached))
		return
	}

	users, err := c.Repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	for i := range users {
		users[i].PasswordHash = ""
	}

	jsonData, _ := json.Marshal(users)
	services.RedisClient.Set(services.Ctx, "users:all", jsonData, 30*time.Second)

	ctx.JSON(http.StatusOK, users)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Returns user by ID
// @Tags users
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
// @Security BearerAuth
func (c *UserController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	user, err := c.Repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	user.PasswordHash = ""

	services.RedisClient.Del(services.Ctx, "users:all")

	ctx.JSON(http.StatusOK, user)
}

// Delete godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
// @Security BearerAuth
func (c *UserController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.Repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	services.RedisClient.Del(services.Ctx, "users:all")

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
