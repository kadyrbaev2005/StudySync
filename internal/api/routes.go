package api

import (
	"net/http"
	_ "net/http/pprof"

	_ "github.com/kadyrbayev2005/studysync/docs"
	"github.com/kadyrbayev2005/studysync/internal/controllers"
	"github.com/kadyrbayev2005/studysync/internal/middleware"
	"github.com/kadyrbayev2005/studysync/internal/repository"
	"github.com/kadyrbayev2005/studysync/internal/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	r.Use(middleware.MetricsMiddleware())

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
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// stress test endpoint
	r.GET("/stress/start", func(c *gin.Context) {
		if GlobalStressManager != nil {
			go GlobalStressManager.StartStress() // Запускаем в горутине чтобы не блокировать ответ
			c.JSON(200, gin.H{
				"message": "CPU stress test started",
				"status":  "running",
			})
		} else {
			c.JSON(500, gin.H{"error": "Stress manager not initialized"})
		}
	})
	r.GET("/stress/stop", func(c *gin.Context) {
		if GlobalStressManager != nil {
			GlobalStressManager.StopStress()
			c.JSON(200, gin.H{
				"message": "CPU stress test stopped",
				"status":  "stopped",
			})
		} else {
			c.JSON(500, gin.H{"error": "Stress manager not initialized"})
		}
	})
	r.GET("/stress/status", func(c *gin.Context) {
		// Импортируем пакет stress для проверки статуса
		// В реальном приложении лучше использовать интерфейс
		c.JSON(200, gin.H{
			"running":   GlobalStressManager != nil,
			"cpu_cores": 0, // Будет заполнено в будущем
		})
	})

	// pprof routes
	pprofGroup := r.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/cmdline", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/profile", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.POST("/symbol", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/symbol", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/trace", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/allocs", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/block", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/goroutine", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/heap", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/mutex", gin.WrapF(http.DefaultServeMux.ServeHTTP))
		pprofGroup.GET("/threadcreate", gin.WrapF(http.DefaultServeMux.ServeHTTP))
	}

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
