package api

import (
	"github.com/kadyrbayev2005/studysync/internal/controllers"
	"github.com/kadyrbayev2005/studysync/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	// Initialize Gin router with default middleware (logger & recovery)
	r := gin.Default()

	// Optionally, you can add global middleware here
	// r.Use(SomeCustomMiddleware())

	// -----------------------
	// Task Routes
	// -----------------------
	taskRepo := repository.NewTaskRepository(db)
	taskController := controllers.NewTaskController(taskRepo)

	taskRoutes := r.Group("/tasks")
	{
		taskRoutes.POST("", taskController.CreateTask)
		taskRoutes.GET("", taskController.GetAllTasks)
		taskRoutes.GET("/:id", taskController.GetTaskByID)
		taskRoutes.PUT("/:id", taskController.UpdateTask)
		taskRoutes.DELETE("/:id", taskController.DeleteTask)
	}

	// -----------------------
	// Subject Routes
	// -----------------------
	subjectRepo := repository.NewSubjectRepository(db)
	subjectController := controllers.NewSubjectController(subjectRepo)

	subjectRoutes := r.Group("/subjects")
	{
		subjectRoutes.POST("", subjectController.CreateSubject)
		subjectRoutes.GET("", subjectController.GetAllSubjects)
		subjectRoutes.GET("/:id", subjectController.GetSubjectByID)
		subjectRoutes.PUT("/:id", subjectController.UpdateSubject)
		subjectRoutes.DELETE("/:id", subjectController.DeleteSubject)
	}

	return r
}
