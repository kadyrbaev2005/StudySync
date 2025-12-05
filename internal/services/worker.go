package services

import (
	"context"
	"fmt"
	"time"

	"github.com/kadyrbayev2005/studysync/internal/models"
	"gorm.io/gorm"
)

func StartReminderWorker(ctx context.Context, db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Reminder worker stopped")
			return
		case <-ticker.C:
			now := time.Now()
			soon := now.Add(15 * time.Minute)

			var due []models.Deadline

			if err := db.Preload("Task").Preload("User").
				Where("due_date > ? AND due_date <= ?", now, soon).
				Find(&due).Error; err != nil {

				fmt.Println("worker query error:", err)
				continue
			}

			for _, d := range due {
				userEmail := d.User.Email
				if userEmail == "" {
					continue
				}

				subject := "StudySync Reminder: Upcoming Deadline"
				body := fmt.Sprintf(
					"Task '%s' is due at %s\nDescription: %s",
					d.Task.Title,
					d.DueDate.Format(time.RFC3339),
					d.Task.Description,
				)

				err := SendEmail(userEmail, subject, body)
				if err != nil {
					fmt.Println("Failed to send email:", err)
				} else {
					fmt.Println("Email sent to", userEmail)
				}
			}
		}
	}
}
