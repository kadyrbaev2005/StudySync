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
    r := gin.New()

    r.Use(middleware.GinLogger())
    r.Use(gin.Recovery())
    r.Use(middleware.CORSMiddleware())

    // repositories
    userRepo := repository.NewUserRepository(db)
    subjectRepo := repository.NewSubjectRepository(db)
    taskRepo := repository.NewTaskRepository(db)
    deadlineRepo := repository.NewDeadlineRepository(db)
    sprintRepo := repository.NewSprintRepository(db)

    // controllers
    userController := controllers.NewUserController(userRepo)
    subjectController := controllers.NewSubjectController(subjectRepo)
    taskController := controllers.NewTaskController(taskRepo)
    deadlineController := controllers.NewDeadlineController(deadlineRepo, taskRepo)
    sprintController := controllers.NewSprintController(sprintRepo)

    // public routes
    r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // auth routes
    auth := r.Group("/auth")
    {
        auth.POST("/register", userController.Register)
        auth.POST("/login", userController.Login)
    }

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
            users.PUT("/:id", userController.Update)
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
        {
            deadlineRoutes.POST("", deadlineController.CreateDeadline)
            deadlineRoutes.GET("", deadlineController.GetAllDeadlines)
            deadlineRoutes.GET("/:id", deadlineController.GetDeadlineByID)
            deadlineRoutes.PUT("/:id", deadlineController.UpdateDeadline)
            deadlineRoutes.DELETE("/:id", deadlineController.DeleteDeadline)
        }

        // Sprints
        sprintRoutes := protected.Group("/sprints")
        {
            sprintRoutes.POST("", sprintController.CreateSprint)
            sprintRoutes.GET("", sprintController.GetAllSprints)
            sprintRoutes.GET("/:id", sprintController.GetSprintByID)
            sprintRoutes.PUT("/:id", sprintController.UpdateSprint)
            sprintRoutes.DELETE("/:id", sprintController.DeleteSprint)
        }
    }

    return r
}