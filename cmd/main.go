package main

import (
	"log"

	"github.com/kadyrbayev2005/studysync/internal/services"

	"github.com/kadyrbayev2005/studysync/internal/api"
)

func main() {
	db, err := services.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	router := api.SetupRouter(db)
	router.Run(":8080")
}
