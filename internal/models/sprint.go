package models

import "time"

type Sprint struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" example:"Sprint 1"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status" example:"active"` // active, completed, planned
	Tasks     []Task    `json:"tasks" gorm:"foreignKey:SprintID"`
	CreatedAt time.Time `json:"created_at"`
}
