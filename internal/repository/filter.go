package repository

import "time"

type TaskFilter struct {
	Page           int
	Limit          int
	Status         string
	SubjectID      *uint
	Search         string
	Sort           string // e.g. "created_at desc"
	DeadlineBefore *time.Time
	DeadlineAfter  *time.Time
}
