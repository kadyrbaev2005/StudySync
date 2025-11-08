package models

import "time"

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Deadline    time.Time `json:"deadline"`
	SubjectID   uint      `json:"subject_id"`
	Subject     Subject   `json:"subject" gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `json:"created_at"`
}
