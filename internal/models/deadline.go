package models

import "time"

type Deadline struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    uint      `json:"task_id"`
	Task      Task      `json:"task" gorm:"constraint:OnDelete:CASCADE"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user"`
	DueDate   time.Time `json:"due_date"`
	CreatedAt time.Time `json:"created_at"`
}
