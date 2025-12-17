package main

import (
	"context"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/kadyrbayev2005/studysync/docs"

	"github.com/joho/godotenv"
	"github.com/kadyrbayev2005/studysync/internal/api"
	"github.com/kadyrbayev2005/studysync/internal/services"
	"github.com/kadyrbayev2005/studysync/internal/stress"
	"github.com/kadyrbayev2005/studysync/internal/utils"
)

// StressManagerImpl реализация интерфейса StressManager
type StressManagerImpl struct{}

func (s *StressManagerImpl) StartStress() {
	stress.StartStress()
}

func (s *StressManagerImpl) StopStress() {
	stress.StopStress()
}

func main() {
	// Флаг для включения нагрузки на CPU
	var debugMode = flag.Bool("debug", false, "Enable CPU stress test for profiling")
	flag.Parse()

	// Запуск нагрузки если указан флаг --debug
	if *debugMode {
		stress.StartStress()
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Это вызовет панику если Logger еще не инициализирован
		// Лучше использовать стандартный вывод для ошибок до инициализации логгера
		println("No .env file found, using system environment variables")
	}

	// Инициализация логгера - ДОЛЖНА БЫТЬ ПЕРВОЙ!
	services.InitLogger()

	services.InitJWT()

	// Теперь можем логировать
	services.Debug("Logger initialized successfully")
	services.Info("Starting StudySync application")

	// Подключение к БД
	db, err := services.ConnectDB()
	if err != nil {
		services.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	services.Info("Database connected successfully")

	services.InitRedis()
	services.InitSMTP() // Добавляем инициализацию SMTP

	// Запуск worker для напоминаний
	workerCtx, workerCancel := context.WithCancel(context.Background())
	go services.StartReminderWorker(workerCtx, db)
	services.Debug("Reminder worker started")

	// Настройка роутера
	router := api.SetupRouter(db)

	// Устанавливаем менеджер стресс-теста
	stressManager := &StressManagerImpl{}
	api.SetStressManager(stressManager)

	port := utils.GetEnv("SERVER_PORT", "8081") // Changed to 8081 to avoid conflicts
	addr := "0.0.0.0:" + port
	srv := api.NewServer(router, addr)

	go func() {
		if err := srv.Run(); err != nil {
			services.Error("Server run error", "error", err)
			os.Exit(1)
		}
	}()
	services.Info("Server started", "port", port)

	// Запуск pprof сервера на отдельном порту
	pprofPort := "6060"
	go func() {
		services.Info("Starting pprof server", "port", pprofPort)
		if err := http.ListenAndServe("127.0.0.1:"+pprofPort, nil); err != nil {
			services.Error("Pprof server error", "error", err)
		}
	}()
	services.Info("Pprof server started", "port", pprofPort)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	services.Info("Shutting down server...")

	// Останавливаем нагрузку если она запущена
	stress.StopStress()
	stress.GlobalStressWG.Wait() // Ждем завершения всех горутин нагрузки

	workerCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		services.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	services.Info("Server stopped gracefully")
}
