package services

import (
	"context"
	"fmt"
	"time"

	"github.com/kadyrbayev2005/studysync/internal/models"
	"gorm.io/gorm"
)

// StartReminderWorker - background worker that periodically scans deadlines and prints reminders.
// In real app replace printing with sending emails/notifications.
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
			soon := now.Add(15 * time.Minute) // example: deadlines in next 15 min
			var due []models.Deadline
			// find deadlines between now and soon
			if err := db.Preload("Task").Where("due_date > ? AND due_date <= ?", now, soon).Find(&due).Error; err != nil {
				fmt.Println("worker query error:", err)
				continue
			}
			for _, d := range due {
				// simple reminder action (log). Replace with real notification.
				fmt.Printf("Reminder: task '%s' is due at %s (deadline id=%d)\n", d.Task.Title, d.DueDate.Format(time.RFC3339), d.ID)
			}
		}
	}
}
