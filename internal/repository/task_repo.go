package repository

import (
	"strings"

	"github.com/kadyrbayev2005/studysync/internal/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) GetAll() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Preload("Subject").Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) GetByID(id uint) (models.Task, error) {
	var task models.Task
	err := r.db.Preload("Subject").First(&task, id).Error
	return task, err
}

func (r *TaskRepository) Update(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Task{}).Where("id = ?", id).Updates(data).Error
}

func (r *TaskRepository) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}

func (r *TaskRepository) GetTasks(filter *TaskFilter) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	tx := r.db.Model(&models.Task{}).Preload("Subject")

	if strings.TrimSpace(filter.Status) != "" {
		tx = tx.Where("status = ?", filter.Status)
	}

	if filter.SubjectID != nil {
		tx = tx.Where("subject_id = ?", *filter.SubjectID)
	}

	if strings.TrimSpace(filter.Search) != "" {
		q := "%" + strings.ToLower(filter.Search) + "%"
		tx = tx.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", q, q)
	}

	if filter.DeadlineBefore != nil {
		tx = tx.Where("deadline <= ?", *filter.DeadlineBefore)
	}
	if filter.DeadlineAfter != nil {
		tx = tx.Where("deadline >= ?", *filter.DeadlineAfter)
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]bool{
		"created_at":      true,
		"created_at desc": true,
		"deadline":        true,
		"deadline desc":   true,
		"title":           true,
		"title desc":      true,
	}
	sort := strings.TrimSpace(filter.Sort)
	if sort == "" {
		sort = "created_at desc"
	}
	if allowedSorts[sort] {
		tx = tx.Order(sort)
	} else {
		// fallback на безопасный сорт
		tx = tx.Order("created_at desc")
	}

	// Пагинация — sane defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	// ограничение верхней границы
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	offset := (filter.Page - 1) * filter.Limit

	if err := tx.Limit(filter.Limit).Offset(offset).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}
