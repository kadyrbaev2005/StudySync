package services

import (
	"fmt"

	"github.com/kadyrbayev2005/studysync/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "user=postgres password=postgres dbname=studysync port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	db.AutoMigrate(&models.User{}, &models.Subject{}, &models.Task{}, &models.Deadline{})
	fmt.Println("✅ Connected to database and migrated successfully")

	return db, nil
}
