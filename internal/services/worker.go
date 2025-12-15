package services

import (
	"context"
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
			Info("Reminder worker stopped")
			return
		case <-ticker.C:
			now := time.Now()
			soon := now.Add(15 * time.Minute)

			var due []models.Deadline
			// find deadlines between now and soon
			if err := db.Preload("Task").Where("due_date > ? AND due_date <= ?", now, soon).Find(&due).Error; err != nil {
				Error("worker query error", "error", err)
				continue
			}

			for _, d := range due {
				// simple reminder action (log). Replace with real notification.
				Info("Upcoming deadline reminder", "task_title", d.Task.Title, "due_date", d.DueDate, "deadline_id", d.ID)
			}
		}
	}
}

func SendEmail(userEmail, subject, body string) any {
	panic("unimplemented")
}
