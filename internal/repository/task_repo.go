package repository

import (
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
