package models

import "time"

type DeadlineRequest struct {
	TaskID  uint      `json:"task_id" binding:"required"`
	DueDate time.Time `json:"due_date" binding:"required"`
}
