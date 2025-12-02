package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kadyrbayev2005/studysync/internal/api"
	"github.com/kadyrbayev2005/studysync/internal/services"
)

func main() {
	db, err := services.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	workerCtx, workerCancel := context.WithCancel(context.Background())
	go services.StartReminderWorker(workerCtx, db)

	router := api.SetupRouter(db)

	srv := api.NewServer(router, ":8080")

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("server run error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	workerCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped gracefully")
}
