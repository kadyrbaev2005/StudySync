package api

import (
	_ "github.com/kadyrbayev2005/studysync/docs"
	"github.com/kadyrbayev2005/studysync/internal/controllers"
	"github.com/kadyrbayev2005/studysync/internal/middleware"
	"github.com/kadyrbayev2005/studysync/internal/repository"
	"github.com/kadyrbayev2005/studysync/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// repositories
	userRepo := repository.NewUserRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	deadlineRepo := repository.NewDeadlineRepository(db)

	// controllers
	userController := controllers.NewUserController(userRepo)
	subjectController := controllers.NewSubjectController(subjectRepo)
	taskController := controllers.NewTaskController(taskRepo)
	deadlineController := controllers.NewDeadlineController(deadlineRepo, taskRepo)

	// auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	// public routes
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// protected routes: require JWT
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Users (admin only)
		users := protected.Group("/users")
		users.Use(middleware.RoleMiddleware(services.RoleAdmin))
		{
			users.GET("", userController.GetAll)
			users.GET("/:id", userController.GetByID)
			users.DELETE("/:id", userController.Delete)
		}

		// Subjects
		subjectRoutes := protected.Group("/subjects")
		{
			subjectRoutes.POST("", subjectController.CreateSubject)
			subjectRoutes.GET("", subjectController.GetAllSubjects)
			subjectRoutes.GET("/:id", subjectController.GetSubjectByID)
			subjectRoutes.PUT("/:id", subjectController.UpdateSubject)
			subjectRoutes.DELETE("/:id", subjectController.DeleteSubject)
		}

		// Tasks
		taskRoutes := protected.Group("/tasks")
		{
			taskRoutes.POST("", taskController.CreateTask)
			taskRoutes.GET("", taskController.GetAllTasks)
			taskRoutes.GET("/:id", taskController.GetTaskByID)
			taskRoutes.PUT("/:id", taskController.UpdateTask)
			taskRoutes.DELETE("/:id", taskController.DeleteTask)
		}

		// Deadlines
		deadlineRoutes := protected.Group("/deadlines")
		deadlineRoutes.Use(middleware.RoleMiddleware(services.RoleAdmin))
		{
			deadlineRoutes.POST("", deadlineController.CreateDeadline)
			deadlineRoutes.GET("", deadlineController.GetAllDeadlines)
			deadlineRoutes.GET("/:id", deadlineController.GetDeadlineByID)
			deadlineRoutes.DELETE("/:id", deadlineController.DeleteDeadline)
		}
	}

	return r
}
