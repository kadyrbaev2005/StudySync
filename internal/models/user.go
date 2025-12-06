package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	Username     string    `json:"username" gorm:"uniqueIndex"`
	PasswordHash string    `json:"-"` // never expose
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
