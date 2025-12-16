package services

import (
    "context"
    "time"

    "github.com/kadyrbayev2005/studysync/internal/models"
    "gorm.io/gorm"
)

func StartReminderWorker(ctx context.Context, db *gorm.DB) {
    ticker := time.NewTicker(1 * time.Minute) // Check every 1 minutes
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            Info("Reminder worker stopped")
            return
        case <-ticker.C:
            now := time.Now()
            soon := now.Add(30 * time.Minute) // Remind 30 minutes before

            var due []models.Deadline
            if err := db.Preload("Task").Preload("User").
                Where("due_date > ? AND due_date <= ?", now, soon).
                Find(&due).Error; err != nil {
                Error("worker query error", "error", err)
                continue
            }

            for _, d := range due {
                if d.User.Email != "" {
                    go func(deadline models.Deadline) {
                        if err := SendDeadlineReminder(
                            deadline.User.Email,
                            deadline.Task.Title,
                            deadline.DueDate.Format("2006-01-02 15:04"),
                        ); err != nil {
                            Error("Failed to send reminder email",
                                "error", err,
                                "user_email", deadline.User.Email,
                                "task", deadline.Task.Title)
                        } else {
                            Info("Deadline reminder sent",
                                "user_email", deadline.User.Email,
                                "task_title", deadline.Task.Title,
                                "due_date", deadline.DueDate)
                        }
                    }(d)
                }
            }
        }
    }
}