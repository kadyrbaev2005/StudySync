package repository

import (
    "github.com/kadyrbayev2005/studysync/internal/models"
    "gorm.io/gorm"
)

type DeadlineRepository struct {
    db *gorm.DB
}

func NewDeadlineRepository(db *gorm.DB) *DeadlineRepository {
    return &DeadlineRepository{db}
}

func (r *DeadlineRepository) Create(d *models.Deadline) error {
    return r.db.Create(d).Error
}

func (r *DeadlineRepository) GetAll() ([]models.Deadline, error) {
    var ds []models.Deadline
    err := r.db.Preload("Task").Preload("User").Find(&ds).Error
    return ds, err
}

func (r *DeadlineRepository) GetByID(id uint) (models.Deadline, error) {
    var d models.Deadline
    err := r.db.Preload("Task").Preload("User").First(&d, id).Error
    return d, err
}

func (r *DeadlineRepository) Update(id uint, data map[string]interface{}) error {
    return r.db.Model(&models.Deadline{}).Where("id = ?", id).Updates(data).Error
}

func (r *DeadlineRepository) Delete(id uint) error {
    return r.db.Delete(&models.Deadline{}, id).Error
}

func (r *DeadlineRepository) GetDueBefore(t string) ([]models.Deadline, error) {
    var ds []models.Deadline
    err := r.db.Preload("Task").Preload("User").Where("due_date <= ?", t).Find(&ds).Error
    return ds, err
}